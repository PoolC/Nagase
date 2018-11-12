package models

import (
	"fmt"
	"time"

	"github.com/graphql-go/graphql"

	"nagase/components/database"
)

type Post struct {
	ID int

	BoardID    int    `gorm:"INDEX"`
	AuthorUUID string `gorm:"type:varchar(40)"`
	Title      string
	Body       string

	CreatedAt time.Time
	UpdatedAt time.Time
}

var postType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Post",
	Fields: graphql.Fields{
		"id": &graphql.Field{Type: graphql.NewNonNull(graphql.Int)},
		"author": &graphql.Field{
			Type: graphql.NewNonNull(memberType),
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				return GetMemberByUUID(params.Source.(Post).AuthorUUID)
			},
		},
		"title": &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"body":  &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"comments": &graphql.Field{
			Type: graphql.NewList(graphql.NewNonNull(commentType)),
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				var comments []Comment
				database.DB.Where(&Comment{PostID: params.Source.(Post).ID}).Find(&comments)
				return comments, nil
			},
		},
		"created_at": &graphql.Field{Type: graphql.NewNonNull(graphql.DateTime)},
		"updated_at": &graphql.Field{Type: graphql.NewNonNull(graphql.DateTime)},
	},
})

func getPosts(boardID int, before int, count int) []Post {
	if count == 0 {
		count = 20
	}

	var posts []Post
	query := database.DB.Where(Post{BoardID: boardID})
	if before != 0 {
		query = query.Where("id < ?", before)
	}
	query.Limit(count).Order("id desc").Find(&posts)

	return posts
}

// Queries
var PostQuery = &graphql.Field{
	Type:        postType,
	Description: "게시물을 조회합니다.",
	Args: graphql.FieldConfigArgument{
		"postID": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		var member *Member
		if memberCtx := params.Context.Value("member"); memberCtx != nil {
			member = memberCtx.(*Member)
		}

		// Get post
		postID, _ := params.Args["postID"].(int)
		var post Post
		database.DB.Where(&Post{ID: postID}).First(&post)
		if post.ID == 0 {
			return nil, fmt.Errorf("bad request")
		}

		// Get board and check permission
		var board Board
		database.DB.Where(&Board{ID: post.BoardID}).First(&board)
		if (board.ReadPermission != "PUBLIC" && member == nil) || (board.ReadPermission == "ADMIN" && !member.IsAdmin) {
			return nil, fmt.Errorf("forbidden")
		}

		return post, nil
	},
}

var PostsQuery = &graphql.Field{
	Type:        graphql.NewList(graphql.NewNonNull(postType)),
	Description: "게시물 목록을 조회합니다. 해당 게시판에 읽기 권한이 있어야합니다. 게시물은 ID의 내림차순으로 반환합니다.",
	Args: graphql.FieldConfigArgument{
		"boardID": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
		"before":  &graphql.ArgumentConfig{Type: graphql.Int},
		"count":   &graphql.ArgumentConfig{Type: graphql.Int},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		var member *Member
		if memberCtx := params.Context.Value("member"); memberCtx != nil {
			member = memberCtx.(*Member)
		}

		// Get board and check permission
		var board Board
		boardID, _ := params.Args["boardID"].(int)
		database.DB.Where(&Board{ID: boardID}).First(&board)
		if (board.ReadPermission != "PUBLIC" && member == nil) || (board.ReadPermission == "ADMIN" && !member.IsAdmin) {
			return nil, fmt.Errorf("forbidden")
		}

		before, _ := params.Args["before"].(int)
		count, _ := params.Args["count"].(int)

		return getPosts(boardID, before, count), nil
	},
}

// Mutations
var postInputType = &graphql.ArgumentConfig{
	Type: graphql.NewNonNull(graphql.NewInputObject(graphql.InputObjectConfig{
		Name:        "PostInput",
		Description: "게시물 작성/수정 InputObject",
		Fields: graphql.InputObjectConfigFieldMap{
			"title": &graphql.InputObjectFieldConfig{Type: graphql.String},
			"body":  &graphql.InputObjectFieldConfig{Type: graphql.String},
		},
	})),
}

var CreatePostMutation = &graphql.Field{
	Type:        postType,
	Description: "게시물을 작성합니다. 해당 게시판에 쓰기 권한이 있어야합니다.",
	Args: graphql.FieldConfigArgument{
		"boardID": &graphql.ArgumentConfig{
			Type:        graphql.NewNonNull(graphql.Int),
			Description: "게시물을 작성할 게시판의 ID",
		},
		"PostInput": postInputType,
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		if params.Context.Value("member") == nil {
			return nil, fmt.Errorf("unauthorized")
		}
		member := params.Context.Value("member").(*Member)

		// Get board and check permission
		boardID, _ := params.Args["boardID"].(int)
		board := new(Board)
		database.DB.Where(&Board{ID: boardID}).First(&board)
		if board.Name == "" {
			return nil, fmt.Errorf("bad request")
		} else if board.WritePermission == "ADMIN" && !member.IsAdmin {
			return nil, fmt.Errorf("forbidden")
		}

		// Create new post
		postInput, _ := params.Args["PostInput"].(map[string]interface{})
		post := Post{
			BoardID:    boardID,
			AuthorUUID: member.UUID,
			Title:      postInput["title"].(string),
			Body:       postInput["body"].(string),
		}

		errs := database.DB.Save(&post).GetErrors()
		if len(errs) > 0 {
			return nil, errs[0]
		}
		return &post, nil
	},
}

var UpdatePostMutation = &graphql.Field{
	Type:        postType,
	Description: "게시물을 수정합니다. 게시물의 작성자이어야 합니다.",
	Args: graphql.FieldConfigArgument{
		"postID": &graphql.ArgumentConfig{
			Type:        graphql.NewNonNull(graphql.Int),
			Description: "수정할 게시물의 ID",
		},
		"PostInput": postInputType,
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		// Check permission to the post.
		if params.Context.Value("member") == nil {
			return nil, fmt.Errorf("unauthorized")
		}
		member := params.Context.Value("member").(*Member)

		var post Post
		postID, _ := params.Args["postID"].(int)
		database.DB.Where(&Post{ID: postID}).First(&post)
		if post.ID == 0 || post.AuthorUUID != member.UUID {
			return nil, fmt.Errorf("unauthorized")
		}

		// Update the post.
		postInput, _ := params.Args["PostInput"].(map[string]interface{})
		if postInput["title"] != nil {
			post.Title = postInput["title"].(string)
		}
		if postInput["body"] != nil {
			post.Body = postInput["body"].(string)
		}
		database.DB.Save(&post)

		return post, nil
	},
}

var DeletePostMutation = &graphql.Field{
	Type:        postType,
	Description: "게시물을 삭제합니다. 게시물의 작성자이거나 관리자 권한이 필요합니다.",
	Args: graphql.FieldConfigArgument{
		"postID": &graphql.ArgumentConfig{
			Type:        graphql.NewNonNull(graphql.Int),
			Description: "삭제할 게시물의 ID",
		},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		// Check permission to the post.
		if params.Context.Value("member") == nil {
			return nil, fmt.Errorf("unauthorized")
		}
		member := params.Context.Value("member").(*Member)

		var post Post
		postID, _ := params.Args["postID"].(int)
		database.DB.Where(&Post{ID: postID}).First(&post)
		if post.ID == 0 || (!member.IsAdmin && post.AuthorUUID != member.UUID) {
			return nil, fmt.Errorf("unauthorized")
		}

		// Delete the post.
		database.DB.Delete(&post)
		return post, nil
	},
}

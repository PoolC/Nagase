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
		"title":      &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"body":       &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"comments":   &graphql.Field{Type: graphql.NewList(graphql.NewNonNull(commentType))},
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
	query.Limit(count).Find(&posts)

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
		member := params.Context.Value("member").(*Member)

		// Get board and check permission
		boardID, _ := params.Args["boardID"].(int)
		board := new(Board)
		database.DB.Where(&Board{ID: boardID}).First(&board)
		if board.Name == "" {
			return nil, fmt.Errorf("bad request")
		} else if (board.ReadPermission != "PUBLIC" && member == nil) || (board.ReadPermission == "ADMIN" && !member.IsAdmin) {
			return nil, fmt.Errorf("forbidden")
		}

		// Get post
		postID, _ := params.Args["postID"].(int)

		post := new(Post)
		database.DB.Where(&Post{ID: postID}).First(&post)
		if post.ID == 0 {
			return nil, fmt.Errorf("bad request")
		}
		return *post, nil
	},
}

var PostsQuery = &graphql.Field{
	Type:        graphql.NewList(graphql.NewNonNull(postType)),
	Description: "게시물 목록을 조회합니다.",
	Args: graphql.FieldConfigArgument{
		"boardID": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
		"before":  &graphql.ArgumentConfig{Type: graphql.Int},
		"count":   &graphql.ArgumentConfig{Type: graphql.Int},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		boardID, _ := params.Args["boardID"].(int)
		before, _ := params.Args["before"].(int)
		count, _ := params.Args["count"].(int)

		return getPosts(boardID, before, count), nil
	},
}

// Mutations
var CreatePostMutation = &graphql.Field{
	Type:        postType,
	Description: "게시글을 작성합니다.",
	Args: graphql.FieldConfigArgument{
		"boardID": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
		"PostInput": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.NewInputObject(graphql.InputObjectConfig{
				Name:        "PostInput",
				Description: "게시물 작성/수정 InputObject",
				Fields: graphql.InputObjectConfigFieldMap{
					"title": &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
					"body":  &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
				},
			})),
		},
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

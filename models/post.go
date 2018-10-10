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

type PostType struct {
	ID int `json:"int"`

	Author   MemberType
	Title    string
	Body     string
	Comments []CommentType

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (post Post) toGraphQLType() PostType {
	author, _ := GetMemberByUUID(post.AuthorUUID)

	var comments []Comment
	database.DB.Where(&Comment{PostID: post.ID}).Find(&comments)
	var commentTypes []CommentType
	for _, v := range comments {
		commentTypes = append(commentTypes, v.toGraphQLType())
	}

	return PostType{
		ID:        post.ID,
		Author:    author.toGraphQLType(),
		Title:     post.Title,
		Body:      post.Body,
		Comments:  commentTypes,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}
}

var postType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Post",
	Fields: graphql.Fields{
		"id":         &graphql.Field{Type: graphql.NewNonNull(graphql.Int)},
		"author":     &graphql.Field{Type: graphql.NewNonNull(memberType)},
		"title":      &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"body":       &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"comments":   &graphql.Field{Type: graphql.NewList(graphql.NewNonNull(commentType))},
		"created_at": &graphql.Field{Type: graphql.NewNonNull(graphql.DateTime)},
		"updated_at": &graphql.Field{Type: graphql.NewNonNull(graphql.DateTime)},
	},
})

// Queries
var PostQuery = &graphql.Field {
	Type: postType,
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
		return post.toGraphQLType(), nil
	},
}

// Mutations
var CreatePostMutation = &graphql.Field{
	Type: postType,
	Description: "게시글을 작성합니다.",
	Args: graphql.FieldConfigArgument{
		"boardID": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
		"title":   &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
		"body":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
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
		title, _ := params.Args["title"].(string)
		body, _ := params.Args["body"].(string)
		post := Post{
			BoardID:    boardID,
			AuthorUUID: member.UUID,
			Title:      title,
			Body:       body,
		}

		errs := database.DB.Save(&post).GetErrors()
		if len(errs) > 0 {
			return nil, errs[0]
		}
		return post.toGraphQLType(), nil
	},
}

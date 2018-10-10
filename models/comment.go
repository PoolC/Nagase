package models

import (
	"fmt"
	"nagase/components/database"
	"time"

	"github.com/graphql-go/graphql"
)

type Comment struct {
	ID int

	PostID     int    `gorm:"INDEX"`
	AuthorUUID string
	Body       string `gorm:"type:varchar(40)"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

type CommentType struct {
	ID        int
	Author    MemberType
	Body      string
	CreatedAt time.Time
}

func (comment Comment) toGraphQLType() CommentType {
	author, _ := GetMemberByUUID(comment.AuthorUUID)
	return CommentType{
		ID:        comment.ID,
		Author:    author.toGraphQLType(),
		Body:      comment.Body,
		CreatedAt: comment.CreatedAt,
	}
}

var commentType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Comment",
	Fields: graphql.Fields{
		"id":         &graphql.Field{Type: graphql.NewNonNull(graphql.Int)},
		"author":     &graphql.Field{Type: graphql.NewNonNull(memberType)},
		"body":       &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"created_at": &graphql.Field{Type: graphql.NewNonNull(graphql.DateTime)},
	},
})

// Mutations
var CreateCommentMutation = &graphql.Field{
	Type:        postType,
	Description: "댓글을 작성합니다.",
	Args: graphql.FieldConfigArgument{
		"postID": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
		"body":   &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		if params.Context.Value("member") == nil {
			return nil, fmt.Errorf("unauthorized")
		}
		member := params.Context.Value("member").(*Member)

		// Get board and check permission.
		// All users who has read permission can create comments.
		postID, _ := params.Args["postID"].(int)
		post := new(Post)
		database.DB.Where(&Post{ID: postID}).First(&post)
		if post.ID == 0 {
			return nil, fmt.Errorf("bad request")
		}

		board := new(Board)
		database.DB.Where(&Board{ID: post.BoardID}).First(&board)
		if board.ID == 0 {
			return nil, fmt.Errorf("bad request")
		} else if board.ReadPermission == "ADMIN" && !member.IsAdmin {
			return nil, fmt.Errorf("forbidden")
		}

		// Create new comment
		body, _ := params.Args["body"].(string)
		comment := Comment{
			PostID:     postID,
			AuthorUUID: member.UUID,
			Body:       body,
		}

		errs := database.DB.Save(&comment).GetErrors()
		if len(errs) > 0 {
			return nil, errs[0]
		}
		return comment.toGraphQLType(), nil
	},
}

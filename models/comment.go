package models

import (
	"fmt"
	"nagase/components/database"
	"time"

	"github.com/graphql-go/graphql"
)

type Comment struct {
	ID int

	PostID     int `gorm:"INDEX"`
	AuthorUUID string
	Body       string `gorm:"type:varchar(40)"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

var commentType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Comment",
	Fields: graphql.Fields{
		"id": &graphql.Field{Type: graphql.NewNonNull(graphql.Int)},
		"author": &graphql.Field{
			Type: graphql.NewNonNull(memberType),
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				return GetMemberByUUID(params.Source.(Comment).AuthorUUID)
			},
		},
		"body":      &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"createdAt": &graphql.Field{Type: graphql.NewNonNull(graphql.DateTime)},
	},
})

// Mutations
var CreateCommentMutation = &graphql.Field{
	Type:        commentType,
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
		return comment, nil
	},
}

var DeleteCommentMutation = &graphql.Field{
	Type:        commentType,
	Description: "댓글을 삭제합니다. 작성자 본인 또는 관리자만 댓글을 삭제할 수 있습니다.",
	Args: graphql.FieldConfigArgument{
		"commentID": &graphql.ArgumentConfig{
			Type:        graphql.NewNonNull(graphql.Int),
			Description: "삭제할 댓글의 ID",
		},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		if params.Context.Value("member") == nil {
			return nil, fmt.Errorf("unauthorized")
		}
		member := params.Context.Value("member").(*Member)

		// Get comment and check permission.
		var comment Comment
		commentID, _ := params.Args["commentID"].(int)
		database.DB.Where(&Comment{ID: commentID}).First(&comment)
		if comment.ID == 0 {
			return nil, fmt.Errorf("bad request")
		} else if comment.AuthorUUID != member.UUID && !member.IsAdmin {
			return nil, fmt.Errorf("forbidden")
		}

		// Delete the comment.
		database.DB.Delete(&comment)
		return comment, nil
	},
}

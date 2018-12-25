package models

import (
	"fmt"
	"strconv"
	"time"

	"github.com/graphql-go/graphql"

	"nagase/components/database"
	"nagase/components/push"
)

type Comment struct {
	ID int

	PostID     int    `gorm:"INDEX"`
	AuthorUUID string `gorm:"type:varchar(40)"`
	Body       string

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
			return nil, fmt.Errorf("ERR401")
		}
		member := params.Context.Value("member").(*Member)

		// 권한을 확인합니다.
		// 해당 게시판에 읽기 권한이 있는 모든 사용자는 댓글을 달 수 있습니다.
		postID, _ := params.Args["postID"].(int)
		post := new(Post)
		database.DB.Where(&Post{ID: postID}).First(&post)
		if post.ID == 0 {
			return nil, fmt.Errorf("ERR400")
		}

		board := new(Board)
		database.DB.Where(&Board{ID: post.BoardID}).First(&board)
		if board.ID == 0 {
			return nil, fmt.Errorf("ERR400")
		} else if board.ReadPermission == "ADMIN" && !member.IsAdmin {
			return nil, fmt.Errorf("ERR403")
		}

		// 댓글을 저장합니다.
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

		// 댓글이 작성된 게시물을 구독하고 있는 유저들에게 푸시를 발송합니다.
		data := make(map[string]string)
		data["boardID"] = strconv.Itoa(board.ID)
		data["postID"] = strconv.Itoa(post.ID)

		var subscriptions []PostSubscription
		database.DB.Where(&PostSubscription{PostID: postID}).Find(&subscriptions)
		for _, s := range subscriptions {
			title := member.Name + " 님이 게시물에 댓글을 남겼습니다."
			body := comment.Body
			go push.SendPush(s.MemberUUID, title, body, data)
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
			return nil, fmt.Errorf("ERR401")
		}
		member := params.Context.Value("member").(*Member)

		// Get comment and check permission.
		var comment Comment
		commentID, _ := params.Args["commentID"].(int)
		database.DB.Where(&Comment{ID: commentID}).First(&comment)
		if comment.ID == 0 {
			return nil, fmt.Errorf("ERR400")
		} else if comment.AuthorUUID != member.UUID && !member.IsAdmin {
			return nil, fmt.Errorf("ERR403")
		}

		// Delete the comment.
		database.DB.Delete(&comment)
		return comment, nil
	},
}

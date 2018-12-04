package models

import (
	"fmt"
	"math"
	"sort"
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
	VoteID     *int

	CreatedAt time.Time
	UpdatedAt time.Time
}

type PostPage struct {
	Posts    []Post
	PageInfo PageInfo
}

var postType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Post",
	Fields: graphql.Fields{
		// To ignore circular dependency, `board` field has initialized lazily on init() method.
		"id": &graphql.Field{Type: graphql.NewNonNull(graphql.Int)},
		"author": &graphql.Field{
			Type: graphql.NewNonNull(memberType),
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				return GetMemberByUUID(params.Source.(Post).AuthorUUID)
			},
		},
		"title": &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"body":  &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"vote": &graphql.Field{
			Type: voteType,
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				if params.Source.(Post).VoteID == nil {
					return nil, nil
				}

				var vote Vote
				database.DB.Where(&Vote{ID: *params.Source.(Post).VoteID}).First(&vote)
				return vote, nil
			},
		},
		"comments": &graphql.Field{
			Type: graphql.NewList(graphql.NewNonNull(commentType)),
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				var comments []Comment
				database.DB.Where(&Comment{PostID: params.Source.(Post).ID}).Order("id asc").Find(&comments)
				return comments, nil
			},
		},
		"createdAt": &graphql.Field{Type: graphql.NewNonNull(graphql.DateTime)},
		"updatedAt": &graphql.Field{Type: graphql.NewNonNull(graphql.DateTime)},
	},
})

var postPageType = graphql.NewObject(graphql.ObjectConfig{
	Name: "PostPage",
	Fields: graphql.Fields{
		"posts":    &graphql.Field{Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(postType)))},
		"pageInfo": &graphql.Field{Type: graphql.NewNonNull(pageInfoType)},
	},
})

func getPostPage(boardID int, pagination *Pagination) PostPage {
	count := 20
	if pagination.Count != 0 {
		count = pagination.Count
	}

	var posts []Post
	query := database.DB.Where(Post{BoardID: boardID})
	if pagination.Before != 0 {
		query = query.Where("id < ?", pagination.Before).Order("id desc")
	} else if pagination.After != 0 {
		query = query.Where("id > ?", pagination.After).Order("id asc")
	} else {
		query = query.Order("id desc")
	}
	query.Limit(count).Find(&posts)
	sort.SliceStable(posts, func(i, j int) bool { return posts[i].ID > posts[j].ID })

	maxID := math.MinInt32
	minID := math.MaxInt32
	for _, p := range posts {
		if p.ID > maxID {
			maxID = p.ID
		}
		if p.ID < minID {
			minID = p.ID
		}
	}

	var prevCount int
	var nextCount int
	database.DB.Model(&Post{}).Where("id > ?", maxID).Count(&prevCount)
	database.DB.Model(&Post{}).Where("id < ?", minID).Count(&nextCount)
	return PostPage{
		Posts: posts,
		PageInfo: PageInfo{
			HasPrevious: prevCount > 0,
			HasNext:     nextCount > 0,
		},
	}
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
			return nil, fmt.Errorf("ERR400")
		}

		// Get board and check permission
		var board Board
		database.DB.Where(&Board{ID: post.BoardID}).First(&board)
		if (board.ReadPermission != "PUBLIC" && member == nil) || (board.ReadPermission == "ADMIN" && !member.IsAdmin) {
			return nil, fmt.Errorf("ERR403")
		}

		return post, nil
	},
}

var PostPageQuery = &graphql.Field{
	Type:        postPageType,
	Description: "게시물 목록을 조회합니다. 해당 게시판에 읽기 권한이 있어야합니다. 게시물은 ID의 내림차순으로 반환합니다.",
	Args: graphql.FieldConfigArgument{
		"boardID": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
		"before":  &graphql.ArgumentConfig{Type: graphql.Int},
		"after":   &graphql.ArgumentConfig{Type: graphql.Int},
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
			return nil, fmt.Errorf("ERR403")
		}

		return getPostPage(boardID, getPaginationFromGraphQLParams(&params)), nil
	},
}

// Mutations
var postInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name:        "PostInput",
	Description: "게시물 작성/수정 InputObject",
	Fields: graphql.InputObjectConfigFieldMap{
		"title": &graphql.InputObjectFieldConfig{Type: graphql.String},
		"body":  &graphql.InputObjectFieldConfig{Type: graphql.String},
	},
})

var CreatePostMutation = &graphql.Field{
	Type:        postType,
	Description: "게시물을 작성합니다. 해당 게시판에 쓰기 권한이 있어야합니다.",
	Args: graphql.FieldConfigArgument{
		"boardID": &graphql.ArgumentConfig{
			Type:        graphql.NewNonNull(graphql.Int),
			Description: "게시물을 작성할 게시판의 ID",
		},
		"PostInput": &graphql.ArgumentConfig{Type: graphql.NewNonNull(postInputType)},
		"VoteInput": &graphql.ArgumentConfig{Type: voteInputType},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		if params.Context.Value("member") == nil {
			return nil, fmt.Errorf("ERR401")
		}
		member := params.Context.Value("member").(*Member)

		// Get board and check permission
		boardID, _ := params.Args["boardID"].(int)
		board := new(Board)
		database.DB.Where(&Board{ID: boardID}).First(&board)
		if board.Name == "" {
			return nil, fmt.Errorf("ERR400")
		} else if board.WritePermission == "ADMIN" && !member.IsAdmin {
			return nil, fmt.Errorf("ERR403")
		}

		// Create new post
		postInput, _ := params.Args["PostInput"].(map[string]interface{})
		post := Post{
			BoardID:    boardID,
			AuthorUUID: member.UUID,
			Title:      postInput["title"].(string),
			Body:       postInput["body"].(string),
		}

		// Create and save new vote (if requested)
		var vote Vote
		if params.Args["VoteInput"] != nil {
			voteInput, _ := params.Args["VoteInput"].(map[string]interface{})
			isMultipleSelectable := false
			if voteInput["isMultipleSelectable"] != nil {
				isMultipleSelectable = voteInput["isMultipleSelectable"].(bool)
			}
			deadline, err := time.Parse(time.RFC3339, voteInput["deadline"].(string))
			if err != nil {
				return nil, fmt.Errorf("ERR400")
			}

			vote = Vote{
				Title:                voteInput["title"].(string),
				IsMultipleSelectable: isMultipleSelectable,
				Deadline:             deadline,
			}
			errs := database.DB.Save(&vote).GetErrors()
			if len(errs) > 0 {
				return nil, errs[0]
			}

			// Create and save vote options.
			for _, t := range voteInput["optionTexts"].([]interface{}) {
				errs = database.DB.Save(&VoteOption{VoteID: vote.ID, Text: t.(string)}).GetErrors()
				if len(errs) > 0 {
					return nil, errs[0]
				}
			}

			post.VoteID = &vote.ID
		}

		errs := database.DB.Save(&post).GetErrors()
		if len(errs) > 0 {
			return nil, errs[0]
		}

		return post, nil
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
		"PostInput": &graphql.ArgumentConfig{Type: graphql.NewNonNull(postInputType)},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		// Check permission to the post.
		if params.Context.Value("member") == nil {
			return nil, fmt.Errorf("ERR401")
		}
		member := params.Context.Value("member").(*Member)

		var post Post
		postID, _ := params.Args["postID"].(int)
		database.DB.Where(&Post{ID: postID}).First(&post)
		if post.ID == 0 || post.AuthorUUID != member.UUID {
			return nil, fmt.Errorf("ERR401")
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
			return nil, fmt.Errorf("ERR401")
		}
		member := params.Context.Value("member").(*Member)

		var post Post
		postID, _ := params.Args["postID"].(int)
		database.DB.Where(&Post{ID: postID}).First(&post)
		if post.ID == 0 || (!member.IsAdmin && post.AuthorUUID != member.UUID) {
			return nil, fmt.Errorf("ERR401")
		}

		// Delete attached vote if exists.
		if post.VoteID != nil {
			database.DB.Where(&Vote{ID: *post.VoteID}).Delete(Vote{})
		}

		// Delete the post.
		database.DB.Delete(&post)
		return post, nil
	},
}

func init() {
	postType.AddFieldConfig("board", &graphql.Field{
		Type: graphql.NewNonNull(boardType),
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			var board Board
			database.DB.Where(&Board{ID: params.Source.(Post).BoardID}).First(&board)
			return board, nil
		},
	})
}

package models

import (
	"fmt"
	"time"

	"github.com/graphql-go/graphql"

	"nagase/components/database"
)

type Board struct {
	ID              int
	Name            string `gorm:"type:varchar(40);UNIQUE_INDEX"`
	URLPath         string `gorm:"type:varchar(40);UNIQUE_INDEX"`
	ReadPermission  string `gorm:"type:varchar(10)"`
	WritePermission string `gorm:"type:varchar(10)"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

var boardType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Board",
	Fields: graphql.Fields{
		"id":              &graphql.Field{Type: graphql.NewNonNull(graphql.Int)},
		"name":            &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"urlPath":         &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"readPermission":  &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"writePermission": &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"postPage": &graphql.Field{
			Type: graphql.NewList(graphql.NewNonNull(postType)),
			Args: graphql.FieldConfigArgument{
				"before": &graphql.ArgumentConfig{Type: graphql.Int},
				"after":  &graphql.ArgumentConfig{Type: graphql.Int},
				"count":  &graphql.ArgumentConfig{Type: graphql.Int},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				return getPostPage(params.Source.(Board).ID, getPaginationFromGraphQLParams(&params)), nil
			},
		},
	},
})

var boardInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name:        "BoardInput",
	Description: "게시판 추가/수정 InputObject",
	Fields: graphql.InputObjectConfigFieldMap{
		"name":            &graphql.InputObjectFieldConfig{Type: graphql.String},
		"urlPath":         &graphql.InputObjectFieldConfig{Type: graphql.String},
		"readPermission":  &graphql.InputObjectFieldConfig{Type: graphql.String},
		"writePermission": &graphql.InputObjectFieldConfig{Type: graphql.String},
	},
})

// Queries
var BoardQuery = &graphql.Field{
	Type:        boardType,
	Description: "게시판을 조회합니다.",
	Args: graphql.FieldConfigArgument{
		"boardID": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		var board Board
		database.DB.Where(&Board{ID: params.Args["boardID"].(int)}).First(&board)
		return board, nil
	},
}

var BoardsQuery = &graphql.Field{
	Type:        graphql.NewList(boardType),
	Description: "게시판 목록을 조회합니다",
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		var boards []Board
		database.DB.Order("id asc").Find(&boards)
		return boards, nil
	},
}

// Mutations
var CreateBoardMutation = &graphql.Field{
	Type:        boardType,
	Description: "게시판을 추가합니다. 관리자 권한이 필요합니다.",
	Args: graphql.FieldConfigArgument{
		"BoardInput": &graphql.ArgumentConfig{Type: boardInputType},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		if member := params.Context.Value("member"); member == nil || !member.(*Member).IsAdmin {
			return nil, fmt.Errorf("ERR401")
		}

		boardInput, _ := params.Args["BoardInput"].(map[string]interface{})
		board := Board{
			Name:            boardInput["name"].(string),
			URLPath:         boardInput["urlPath"].(string),
			ReadPermission:  boardInput["readPermission"].(string),
			WritePermission: boardInput["writePermission"].(string),
		}
		errs := database.DB.Save(&board).GetErrors()
		if len(errs) > 0 {
			return nil, errs[0]
		}

		return board, nil
	},
}

var UpdateBoardMutation = &graphql.Field{
	Type:        boardType,
	Description: "게시판을 수정합니다. 관리자 권한이 필요합니다.",
	Args: graphql.FieldConfigArgument{
		"boardID":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
		"BoardInput": &graphql.ArgumentConfig{Type: boardInputType},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		if member := params.Context.Value("member"); member == nil || !member.(*Member).IsAdmin {
			return nil, fmt.Errorf("ERR401")
		}

		var board Board
		database.DB.Where(&Board{ID: params.Args["boardID"].(int)}).First(&board)
		fmt.Println(board)
		if board.ID == 0 {
			return nil, fmt.Errorf("ERR400")
		}

		boardInput, _ := params.Args["BoardInput"].(map[string]interface{})
		if boardInput["name"] != nil {
			board.Name = boardInput["name"].(string)
		}
		if boardInput["urlPath"] != nil {
			board.URLPath = boardInput["urlPath"].(string)
		}
		if boardInput["readPermission"] != nil {
			board.ReadPermission = boardInput["readPermission"].(string)
		}
		if boardInput["writePermission"] != nil {
			board.WritePermission = boardInput["writePermission"].(string)
		}

		errs := database.DB.Save(&board).GetErrors()
		if len(errs) > 0 {
			return nil, errs[0]
		}

		return board, nil
	},
}

var DeleteBoardMutation = &graphql.Field{
	Type:        boardType,
	Description: "게시판을 삭제합니다. 관리자 권한이 필요합니다.",
	Args: graphql.FieldConfigArgument{
		"boardID": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		if member := params.Context.Value("member"); member == nil || !member.(*Member).IsAdmin {
			return nil, fmt.Errorf("ERR401")
		}

		var board Board
		database.DB.Where(&Board{ID: params.Args["boardID"].(int)}).First(&board)
		if board.ID == 0 {
			return nil, fmt.Errorf("ERR400")
		}

		errs := database.DB.Delete(&board).GetErrors()
		if len(errs) > 0 {
			return nil, errs[0]
		}

		return board, nil
	},
}

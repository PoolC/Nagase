package main

import (
	"fmt"
	"nagase"
	"nagase/components/database"

	"github.com/graphql-go/graphql"
)

var boardType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Board",
	Fields: graphql.Fields{
		"id":              &graphql.Field{Type: graphql.NewNonNull(graphql.Int)},
		"name":            &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"urlPath":         &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"readPermission":  &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"writePermission": &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"posts":           postsQuery,
		"isSubscribed":    nil, // TODO - implement
	},
})

var boardQuery = &graphql.Field{
	Type:        boardType,
	Description: "게시판을 조회합니다.",
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		var board nagase.Board
		database.DB.Where(&nagase.Board{ID: p.Args["id"].(int)}).First(&board)
		if board.ID != p.Args["id"].(int) {
			return nil, nil
		}
		return &board, nil
	},
}

var boardsQuery = &graphql.Field{
	Type:        graphql.NewList(graphql.NewNonNull(boardType)),
	Description: "게시판 목록을 조회합니다",
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		var boards []*nagase.Board
		database.DB.Order("id asc").Find(&boards)
		return boards, nil
	},
}

var createBoardMutation = &graphql.Field{
	Type:        boardType,
	Description: "게시판을 수정합니다. 관리자 권한이 필요합니다.",
	Args: graphql.FieldConfigArgument{
		"input": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.NewInputObject(graphql.InputObjectConfig{
				Name:        "CreateBoardInput",
				Description: "게시판 수정 InputObject",
				Fields: graphql.InputObjectConfigFieldMap{
					"name":            &graphql.InputObjectFieldConfig{Type: graphql.String},
					"urlPath":         &graphql.InputObjectFieldConfig{Type: graphql.String},
					"readPermission":  &graphql.InputObjectFieldConfig{Type: graphql.String},
					"writePermission": &graphql.InputObjectFieldConfig{Type: graphql.String},
				},
			})),
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		input := p.Args["input"].(map[string]interface{})
		board := nagase.Board{
			Name:            input["name"].(string),
			URLPath:         input["urlPath"].(string),
			ReadPermission:  input["readPermission"].(string),
			WritePermission: input["writePermission"].(string),
		}
		database.DB.Save(&board)
		return &board, nil
	},
}

var updateBoardMutation = &graphql.Field{
	Type:        boardType,
	Description: "게시판을 수정합니다. 관리자 권한이 필요합니다.",
	Args: graphql.FieldConfigArgument{
		"input": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.NewInputObject(graphql.InputObjectConfig{
				Name:        "UpdateBoardInput",
				Description: "게시판 수정 InputObject",
				Fields: graphql.InputObjectConfigFieldMap{
					"id":              &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.Int)},
					"name":            &graphql.InputObjectFieldConfig{Type: graphql.String},
					"urlPath":         &graphql.InputObjectFieldConfig{Type: graphql.String},
					"readPermission":  &graphql.InputObjectFieldConfig{Type: graphql.String},
					"writePermission": &graphql.InputObjectFieldConfig{Type: graphql.String},
				},
			})),
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		input := p.Args["input"].(map[string]interface{})

		var board nagase.Board
		database.DB.Where(&nagase.Board{ID: input["id"].(int)}).First(&board)
		if board.ID == 0 {
			return nil, fmt.Errorf("ERR400")
		}

		if input["name"] != nil {
			board.Name = input["name"].(string)
		}
		if input["urlPath"] != nil {
			board.URLPath = input["urlPath"].(string)
		}
		if input["readPermission"] != nil {
			board.ReadPermission = input["readPermission"].(string)
		}
		if input["writePermission"] != nil {
			board.WritePermission = input["writePermission"].(string)
		}

		database.DB.Save(&board)
		return &board, nil
	},
}

var deleteBoardMutation = &graphql.Field{
	Type:        boardType,
	Description: "게시판을 삭제합니다. 관리자 권한이 필요합니다.",
	Args: graphql.FieldConfigArgument{
		"input": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.NewInputObject(graphql.InputObjectConfig{
				Name:        "DeleteBoardInput",
				Description: "게시판 수정 InputObject",
				Fields: graphql.InputObjectConfigFieldMap{
					"id": &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.Int)},
				},
			})),
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		input := p.Args["input"].(map[string]interface{})

		var board nagase.Board
		database.DB.Where(&nagase.Board{ID: input["id"].(int)}).First(&board)
		if board.ID == 0 {
			return nil, fmt.Errorf("ERR400")
		}

		database.DB.Delete(&board)
		return &board, nil
	},
}

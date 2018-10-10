package models

import (
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

type BoardType struct {
	Name            string
	URLPath         string
	ReadPermission  string
	WritePermission string
}

func (board Board) toGraphQLType() BoardType {
	return BoardType{
		Name:            board.Name,
		URLPath:         board.URLPath,
		ReadPermission:  board.ReadPermission,
		WritePermission: board.WritePermission,
	}
}

var boardType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Board",
	Fields: graphql.Fields{
		"name":            &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"urlPath":         &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"readPermission":  &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"writePermission": &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
	},
})

// Queries
var BoardsQuery = &graphql.Field{
	Type: graphql.NewList(boardType),
	Description: "게시판 목록을 조회합니다",
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		var boards []Board
		database.DB.Find(&boards)

		var boardTypes []BoardType
		for _, b := range boards {
			boardTypes = append(boardTypes, b.toGraphQLType())
		}

		return boardTypes, nil
	},
}

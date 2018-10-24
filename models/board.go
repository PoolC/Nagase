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

var boardType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Board",
	Fields: graphql.Fields{
		"id":              &graphql.Field{Type: graphql.NewNonNull(graphql.Int)},
		"name":            &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"urlPath":         &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"readPermission":  &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"writePermission": &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"posts": &graphql.Field{
			Type: graphql.NewList(graphql.NewNonNull(postType)),
			Args: graphql.FieldConfigArgument{
				"before": &graphql.ArgumentConfig{Type: graphql.Int},
				"count":  &graphql.ArgumentConfig{Type: graphql.Int},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				before, _ := params.Args["before"].(int)
				count, _ := params.Args["count"].(int)

				return getPosts(params.Source.(Board).ID, before, count), nil
			},
		},
	},
})

// Queries
var BoardsQuery = &graphql.Field{
	Type:        graphql.NewList(boardType),
	Description: "게시판 목록을 조회합니다",
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		var boards []Board
		database.DB.Find(&boards)
		return boards, nil
	},
}

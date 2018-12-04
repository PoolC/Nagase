package models

import (
	"nagase/components/database"

	"github.com/graphql-go/graphql"
)

/// 페이지네이션과 관련된 type 및 함수.
type Pagination struct {
	Before int
	After  int
	Count  int
}

type PageInfo struct {
	HasPrevious bool
	HasNext     bool
}

var pageInfoType = graphql.NewObject(graphql.ObjectConfig{
	Name: "PageInfo",
	Fields: graphql.Fields{
		"hasPrevious": &graphql.Field{Type: graphql.NewNonNull(graphql.Boolean)},
		"hasNext":     &graphql.Field{Type: graphql.NewNonNull(graphql.Boolean)},
	},
})

func getPaginationFromGraphQLParams(params *graphql.ResolveParams) *Pagination {
	var pagination Pagination
	if params.Args["before"] != nil {
		before, _ := params.Args["before"].(int)
		pagination.Before = before
	}
	if params.Args["after"] != nil {
		after, _ := params.Args["after"].(int)
		pagination.After = after
	}
	if params.Args["count"] != nil {
		count, _ := params.Args["count"].(int)
		pagination.Count = count
	}
	return &pagination
}

func init() {
	database.DB.AutoMigrate(
		&Board{},
		&Member{},
		&Post{},
		&Comment{},
		&Vote{},
		&VoteOption{},
		&VoteSelection{},
		&Project{},
	)
}

package main

import (
	"encoding/base64"
	"nagase/components/database"
	"strconv"

	"github.com/graphql-go/graphql"

	"nagase"
)

var postType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Post",
	Fields: graphql.Fields{
		"id": &graphql.Field{Type: graphql.NewNonNull(graphql.Int)},
		"author": &graphql.Field{
			Type: graphql.NewNonNull(memberType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				var member nagase.Member
				database.DB.Where(&nagase.Member{UUID: p.Source.(*nagase.Post).AuthorUUID}).First(&member)
				return &member, nil
			},
		},
		"title":     &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"body":      &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"createdAt": &graphql.Field{Type: graphql.NewNonNull(graphql.DateTime)},
	},
})

var pageInfoType = graphql.NewObject(graphql.ObjectConfig{
	Name: "PageInfo",
	Fields: graphql.Fields{
		"startCursor": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				posts := p.Source.([]*nagase.Post)
				if len(posts) == 0 {
					return nil, nil
				}
				return base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(posts[0].ID))), nil
			},
		},
		"endCursor": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				posts := p.Source.([]*nagase.Post)
				if len(posts) == 0 {
					return nil, nil
				}
				return base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(posts[len(posts)-1].ID))), nil
			},
		},
		"hasPreviousPage": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Boolean),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				posts := p.Source.([]*nagase.Post)
				if len(posts) == 0 {
					return false, nil
				}

				startPost := posts[0]
				var count int
				database.DB.Model(&nagase.Post{}).Where("id > ?", startPost.ID).Count(&count)
				return count > 0, nil
			},
		},
		"hasNextPage": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Boolean),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				posts := p.Source.([]*nagase.Post)
				if len(posts) == 0 {
					return false, nil
				}

				endPost := posts[len(posts)-1]
				var count int
				database.DB.Model(&nagase.Post{}).Where("id < ?", endPost.ID).Count(&count)
				return count > 0, nil
			},
		},
	},
})

var postEdgeType = graphql.NewObject(graphql.ObjectConfig{
	Name: "PostEdge",
	Fields: graphql.Fields{
		"cursor": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id := p.Source.(*nagase.Post).ID
				return base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(id))), nil
			},
		},
		"node": &graphql.Field{
			Type: graphql.NewNonNull(postType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source, nil
			},
		},
	},
})

var postConnectionType = graphql.NewObject(graphql.ObjectConfig{
	Name: "PostConnection",
	Fields: graphql.Fields{
		"edges": &graphql.Field{
			Type: graphql.NewList(graphql.NewNonNull(postEdgeType)),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source, nil
			},
		},
		"pageInfo": &graphql.Field{
			Type: graphql.NewNonNull(pageInfoType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source, nil
			},
		},
		"totalCount": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Int),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				var count int
				database.DB.Model(&nagase.Post{}).Count(&count)
				return &count, nil
			},
		},
	},
})

var postsQuery = &graphql.Field{
	Type:        graphql.NewNonNull(postConnectionType),
	Description: "게시글 목록을 조회합니다.",
	Args: graphql.FieldConfigArgument{
		"after":  &graphql.ArgumentConfig{Type: graphql.String},
		"before": &graphql.ArgumentConfig{Type: graphql.String},
		"limit":  &graphql.ArgumentConfig{Type: graphql.Int},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		board := p.Source.(*nagase.Board)

		// Resolve pagination parameters
		afterCursor := ""
		if afterArgs, ok := p.Args["after"].(string); ok {
			afterCursor = afterArgs
		}
		beforeCursor := ""
		if beforeArgs, ok := p.Args["before"].(string); ok {
			beforeCursor = beforeArgs
		}
		limit := 20
		if argsLimit, ok := p.Args["limit"].(int); ok {
			limit = argsLimit
		}

		query := database.DB.Where(&nagase.Post{BoardID: board.ID})
		if afterCursor != "" {
			id, err := base64.StdEncoding.DecodeString(afterCursor)
			if err != nil {
				return nil, err
			}
			query = query.Where("id < ?", id)
		}
		if beforeCursor != "" {
			id, err := base64.StdEncoding.DecodeString(beforeCursor)
			if err != nil {
				return nil, err
			}
			query = query.Where("id > ?", id)
		}

		var posts []*nagase.Post
		query.Order("id desc").Limit(limit).Find(&posts)
		return posts, nil
	},
}

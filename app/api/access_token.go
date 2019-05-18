package main

import (
	"fmt"

	"github.com/graphql-go/graphql"

	"nagase"
	"nagase/components/database"
)

var accessTokenType = graphql.NewObject(graphql.ObjectConfig{
	Name: "AccessToken",
	Fields: graphql.Fields{
		"key": &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
	},
})

var createAccessTokenMutation = &graphql.Field{
	Type:        graphql.NewNonNull(accessTokenType),
	Description: "Access Token을 발급합니다.",
	Args: graphql.FieldConfigArgument{
		"input": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.NewInputObject(graphql.InputObjectConfig{
				Name:        "LoginInput",
				Description: "로그인 정보 InputObject",
				Fields: graphql.InputObjectConfigFieldMap{
					"loginID":  &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
					"password": &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
				},
			})),
		},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		input := params.Args["input"].(map[string]interface{})

		// Get the member by login id and password.
		var member nagase.Member
		database.DB.Where(&nagase.Member{LoginID: input["loginID"].(string)}).First(&member)
		if member.UUID == "" || !member.ValidatePassword(input["password"].(string)) {
			return nil, fmt.Errorf("TKN000")
		}
		if !member.IsActivated {
			return nil, fmt.Errorf("TKN002")
		}

		// Generate new token and return.
		key, err := member.GenerateAccessToken()
		if err != nil {
			return nil, fmt.Errorf("ERR500")
		}

		return map[string]string{"key": key}, nil
	},
}

var refreshAccessTokenMutation = &graphql.Field{
	Type:        graphql.NewNonNull(accessTokenType),
	Description: "Access Token을 갱신합니다.",
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		if params.Context.Value("member") == nil {
			return nil, fmt.Errorf("ERR401")
		}
		member := params.Context.Value("member").(*nagase.Member)

		// Generate new token and return.
		key, err := member.GenerateAccessToken()
		if err != nil {
			return nil, fmt.Errorf("ERR500")
		}

		return map[string]string{"key": key}, nil
	},
}

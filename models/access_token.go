package models

import (
	"fmt"

	"github.com/graphql-go/graphql"

	"nagase/components/auth"
	"nagase/components/database"
)

type AccessToken struct {
	Key string `json:"key"`
}

var accessTokenType = graphql.NewObject(graphql.ObjectConfig{
	Name: "AccessToken",
	Fields: graphql.Fields{
		"key": &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
	},
})

var CreateAccessTokenMutation = &graphql.Field{
	Type:        accessTokenType,
	Description: "Access Token을 발급합니다.",
	Args: graphql.FieldConfigArgument{
		"LoginInput": &graphql.ArgumentConfig{
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
		loginInput := params.Args["LoginInput"].(map[string]interface{})

		// Get the member by login id and password.
		var member Member
		database.DB.Where(&Member{LoginID: loginInput["loginID"].(string)}).First(&member)
		if member.UUID == "" || !member.ValidatePassword(loginInput["password"].(string)) || !member.IsActivated {
			return nil, fmt.Errorf("TKN000")
		}

		// Generate new token and return.
		key, err := auth.GenerateToken(member.UUID, member.IsAdmin)
		if err != nil {
			return nil, fmt.Errorf("ERR500")
		}

		return AccessToken{Key: key}, nil
	},
}

var RefreshAccessTokenMutation = &graphql.Field{
	Type:        accessTokenType,
	Description: "Access Token을 갱신합니다.",
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		if params.Context.Value("member") == nil {
			return nil, fmt.Errorf("ERR401")
		}
		member := params.Context.Value("member").(*Member)

		// Generate new token and return.
		key, err := auth.GenerateToken(member.UUID, member.IsAdmin)
		if err != nil {
			return nil, fmt.Errorf("ERR500")
		}

		return AccessToken{Key: key}, nil
	},
}

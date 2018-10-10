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
		"loginID":  &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
		"password": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		loginID, _ := params.Args["loginID"].(string)
		password, _ := params.Args["password"].(string)

		// Get the member by login id and password.
		member := new(Member)
		database.DB.Where(&Member{LoginID: loginID}).First(&member)
		if member.UUID == "" || !member.ValidatePassword(password) || !member.IsActivated {
			return nil, fmt.Errorf("invalid login id or password")
		}

		// If token not exists or expired, generate new token.
		key, err := auth.GenerateToken(member.UUID, member.IsAdmin)
		if err != nil {
			return nil, fmt.Errorf("unknown error")
		}
		
		return AccessToken{Key: key}, nil
	},
}

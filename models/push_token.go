package models

import (
	"fmt"

	"github.com/graphql-go/graphql"

	"nagase/components/push"
)

type PushToken struct {
	MemberUUID string
	Token      string
}

var pushTokenType = graphql.NewObject(graphql.ObjectConfig{
	Name: "PushToken",
	Fields: graphql.Fields{
		"memberUUID": &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"token":      &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
	},
})

// Mutations
var RegisterPushTokenMutation = &graphql.Field{
	Type:        pushTokenType,
	Description: "푸시 토큰을 등록합니다.",
	Args: graphql.FieldConfigArgument{
		"token": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		if params.Context.Value("member") == nil {
			return nil, fmt.Errorf("ERR401")
		}
		member := params.Context.Value("member").(*Member)

		err := push.RegisterToken(member.UUID, params.Args["token"].(string))
		if err != nil {
			return nil, err
		}

		token := PushToken{MemberUUID: member.UUID, Token: params.Args["token"].(string)}
		return token, nil
	},
}

var DeregisterPushTokenMutation = &graphql.Field{
	Type:        pushTokenType,
	Description: "푸시 토큰을 해제합니다.",
	Args: graphql.FieldConfigArgument{
		"token": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		if params.Context.Value("member") == nil {
			return nil, fmt.Errorf("ERR401")
		}
		member := params.Context.Value("member").(*Member)

		err := push.DeregisterToken(member.UUID, params.Args["token"].(string))
		if err != nil {
			return nil, err
		}

		token := PushToken{MemberUUID: member.UUID, Token: params.Args["token"].(string)}
		return token, nil
	},
}

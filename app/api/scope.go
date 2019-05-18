package main

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/graphql-go/graphql"

	"nagase"
	"nagase/components/database"
)

type scope int8
type request int8

const (
	ctxScopeAdmin scope = iota
	ctxScopeMember

	ctxHTTPReq request = iota
)

var (
	errForbidden = errors.New("ERR401")
)

func withAdminScope(field *graphql.Field) *graphql.Field {
	resolver := field.Resolve
	field.Resolve = func(p graphql.ResolveParams) (interface{}, error) {
		req := p.Context.Value(ctxHTTPReq).(*http.Request)

		authorization := req.Header.Get("Authorization")
		if authorization == "" {
			return nil, errForbidden
		}

		memberUUID, err := nagase.GetMemberUUIDFromToken(strings.Split(authorization, " ")[1])
		if err != nil {
			return nil, errForbidden
		}

		var member nagase.Member
		database.DB.Where(&nagase.Member{UUID: memberUUID}).First(&member)
		if member.UUID == "" || !member.IsActivated || !member.IsAdmin {
			return nil, errForbidden
		}

		p.Context = context.WithValue(p.Context, ctxScopeAdmin, &member)
		return resolver(p)
	}

	return field
}

func withMemberScope(field *graphql.Field) *graphql.Field {
	resolver := field.Resolve
	field.Resolve = func(p graphql.ResolveParams) (interface{}, error) {
		req := p.Context.Value(ctxHTTPReq).(*http.Request)

		authorization := req.Header.Get("Authorization")
		if authorization == "" {
			return nil, errForbidden
		}

		memberUUID, err := nagase.GetMemberUUIDFromToken(strings.Split(authorization, " ")[1])
		if err != nil {
			return nil, errForbidden
		}

		var member nagase.Member
		database.DB.Where(&nagase.Member{UUID: memberUUID}).First(&member)
		if member.UUID == "" || !member.IsActivated {
			return nil, errForbidden
		}

		p.Context = context.WithValue(p.Context, ctxScopeAdmin, &member)
		return resolver(p)
	}

	return field
}

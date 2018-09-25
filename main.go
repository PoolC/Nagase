package main

import (
	"context"
	"net/http"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"

	"nagase/components/auth"
	"nagase/models"
)

func main() {
	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name:   "RootQuery",
			Fields: graphql.Fields{
				"me": models.MeQuery,
			},
		}),
		Mutation: graphql.NewObject(graphql.ObjectConfig{
			Name: "RootMutation",
			Fields: graphql.Fields{
				"createAccessToken": models.CreateAccessTokenMutation,
				"createMember": models.CreateMemberMutation,
			},
		}),
	})

	// Set GraphQL endpoint.
	h := handler.New(&handler.Config{Schema: &schema})

	http.Handle("/graphql", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// CORS configurations.
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")

		// Set context.
		ctx := context.Background()
		authorization := r.Header.Get("Authorization")
		if authorization != "" {
			memberUUID, _, err := auth.ValidatedToken(strings.Split(authorization, " ")[1])
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return 
			}

			member, err := models.GetMemberByUUID(memberUUID)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return 
			}

			ctx = context.WithValue(ctx, "member", member)
		}

		h.ContextHandler(ctx, w, r)
	}))

	http.ListenAndServe(":8080", nil)
}

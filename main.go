package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"

	"nagase/components/auth"
	"nagase/models"
)

func main() {
	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "RootQuery",
			Fields: graphql.Fields{
				"me":       models.MeQuery,
				"members":  models.MembersQuery,
				"board":    models.BoardQuery,
				"boards":   models.BoardsQuery,
				"post":     models.PostQuery,
				"postPage": models.PostPageQuery,
				"project":  models.ProjectQuery,
				"projects": models.ProjectsQuery,
				"vote":     models.VoteQuery,
			},
		}),
		Mutation: graphql.NewObject(graphql.ObjectConfig{
			Name: "RootMutation",
			Fields: graphql.Fields{
				// Access tokens
				"createAccessToken":  models.CreateAccessTokenMutation,
				"refreshAccessToken": models.RefreshAccessTokenMutation,

				// Boards
				"createBoard": models.CreateBoardMutation,
				"updateBoard": models.UpdateBoardMutation,
				"deleteBoard": models.DeleteBoardMutation,

				// Comments
				"createComment": models.CreateCommentMutation,
				"deleteComment": models.DeleteCommentMutation,

				// Members
				"createMember":               models.CreateMemberMutation,
				"updateMember":               models.UpdateMemberMutation,
				"updateMemberPassword":       models.UpdateMemberPasswordMutation,
				"deleteMember":               models.DeleteMemberMutation,
				"toggleMemberIsActivated":    models.ToggleMemberIsActivatedMutation,
				"toggleMemberIsAdmin":        models.ToggleMemberIsAdminMutation,
				"requestMemberPasswordReset": models.RequestPasswordResetMutation,

				// Posts
				"createPost": models.CreatePostMutation,
				"deletePost": models.DeletePostMutation,
				"updatePost": models.UpdatePostMutation,

				// Projects
				"createProject": models.CreateProjectMutation,
				"updateProject": models.UpdateProjectMutation,
				"deleteProject": models.DeleteProjectMutation,

				// Votes
				"selectVoteOption": models.SelectVoteOptionMutation,

				// Push tokens & subscriptions
				"registerPushToken":   models.RegisterPushTokenMutation,
				"deregisterPushToken": models.DeregisterPushTokenMutation,
				"subscribeBoard":      models.SubscribeBoardMutation,
				"unsubscribeBoard":    models.UnsubscribeBoardMutation,
				"subscribePost":       models.SubscribePostMutation,
				"unsubscribePost":     models.UnsubscribePostMutation,
			},
		}),
	})

	// Set GraphQL endpoint.
	h := handler.New(&handler.Config{
		Schema:     &schema,
		Playground: true,
	})

	server := http.NewServeMux()
	server.Handle("/graphql", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// CORS configurations.
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")

		// Set context.
		ctx := context.Background()
		authorization := r.Header.Get("Authorization")
		if authorization != "" {
			memberUUID, err := auth.ValidatedToken(strings.Split(authorization, " ")[1])
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			member, err := models.GetMemberByUUID(memberUUID)
			if err != nil || !member.IsActivated {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx = context.WithValue(ctx, "member", member)
		}

		h.ContextHandler(ctx, w, r)
	})))
	server.Handle("/files/", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse fileName.
		paths := strings.Split(r.URL.Path, "/")
		if len(paths) != 3 {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		fileName := paths[2]

		// CORS configurations.
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")

		if r.Method == "GET" {
			buffer, contentType, err := models.GetFile(fileName)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if contentType != "" {
				w.Header().Set("Content-Type", contentType)
			}
			io.Copy(w, buffer)
		} else if r.Method == "POST" {
			// Only logged-in member can upload new files.
			authorization := r.Header.Get("Authorization")
			if authorization == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			memberUUID, err := auth.ValidatedToken(strings.Split(authorization, " ")[1])
			if memberUUID == "" || err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// Limit file size to 5MB.
			r.Body = http.MaxBytesReader(w, r.Body, 5*1024*1024)

			// Get bytes from the HTTP request.
			var buffer bytes.Buffer
			upload, _, err := r.FormFile("upload")
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			io.Copy(&buffer, upload)

			// Save file
			err = models.SaveFile(&buffer, fileName)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusOK)
		} else if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})))

	fmt.Println("Server listening port 8080...")
	http.ListenAndServe(":8080", handlers.CompressHandler(server))
}

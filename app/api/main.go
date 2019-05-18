package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"nagase"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"

	"nagase/models"
)

func main() {
	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "RootQuery",
			Fields: graphql.Fields{
				// 회원 정보
				"me":      meQuery,
				"members": withAdminScope(membersQuery),

				// 게시판
				"board":  boardQuery,
				"boards": boardsQuery,
			},
		}),
		Mutation: graphql.NewObject(graphql.ObjectConfig{
			Name: "RootMutation",
			Fields: graphql.Fields{
				// 회원 정보
				"createMember":            createMemberMutation,
				"updateMember":            withMemberScope(updateMemberMutation),
				"updateMemberPassword":    withMemberScope(updateMemberPasswordMutation),
				"deleteMember":            withAdminScope(deleteMemberMutation),
				"requestPasswordReset":    requestPasswordResetMutation,
				"toggleMemberIsActivated": withAdminScope(toggleMemberIsActivatedMutation),
				"toggleMemberIsAdmin":     withAdminScope(toggleMemberIsAdminMutation),

				// 게시판
				"createBoard": withAdminScope(createBoardMutation),
				"updateBoard": withAdminScope(updateBoardMutation),
				"deleteBoard": withAdminScope(deleteBoardMutation),
			},
		}),
	})

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
		ctx = context.WithValue(ctx, ctxHTTPReq, r)

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
			memberUUID, err := nagase.GetMemberUUIDFromToken(strings.Split(authorization, " ")[1])
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

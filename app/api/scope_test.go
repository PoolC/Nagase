package main

import (
	"context"
	"net/http"

	"github.com/graphql-go/graphql"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("withAdminScope", func() {
	field := &graphql.Field{
		Resolve: func(p graphql.ResolveParams) (interface{}, error) { return true, nil },
	}

	It("success for admin", func() {
		token, _ := dummyAdminMember().GenerateAccessToken()
		req, _ := http.NewRequest("GET", "dummy", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		ctx := context.Background()
		p := graphql.ResolveParams{
			Context: context.WithValue(ctx, ctxHTTPReq, req),
		}
		result, err := withAdminScope(field).Resolve(p)

		Expect(err).To(BeNil())
		Expect(result.(bool)).To(BeTrue())
	})

	It("error for non-admin", func() {
		token, _ := dummyMember().GenerateAccessToken()
		req, _ := http.NewRequest("GET", "dummy", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		ctx := context.Background()
		p := graphql.ResolveParams{
			Context: context.WithValue(ctx, ctxHTTPReq, req),
		}
		result, err := withAdminScope(field).Resolve(p)

		Expect(err).NotTo(BeNil())
		Expect(result).To(BeNil())
	})

	It("error for non-member", func() {
		req, _ := http.NewRequest("GET", "dummy", nil)

		ctx := context.Background()
		p := graphql.ResolveParams{
			Context: context.WithValue(ctx, ctxHTTPReq, req),
		}
		result, err := withAdminScope(field).Resolve(p)

		Expect(err).NotTo(BeNil())
		Expect(result).To(BeNil())
	})
})

var _ = Describe("withMemberScope", func() {
	field := &graphql.Field{
		Resolve: func(p graphql.ResolveParams) (interface{}, error) { return true, nil },
	}

	It("success for member", func() {
		token, _ := dummyMember().GenerateAccessToken()
		req, _ := http.NewRequest("GET", "dummy", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		ctx := context.Background()
		p := graphql.ResolveParams{
			Context: context.WithValue(ctx, ctxHTTPReq, req),
		}
		result, err := withMemberScope(field).Resolve(p)

		Expect(err).To(BeNil())
		Expect(result.(bool)).To(BeTrue())
	})

	It("error for non-member", func() {
		req, _ := http.NewRequest("GET", "dummy", nil)

		ctx := context.Background()
		p := graphql.ResolveParams{
			Context: context.WithValue(ctx, ctxHTTPReq, req),
		}
		result, err := withMemberScope(field).Resolve(p)

		Expect(err).NotTo(BeNil())
		Expect(result).To(BeNil())
	})
})

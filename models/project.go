package models

import (
	"fmt"
	"time"

	"github.com/graphql-go/graphql"

	"nagase/components/database"
)

type Project struct {
	ID int

	Name         string `gorm:"type:varchar(255)"`
	Genre        string `gorm:"type:varchar(255)"`
	ThumbnailURL string `gorm:"type:varchar(255)"`
	Body         string

	CreatedAt time.Time
	UpdatedAt time.Time
}

var projectType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Project",
	Fields: graphql.Fields{
		"id":           &graphql.Field{Type: graphql.NewNonNull(graphql.Int)},
		"name":         &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"genre":        &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"thumbnailURL": &graphql.Field{Type: graphql.String},
		"body":         &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
	},
})

var projectInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name:        "ProjectInput",
	Description: "프로젝트 추가/수정 InputObject",
	Fields: graphql.InputObjectConfigFieldMap{
		"name":         &graphql.InputObjectFieldConfig{Type: graphql.String},
		"genre":        &graphql.InputObjectFieldConfig{Type: graphql.String},
		"thumbnailURL": &graphql.InputObjectFieldConfig{Type: graphql.String},
		"body":         &graphql.InputObjectFieldConfig{Type: graphql.String},
	},
})

// Queries
var ProjectsQuery = &graphql.Field{
	Type:        graphql.NewList(graphql.NewNonNull(projectType)),
	Description: "프로젝트 목록을 조회합니다.",
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		var prjs []Project
		database.DB.Order("id desc").Find(&prjs)
		return prjs, nil
	},
}

var ProjectQuery = &graphql.Field{
	Type:        projectType,
	Description: "프로젝트를 조회합니다.",
	Args: graphql.FieldConfigArgument{
		"projectID": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		var prj Project
		database.DB.Where(&Project{ID: params.Args["projectID"].(int)}).First(&prj)
		if prj.ID == 0 {
			return nil, fmt.Errorf("ERR400")
		}
		return prj, nil
	},
}

// Mutations
var CreateProjectMutation = &graphql.Field{
	Type:        projectType,
	Description: "프로젝트를 추가합니다. 관리자 권한이 필요합니다.",
	Args: graphql.FieldConfigArgument{
		"ProjectInput": &graphql.ArgumentConfig{Type: projectInputType},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		if member := params.Context.Value("member"); member == nil || !member.(*Member).IsAdmin {
			return nil, fmt.Errorf("ERR401")
		}

		prjInput, _ := params.Args["ProjectInput"].(map[string]interface{})
		prj := Project{
			Name:         prjInput["name"].(string),
			Genre:        prjInput["genre"].(string),
			ThumbnailURL: prjInput["thumbnailURL"].(string),
			Body:         prjInput["body"].(string),
		}
		errs := database.DB.Save(&prj).GetErrors()
		if len(errs) > 0 {
			return nil, errs[0]
		}

		return prj, nil
	},
}

var UpdateProjectMutation = &graphql.Field{
	Type:        projectType,
	Description: "프로젝트를 수정합니다. 관리자 권한이 필요합니다.",
	Args: graphql.FieldConfigArgument{
		"projectID":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
		"ProjectInput": &graphql.ArgumentConfig{Type: projectInputType},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		if member := params.Context.Value("member"); member == nil || !member.(*Member).IsAdmin {
			return nil, fmt.Errorf("ERR401")
		}

		var prj Project
		database.DB.Where(&Project{ID: params.Args["projectID"].(int)}).First(&prj)
		if prj.ID == 0 {
			return nil, fmt.Errorf("ERR400")
		}

		prjInput, _ := params.Args["ProjectInput"].(map[string]interface{})
		if prjInput["name"] != nil {
			prj.Name = prjInput["name"].(string)
		}
		if prjInput["genre"] != nil {
			prj.Genre = prjInput["genre"].(string)
		}
		if prjInput["thumbnailURL"] != nil {
			prj.ThumbnailURL = prjInput["thumbnailURL"].(string)
		}
		if prjInput["body"] != nil {
			prj.Body = prjInput["body"].(string)
		}

		errs := database.DB.Save(&prj).GetErrors()
		if len(errs) > 0 {
			return nil, errs[0]
		}

		return prj, nil
	},
}

var DeleteProjectMutation = &graphql.Field{
	Type:        projectType,
	Description: "프로젝트를 삭제합니다. 관리자 권한이 필요합니다.",
	Args: graphql.FieldConfigArgument{
		"projectID": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		if member := params.Context.Value("member"); member == nil || !member.(*Member).IsAdmin {
			return nil, fmt.Errorf("ERR401")
		}

		var prj Project
		database.DB.Where(&Project{ID: params.Args["projectID"].(int)}).First(&prj)
		if prj.ID == 0 {
			return nil, fmt.Errorf("ERR400")
		}

		errs := database.DB.Delete(&prj).GetErrors()
		if len(errs) > 0 {
			return nil, errs[0]
		}

		return prj, nil
	},
}

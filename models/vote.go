package models

import (
	"fmt"
	"time"

	"github.com/graphql-go/graphql"

	"nagase/components/database"
)

type Vote struct {
	ID int

	Title                string `gorm:"type:varchar(255)"`
	IsMultipleSelectable bool   `gorm:"default:false"`
	Deadline             time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

type VoteOption struct {
	ID int

	VoteID int
	Text   string `gorm:"type:varchar(255)"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

type VoteSelection struct {
	MemberUUID   string `gorm:"type:varchar(40)"`
	VoteID       int    `gorm:"INDEX"`
	VoteOptionID int

	CreatedAt time.Time
	UpdatedAt time.Time
}

var voteType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Vote",
	Fields: graphql.Fields{
		"id":                   &graphql.Field{Type: graphql.NewNonNull(graphql.Int)},
		"title":                &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"isMultipleSelectable": &graphql.Field{Type: graphql.NewNonNull(graphql.Boolean)},
		"deadline":             &graphql.Field{Type: graphql.NewNonNull(graphql.DateTime)},
		"options": &graphql.Field{
			Type: graphql.NewList(graphql.NewNonNull(voteOptionType)),
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				var options []VoteOption
				database.DB.Where(&VoteOption{VoteID: params.Source.(Vote).ID}).Find(&options)
				return options, nil
			},
		},
		"totalVotersCount": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Int),
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				var voters []VoteSelection
				database.DB.Select("DISTINCT(member_uuid)").Where(&VoteSelection{VoteID: params.Source.(Vote).ID}).Find(&voters)
				return len(voters), nil
			},
		},
	},
})

var voteOptionType = graphql.NewObject(graphql.ObjectConfig{
	Name: "VoteOption",
	Fields: graphql.Fields{
		"id":   &graphql.Field{Type: graphql.NewNonNull(graphql.Int)},
		"text": &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"votersCount": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Int),
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				var voters []VoteSelection
				option := params.Source.(VoteOption)
				database.DB.Where(&VoteSelection{VoteID: option.VoteID, VoteOptionID: option.ID}).Find(&voters)
				return len(voters), nil
			},
		},
		"voters": &graphql.Field{
			Type: graphql.NewList(graphql.NewNonNull(memberType)),
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				var voters []VoteSelection
				option := params.Source.(VoteOption)
				database.DB.Where(&VoteSelection{VoteID: option.VoteID, VoteOptionID: option.ID}).Find(&voters)

				var members []Member
				for _, v := range voters {
					member, _ := GetMemberByUUID(v.MemberUUID)
					members = append(members, *member)
				}
				return members, nil
			},
		},
	},
})

var voteInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name:        "VoteInput",
	Description: "투표 InputObject",
	Fields: graphql.InputObjectConfigFieldMap{
		"title":                &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
		"optionTexts":          &graphql.InputObjectFieldConfig{Type: graphql.NewList(graphql.NewNonNull(graphql.String))},
		"isMultipleSelectable": &graphql.InputObjectFieldConfig{Type: graphql.Boolean},
		"deadline":             &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.DateTime)},
	},
})

// Queries
var VoteQuery = &graphql.Field{
	Type:        voteType,
	Description: "투표를 조회합니다.",
	Args: graphql.FieldConfigArgument{
		"voteID": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		if params.Context.Value("member") == nil {
			return nil, fmt.Errorf("ERR401")
		}

		// Get vote
		voteID, _ := params.Args["voteID"].(int)
		var vote Vote
		database.DB.Where(&Vote{ID: voteID}).First(&vote)
		if vote.ID == 0 {
			return nil, fmt.Errorf("ERR400")
		}
		return vote, nil
	},
}

// Mutations
var SelectVoteOptionMutation = &graphql.Field{
	Type:        voteType,
	Description: "투표 선택지를 선택합니다. 이미 투표한 경우, 기존의 선택을 무르고 새로운 요청으로 덮어씌웁니다.",
	Args: graphql.FieldConfigArgument{
		"voteID":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
		"optionIDs": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(graphql.Int)))},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		if params.Context.Value("member") == nil {
			return nil, fmt.Errorf("ERR401")
		}
		member := params.Context.Value("member").(*Member)

		voteID := params.Args["voteID"].(int)
		var vote Vote
		database.DB.Where(&Vote{ID: voteID}).First(&vote)
		if vote.ID == 0 {
			return nil, fmt.Errorf("ERR400")
		}

		// Validate vote and selection(s).
		if vote.Deadline.Before(time.Now()) ||
			(len(params.Args["optionIDs"].([]interface{})) > 1 && !vote.IsMultipleSelectable) {
			return nil, fmt.Errorf("ERR400")
		}

		// If the member has voted already, remove selections.
		database.DB.Where(&VoteSelection{VoteID: voteID, MemberUUID: member.UUID}).Delete(&VoteSelection{})

		// Save selections.
		for _, id := range params.Args["optionIDs"].([]interface{}) {
			database.DB.Save(&VoteSelection{VoteID: voteID, VoteOptionID: id.(int), MemberUUID: member.UUID})
		}

		return vote, nil
	},
}

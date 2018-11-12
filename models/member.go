package models

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/graphql-go/graphql"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/argon2"

	"nagase/components/database"
)

type Member struct {
	UUID string `gorm:"type:varchar(40);PRIMARY_KEY"`

	LoginID      string `gorm:"type:varchar(40);UNIQUE_INDEX"`
	PasswordHash []byte `gorm:"NOT NULL" json:"-"`
	PasswordSalt []byte `gorm:"NOT NULL" json:"-"`
	Email        string `gorm:"type:varchar(255);UNIQUE_INDEX"`
	Name         string `gorm:"type:varchar(40)"`
	Department   string `gorm:"type:varchar(40)"`
	StudentID    string `gorm:"type:varchar(40);UNIQUE_INDEX"`

	IsActivated bool `gorm:"default:false"`
	IsAdmin     bool `gorm:"default:false"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (member Member) ValidatePassword(password string) bool {
	hash := argon2.IDKey([]byte(password), member.PasswordSalt, 1, 8*1024, 4, 32)
	return bytes.Compare(hash, member.PasswordHash) == 0
}

var memberType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Member",
	Fields: graphql.Fields{
		"uuid":        &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"loginID":     &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"email":       &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"name":        &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"department":  &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"studentID":   &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"isActivated": &graphql.Field{Type: graphql.NewNonNull(graphql.Boolean)},
		"isAdmin":     &graphql.Field{Type: graphql.NewNonNull(graphql.Boolean)},
	},
})

var memberInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name:        "MemberInput",
	Description: "회원가입 InputObject",
	Fields: graphql.InputObjectConfigFieldMap{
		"loginID":    &graphql.InputObjectFieldConfig{Type: graphql.String},
		"password":   &graphql.InputObjectFieldConfig{Type: graphql.String},
		"email":      &graphql.InputObjectFieldConfig{Type: graphql.String},
		"name":       &graphql.InputObjectFieldConfig{Type: graphql.String},
		"department": &graphql.InputObjectFieldConfig{Type: graphql.String},
		"studentID":  &graphql.InputObjectFieldConfig{Type: graphql.String},
	},
})

// Queries
var MeQuery = &graphql.Field{
	Type:        memberType,
	Description: "자신의 회원 정보를 조회합니다.",
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		return params.Context.Value("member").(*Member), nil
	},
}

// Mutations
var CreateMemberMutation = &graphql.Field{
	Type:        memberType,
	Description: "회원을 추가합니다.",
	Args: graphql.FieldConfigArgument{
		"MemberInput": &graphql.ArgumentConfig{Type: graphql.NewNonNull(memberInputType)},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		memberInput := params.Args["MemberInput"].(map[string]interface{})

		// Create member model.
		salt := make([]byte, 32)
		rand.Read(salt)

		hash := argon2.IDKey([]byte(memberInput["password"].(string)), salt, 1, 8*1024, 4, 32)
		member := Member{
			UUID:         uuid.NewV4().String(),
			LoginID:      memberInput["loginID"].(string),
			PasswordHash: hash,
			PasswordSalt: salt,
			Email:        memberInput["email"].(string),
			Name:         memberInput["name"].(string),
			Department:   memberInput["department"].(string),
			StudentID:    memberInput["studentID"].(string),
			IsActivated:  false,
			IsAdmin:      false,
		}

		// Save record on DB.
		errs := database.DB.Create(&member).GetErrors()
		if len(errs) > 0 {
			return nil, fmt.Errorf("failed to create member")
		}

		return member, nil
	},
}

var UpdateMemberMutation = &graphql.Field{
	Type:        memberType,
	Description: "회원 정보를 수정합니다. 본인만 수정할 수 있습니다.",
	Args: graphql.FieldConfigArgument{
		"MemberInput": &graphql.ArgumentConfig{Type: graphql.NewNonNull(memberInputType)},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		if params.Context.Value("member") == nil {
			return nil, fmt.Errorf("unauthorized")
		}

		// Read member from database to keep persistence.
		member, _ := GetMemberByUUID(params.Context.Value("member").(*Member).UUID)

		// Updated fields except password.
		memberInput := params.Args["MemberInput"].(map[string]interface{})
		if memberInput["loginID"] != nil {
			member.LoginID = memberInput["loginID"].(string)
		}
		if memberInput["email"] != nil {
			member.Email = memberInput["email"].(string)
		}
		if memberInput["name"] != nil {
			member.Name = memberInput["name"].(string)
		}
		if memberInput["department"] != nil {
			member.Department = memberInput["department"].(string)
		}
		if memberInput["studentID"] != nil {
			member.StudentID = memberInput["studentID"].(string)
		}

		// Update password (if requested).
		if memberInput["password"] != nil {
			hash := argon2.IDKey([]byte(memberInput["password"].(string)), member.PasswordSalt, 1, 8*1024, 4, 32)
			member.PasswordHash = hash
		}

		errs := database.DB.Save(&member).GetErrors()
		if len(errs) > 0 {
			return nil, errs[0]
		}
		return member, nil
	},
}

// Common functions
func GetMemberByUUID(uuid string) (*Member, error) {
	var member Member
	database.DB.Where(&Member{UUID: uuid}).First(&member)
	if member.LoginID == "" {
		return nil, fmt.Errorf("invalid member uuid")
	}
	return &member, nil
}

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

	LoginID      string `gorm:"type:varchar(40);UNIQUE_INDEX`
	PasswordHash []byte `gorm:"NOT NULL"`
	PasswordSalt []byte `gorm:"NOT NULL"`
	Email        string `gorm:"type:varchar(255);UNIQUE_INDEX"`
	Name         string `gorm:"type:varchar(40)"`
	Department   string `gorm:"type:varchar(40)"`
	StudentID    string `gorm:"type:varchar(40);UNIQUE_INDEX"`

	IsActivated bool `gorm:"default:false"`
	IsAdmin     bool `gorm:"default:false"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

type MemberType struct {
	UUID string

	LoginID    string
	Email      string
	Name       string
	Department string
	StudentID  string

	IsActivated bool
	IsAdmin     bool

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (member Member) toGraphQLType() MemberType {
	return MemberType {
		UUID:        member.UUID,
		LoginID:     member.LoginID,
		Email:       member.Email,
		Name:        member.Name,
		Department:  member.Department,
		StudentID:   member.StudentID,
		IsActivated: member.IsActivated,
		IsAdmin:     member.IsAdmin,
		CreatedAt:   member.CreatedAt,
		UpdatedAt:   member.UpdatedAt,
	}
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

// Queries
var MeQuery = &graphql.Field{
	Type:        memberType,
	Description: "자신의 회원 정보를 조회합니다.",
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		return params.Context.Value("member").(*Member).toGraphQLType(), nil
	},
}

// Mutations
var CreateMemberMutation = &graphql.Field{
	Type:        memberType,
	Description: "회원을 추가합니다.",
	Args: graphql.FieldConfigArgument{
		"loginID":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
		"password":   &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
		"email":      &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
		"name":       &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
		"department": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
		"studentID":  &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		loginID, _ := params.Args["loginID"].(string)
		password, _ := params.Args["password"].(string)
		email, _ := params.Args["email"].(string)
		name, _ := params.Args["name"].(string)
		department, _ := params.Args["department"].(string)
		studentID, _ := params.Args["studentID"].(string)

		// Create member model.
		salt := make([]byte, 32)
		rand.Read(salt)

		hash := argon2.IDKey([]byte(password), salt, 1, 8*1024, 4, 32)
		member := Member{
			UUID:         uuid.NewV4().String(),
			LoginID:      loginID,
			PasswordHash: hash,
			PasswordSalt: salt,
			Email:        email,
			Name:         name,
			Department:   department,
			StudentID:    studentID,
			IsActivated:  false,
			IsAdmin:      false,
		}

		// Save record on DB.
		errs := database.DB.Create(&member).GetErrors()
		if len(errs) > 0 {
			return nil, fmt.Errorf("failed to create member")
		}

		return member.toGraphQLType(), nil
	},
}

func GetMemberByUUID(uuid string) (*Member, error) {
	member := new(Member)
	database.DB.Where(&Member{UUID: uuid}).First(&member)
	if member.LoginID == "" {
		return nil, fmt.Errorf("invalid member uuid")
	}
	return member, nil
}

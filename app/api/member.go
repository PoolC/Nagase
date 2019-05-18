package main

import (
	"crypto/rand"
	"fmt"
	"nagase/components/database"
	"nagase/components/email"
	"nagase/components/random"
	"time"

	"github.com/graphql-go/graphql"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/argon2"

	"nagase"
)

var memberType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Member",
	Fields: graphql.Fields{
		"uuid":        &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"loginID":     &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"email":       &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"phoneNumber": &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"name":        &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"department":  &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"studentID":   &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"isActivated": &graphql.Field{Type: graphql.NewNonNull(graphql.Boolean)},
		"isAdmin":     &graphql.Field{Type: graphql.NewNonNull(graphql.Boolean)},
	},
})

var meQuery = &graphql.Field{
	Type:        graphql.NewNonNull(memberType),
	Description: "자신의 회원 정보를 조회합니다.",
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		return p.Context.Value(ctxScopeMember).(*nagase.Member), nil
	},
}

var membersQuery = &graphql.Field{
	Type:        graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(memberType))),
	Description: "회원 목록을 조회합니다. 관리자 권한이 필요합니다.",
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		var members []nagase.Member
		database.DB.Order("created_at desc").Find(&members)
		return members, nil
	},
}

var createMemberMutation = &graphql.Field{
	Type:        graphql.NewNonNull(memberType),
	Description: "회원을 추가합니다.",
	Args: graphql.FieldConfigArgument{
		"input": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.NewInputObject(graphql.InputObjectConfig{
				Name:        "CreateMemberInput",
				Description: "회원가입 InputObject",
				Fields: graphql.InputObjectConfigFieldMap{
					"loginID":     &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
					"password":    &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
					"email":       &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
					"phoneNumber": &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
					"name":        &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
					"department":  &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
					"studentID":   &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
				},
			})),
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		input := p.Args["input"].(map[string]interface{})

		// Find duplicated member info, and return error if exists.
		var existingMember nagase.Member
		database.DB.Where(&nagase.Member{LoginID: input["loginID"].(string)}).First(&existingMember)
		if existingMember.UUID != "" {
			return nil, fmt.Errorf("MEM000")
		}
		database.DB.Where(&nagase.Member{Email: input["email"].(string)}).First(&existingMember)
		if existingMember.UUID != "" {
			return nil, fmt.Errorf("MEM001")
		}

		// Create member model.
		salt := make([]byte, 32)
		rand.Read(salt)

		hash := argon2.IDKey([]byte(input["password"].(string)), salt, 1, 8*1024, 4, 32)
		member := nagase.Member{
			UUID:         uuid.NewV4().String(),
			LoginID:      input["loginID"].(string),
			PasswordHash: hash,
			PasswordSalt: salt,
			Email:        input["email"].(string),
			PhoneNumber:  input["phoneNumber"].(string),
			Name:         input["name"].(string),
			Department:   input["department"].(string),
			StudentID:    input["studentID"].(string),
			IsActivated:  false,
			IsAdmin:      false,
		}

		// Save record on DB.
		errs := database.DB.Create(&member).GetErrors()
		if len(errs) > 0 {
			return nil, fmt.Errorf("failed to create member")
		}

		return &member, nil
	},
}

var updateMemberMutation = &graphql.Field{
	Type:        graphql.NewNonNull(memberType),
	Description: "회원 정보를 수정합니다. 본인만 수정할 수 있습니다.",
	Args: graphql.FieldConfigArgument{
		"input": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.NewInputObject(graphql.InputObjectConfig{
				Name:        "UpdateMemberInput",
				Description: "회원가입 InputObject",
				Fields: graphql.InputObjectConfigFieldMap{
					"email":       &graphql.InputObjectFieldConfig{Type: graphql.String},
					"phoneNumber": &graphql.InputObjectFieldConfig{Type: graphql.String},
					"department":  &graphql.InputObjectFieldConfig{Type: graphql.String},
				},
			})),
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		// Read member from database to keep persistence.
		var member nagase.Member
		database.DB.Where(&nagase.Member{UUID: p.Context.Value(ctxScopeMember).(*nagase.Member).UUID}).First(&member)

		input := p.Args["input"].(map[string]interface{})
		if input["email"] != nil {
			member.Email = input["email"].(string)
		}
		if input["phoneNumber"] != nil {
			member.PhoneNumber = input["phoneNumber"].(string)
		}
		if input["password"] != nil {
			hash := argon2.IDKey([]byte(input["password"].(string)), member.PasswordSalt, 1, 8*1024, 4, 32)
			member.PasswordHash = hash
		}

		errs := database.DB.Save(&member).GetErrors()
		if len(errs) > 0 {
			return nil, errs[0]
		}
		return &member, nil
	},
}

var updateMemberPasswordMutation = &graphql.Field{
	Type:        graphql.NewNonNull(memberType),
	Description: "회원 비밀번호를 수정합니다.",
	Args: graphql.FieldConfigArgument{
		"input": &graphql.ArgumentConfig{
			Type: graphql.NewInputObject(graphql.InputObjectConfig{
				Name:        "UpdateMemberPasswordInput",
				Description: "회원 비밀번호 수정 InputObject",
				Fields: graphql.InputObjectConfigFieldMap{
					"token":    &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
					"password": &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
				},
			}),
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		input := p.Args["input"].(map[string]interface{})

		var member nagase.Member
		database.DB.Where(&nagase.Member{PasswordResetToken: input["token"].(string)}).First(&member)
		if member.UUID == "" || member.PasswordResetTokenValidUntil.Before(time.Now()) {
			return nil, fmt.Errorf("TKN001")
		}

		salt := make([]byte, 32)
		hash := argon2.IDKey([]byte(input["password"].(string)), salt, 1, 8*1024, 4, 32)
		member.PasswordHash = hash
		member.PasswordSalt = salt
		errs := database.DB.Save(&member).GetErrors()
		if len(errs) > 0 {
			return nil, errs[0]
		}

		return &member, nil
	},
}

var passwordResetEmailTitle = "비밀번호 초기화 안내"
var passwordResetEmailBody = `
안녕하세요,
PoolC 홈페이지 비밀번호 초기화 안내 메일입니다.

아래 링크를 눌러 비밀번호 초기화를 진행해주세요.
<a href="https://poolc.org/accounts/password-reset?token=%s">https://poolc.org/accounts/password-reset?token=%s</a>
링크는 24시간 동안 유효합니다.

본인이 비밀번호 초기화를 요청하지 않은 경우, 즉시 관리자에게 알려주세요.
감사합니다.
`

var requestPasswordResetMutation = &graphql.Field{
	Type:        graphql.NewNonNull(graphql.Boolean),
	Description: "비밀번호 초기화 이메일 발송을 요청합니다.",
	Args: graphql.FieldConfigArgument{
		"input": &graphql.ArgumentConfig{
			Type: graphql.NewInputObject(graphql.InputObjectConfig{
				Name:        "RequestPasswordResetInput",
				Description: "비밀번호 초기화 요청 InputObject",
				Fields: graphql.InputObjectConfigFieldMap{
					"email": &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
				},
			}),
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		address := p.Args["input"].(map[string]interface{})["address"].(string)

		var member nagase.Member
		database.DB.Where(&nagase.Member{Email: address}).First(&member)
		if member.UUID == "" {
			return true, nil
		}

		token := random.GenerateRandomString(40)
		member.PasswordResetToken = token
		member.PasswordResetTokenValidUntil = time.Now().AddDate(0, 0, 1)
		database.DB.Save(&member)

		mail := email.Email{
			Title: passwordResetEmailTitle,
			Body:  fmt.Sprintf(passwordResetEmailBody, token, token),
			To:    address,
		}
		go func() { mail.Send() }()

		return true, nil
	},
}

var deleteMemberMutation = &graphql.Field{
	Type:        memberType,
	Description: "회원을 삭제합니다. 관리자 권한이 필요합니다.",
	Args: graphql.FieldConfigArgument{
		"input": &graphql.ArgumentConfig{
			Type: graphql.NewInputObject(graphql.InputObjectConfig{
				Name:        "DeleteMemberInput",
				Description: "회원 삭제 InputObject",
				Fields: graphql.InputObjectConfigFieldMap{
					"uuid": &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
				},
			}),
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		var member nagase.Member
		database.DB.Where(&nagase.Member{UUID: p.Args["input"].(map[string]interface{})["uuid"].(string)}).First(&member)
		if member.UUID == "" {
			return nil, nil
		}

		database.DB.Delete(&member)
		return &member, nil
	},
}

var toggleMemberIsActivatedMutation = &graphql.Field{
	Type:        memberType,
	Description: "회원을 활성화/비활성화합니다. 관리자 권한이 필요합니다.",
	Args: graphql.FieldConfigArgument{
		"input": &graphql.ArgumentConfig{
			Type: graphql.NewInputObject(graphql.InputObjectConfig{
				Name:        "ToggleMemberIsActivated",
				Description: "회원 활성화/비활성화 InputObject",
				Fields: graphql.InputObjectConfigFieldMap{
					"uuid": &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
				},
			}),
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		var member nagase.Member
		database.DB.Where(&nagase.Member{UUID: p.Args["input"].(map[string]interface{})["uuid"].(string)}).First(&member)
		member.IsActivated = !member.IsActivated
		errs := database.DB.Save(&member).GetErrors()
		if len(errs) > 0 {
			return nil, errs[0]
		}
		return &member, nil
	},
}

var toggleMemberIsAdminMutation = &graphql.Field{
	Type:        memberType,
	Description: "회원을 관리자로 만들거나, 관리자 권한을 해제합니다. 관리자 권한이 필요합니다.",
	Args: graphql.FieldConfigArgument{
		"input": &graphql.ArgumentConfig{
			Type: graphql.NewInputObject(graphql.InputObjectConfig{
				Name:        "ToggleMemberIsAdmin",
				Description: "회원 활성화/비활성화 InputObject",
				Fields: graphql.InputObjectConfigFieldMap{
					"uuid": &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
				},
			}),
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		var member nagase.Member
		database.DB.Where(&nagase.Member{UUID: p.Args["input"].(map[string]interface{})["uuid"].(string)}).First(&member)
		member.IsAdmin = !member.IsAdmin
		errs := database.DB.Save(&member).GetErrors()
		if len(errs) > 0 {
			return nil, errs[0]
		}
		return &member, nil
	},
}

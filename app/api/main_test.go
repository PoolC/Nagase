package main

import (
	"nagase"
	"nagase/components/database"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	uuid "github.com/satori/go.uuid"
)

func TestAPI(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "API Suite")
}

func dummyMember() *nagase.Member {
	member := nagase.Member{
		UUID:        uuid.NewV4().String(),
		LoginID:     uuid.NewV4().String(),
		Email:       uuid.NewV4().String(),
		StudentID:   uuid.NewV4().String(),
		IsActivated: true,
	}
	database.DB.Save(&member)
	return &member
}

func dummyAdminMember() *nagase.Member {
	member := nagase.Member{
		UUID:        uuid.NewV4().String(),
		LoginID:     uuid.NewV4().String(),
		Email:       uuid.NewV4().String(),
		StudentID:   uuid.NewV4().String(),
		IsAdmin:     true,
		IsActivated: true,
	}
	database.DB.Save(&member)
	return &member
}

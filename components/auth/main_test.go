package auth

import (
	"testing"
	"time"

	"github.com/bouk/monkey"
)

var baseTime = time.Date(2018, 9, 24, 0, 39, 39, 0, time.UTC)

func TestGenerateToken(t *testing.T) {
	patch := monkey.Patch(time.Now, func() time.Time { return baseTime })

	token, err := GenerateToken("00000000-0000-0000-0000-000000000000", false)
	expected := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJtZW1iZXJfdXVpZCI6IjAwMDAwMDAwLTAwMDAtMDAwMC0wMDAwLTAwMDAwMDAwMDAwMCIsImlzX2FkbWluIjpmYWxzZSwiZXhwIjoxNTM4MzU0Mzc5LCJpYXQiOjE1Mzc3NDk1NzksImlzcyI6IlBvb2xDL05hZ2FzZSJ9.na66QDHHdyIFodsCi8cguPirT9YSqogqR3mcThRylVA"
	if err != nil || token != expected {
		t.Fail()
	}

	defer patch.Unpatch()
}

func TestValidatedToken(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJtZW1iZXJfdXVpZCI6IjAwMDAwMDAwLTAwMDAtMDAwMC0wMDAwLTAwMDAwMDAwMDAwMCIsImlzX2FkbWluIjpmYWxzZSwiZXhwIjoxNTM4MzU0Mzc5LCJpYXQiOjE1Mzc3NDk1NzksImlzcyI6IlBvb2xDL05hZ2FzZSJ9.na66QDHHdyIFodsCi8cguPirT9YSqogqR3mcThRylVA"

	// Test should be failed : token used before issued
	patch := monkey.Patch(time.Now, func() time.Time { return baseTime.AddDate(0, 0, -1) })
	_, _, err := ValidatedToken(token)
	if err == nil {
		t.Fail()
	}

	// Test should be failed : token is expired
	patch = monkey.Patch(time.Now, func() time.Time { return baseTime.AddDate(0, 0, 10) })
	_, _, err = ValidatedToken(token)
	if err == nil {
		t.Fail()
	}

	// Test should be passed
	patch = monkey.Patch(time.Now, func() time.Time { return baseTime.AddDate(0, 0, 1) })
	memberUUID, isAdmin, err := ValidatedToken(token)
	if err != nil || memberUUID != "00000000-0000-0000-0000-000000000000" || isAdmin {
		t.Fail()
	}

	defer patch.Unpatch()
}

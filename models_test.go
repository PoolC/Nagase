package nagase_test

import (
	"time"

	"bou.ke/monkey"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "nagase"
)

var _ = Describe("Member", func() {
	Describe("GenerateAccessToken", func() {
		var patch *monkey.PatchGuard
		member := Member{UUID: "00000000-0000-0000-0000-000000000000"}

		BeforeEach(func() {
			patch = monkey.Patch(time.Now, func() time.Time { return time.Date(2018, 9, 24, 0, 39, 39, 0, time.UTC) })
		})

		AfterEach(func() {
			patch.Unpatch()
		})

		It("valid", func() {
			expected := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJtZW1iZXJfdXVpZCI6IjAwMDAwMDAwLTAwMDAtMDAwMC0wMDAwLTAwMDAwMDAwMDAwMCIsImlzX2FkbWluIjpmYWxzZSwiZXhwIjoxNTM4MzU0Mzc5LCJpYXQiOjE1Mzc3NDk1NzksImlzcyI6IlBvb2xDL05hZ2FzZSJ9.na66QDHHdyIFodsCi8cguPirT9YSqogqR3mcThRylVA"
			Expect(member.GenerateAccessToken()).To(Equal(expected))
		})
	})

	Describe("GetMemberUUIDFromToken", func() {
		var patch *monkey.PatchGuard
		baseTime := time.Date(2018, 9, 24, 0, 39, 39, 0, time.UTC)
		token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJtZW1iZXJfdXVpZCI6IjAwMDAwMDAwLTAwMDAtMDAwMC0wMDAwLTAwMDAwMDAwMDAwMCIsImlzX2FkbWluIjpmYWxzZSwiZXhwIjoxNTM4MzU0Mzc5LCJpYXQiOjE1Mzc3NDk1NzksImlzcyI6IlBvb2xDL05hZ2FzZSJ9.na66QDHHdyIFodsCi8cguPirT9YSqogqR3mcThRylVA"

		AfterEach(func() { patch.Unpatch() })

		It("token used before issued", func() {
			patch = monkey.Patch(time.Now, func() time.Time { return baseTime.AddDate(0, 0, -1) })
			_, err := GetMemberUUIDFromToken(token)

			Expect(err).NotTo(BeNil())
		})

		It("token is expired", func() {
			patch = monkey.Patch(time.Now, func() time.Time { return baseTime.AddDate(0, 0, 10) })
			_, err := GetMemberUUIDFromToken(token)

			Expect(err).NotTo(BeNil())
		})

		It("Test should be passed", func() {
			patch = monkey.Patch(time.Now, func() time.Time { return baseTime.AddDate(0, 0, 1) })
			memberUUID, _ := GetMemberUUIDFromToken(token)

			Expect(memberUUID).To(Equal("00000000-0000-0000-0000-000000000000"))
		})
	})
})

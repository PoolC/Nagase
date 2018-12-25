// Script to migrate database records from Yuzuki.
//
// Requirements =>
// Export tables as csv form, and place it on same directory with this script.

package main

import (
	"encoding/csv"
	"os"
	"strconv"
	"strings"
	"time"

	"nagase/components/database"
	"nagase/models"
)

func idToUUID(idStr string) string {
	return "00000000-0000-0000-0000-" + strings.Repeat("0", 12-len(idStr)) + idStr
}

func parseTime(timeStr string) time.Time {
	t, _ := time.Parse("2006-01-02 15:04:05", timeStr)
	return t
}

func readCSV(tableName string) [][]string {
	f, err := os.Open("./scripts/" + tableName + ".csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		panic(err)
	}

	return lines
}

func migrateMembers() {
	for _, src := range readCSV("user") {
		email := src[5]
		if email == "NULL" {
			email = "dummy" + src[0] + "@dummy.com"
		}

		database.DB.Save(&models.Member{
			UUID:        idToUUID(src[0]),
			LoginID:     src[1],
			Email:       "_" + email, // Add underscore to prevent duplication with the Nagase account.
			PhoneNumber: src[7],
			Name:        src[4],
			Department:  "",
			StudentID:   strings.Repeat("0", 10-len(src[0])) + src[0],
			IsActivated: false,
			IsAdmin:     false,
			CreatedAt:   parseTime(src[1]),
			UpdatedAt:   time.Now(),
		})
	}
}

func migratePosts() {
	legacyBoards := readCSV("board")

	for _, src := range readCSV("article") {
		boardID := -1
		titlePrefix := ""
		legacyBoardId, _ := strconv.Atoi(src[1])
		if legacyBoardId == 1 { // 공지사항
			boardID = 1
		} else if legacyBoardId == 19 { // 학술부
			boardID = 4
		} else if legacyBoardId == 20 { // 게임제작부
			boardID = 5
		} else if legacyBoardId == 21 { // 자유게시판
			boardID = 2
		} else if legacyBoardId == 23 { // 버그레포트
			continue
		} else { // 기타 게시물은 모두 자유게시판으로 옮기되, 제목 앞에 게시판 이름을 붙힘
			boardID = 2
			for _, legacyBoard := range legacyBoards {
				i, _ := strconv.Atoi(legacyBoard[0])
				if i == legacyBoardId {
					titlePrefix = "[" + legacyBoard[2] + "] "
					break
				}
			}
		}

		id, _ := strconv.Atoi(src[0])
		post := models.Post{
			ID:         id,
			BoardID:    boardID,
			AuthorUUID: idToUUID(src[2]),
			Title:      titlePrefix + src[3],
			Body:       src[4],
			VoteID:     nil,
			CreatedAt:  parseTime(src[9]),
			UpdatedAt:  parseTime(src[8]),
		}

		database.DB.Save(&post)
	}
}

func migrateComments() {
	for _, src := range readCSV("reply") {
		id, _ := strconv.Atoi(src[0])
		postID, _ := strconv.Atoi(src[1])
		comment := models.Comment{
			ID:         id,
			PostID:     postID,
			AuthorUUID: idToUUID(src[2]),
			Body:       src[3],
			CreatedAt:  parseTime(src[7]),
			UpdatedAt:  parseTime(src[6]),
		}

		database.DB.Save(&comment)
	}
}

func main() {
	// migrateMembers()
	// migratePosts()
	migrateComments()
}

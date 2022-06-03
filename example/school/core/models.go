package core

import (
	"crypto/md5"
	"encoding/hex"
	"time"
)

type Student struct {
	ID           int    `json:"id"`
	EmailAddress string `json:"email_address"`
	Name         string `json:"name"`
}

type Book struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Borrowing struct {
	ID         string    `json:"id"`
	StudentID  int       `json:"student_id"`
	BookID     int       `json:"book_id"`
	BorrowDate time.Time `json:"borrow_date"`
}

func GenerateID(b []byte) string {
	return md5HashHexStr(b)
}

func md5Hash(input []byte) []byte {
	h := md5.New()
	h.Write(input)
	return h.Sum(nil)
}
func md5HashHexStr(data []byte) string {
	return hex.EncodeToString(md5Hash(data))
}

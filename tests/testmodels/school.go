package testmodels

import "encoding/json"

type Book struct {
	ID int `json:"book_id"`
}

func (s *Book) SayHello() string {
	return "hello book"
}

func (s Book) SayBye() string {
	return "bye book"
}

func (b *Book) String() string {
	byt, _ := json.Marshal(b)
	return string(byt)
}

type Student struct {
	ID           int    `json:"id"`
	EmailAddress string `json:"email_address"`
	Name         string `json:"name"`
	Borrowed     bool   `json:"borrowed"`
	*Book
}

func (s *Student) SayHello() string {
	return "hello"
}

func (s Student) SayBye() string {
	return "bye"
}

func (s *Student) String() string {
	byt, _ := json.Marshal(s)
	return string(byt)
}

type StudentList []Student

func (sl *StudentList) String() string {
	byt, _ := json.Marshal(sl)
	return string(byt)
}

package model

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

// Request user
type User struct {
	gorm.Model
	FirstName    string     `json:"firstname" validate:"min=1,max=15"`
	LastName     string     `json:"lastname" validate:"min=0,max=15"`
	Email        string     `json:"email" gorm:"unique" validate:"min=1,max=20,email"`
	Password     string     `json:"password,omitempty" validate:"min=8,max=20"`
	DateOfBirth  *time.Time `json:dateofbirth`
	LastAccessAt *time.Time `json:"lastaccessedat"`
	Archived     bool       `json:"archived" `
}

// Response user
type ResUser struct {
	Id           uint       `json:"ID"`
	CreatedAt    *time.Time `json:"CreatedAt"`
	UpdatedAt    *time.Time `json:"UpdatedAt"`
	FirstName    string     `json:"firstname" `
	LastName     string     `json:"lastname" `
	Email        string     `json:"email" `
	DateOfBirth  *time.Time `json:"dateofbirth"`
	Age          int        `json:"age"`
	LastAccessAt *time.Time `json:"lastaccessedat"`
}
//For Login
type LogUser struct {
	Email    string
	Password string
}
// This function is mapping requset user and response user
func Mapping(u User) ResUser {
	var res ResUser
	res.Id = u.ID
	res.CreatedAt = &u.CreatedAt
	res.UpdatedAt = &u.UpdatedAt
	res.FirstName = u.FirstName
	res.LastName = u.LastName
	res.Email = u.Email
	res.DateOfBirth = u.DateOfBirth
	res.LastAccessAt = u.LastAccessAt
	return res

}
// This type is storing credential information  while login 
type Cred struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}
// This is for paginating the response
type Paginate struct{
	Data []ResUser    `json:"data"`
	Total int64        `json:"total"`
	Page  int         `json:"page"`
	Last_page float64     `json:"lastpage"`
}
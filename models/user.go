package models

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"html"
	"schedule/utils/token"
	"strings"
)

type User struct {
	gorm.Model
	Username  string `gorm:"size:255;not null;unique" json:"username"`
	Password  string `gorm:"size:255;not null;" json:"password"`
	Email     string `gorm:"size:255;not null;" json:"email"`
	FirstName string `gorm:"size:255;not null;" json:"first_name"`
	LastName  string `gorm:"size:255;not null;" json:"last_name"`
}

func GetUserByID(uid uint) (User, error) {
	var u User
	if err := DB.First(&u, uid).Error; err != nil {
		return u, errors.New("User not found!")
	}
	u.PrepareGive()
	return u, nil
}

func GetUserByUsername(username string) (User, error) {
	var u User
	err := DB.Model(User{}).Where("username = ?", username).Take(&u).Error
	if err != nil {
		return u, errors.New("User not found!")
	}
	fmt.Println(u)
	u.PrepareGive()
	return u, nil
}

func GetUserByEmail(email string) (User, error) {
	var u User
	err := DB.Model(User{}).Where("email = ?", email).Take(&u).Error
	if err != nil {
		return u, errors.New("User not found!")
	}
	fmt.Println(u)
	u.PrepareGive()
	return u, nil
}

func (u *User) PrepareGive() {
	u.Password = ""
}

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func LoginCheck(username string, password string) (string, error) {

	var err error

	u := User{}

	err = DB.Model(User{}).Where("username = ?", username).Take(&u).Error

	if err != nil {
		return "", err
	}

	err = VerifyPassword(password, u.Password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}

	token, err := token.GenerateToken(u.ID)

	if err != nil {
		return "", err
	}

	return token, nil

}

func (u *User) SaveUser() (*User, error) {

	var err error
	err = DB.Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) BeforeSave() error {

	u.Username = html.EscapeString(strings.TrimSpace(u.Username))

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)

	return nil

}

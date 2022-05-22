package database

import (
	"errors"
	Logger "example/server/logger"
	model "example/server/model"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error
var DSN = "aayush:Password@123@tcp(127.0.0.1:3306)/USER?parseTime=true"

func DbMigrator() {
	DB, err = gorm.Open(mysql.Open(DSN), &gorm.Config{})
	if err != nil {
		Logger.ErrorLog.Println(err.Error())
		panic("the database is not connected")
	}
	DB.AutoMigrate(&model.User{})

}

//To get Database connection
func GetDb() *gorm.DB {
	return DB
}

// Create user
func CreateUser(u *model.User) model.User {
	DB.Create(&u)
	return *u
}

// To get All user
func Getuser() []model.User {
	var user []model.User
	DB.Find(&user)
	return user
}

// To get User By Id
func GetUserById(id string) model.User {
	var user model.User
	DB.First(&user, id)
	return user
}

// To update user
func UpdateUser(user *model.User) model.User {
	DB.Save(&user)
	return *user
}

// To Delete User By Id
func DeleteUser(id string) error {
	var user model.User
	DB.First(&user, id)
	fmt.Println("The Id :", id)
	fmt.Println("The user :", user)

	if user.ID == 0 {
		err := errors.New("NO DATA FOUND TO DELTE")
		return err
	}
	fmt.Println("the user is", user)
	user.Archived = true
	DB.Save(&user)
	DB.Delete(&user)
	return nil
}

package config

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	VERSION string = "0.0.1"
)

var TOKEN string
var URL_BOT string
var DB_DSN string
var DB *gorm.DB
var DEBUG bool = false

type User struct {
	ID      uint   `gorm:"primarykey"`
	Tid     int64  `gorm:"unique_index"`
	Name    string `gorm:"size:255"`
	Active  bool   `gorm:"type:bool"`
	Sending bool   `gorm:"type:bool"`
	Todo    []Todo
}

type Todo struct {
	gorm.Model
	Num         uint   `gorm:"type:int"`
	UserID      uint   `gorm:"type:bigint"`
	Description string `gorm:"type:text"`
	Active      bool   `gorm:"type:bool"`
}

// get data from .env
func init() {
	file, err := os.Open(".env")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := scanner.Text()
		if strings.Contains(s, "TOKEN=") {
			TOKEN = strings.ReplaceAll(s, "TOKEN=", "")
		}

		if strings.Contains(s, "URL_BOT=") {
			URL_BOT = strings.ReplaceAll(s, "URL_BOT=", "")
		}

		if strings.Contains(s, "DB_DSN=") {
			DB_DSN = strings.ReplaceAll(s, "DB_DSN=", "")
		}

		if strings.Contains(s, "DEBUG=") {
			DEBUG, _ = strconv.ParseBool(strings.ReplaceAll(s, "DEBUG=", ""))
		}
	}

	db, err := gorm.Open(postgres.Open(DB_DSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	db.AutoMigrate(&User{}, &Todo{})
	DB = db
}

func SetUser(tid int64, userName string) User {
	var us User
	us.Tid = tid
	us.Name = userName
	us.Sending = true
	us.Active = true
	DB.FirstOrCreate(&us, User{Tid: tid})

	return us
}

func GetListUsers() []User {
	var users []User
	DB.Select("id", "tid", "name").Where("active = ?", true).Where("sending = ?", true).Find(&users)

	return users
}

func SetUserSending(uid uint, sending bool) User {
	var us User
	us.ID = uid
	DB.Model(&us).Update("sending", sending)

	return us
}

func getNum(uid uint) uint {
	var td Todo
	DB.Select("id", "num").Where("user_id = ?", uid).Order("id DESC").First(&td)
	return td.Num + 1
}

func AddTodo(us *User, Text string) Todo {
	var td Todo
	td.Active = true
	td.Num = getNum(us.ID)
	td.UserID = us.ID

	DB.FirstOrCreate(&td, Todo{Description: Text})

	return td
}

func GetTodoList(uid uint, page int, pageSize int) ([]Todo, int64, int, bool) {
	var count int64
	var tds []Todo

	DB.Model(&Todo{}).Where("user_id = ?", uid).Order("id DESC").Count(&count)

	offset := (page - 1) * pageSize
	DB.Where("user_id = ?", uid).Order("id DESC").Offset(offset).Limit(pageSize).Find(&tds)

	return tds, count, page, count > int64(pageSize)
}

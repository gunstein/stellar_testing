package models

import (
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Art struct {
	gorm.Model
	Name          string
	Description   string
	SmallFileName string
	BigFileName   string
	Order         []Order
}

type Order struct {
	gorm.Model
	Email string
	Paid  bool
	ArtId uint
}

var DB *gorm.DB

func ConnectDatabase() {
	os.Remove("test.db")
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Art{}, &Order{})

	art1 := Art{Name: "Aurora Borealis", Description: "bla bla bla", SmallFileName: "small_aurora_borealis.jpg", BigFileName: "big_aurora_borealis.jpg"}
	db.Create(&art1)

	//order := Order{Email: "gvtest", ArtId: 1}
	//db.Create(&order)

	DB = db
}

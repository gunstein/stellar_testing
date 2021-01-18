package models

import (
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Art struct {
	gorm.Model
	Title        string
	Description  string
	Artist       string
	Price        float32
	SmallFileUrl string
	BigFileUrl   string
	Order        []Order
}

type Order struct {
	gorm.Model
	Email string
	Paid  bool
	ArtId uint
}

var DB *gorm.DB

func ConnectDatabase() {
	os.Remove("stellar_art_gallery.db")
	db, err := gorm.Open(sqlite.Open("stellar_art_gallery.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Art{}, &Order{})

	art1 := Art{Title: "Sushi", Description: "Made in 2020 during the Lockdown.", Artist: "Anneli", Price: 1, SmallFileUrl: "https://drive.google.com/uc?export=view&id=1eU1DLGWLha8KfJct5DWS2WadJz0RcifC", BigFileUrl: "big_aurora_borealis.jpg"}
	db.Create(&art1)
	art2 := Art{Title: "ScoobyDoo", Description: "Made in 2020 during the Lockdown.", Artist: "Anneli", Price: 1, SmallFileUrl: "https://drive.google.com/uc?export=view&id=1TciUqcac_FN1P-c-h1gnhYlnT2cnK1qa", BigFileUrl: "big_aurora_borealis.jpg"}
	db.Create(&art2)
	art3 := Art{Title: "Paperbag", Description: "Made in 2020 during the Lockdown.", Artist: "Anneli", Price: 1, SmallFileUrl: "https://drive.google.com/uc?export=view&id=1IxDAhNTq5V7TTRE3s0xfwvvb0GaLk3Fh", BigFileUrl: "big_aurora_borealis.jpg"}
	db.Create(&art3)
	art4 := Art{Title: "LionKing", Description: "Made in 2020 during the Lockdown.", Artist: "Anneli", Price: 1, SmallFileUrl: "https://drive.google.com/uc?export=view&id=1VgD7J5lYPoDLtJ6d9-sJH4b9u6eCyGXD", BigFileUrl: "big_aurora_borealis.jpg"}
	db.Create(&art4)
	DB = db
}

func GetOrderForId(OrderId uint) (order Order, err error) {
	err = DB.First(&order, OrderId).Error
	return order, err
}

func GetArtForId(ArtId uint) (art Art, err error) {
	err = DB.First(&art, ArtId).Error
	return art, err
}

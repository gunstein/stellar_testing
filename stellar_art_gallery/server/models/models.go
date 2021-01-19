package models

import (
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Art struct {
	gorm.Model
	Title        string `gorm:"not null"`
	Description  string
	Artist       string 
	Price        float32 `gorm:"not null"`
	SmallFileUrl string `gorm:"not null"`
	BigFileUrl   string `gorm:"not null"`
	Order        []Order
}

type Order struct {
	gorm.Model
	Email string `gorm:"not null"`
	Paid  bool `gorm:"default:false"`
	EmailSent bool `gorm:"default:false"`
	ArtId uint `gorm:"not null"`
}

var DB *gorm.DB

func ConnectDatabase() {
	os.Remove("stellar_art_gallery.db")
	db, err := gorm.Open(sqlite.Open("stellar_art_gallery.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Art{}, &Order{})

	art1 := Art{Title: "Sushi", Description: "Made in 2020 during the Lockdown.", Artist: "Anneli", Price: 1, SmallFileUrl: "https://drive.google.com/uc?export=view&id=1eU1DLGWLha8KfJct5DWS2WadJz0RcifC", BigFileUrl: "https://drive.google.com/uc?export=view&id=1kZFnq6rwmKjJ0P8B9ClEWXcdOlkBz98s"}
	db.Create(&art1)
	art2 := Art{Title: "ScoobyDoo", Description: "Made in 2020 during the Lockdown.", Artist: "Anneli", Price: 1, SmallFileUrl: "https://drive.google.com/uc?export=view&id=1TciUqcac_FN1P-c-h1gnhYlnT2cnK1qa", BigFileUrl: "https://drive.google.com/uc?export=view&id=1hKwCJ-x3UKwJEfS-HRRPL9OkgMPv5YSn"}
	db.Create(&art2)
	art3 := Art{Title: "Paperbag", Description: "Made in 2020 during the Lockdown.", Artist: "Anneli", Price: 1, SmallFileUrl: "https://drive.google.com/uc?export=view&id=1IxDAhNTq5V7TTRE3s0xfwvvb0GaLk3Fh", BigFileUrl: "https://drive.google.com/uc?export=view&id=1a4U3D_znz-_vRccPCB_yRw5XGRSWEoQH"}
	db.Create(&art3)
	art4 := Art{Title: "LionKing", Description: "Made in 2020 during the Lockdown.", Artist: "Anneli", Price: 1, SmallFileUrl: "https://drive.google.com/uc?export=view&id=1VgD7J5lYPoDLtJ6d9-sJH4b9u6eCyGXD", BigFileUrl: "https://drive.google.com/uc?export=view&id=1ys0KJ200qUiralOGfHWYL6dRIaXGPLf3"}
	db.Create(&art4)
	DB = db
}

func GetArtForId(ArtId uint) (art Art, err error) {
	err = DB.First(&art, ArtId).Error
	return art, err
}

func GetOrderForId(OrderId uint) (order Order, err error) {
	err = DB.First(&order, OrderId).Error
	return order, err
}

func UpdateOrderToPaid(OrderId uint) (order Order, err error) {
	err = DB.First(&order, OrderId).Error
	if err != nil {
		return order, err
	}
	DB.Model(&order).Update("paid", true)	
	return order, err
}

func UpdateOrderToEmailSent(OrderId uint) (order Order, err error) {
	err = DB.First(&order, OrderId).Error
	if err != nil {
		return order, err
	}
	DB.Model(&order).Update("email_sent", true)	
	return order, err
}
package models

import (
	"gorm.io/driver/postgres"
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
	AlternativeBigFileUrl   string `gorm:"not null"`
	Order        []Order
}

type Order struct {
	gorm.Model
	Paid  bool `gorm:"default:false"`
	ArtId uint `gorm:"not null"`
	DownloadKey string
}

var DB *gorm.DB

func ConnectDatabase() {
	dsn := "host=postgres user=postgres password=postgres dbname=postgres port=5432"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.Migrator().DropTable(&Art{})
	db.Migrator().DropTable(&Order{})
	db.AutoMigrate(&Art{}, &Order{})

	art1 := Art{Title: "Sushi", Description: "Made in 2020 during the Lockdown.", Artist: "Anneli", Price: 1, SmallFileUrl: "/assets/images/sushi_small.jpg", BigFileUrl: "https://storageforgv.blob.core.windows.net/stellargallery/sushi_big.jpg", AlternativeBigFileUrl: "https://drive.google.com/uc?export=view&id=1kZFnq6rwmKjJ0P8B9ClEWXcdOlkBz98s"}
	db.Create(&art1)
	art2 := Art{Title: "ScoobyDoo", Description: "Made in 2020 during the Lockdown.", Artist: "Anneli", Price: 1, SmallFileUrl: "/assets/images/scoobydoo_small.jpg", BigFileUrl: "https://storageforgv.blob.core.windows.net/stellargallery/scoobydoo_big.jpg", AlternativeBigFileUrl: "https://drive.google.com/uc?export=view&id=1hKwCJ-x3UKwJEfS-HRRPL9OkgMPv5YSn"}
	db.Create(&art2)
	art3 := Art{Title: "Paperbag", Description: "Made in 2020 during the Lockdown.", Artist: "Anneli", Price: 1, SmallFileUrl: "/assets/images/paperbag_small.jpg", BigFileUrl: "https://storageforgv.blob.core.windows.net/stellargallery/paperbag_big.jpg", AlternativeBigFileUrl: "https://drive.google.com/uc?export=view&id=1a4U3D_znz-_vRccPCB_yRw5XGRSWEoQH"}
	db.Create(&art3)
	art4 := Art{Title: "LionKing", Description: "Made in 2020 during the Lockdown.", Artist: "Anneli", Price: 1, SmallFileUrl: "/assets/images/lionking_small.jpg", BigFileUrl: "https://storageforgv.blob.core.windows.net/stellargallery/lionking_big.jpg", AlternativeBigFileUrl: "https://drive.google.com/uc?export=view&id=1ys0KJ200qUiralOGfHWYL6dRIaXGPLf3"}
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

func UpdateOrderToSSESent(OrderId uint) (order Order, err error) {
	err = DB.First(&order, OrderId).Error
	if err != nil {
		return order, err
	}
	DB.Model(&order).Update("sse_sent", true)	
	return order, err
}
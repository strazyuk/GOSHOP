package database
import (
	"log"
	"backend/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)
var DB *gorm.DB

func Connect(dbPath string){
	var err error
	DB,err = gorm.Open(sqlite.Open(dbPath)
	 , &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),


		})
	if err != nil {
		log.Fatal("Failed to connect to database:",err)
	}
	log.Fatal("Database connected:",dbPath)

}

err = DB.AutoMigrate(
	&models.User{},
	&models,Category{},
    &models.Product{},
    &models.Cart{},
    &models.CartItem{},
    &models.Order{},
    &models.OrderItem{},

)

if err != nil {
	log.Fatal("Migration failed:",err)
 }
 log.Println("Database migration complete")
}

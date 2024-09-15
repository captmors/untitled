package cfg

import (
	"log"
	"untitled/internal/interfaces"
	"untitled/internal/users"
	"untitled/internal/users/mdl"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init() *gin.Engine {
	r := gin.Default()

	db := initDB()
	initApps(r, db)

	return r
}

var (
	apps []interfaces.App
)

func initDB() *gorm.DB {
	DB, err := gorm.Open(postgres.Open(DatabaseUrl), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	if err := DB.AutoMigrate(&mdl.User{}); err != nil {
		log.Fatal("Error migrating database:", err)
	}

	return DB
}

func initApps(r *gin.Engine, db *gorm.DB) {
	usersApp := users.NewUserApp(db)
	apps = append(apps, usersApp)

	for _, app := range apps {
		app.Init(r)
	}
}

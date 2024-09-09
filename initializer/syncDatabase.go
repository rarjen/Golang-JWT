package initializer

import "go-jwt/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
}

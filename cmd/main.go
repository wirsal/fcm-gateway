package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/wirsal/fcm-gateway/api"
	"github.com/wirsal/fcm-gateway/fcm"
	"github.com/wirsal/fcm-gateway/internal/config"
)

func main() {

	cfg, err := config.LoadConfig("configs/")
	if err != nil {
		log.Fatalf("Gagal memuat konfigurasi: %v", err)
	}

	ctx := context.Background()

	fcmService, err := fcm.NewService(ctx, cfg.FCM.CredentialsFile, cfg.FCM.Scopes, cfg.FCM.EndpointURL)
	if err != nil {
		log.Fatalf("Gagal inisialisasi service FCM: %v", err)
	}

	apiHandler := api.NewHandler(fcmService)

	router := gin.Default()
	router.GET("/", apiHandler.Welcome)
	router.POST("/send", apiHandler.SendNotification)

	log.Printf("Server Gin berjalan di http://localhost:%s", cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Gagal menjalankan server Gin: %v", err)
	}

}

package main

import (
	"crypto-watcher/pkg/api"
	"crypto-watcher/pkg/db"
	"crypto-watcher/pkg/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	_ "crypto-watcher/docs" // ← мы это сгенерируем

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Crypto Watcher API
// @version 1.0
// @description API для отслеживания цен криптовалют

// @host localhost:8080
// @BasePath /

func main() {
	db.Init()
	go service.StartScheduler()

	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	api.SetupRoutes(r)
	// Регистрируем frontend
	r.Static("/ui", "/app/frontend")

	// Редирект с корня на UI
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/ui/")
	})
	log.Fatal(r.Run(":8080"))
}

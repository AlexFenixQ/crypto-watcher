package api

import (
	"crypto-watcher/pkg/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/currency/add", AddCurrency)
	r.POST("/currency/remove", RemoveCurrency)
	r.POST("/currency/price", GetPrice)
}

// AddCurrency godoc
// @Summary Добавить валюту для отслеживания
// @Description Добавляет новую криптовалюту в список отслеживания
// @Tags currency
// @Accept json
// @Produce json
// @Param request body model.CurrencyRequest true "Название валюты"
// @Success 200 {object} model.StatusResponse
// @Failure 400 {object} model.ErrorResponse
// @Router /currency/add [post]
func AddCurrency(c *gin.Context) {
	var body struct {
		Coin string `json:"coin"`
	}
	if c.BindJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	go db.TrackCurrency(body.Coin)
	c.JSON(http.StatusOK, gin.H{"status": "tracking started"})
}

// RemoveCurrency godoc
// @Summary Удалить валюту
// @Description Удаляет криптовалюту из отслеживания
// @Tags currency
// @Accept json
// @Produce json
// @Param request body model.CurrencyRequest true "Название валюты"
// @Success 200 {object} model.StatusResponse
// @Failure 400 {object} model.ErrorResponse
// @Router /currency/remove [post]
func RemoveCurrency(c *gin.Context) {
	var body struct {
		Coin string `json:"coin"`
	}
	if c.BindJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	db.UntrackCurrency(body.Coin)
	c.JSON(http.StatusOK, gin.H{"status": "tracking stopped"})
}

// GetPrice godoc
// @Summary Получить цену криптовалюты
// @Description Возвращает цену криптовалюты на заданный timestamp
// @Tags currency
// @Accept json
// @Produce json
// @Param request body model.PriceRequest true "Монета и метка времени"
// @Success 200 {object} model.PriceResponse
// @Failure 400 {object} model.ErrorResponse
// @Router /currency/price [post]
func GetPrice(c *gin.Context) {
	var body struct {
		Coin      string `json:"coin"`
		Timestamp int64  `json:"timestamp"`
	}
	if c.BindJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	price, ts := db.FindClosestPrice(body.Coin, body.Timestamp)
	c.JSON(http.StatusOK, gin.H{
		"coin":      body.Coin,
		"price":     price,
		"timestamp": ts,
	})
}

package main

import (
	"github.com/MJarmolkiewicz/kryptonim/handlers"
	"github.com/MJarmolkiewicz/kryptonim/services"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	curHandler := handlers.NewCurrencyHandler(services.NewRateService())

	router.GET("/rates", curHandler.GetFIATRates)
	router.GET("/exchange", curHandler.ExchangeCrypto)

	router.Run(":8080")
}

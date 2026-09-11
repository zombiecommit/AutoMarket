package main

import (
	"net/http"

	"AutoMarket/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "AutoMarket funcionando")
	})

	router.POST("/usuarios", handlers.RegistrarUsuario)
	router.POST("/login", handlers.IniciarSesion)

	router.Run(":8080")
}

package main

import (
	"AutoMarket/middleware"
	"AutoMarket/services"
	"net/http"

	"AutoMarket/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "AutoMarket funcionando")
	})

	router.POST(
		"/usuarios",
		middleware.RegistrarAccounting(
			services.OpRegistrarUsuario,
			"usuario",
		),
		handlers.RegistrarUsuario,
	)

	router.POST(
		"/login",
		middleware.RegistrarAccounting(
			services.OpIniciarSesion,
			"autenticacion",
		),
		handlers.IniciarSesion,
	)

	router.Run(":8080")
}

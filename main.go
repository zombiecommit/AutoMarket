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

	// NOTA: las rutas de registro (POST /usuarios) y login (POST /login) las
	// debe agregar el equipo de AAA aquí, apuntando a su propio handler que
	// use services.AutenticarUsuario y storage.GuardarUsuarios.

	// Visitante: acceso público, sin autenticación
	router.GET("/catalogo", handlers.ConsultarCatalogo)
	router.GET("/contacto", handlers.ConsultarContacto)

	// Vendedor: publicar y gestionar sus propios vehículos
	router.POST("/vehiculos", handlers.PublicarVehiculo)
	router.GET("/vendedores/:vendedorId/publicaciones", handlers.MisPublicaciones)
	router.PUT("/vehiculos/:id/vendido", handlers.ReportarVendido)

	// Administrador: moderar publicaciones y usuarios
	router.PUT("/vehiculos/:id/autorizar", handlers.AutorizarPublicacion)
	router.DELETE("/vehiculos/:id", handlers.EliminarPublicacion)
	router.DELETE("/usuarios/:id", handlers.EliminarUsuario)

	router.Run(":8080")
}

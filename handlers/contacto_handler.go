package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// BLOQUE: consultar informacion de contacto (visitante)
func ConsultarContacto(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"nombre_negocio": "AutoMarket",
		"correo":         "contacto@automarket.com",
		"telefono":       "+57 300 000 0000",
		"direccion":      "Cali, Valle del Cauca, Colombia",
	})
}

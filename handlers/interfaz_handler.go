package handlers

import (
	"net/http"

	"AutoMarket/services"

	"github.com/gin-gonic/gin"
)

func MostrarCatalogo(c *gin.Context) {
	catalogo, err := services.ObtenerCatalogo()

	if err != nil {
		c.HTML(http.StatusInternalServerError, "catalogo.html", gin.H{
			"error": "No se pudo obtener el catálogo",
		})
		return
	}

	c.HTML(http.StatusOK, "catalogo.html", gin.H{
		"vehiculos": catalogo,
	})
}

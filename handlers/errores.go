package handlers

import (
	"errors"
	"net/http"

	"AutoMarket/services"

	"github.com/gin-gonic/gin"
)

// BLOQUE: traducir un error de negocio a su código HTTP

func responderError(c *gin.Context, err error, mensajeInterno string) {
	switch {
	case services.EsErrorDeValidacion(err):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

	case errors.Is(err, services.ErrCredencialesInvalidas):
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})

	case errors.Is(err, services.ErrPublicacionAjena):
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})

	case errors.Is(err, services.ErrUsuarioNoEncontrado),
		errors.Is(err, services.ErrPublicacionNoEncontrada):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

	case errors.Is(err, services.ErrCorreoDuplicado),
		errors.Is(err, services.ErrPublicacionYaVendida),
		errors.Is(err, services.ErrPublicacionNoPendiente),
		errors.Is(err, services.ErrPublicacionNoPublicada),
		errors.Is(err, services.ErrUsuarioConPublicaciones):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})

	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": mensajeInterno})
	}
}

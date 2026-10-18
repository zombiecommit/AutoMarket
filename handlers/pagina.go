package handlers

import (
	"AutoMarket/middleware"
	"AutoMarket/services"

	"github.com/gin-gonic/gin"
)

// BLOQUE: datos comunes de toda pagina (rol y menu de navegacion)
func DatosDePagina(c *gin.Context, extra gin.H) gin.H {
	rol := middleware.RolDeContexto(c)

	datos := gin.H{
		"Rol":  rol,
		"Menu": services.MenuParaRol(rol),
	}

	for clave, valor := range extra {
		datos[clave] = valor
	}

	return datos
}

package middleware

import (
	"net/http"

	"AutoMarket/services"

	"github.com/gin-gonic/gin"
)

// BLOQUE: autorizacion por permiso
func RequerirPermiso(operacion string) gin.HandlerFunc {
	return func(c *gin.Context) {
		usuario, autenticado := ObtenerUsuarioDeContexto(c)

		if !autenticado {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "se requiere autenticación para realizar esta operación",
			})
			return
		}

		if !services.TienePermiso(usuario.Rol, operacion) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "el rol '" + usuario.Rol + "' no tiene permiso para realizar esta operación",
			})
			return
		}

		c.Next()
	}
}

// BLOQUE: autorizacion por rol
func RequerirRol(rolesPermitidos ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		usuario, autenticado := ObtenerUsuarioDeContexto(c)

		if !autenticado {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "se requiere autenticación para realizar esta operación",
			})
			return
		}

		for _, rol := range rolesPermitidos {
			if usuario.Rol == rol {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "el rol '" + usuario.Rol + "' no está autorizado para acceder a este recurso",
		})
	}
}

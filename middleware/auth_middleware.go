package middleware

import (
	"net/http"
	"strings"

	"AutoMarket/services"

	"github.com/gin-gonic/gin"
)

func Autenticar() gin.HandlerFunc {
	return func(c *gin.Context) {
		encabezado := c.GetHeader("Authorization")

		if encabezado == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "se requiere autenticación",
			})
			c.Abort()
			return
		}

		partes := strings.SplitN(encabezado, " ", 2)

		if len(partes) != 2 || partes[0] != "Bearer" || partes[1] == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "token de autenticación inválido",
			})
			c.Abort()
			return
		}

		token := partes[1]

		sesion, existe := services.ObtenerSesion(token)

		if !existe {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "token de autenticación inválido o sesión inexistente",
			})
			c.Abort()
			return
		}

		c.Set("usuario", sesion.Usuario)
		c.Next()
	}
}

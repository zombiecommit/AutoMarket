package middleware

import (
	"AutoMarket/models"

	"github.com/gin-gonic/gin"
)

// BLOQUE: llave del usuario en el contexto
const ClaveUsuarioContexto = "usuario"

// BLOQUE: obtener usuario autenticado del contexto
func ObtenerUsuarioDeContexto(c *gin.Context) (models.Usuario, bool) {
	valor, existe := c.Get(ClaveUsuarioContexto)
	if !existe {
		return models.Usuario{}, false
	}

	usuario, esUsuario := valor.(models.Usuario)
	if !esUsuario {
		return models.Usuario{}, false
	}

	return usuario, true
}

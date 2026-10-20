package middleware

import (
	"AutoMarket/models"
	"AutoMarket/services"

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

// BLOQUE: rol de la peticion actual
func RolDeContexto(c *gin.Context) string {
	usuario, autenticado := ObtenerUsuarioDeContexto(c)

	if !autenticado {
		return services.RolVisitante
	}

	return usuario.Rol
}

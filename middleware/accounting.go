package middleware

import (
	"AutoMarket/services"

	"github.com/gin-gonic/gin"
)

func RegistrarAccounting(operacion, recurso string) gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Next()

		usuario, autenticado := ObtenerUsuarioDeContexto(c)

		recursoID := c.Param("id")

		exito := c.Writer.Status() >= 200 && c.Writer.Status() < 400

		var err error

		if autenticado {
			err = services.RegistrarOperacion(
				usuario,
				operacion,
				recurso,
				recursoID,
				exito,
			)
		} else {
			err = services.RegistrarOperacionSinAutenticar(
				operacion,
				recurso,
				recursoID,
				exito,
			)
		}

		if err != nil {
			_ = c.Error(err)
		}
	}
}

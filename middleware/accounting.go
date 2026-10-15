package middleware

import (
	"fmt"

	"AutoMarket/services"

	"github.com/gin-gonic/gin"
)

func RegistrarAccounting(operacion, recurso string) gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Next()

		usuario, autenticado := ObtenerUsuarioDeContexto(c)

		recursoID := c.Param("id")

		if recursoID == "" {
			if id, existe := c.Get("accounting_recurso_id"); existe {
				recursoID, _ = id.(string)
			}
		}

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
			c.Error(fmt.Errorf("falló el registro de Accounting: %w", err))
		}
	}
}

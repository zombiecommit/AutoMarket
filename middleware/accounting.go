package middleware

import (
	"AutoMarket/services"

	"github.com/gin-gonic/gin"
)

func RegistrarAccounting(operacion, recurso string) gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Next() // Dejar que la petición continúe

		usuario, autenticado := ObtenerUsuarioDeContexto(c) // Intentamos obtener el usuario autenticado

		exito := c.Writer.Status() >= 200 && c.Writer.Status() < 400 // Considera exitosa una respuesta entre 200 y 399

		var err error

		if autenticado {
			err = services.RegistrarOperacion(
				usuario,
				operacion,
				recurso,
				exito,
			)
		} else {
			err = services.RegistrarOperacionSinAutenticar(
				operacion,
				recurso,
				exito,
			)
		}

		if err != nil {
			_ = c.Error(err)
		}
	}
}

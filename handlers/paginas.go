package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// BLOQUE: paginas HTML del frontend de vehiculos
//
// Las paginas son "cascarones": el navegador las pide sin token y luego
// su JavaScript (static/js/app.js) consulta la API REST con el token
// guardado en sessionStorage. Por eso estas rutas no usan Autenticar().

func PaginaCatalogo(c *gin.Context) {
	c.HTML(http.StatusOK, "catalogo.html", DatosDePagina(c, gin.H{
		"Titulo": "Catálogo",
	}))
}

func PaginaDetalleVehiculo(c *gin.Context) {
	c.HTML(http.StatusOK, "detalle_vehiculo.html", DatosDePagina(c, gin.H{
		"Titulo": "Detalle del vehículo",
	}))
}

func PaginaMisVehiculos(c *gin.Context) {
	c.HTML(http.StatusOK, "mis_vehiculos.html", DatosDePagina(c, gin.H{
		"Titulo": "Mis vehículos",
	}))
}

func PaginaPublicar(c *gin.Context) {
	c.HTML(http.StatusOK, "publicar.html", DatosDePagina(c, gin.H{
		"Titulo": "Publicar vehículo",
	}))
}

// BLOQUE: negociacion de contenido para GET /catalogo
// La misma ruta sirve JSON (API) y HTML (navegador). Si quien pide acepta
// text/html (navegacion normal del navegador) se devuelve la pagina y se
// corta la cadena, de modo que no se registra en Accounting como consulta
// de API. Cualquier otro cliente (fetch con Accept: application/json,
// curl, tests) sigue recibiendo JSON como antes.
func CatalogoComoPaginaSiEsNavegador() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.Contains(c.GetHeader("Accept"), "text/html") {
			PaginaCatalogo(c)
			c.Abort()
			return
		}

		c.Next()
	}
}

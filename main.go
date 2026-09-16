package main

import (
	"AutoMarket/middleware"
	"AutoMarket/services"
	"net/http"

	"AutoMarket/handlers"

	"github.com/gin-gonic/gin"
)

// BLOQUE: construir el router con todas las rutas de la aplicacion
// (se separa de main() para poder probarlo con httptest en los tests)
func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "AutoMarket funcionando")
	})

	router.POST(
		"/usuarios",
		middleware.RegistrarAccounting(
			services.OpRegistrarUsuario,
			"usuario",
		),
		handlers.RegistrarUsuario,
	)

	router.POST(
		"/login",
		middleware.RegistrarAccounting(
			services.OpIniciarSesion,
			"autenticacion",
		),
		handlers.IniciarSesion,
	)

	// BLOQUE: rutas de visitante (publicas, sin autenticacion)
	router.GET(
		"/catalogo",
		middleware.RegistrarAccounting(
			services.OpConsultarCatalogo,
			"vehiculo",
		),
		handlers.ConsultarCatalogo,
	)

	router.GET(
		"/contacto",
		middleware.RegistrarAccounting(
			services.OpConsultarContacto,
			"contacto",
		),
		handlers.ConsultarContacto,
	)

	// BLOQUE: rutas de vendedor (requieren autenticacion + permiso)
	router.POST(
		"/vehiculos",
		middleware.Autenticar(),
		middleware.RequerirPermiso(services.OpPublicarVehiculo),
		middleware.RegistrarAccounting(
			services.OpPublicarVehiculo,
			"vehiculo",
		),
		handlers.PublicarVehiculo,
	)

	router.GET(
		"/vehiculos/mios",
		middleware.Autenticar(),
		middleware.RequerirPermiso(services.OpVerPublicacionesPropias),
		middleware.RegistrarAccounting(
			services.OpVerPublicacionesPropias,
			"vehiculo",
		),
		handlers.VerMisPublicaciones,
	)

	router.PATCH(
		"/vehiculos/:id/vender",
		middleware.Autenticar(),
		middleware.RequerirPermiso(services.OpReportarVenta),
		middleware.RegistrarAccounting(
			services.OpReportarVenta,
			"vehiculo",
		),
		handlers.ReportarVenta,
	)

	// BLOQUE: rutas de administrador (requieren autenticacion + permiso)
	router.PATCH(
		"/vehiculos/:id/autorizar",
		middleware.Autenticar(),
		middleware.RequerirPermiso(services.OpAutorizarPublicacion),
		middleware.RegistrarAccounting(
			services.OpAutorizarPublicacion,
			"vehiculo",
		),
		handlers.AutorizarPublicacion,
	)

	router.DELETE(
		"/vehiculos/:id",
		middleware.Autenticar(),
		middleware.RequerirPermiso(services.OpEliminarPublicacion),
		middleware.RegistrarAccounting(
			services.OpEliminarPublicacion,
			"vehiculo",
		),
		handlers.EliminarPublicacion,
	)

	router.DELETE(
		"/usuarios/:id",
		middleware.Autenticar(),
		middleware.RequerirPermiso(services.OpEliminarUsuario),
		middleware.RegistrarAccounting(
			services.OpEliminarUsuario,
			"usuario",
		),
		handlers.EliminarUsuario,
	)

	return router
}

func main() {
	router := SetupRouter()
	router.Run(":8080")
}

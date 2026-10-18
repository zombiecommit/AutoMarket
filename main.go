package main

import (
	"AutoMarket/middleware"
	"AutoMarket/services"

	"AutoMarket/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.LoadHTMLGlob("templates/*.html")
	router.Static("/static", "./static")

	router.GET("/login", func(c *gin.Context) {
		c.HTML(200, "login.html", gin.H{})
	})

	router.GET("/registro", func(c *gin.Context) {
		c.HTML(200, "registro.html", gin.H{})
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

	// BLOQUE: rutas de visitante (públicas, sin autenticacion)
	router.GET(
		"/catalogo",
		handlers.CatalogoComoPaginaSiEsNavegador(),
		middleware.RegistrarAccounting(
			services.OpConsultarCatalogo,
			"vehiculo",
		),
		handlers.ConsultarCatalogo,
	)

	// BLOQUE: paginas HTML de vehiculos (los datos llegan por la API REST)
	router.GET("/catalogo/:id", handlers.PaginaDetalleVehiculo)
	router.GET("/mis-vehiculos", handlers.PaginaMisVehiculos)
	router.GET("/publicar", handlers.PaginaPublicar)

	// BLOQUE: rutas de vendedor (requieren autenticacion + permiso)
	router.POST(
		"/vehiculos",
		middleware.RegistrarAccounting(
			services.OpPublicarVehiculo,
			"vehiculo",
		),
		middleware.Autenticar(),
		middleware.RequerirPermiso(services.OpPublicarVehiculo),
		handlers.PublicarVehiculo,
	)

	router.GET(
		"/vehiculos/mios",
		middleware.RegistrarAccounting(
			services.OpVerPublicacionesPropias,
			"vehiculo",
		),
		middleware.Autenticar(),
		middleware.RequerirPermiso(services.OpVerPublicacionesPropias),
		handlers.VerMisPublicaciones,
	)

	router.PATCH(
		"/vehiculos/:id/vender",
		middleware.RegistrarAccounting(
			services.OpReportarVenta,
			"vehiculo",
		),
		middleware.Autenticar(),
		middleware.RequerirPermiso(services.OpReportarVenta),
		handlers.ReportarVenta,
	)

	// BLOQUE: rutas de administrador (requieren autenticacion + permiso)
	router.PATCH(
		"/vehiculos/:id/autorizar",
		middleware.RegistrarAccounting(
			services.OpAutorizarPublicacion,
			"vehiculo",
		),
		middleware.Autenticar(),
		middleware.RequerirPermiso(services.OpAutorizarPublicacion),
		handlers.AutorizarPublicacion,
	)

	router.DELETE(
		"/vehiculos/:id",
		middleware.RegistrarAccounting(
			services.OpEliminarPublicacion,
			"vehiculo",
		),
		middleware.Autenticar(),
		middleware.RequerirPermiso(services.OpEliminarPublicacion),
		handlers.EliminarPublicacion,
	)

	router.DELETE(
		"/usuarios/:id",
		middleware.RegistrarAccounting(
			services.OpEliminarUsuario,
			"usuario",
		),
		middleware.Autenticar(),
		middleware.RequerirPermiso(services.OpEliminarUsuario),
		handlers.EliminarUsuario,
	)

	return router
}

func main() {
	router := SetupRouter()
	router.Run(":8080")
}

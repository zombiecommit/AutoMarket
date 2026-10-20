package services

// BLOQUE: rutas de las paginas del menu
const (
	RutaCatalogo               = "/catalogo"
	RutaLogin                  = "/login"
	RutaRegistro               = "/registro"
	RutaMisVehiculos           = "/mis-vehiculos"
	RutaPublicarVehiculo       = "/publicar"
	RutaReportarVenta          = "/ventas"
	RutaAutorizarPublicaciones = "/admin/publicaciones"
	RutaEliminarPublicaciones  = "/admin/vehiculos"
	RutaEliminarUsuarios       = "/admin/usuarios"
)

// BLOQUE: elemento del menu de navegacion
type ItemMenu struct {
	Texto         string
	Ruta          string
	Operacion     string
	SoloVisitante bool
}

// BLOQUE: elementos del menu, en el orden en que se muestran
var itemsDelMenu = []ItemMenu{
	{Texto: "Catálogo", Ruta: RutaCatalogo},
	{Texto: "Iniciar sesión", Ruta: RutaLogin, SoloVisitante: true},
	{Texto: "Registrarse", Ruta: RutaRegistro, SoloVisitante: true},
	{Texto: "Mis vehículos", Ruta: RutaMisVehiculos, Operacion: OpVerPublicacionesPropias},
	{Texto: "Publicar vehículo", Ruta: RutaPublicarVehiculo, Operacion: OpPublicarVehiculo},
	{Texto: "Reportar venta", Ruta: RutaReportarVenta, Operacion: OpReportarVenta},
	{Texto: "Autorizar publicaciones", Ruta: RutaAutorizarPublicaciones, Operacion: OpAutorizarPublicacion},
	{Texto: "Eliminar publicaciones", Ruta: RutaEliminarPublicaciones, Operacion: OpEliminarPublicacion},
	{Texto: "Eliminar usuarios", Ruta: RutaEliminarUsuarios, Operacion: OpEliminarUsuario},
}

// BLOQUE: menu que corresponde a un rol
func MenuParaRol(rol string) []ItemMenu {
	menu := []ItemMenu{}

	for _, item := range itemsDelMenu {
		if itemVisiblePara(item, rol) {
			menu = append(menu, item)
		}
	}

	return menu
}

// BLOQUE: decidir si un elemento del menu se muestra a un rol
func itemVisiblePara(item ItemMenu, rol string) bool {
	if item.SoloVisitante {
		return rol == RolVisitante
	}

	if item.Operacion == "" {
		return true
	}

	return TienePermiso(rol, item.Operacion)
}

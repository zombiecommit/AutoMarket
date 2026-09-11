package services

//BLOQUE: roles del sistema
const (
	RolVisitante     = "visitante"
	RolVendedor      = "vendedor"
	RolAdministrador = "administrador"
)

//BLOQUE: operaciones protegidas
const (
	OpConsultarCatalogo       = "CONSULTAR_CATALOGO"
	OpConsultarContacto       = "CONSULTAR_CONTACTO"
	OpPublicarVehiculo        = "PUBLICAR_VEHICULO"
	OpVerPublicacionesPropias = "VER_PUBLICACIONES_PROPIAS"
	OpReportarVenta           = "REPORTAR_VENTA"
	OpAutorizarPublicacion    = "AUTORIZAR_PUBLICACION"
	OpEliminarPublicacion     = "ELIMINAR_PUBLICACION"
	OpEliminarUsuario         = "ELIMINAR_USUARIO"
)

//BLOQUE: mapa de permisos por rol
var permisosPorRol = map[string]map[string]bool{
	RolVisitante: {
		OpConsultarCatalogo: true,
		OpConsultarContacto: true,
	},
	RolVendedor: {
		OpPublicarVehiculo:        true,
		OpVerPublicacionesPropias: true,
		OpReportarVenta:           true,
	},
	RolAdministrador: {
		OpAutorizarPublicacion: true,
		OpEliminarPublicacion:  true,
		OpEliminarUsuario:      true,
	},
}

//BLOQUE: verificar permiso
func TienePermiso(rol, operacion string) bool {
	operacionesDelRol, existeRol := permisosPorRol[rol]
	if !existeRol {
		return false
	}

	return operacionesDelRol[operacion]
}

//BLOQUE: verificar rol valido
func EsRolValido(rol string) bool {
	_, existe := permisosPorRol[rol]
	return existe
}

package services_test

import (
	"testing"

	"AutoMarket/services"
)

// BLOQUE: test vendedor sin permiso de administrador
func TestVendedorNoTienePermisoDeAdministrador(t *testing.T) {
	if services.TienePermiso(services.RolVendedor, services.OpEliminarUsuario) {
		t.Error("un vendedor no debería tener permiso para eliminar usuarios")
	}
}

// BLOQUE: test administrador con permiso de administrador
func TestAdministradorTienePermisoDeAdministrador(t *testing.T) {
	if !services.TienePermiso(services.RolAdministrador, services.OpEliminarUsuario) {
		t.Error("un administrador debería tener permiso para eliminar usuarios")
	}
}

// BLOQUE: test visitante solo con permisos publicos
func TestVisitanteSoloTienePermisosPublicos(t *testing.T) {
	if !services.TienePermiso(services.RolVisitante, services.OpConsultarCatalogo) {
		t.Error("un visitante debería poder consultar el catálogo")
	}

	if services.TienePermiso(services.RolVisitante, services.OpPublicarVehiculo) {
		t.Error("un visitante no debería poder publicar un vehículo")
	}
}

// BLOQUE: test rol inexistente sin permisos
func TestRolInexistenteNoTieneNingunPermiso(t *testing.T) {
	if services.TienePermiso("rol_que_no_existe", services.OpConsultarCatalogo) {
		t.Error("un rol inexistente no debería tener ningún permiso")
	}

	if services.EsRolValido("rol_que_no_existe") {
		t.Error("'rol_que_no_existe' no debería considerarse un rol válido")
	}
}

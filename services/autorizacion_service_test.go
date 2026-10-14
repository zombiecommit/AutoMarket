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

}

// BLOQUE: matriz completa de permisos (3 roles x 7 operaciones)

func TestMatrizDePermisosPorRol(t *testing.T) {
	operaciones := []string{
		services.OpConsultarCatalogo,
		services.OpPublicarVehiculo,
		services.OpVerPublicacionesPropias,
		services.OpReportarVenta,
		services.OpAutorizarPublicacion,
		services.OpEliminarPublicacion,
		services.OpEliminarUsuario,
	}

	// true = el rol SI puede ejecutar la operacion; lo que no aparece = denegado
	esperado := map[string]map[string]bool{
		services.RolVisitante: {
			services.OpConsultarCatalogo: true,
		},
		services.RolVendedor: {
			services.OpPublicarVehiculo:        true,
			services.OpVerPublicacionesPropias: true,
			services.OpReportarVenta:           true,
		},
		services.RolAdministrador: {
			services.OpAutorizarPublicacion: true,
			services.OpEliminarPublicacion:  true,
			services.OpEliminarUsuario:      true,
		},
	}

	for rol, permitidas := range esperado {
		for _, operacion := range operaciones {
			t.Run(rol+" / "+operacion, func(t *testing.T) {
				obtenido := services.TienePermiso(rol, operacion)

				if obtenido != permitidas[operacion] {
					t.Errorf(
						"rol %q, operación %q: se esperaba %v y se obtuvo %v",
						rol, operacion, permitidas[operacion], obtenido,
					)
				}
			})
		}
	}
}

package services_test

import (
	"strings"
	"testing"

	"AutoMarket/services"
)

// BLOQUE: obtener solo los textos de un menu
func textosDelMenu(menu []services.ItemMenu) []string {
	textos := []string{}

	for _, item := range menu {
		textos = append(textos, item.Texto)
	}

	return textos
}

// BLOQUE: test menu exacto segun el rol
func TestMenuSegunRol(t *testing.T) {
	esperado := map[string][]string{
		services.RolVisitante: {
			"Catálogo", "Iniciar sesión", "Registrarse",
		},
		services.RolVendedor: {
			"Catálogo", "Mis vehículos", "Publicar vehículo", "Reportar venta",
		},
		services.RolAdministrador: {
			"Catálogo", "Autorizar publicaciones", "Eliminar publicaciones", "Eliminar usuarios",
		},
	}

	for rol, textos := range esperado {
		t.Run(rol, func(t *testing.T) {
			obtenidos := textosDelMenu(services.MenuParaRol(rol))

			if strings.Join(obtenidos, " | ") != strings.Join(textos, " | ") {
				t.Errorf("rol %q: se esperaba %v y se obtuvo %v", rol, textos, obtenidos)
			}
		})
	}
}

// BLOQUE: test rol desconocido solo ve el catalogo publico
func TestRolDesconocidoSoloVeElCatalogoEnElMenu(t *testing.T) {
	obtenidos := textosDelMenu(services.MenuParaRol("rol_inventado"))

	if len(obtenidos) != 1 || obtenidos[0] != "Catálogo" {
		t.Errorf("un rol desconocido solo debería ver el catálogo y se obtuvo %v", obtenidos)
	}
}

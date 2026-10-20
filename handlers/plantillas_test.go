package handlers_test

import (
	"bytes"
	"html/template"
	"net/http/httptest"
	"strings"
	"testing"

	"AutoMarket/handlers"
	"AutoMarket/middleware"
	"AutoMarket/models"
	"AutoMarket/services"

	"github.com/gin-gonic/gin"
)

// BLOQUE: dibujar el menu de navegacion para un rol
func renderizarNav(t *testing.T, rol string) string {
	t.Helper()

	plantilla, err := template.ParseFiles("../templates/nav.html")
	if err != nil {
		t.Fatalf("no se pudo cargar templates/nav.html: %v", err)
	}

	datos := map[string]any{
		"Rol":  rol,
		"Menu": services.MenuParaRol(rol),
	}

	var salida bytes.Buffer

	if err := plantilla.ExecuteTemplate(&salida, "nav", datos); err != nil {
		t.Fatalf("no se pudo dibujar el menú: %v", err)
	}

	return salida.String()
}

// BLOQUE: test del menu dibujado segun el rol
func TestMenuDeNavegacionSegunRol(t *testing.T) {
	casos := []struct {
		rol           string
		debeMostrar   []string
		noDebeMostrar []string
	}{
		{
			rol:           services.RolVisitante,
			debeMostrar:   []string{"Catálogo", "Iniciar sesión", "Registrarse"},
			noDebeMostrar: []string{"Publicar vehículo", "Autorizar publicaciones", "Eliminar usuarios"},
		},
		{
			rol:           services.RolVendedor,
			debeMostrar:   []string{"Catálogo", "Mis vehículos", "Publicar vehículo", "Reportar venta"},
			noDebeMostrar: []string{"Iniciar sesión", "Registrarse", "Autorizar publicaciones", "Eliminar usuarios"},
		},
		{
			rol:           services.RolAdministrador,
			debeMostrar:   []string{"Catálogo", "Autorizar publicaciones", "Eliminar publicaciones", "Eliminar usuarios"},
			noDebeMostrar: []string{"Iniciar sesión", "Registrarse", "Publicar vehículo", "Reportar venta"},
		},
	}

	for _, caso := range casos {
		t.Run(caso.rol, func(t *testing.T) {
			pagina := renderizarNav(t, caso.rol)

			for _, texto := range caso.debeMostrar {
				if !strings.Contains(pagina, texto) {
					t.Errorf("el menú de %q debería mostrar %q", caso.rol, texto)
				}
			}

			for _, texto := range caso.noDebeMostrar {
				if strings.Contains(pagina, texto) {
					t.Errorf("el menú de %q no debería mostrar %q", caso.rol, texto)
				}
			}
		})
	}
}

// BLOQUE: test datos de pagina segun el usuario del contexto
func TestDatosDePaginaUsaElRolDelContexto(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("sin usuario autenticado el rol es visitante", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())

		datos := handlers.DatosDePagina(c, nil)

		if datos["Rol"] != services.RolVisitante {
			t.Errorf("se esperaba el rol visitante y se obtuvo %v", datos["Rol"])
		}
	})

	t.Run("con usuario autenticado toma su rol y conserva los datos extra", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Set(middleware.ClaveUsuarioContexto, models.Usuario{ID: 1, Nombre: "Ana", Rol: services.RolVendedor})

		datos := handlers.DatosDePagina(c, gin.H{"Titulo": "Catálogo"})

		if datos["Rol"] != services.RolVendedor {
			t.Errorf("se esperaba el rol vendedor y se obtuvo %v", datos["Rol"])
		}

		if datos["Titulo"] != "Catálogo" {
			t.Errorf("los datos extra deberían conservarse y se obtuvo %v", datos["Titulo"])
		}

		menu, esMenu := datos["Menu"].([]services.ItemMenu)
		if !esMenu || len(menu) == 0 {
			t.Errorf("los datos deberían incluir el menú del rol y se obtuvo %v", datos["Menu"])
		}
	})
}

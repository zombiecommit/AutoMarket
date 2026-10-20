package main

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"AutoMarket/services"

	"github.com/gin-gonic/gin"
)

// ---------------------------------------------------------------------
// BLOQUE: helpers propios de estas pruebas (reutilizan los de main_test.go)
// ---------------------------------------------------------------------

func publicarVehiculoDePrueba(t *testing.T, router *gin.Engine, token string) int {
	t.Helper()

	respuesta := hacerPeticion(router, http.MethodPost, "/vehiculos", token, map[string]any{
		"marca": "Mazda", "modelo": "3", "anio": 2020, "precio": 60000000,
	})

	if respuesta.Code != http.StatusCreated {
		t.Fatalf("no se pudo publicar el vehículo: %d %s", respuesta.Code, respuesta.Body.String())
	}

	return int(decodificar(t, respuesta)["vehiculo"].(map[string]any)["id"].(float64))
}

func misPublicaciones(t *testing.T, router *gin.Engine, token string) []any {
	t.Helper()

	respuesta := hacerPeticion(router, http.MethodGet, "/vehiculos/mios", token, nil)

	if respuesta.Code != http.StatusOK {
		t.Fatalf("no se pudieron consultar las publicaciones: %d %s", respuesta.Code, respuesta.Body.String())
	}

	return decodificar(t, respuesta)["vehiculos"].([]any)
}

func estadoDePublicacionPropia(t *testing.T, router *gin.Engine, token string, id int) string {
	t.Helper()

	for _, item := range misPublicaciones(t, router, token) {
		vehiculo := item.(map[string]any)

		if vehiculo["id"] == float64(id) {
			return vehiculo["estado"].(string)
		}
	}

	t.Fatalf("la publicación %d no está entre las publicaciones propias", id)

	return ""
}

func buscarEnCatalogo(t *testing.T, router *gin.Engine, id int) (map[string]any, bool) {
	t.Helper()

	respuesta := hacerPeticion(router, http.MethodGet, "/catalogo", "", nil)

	if respuesta.Code != http.StatusOK {
		t.Fatalf("no se pudo consultar el catálogo: %d %s", respuesta.Code, respuesta.Body.String())
	}

	for _, item := range decodificar(t, respuesta)["vehiculos"].([]any) {
		vehiculo := item.(map[string]any)

		if vehiculo["id"] == float64(id) {
			return vehiculo, true
		}
	}

	return nil, false
}

func rutaVehiculo(id int, accion string) string {
	if accion == "" {
		return fmt.Sprintf("/vehiculos/%d", id)
	}

	return fmt.Sprintf("/vehiculos/%d/%s", id, accion)
}

// ---------------------------------------------------------------------
// BLOQUE: transiciones de estado, validaciones y coherencia usuarios/publicaciones
// ---------------------------------------------------------------------

func TestReglasDeNegocioDeVehiculos(t *testing.T) {
	gin.SetMode(gin.TestMode)
	snapshotYRestaurarAlFinal(t)

	router := SetupRouter()
	pid := os.Getpid()

	tokenVendedor, _ := registrarYLoguear(t, router, fmt.Sprintf("reglas.vendedor.%d@automarket.test", pid))
	tokenOtroVendedor, _ := registrarYLoguear(t, router, fmt.Sprintf("reglas.otro.%d@automarket.test", pid))

	correoAdmin := fmt.Sprintf("reglas.admin.%d@automarket.test", pid)
	registrarYLoguear(t, router, correoAdmin)
	tokenAdmin := promoverAAdministrador(t, router, correoAdmin)

	// --- Punto 1: transiciones de estado ---

	t.Run("no se puede reportar como vendida una publicación pendiente (409)", func(t *testing.T) {
		id := publicarVehiculoDePrueba(t, router, tokenVendedor)

		respuesta := hacerPeticion(router, http.MethodPatch, rutaVehiculo(id, "vender"), tokenVendedor, nil)

		if respuesta.Code != http.StatusConflict {
			t.Fatalf("se esperaba 409, se obtuvo %d: %s", respuesta.Code, respuesta.Body.String())
		}

		comprobarUltimoAccounting(t, services.OpReportarVenta, "vehiculo", fmt.Sprintf("%d", id), false)

		if estado := estadoDePublicacionPropia(t, router, tokenVendedor, id); estado != services.EstadoPendiente {
			t.Errorf("la publicación debía seguir pendiente, quedó en %q", estado)
		}
	})

	t.Run("reportar la venta de una publicación ajena pendiente sigue siendo 403", func(t *testing.T) {
		id := publicarVehiculoDePrueba(t, router, tokenVendedor)

		respuesta := hacerPeticion(router, http.MethodPatch, rutaVehiculo(id, "vender"), tokenOtroVendedor, nil)

		if respuesta.Code != http.StatusForbidden {
			t.Fatalf("se esperaba 403, se obtuvo %d: %s", respuesta.Code, respuesta.Body.String())
		}
	})

	t.Run("flujo completo: pendiente, autorizada y recién entonces vendida", func(t *testing.T) {
		id := publicarVehiculoDePrueba(t, router, tokenVendedor)

		if respuesta := hacerPeticion(router, http.MethodPatch, rutaVehiculo(id, "autorizar"), tokenAdmin, nil); respuesta.Code != http.StatusOK {
			t.Fatalf("autorizar: se esperaba 200, se obtuvo %d: %s", respuesta.Code, respuesta.Body.String())
		}

		if respuesta := hacerPeticion(router, http.MethodPatch, rutaVehiculo(id, "vender"), tokenVendedor, nil); respuesta.Code != http.StatusOK {
			t.Fatalf("vender: se esperaba 200, se obtuvo %d: %s", respuesta.Code, respuesta.Body.String())
		}

		if estado := estadoDePublicacionPropia(t, router, tokenVendedor, id); estado != services.EstadoVendido {
			t.Errorf("se esperaba estado %q, se obtuvo %q", services.EstadoVendido, estado)
		}
	})

	t.Run("autorizar una publicación ya publicada devuelve 409", func(t *testing.T) {
		id := publicarVehiculoDePrueba(t, router, tokenVendedor)

		hacerPeticion(router, http.MethodPatch, rutaVehiculo(id, "autorizar"), tokenAdmin, nil)
		respuesta := hacerPeticion(router, http.MethodPatch, rutaVehiculo(id, "autorizar"), tokenAdmin, nil)

		if respuesta.Code != http.StatusConflict {
			t.Fatalf("se esperaba 409, se obtuvo %d: %s", respuesta.Code, respuesta.Body.String())
		}
	})

	// --- Punto 4: validaciones de entrada ---

	t.Run("publicar con datos inválidos devuelve 400 y no guarda nada", func(t *testing.T) {
		antes := len(misPublicaciones(t, router, tokenVendedor))

		casos := []struct {
			nombre string
			cuerpo map[string]any
		}{
			{"marca vacía", map[string]any{"marca": "", "modelo": "3", "anio": 2020, "precio": 1000}},
			{"marca solo con espacios", map[string]any{"marca": "   ", "modelo": "3", "anio": 2020, "precio": 1000}},
			{"modelo vacío", map[string]any{"marca": "Mazda", "modelo": "", "anio": 2020, "precio": 1000}},
			{"año cero", map[string]any{"marca": "Mazda", "modelo": "3", "anio": 0, "precio": 1000}},
			{"año negativo", map[string]any{"marca": "Mazda", "modelo": "3", "anio": -5, "precio": 1000}},
			{"precio cero", map[string]any{"marca": "Mazda", "modelo": "3", "anio": 2020, "precio": 0}},
			{"precio negativo", map[string]any{"marca": "Mazda", "modelo": "3", "anio": 2020, "precio": -1000}},
		}

		for _, caso := range casos {
			respuesta := hacerPeticion(router, http.MethodPost, "/vehiculos", tokenVendedor, caso.cuerpo)

			if respuesta.Code != http.StatusBadRequest {
				t.Errorf("%s: se esperaba 400, se obtuvo %d: %s", caso.nombre, respuesta.Code, respuesta.Body.String())
			}
		}

		if despues := len(misPublicaciones(t, router, tokenVendedor)); despues != antes {
			t.Errorf("las solicitudes inválidas no deben guardarse: había %d publicaciones y ahora hay %d", antes, despues)
		}
	})

	t.Run("un identificador no numérico devuelve 400", func(t *testing.T) {
		casos := []struct {
			nombre string
			metodo string
			ruta   string
			token  string
		}{
			{"vender", http.MethodPatch, "/vehiculos/abc/vender", tokenVendedor},
			{"autorizar", http.MethodPatch, "/vehiculos/abc/autorizar", tokenAdmin},
			{"eliminar publicación", http.MethodDelete, "/vehiculos/abc", tokenAdmin},
			{"eliminar usuario", http.MethodDelete, "/usuarios/abc", tokenAdmin},
		}

		for _, caso := range casos {
			respuesta := hacerPeticion(router, caso.metodo, caso.ruta, caso.token, nil)

			if respuesta.Code != http.StatusBadRequest {
				t.Errorf("%s: se esperaba 400, se obtuvo %d: %s", caso.nombre, respuesta.Code, respuesta.Body.String())
			}
		}
	})

	// --- Punto 2 y 3: usuarios, publicaciones y datos de contacto ---

	t.Run("el catálogo incluye nombre y correo del vendedor", func(t *testing.T) {
		correo := fmt.Sprintf("reglas.contacto.%d@automarket.test", pid)
		token, _ := registrarYLoguear(t, router, correo)

		id := publicarVehiculoDePrueba(t, router, token)
		hacerPeticion(router, http.MethodPatch, rutaVehiculo(id, "autorizar"), tokenAdmin, nil)

		vehiculo, encontrado := buscarEnCatalogo(t, router, id)
		if !encontrado {
			t.Fatal("el vehículo autorizado debería estar en el catálogo")
		}

		vendedor := vehiculo["vendedor"].(map[string]any)

		if vendedor["correo"] != correo {
			t.Errorf("correo de contacto esperado %q, obtenido %v", correo, vendedor["correo"])
		}

		if vendedor["nombre"] != "Usuario de Prueba" {
			t.Errorf("nombre de contacto esperado %q, obtenido %v", "Usuario de Prueba", vendedor["nombre"])
		}
	})

	t.Run("no se puede eliminar un usuario con publicaciones (409) y el catálogo no queda huérfano", func(t *testing.T) {
		correo := fmt.Sprintf("reglas.eliminar.%d@automarket.test", pid)
		token, idUsuario := registrarYLoguear(t, router, correo)
		rutaUsuario := fmt.Sprintf("/usuarios/%d", int(idUsuario))

		id := publicarVehiculoDePrueba(t, router, token)
		hacerPeticion(router, http.MethodPatch, rutaVehiculo(id, "autorizar"), tokenAdmin, nil)

		// Con una publicación publicada: se rechaza.
		respuesta := hacerPeticion(router, http.MethodDelete, rutaUsuario, tokenAdmin, nil)

		if respuesta.Code != http.StatusConflict {
			t.Fatalf("se esperaba 409, se obtuvo %d: %s", respuesta.Code, respuesta.Body.String())
		}

		comprobarUltimoAccounting(t, services.OpEliminarUsuario, "usuario", fmt.Sprintf("%d", int(idUsuario)), false)

		// El usuario no se tocó: su sesión sigue válida y el catálogo conserva su contacto.
		misPublicaciones(t, router, token)

		vehiculo, encontrado := buscarEnCatalogo(t, router, id)
		if !encontrado {
			t.Fatal("la publicación debía seguir en el catálogo")
		}

		if vehiculo["vendedor"].(map[string]any)["correo"] != correo {
			t.Error("la publicación perdió los datos de contacto del vendedor")
		}

		// Una publicación vendida también cuenta: sigue asociada al usuario.
		hacerPeticion(router, http.MethodPatch, rutaVehiculo(id, "vender"), token, nil)

		respuesta = hacerPeticion(router, http.MethodDelete, rutaUsuario, tokenAdmin, nil)

		if respuesta.Code != http.StatusConflict {
			t.Fatalf("con una publicación vendida se esperaba 409, se obtuvo %d: %s", respuesta.Code, respuesta.Body.String())
		}

		// Eliminadas sus publicaciones, ya se puede eliminar al usuario.
		if respuesta = hacerPeticion(router, http.MethodDelete, rutaVehiculo(id, ""), tokenAdmin, nil); respuesta.Code != http.StatusOK {
			t.Fatalf("eliminar publicación: se esperaba 200, se obtuvo %d: %s", respuesta.Code, respuesta.Body.String())
		}

		if respuesta = hacerPeticion(router, http.MethodDelete, rutaUsuario, tokenAdmin, nil); respuesta.Code != http.StatusOK {
			t.Fatalf("eliminar usuario: se esperaba 200, se obtuvo %d: %s", respuesta.Code, respuesta.Body.String())
		}

		if _, encontrado := buscarEnCatalogo(t, router, id); encontrado {
			t.Error("la publicación eliminada no debería estar en el catálogo")
		}
	})
}

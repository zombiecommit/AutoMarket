package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func pedirConAccept(router *gin.Engine, ruta string, accept string) *httptest.ResponseRecorder {
	peticion := httptest.NewRequest(http.MethodGet, ruta, nil)

	if accept != "" {
		peticion.Header.Set("Accept", accept)
	}

	respuesta := httptest.NewRecorder()
	router.ServeHTTP(respuesta, peticion)

	return respuesta
}

// BLOQUE: las paginas HTML de vehiculos se sirven y cada una tiene su contenido
func TestPaginasHTMLDeVehiculos(t *testing.T) {
	gin.SetMode(gin.TestMode)
	snapshotYRestaurarAlFinal(t)

	router := SetupRouter()

	casos := []struct {
		nombre   string
		ruta     string
		accept   string
		contiene string
	}{
		{"catálogo para el navegador", "/catalogo", "text/html,application/xhtml+xml", "Catálogo de vehículos"},
		{"detalle de un vehículo", "/catalogo/1", "text/html", "Volver al catálogo"},
		{"mis vehículos", "/mis-vehiculos", "text/html", "Mis vehículos"},
		{"publicar vehículo", "/publicar", "text/html", "Publicar vehículo"},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			respuesta := pedirConAccept(router, caso.ruta, caso.accept)

			if respuesta.Code != http.StatusOK {
				t.Fatalf("se esperaba 200, se obtuvo %d", respuesta.Code)
			}

			if !strings.Contains(respuesta.Header().Get("Content-Type"), "text/html") {
				t.Errorf("se esperaba text/html, se obtuvo %q", respuesta.Header().Get("Content-Type"))
			}

			if !strings.Contains(respuesta.Body.String(), caso.contiene) {
				t.Errorf("la página debería contener %q", caso.contiene)
			}
		})
	}
}

// BLOQUE: GET /catalogo sigue devolviendo JSON a quien no es un navegador
func TestCatalogoSigueDevolviendoJSONParaLaAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)
	snapshotYRestaurarAlFinal(t)

	router := SetupRouter()

	for _, accept := range []string{"", "application/json", "*/*"} {
		respuesta := pedirConAccept(router, "/catalogo", accept)

		if respuesta.Code != http.StatusOK {
			t.Fatalf("Accept %q: se esperaba 200, se obtuvo %d", accept, respuesta.Code)
		}

		if !strings.Contains(respuesta.Header().Get("Content-Type"), "application/json") {
			t.Errorf("Accept %q: se esperaba JSON, se obtuvo %q", accept, respuesta.Header().Get("Content-Type"))
		}

		if _, existe := decodificar(t, respuesta)["vehiculos"]; !existe {
			t.Errorf("Accept %q: la respuesta debería incluir la clave \"vehiculos\"", accept)
		}
	}
}

// BLOQUE: los archivos estaticos (CSS, JS, imagen por defecto) se sirven
func TestArchivosEstaticos(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := SetupRouter()

	for _, ruta := range []string{
		"/static/css/style.css",
		"/static/js/app.js",
		"/static/images/sin-foto.svg",
	} {
		if respuesta := pedirConAccept(router, ruta, ""); respuesta.Code != http.StatusOK {
			t.Errorf("%s: se esperaba 200, se obtuvo %d", ruta, respuesta.Code)
		}
	}
}

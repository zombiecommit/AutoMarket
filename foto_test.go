package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"AutoMarket/services"

	"github.com/gin-gonic/gin"
)

// PNG válido de 1x1 píxel.
const pngMinimoBase64 = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg=="

func publicarConFoto(
	router *gin.Engine,
	token string,
	campos map[string]string,
	nombreArchivo string,
	contenido []byte,
) *httptest.ResponseRecorder {
	var cuerpo bytes.Buffer
	escritor := multipart.NewWriter(&cuerpo)

	for clave, valor := range campos {
		_ = escritor.WriteField(clave, valor)
	}

	if nombreArchivo != "" {
		parte, _ := escritor.CreateFormFile("foto", nombreArchivo)
		_, _ = parte.Write(contenido)
	}

	_ = escritor.Close()

	peticion := httptest.NewRequest(http.MethodPost, "/vehiculos", &cuerpo)
	peticion.Header.Set("Content-Type", escritor.FormDataContentType())
	peticion.Header.Set("Authorization", "Bearer "+token)

	respuesta := httptest.NewRecorder()
	router.ServeHTTP(respuesta, peticion)

	return respuesta
}

func camposValidos() map[string]string {
	return map[string]string{
		"marca": "Mazda", "modelo": "3", "anio": "2020",
		"precio": "60000000", "descripcion": "Full equipo",
	}
}

// BLOQUE: publicar un vehículo subiendo su foto
func TestPublicarVehiculoConFoto(t *testing.T) {
	gin.SetMode(gin.TestMode)
	snapshotYRestaurarAlFinal(t)

	router := SetupRouter()
	token, _ := registrarYLoguear(t, router, fmt.Sprintf("foto.vendedor.%d@automarket.test", os.Getpid()))

	png, err := base64.StdEncoding.DecodeString(pngMinimoBase64)
	if err != nil {
		t.Fatalf("PNG de prueba inválido: %v", err)
	}

	t.Run("foto válida: 201, queda guardada y se sirve por /static", func(t *testing.T) {
		respuesta := publicarConFoto(router, token, camposValidos(), "mi-carro.png", png)

		if respuesta.Code != http.StatusCreated {
			t.Fatalf("se esperaba 201, se obtuvo %d: %s", respuesta.Code, respuesta.Body.String())
		}

		vehiculo := decodificar(t, respuesta)["vehiculo"].(map[string]any)
		id := int(vehiculo["id"].(float64))
		rutaFoto, _ := vehiculo["foto"].(string)

		esperada := fmt.Sprintf("/static/images/vehiculos/%d.png", id)
		if rutaFoto != esperada {
			t.Fatalf("foto esperada %q, obtenida %q", esperada, rutaFoto)
		}

		archivo := filepath.Join("static", "images", "vehiculos", fmt.Sprintf("%d.png", id))
		t.Cleanup(func() { _ = os.Remove(archivo) })

		if _, err := os.Stat(archivo); err != nil {
			t.Fatalf("el archivo de la foto no se guardó: %v", err)
		}

		if servida := pedirConAccept(router, rutaFoto, ""); servida.Code != http.StatusOK {
			t.Errorf("la foto debería servirse en %s, código %d", rutaFoto, servida.Code)
		}

		// La foto aparece entre las publicaciones propias.
		for _, item := range misPublicaciones(t, router, token) {
			propio := item.(map[string]any)
			if propio["id"] == float64(id) && propio["foto"] != esperada {
				t.Errorf("la publicación propia debería tener foto %q, tiene %v", esperada, propio["foto"])
			}
		}
	})

	t.Run("sin foto: 400 y no se guarda la publicación", func(t *testing.T) {
		antes := len(misPublicaciones(t, router, token))

		respuesta := publicarConFoto(router, token, camposValidos(), "", nil)

		if respuesta.Code != http.StatusBadRequest {
			t.Fatalf("se esperaba 400, se obtuvo %d: %s", respuesta.Code, respuesta.Body.String())
		}

		if despues := len(misPublicaciones(t, router, token)); despues != antes {
			t.Errorf("no debía guardarse nada: había %d y ahora hay %d", antes, despues)
		}
	})

	t.Run("archivo que no es imagen: 400 y no se guarda la publicación", func(t *testing.T) {
		antes := len(misPublicaciones(t, router, token))

		respuesta := publicarConFoto(router, token, camposValidos(), "foto.png", []byte("esto es texto, no una imagen"))

		if respuesta.Code != http.StatusBadRequest {
			t.Fatalf("se esperaba 400, se obtuvo %d: %s", respuesta.Code, respuesta.Body.String())
		}

		if despues := len(misPublicaciones(t, router, token)); despues != antes {
			t.Errorf("no debía guardarse nada: había %d y ahora hay %d", antes, despues)
		}
	})

	t.Run("foto de más de 5 MB: 400", func(t *testing.T) {
		grande := append([]byte{}, png...)
		grande = append(grande, bytes.Repeat([]byte{0}, services.MaxFotoBytes)...)

		respuesta := publicarConFoto(router, token, camposValidos(), "grande.png", grande)

		if respuesta.Code != http.StatusBadRequest {
			t.Fatalf("se esperaba 400, se obtuvo %d: %s", respuesta.Code, respuesta.Body.String())
		}
	})

	t.Run("datos inválidos con foto válida: 400 y no queda archivo huérfano", func(t *testing.T) {
		campos := camposValidos()
		campos["marca"] = ""

		antes := len(misPublicaciones(t, router, token))

		respuesta := publicarConFoto(router, token, campos, "mi-carro.png", png)

		if respuesta.Code != http.StatusBadRequest {
			t.Fatalf("se esperaba 400, se obtuvo %d: %s", respuesta.Code, respuesta.Body.String())
		}

		if despues := len(misPublicaciones(t, router, token)); despues != antes {
			t.Errorf("no debía guardarse nada: había %d y ahora hay %d", antes, despues)
		}
	})

	t.Run("sin token no se puede publicar con foto", func(t *testing.T) {
		respuesta := publicarConFoto(router, "", camposValidos(), "mi-carro.png", png)

		if respuesta.Code != http.StatusUnauthorized {
			t.Fatalf("se esperaba 401, se obtuvo %d", respuesta.Code)
		}
	})

	t.Run("al eliminar la publicación se borra su foto", func(t *testing.T) {
		respuesta := publicarConFoto(router, token, camposValidos(), "mi-carro.png", png)
		if respuesta.Code != http.StatusCreated {
			t.Fatalf("no se pudo publicar: %d %s", respuesta.Code, respuesta.Body.String())
		}

		id := int(decodificar(t, respuesta)["vehiculo"].(map[string]any)["id"].(float64))
		archivo := filepath.Join("static", "images", "vehiculos", fmt.Sprintf("%d.png", id))
		t.Cleanup(func() { _ = os.Remove(archivo) })

		correoAdmin := fmt.Sprintf("foto.admin.%d@automarket.test", os.Getpid())
		registrarYLoguear(t, router, correoAdmin)
		tokenAdmin := promoverAAdministrador(t, router, correoAdmin)

		borrado := hacerPeticion(router, http.MethodDelete, rutaVehiculo(id, ""), tokenAdmin, nil)
		if borrado.Code != http.StatusOK {
			t.Fatalf("eliminar: se esperaba 200, se obtuvo %d: %s", borrado.Code, borrado.Body.String())
		}

		if _, err := os.Stat(archivo); !os.IsNotExist(err) {
			t.Errorf("el archivo de la foto debería haberse borrado (%s)", strings.TrimSpace(archivo))
		}
	})
}

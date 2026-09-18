package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

// ---------------------------------------------------------------------
// BLOQUE: utilidades para no dañar los datos reales de tus compañeros.
// Antes de correr los tests se guarda una "foto" de los archivos de
// almacenamiento, y al terminar (con defer) se restauran tal cual estaban.
// ---------------------------------------------------------------------

const (
	rutaUsuariosJSON   = "storage/usuarios.json"
	rutaVehiculosJSON  = "storage/vehiculos.json"
	rutaAccountingJSON = "accounting.json" // así está definido en storage/accounting.go
)

func snapshotArchivo(t *testing.T, ruta string) (existia bool, contenido []byte) {
	t.Helper()

	datos, err := os.ReadFile(ruta)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		t.Fatalf("no se pudo leer %s antes del test: %v", ruta, err)
	}

	return true, datos
}

func restaurarArchivo(t *testing.T, ruta string, existia bool, contenido []byte) {
	t.Helper()

	if existia {
		if err := os.WriteFile(ruta, contenido, 0644); err != nil {
			t.Errorf("no se pudo restaurar %s: %v", ruta, err)
		}
		return
	}

	// El archivo no existía antes del test (p. ej. accounting.json en un
	// clon nuevo del repo): lo eliminamos para dejar el repo como estaba.
	_ = os.Remove(ruta)
}

// snapshotYRestaurarAlFinal deja todo el almacenamiento (usuarios,
// vehiculos y accounting) exactamente como estaba antes de correr el test.
func snapshotYRestaurarAlFinal(t *testing.T) {
	t.Helper()

	usuariosExistia, usuariosContenido := snapshotArchivo(t, rutaUsuariosJSON)
	vehiculosExistia, vehiculosContenido := snapshotArchivo(t, rutaVehiculosJSON)
	accountingExistia, accountingContenido := snapshotArchivo(t, rutaAccountingJSON)

	t.Cleanup(func() {
		restaurarArchivo(t, rutaUsuariosJSON, usuariosExistia, usuariosContenido)
		restaurarArchivo(t, rutaVehiculosJSON, vehiculosExistia, vehiculosContenido)
		restaurarArchivo(t, rutaAccountingJSON, accountingExistia, accountingContenido)
	})
}

// ---------------------------------------------------------------------
// BLOQUE: helpers HTTP para hablar con el router real (httptest)
// ---------------------------------------------------------------------

func hacerPeticion(
	router *gin.Engine,
	metodo string,
	ruta string,
	token string,
	cuerpo any,
) *httptest.ResponseRecorder {
	var lector *bytes.Reader

	if cuerpo != nil {
		datos, _ := json.Marshal(cuerpo)
		lector = bytes.NewReader(datos)
	} else {
		lector = bytes.NewReader(nil)
	}

	peticion := httptest.NewRequest(metodo, ruta, lector)
	peticion.Header.Set("Content-Type", "application/json")

	if token != "" {
		peticion.Header.Set("Authorization", "Bearer "+token)
	}

	respuesta := httptest.NewRecorder()
	router.ServeHTTP(respuesta, peticion)

	return respuesta
}

func decodificar(t *testing.T, respuesta *httptest.ResponseRecorder) map[string]any {
	t.Helper()

	var cuerpo map[string]any

	if err := json.Unmarshal(respuesta.Body.Bytes(), &cuerpo); err != nil {
		t.Fatalf("la respuesta no es JSON válido: %v (cuerpo: %s)", err, respuesta.Body.String())
	}

	return cuerpo
}

// registrarYLoguear crea un usuario nuevo (rol "vendedor" por defecto en el
// registro) usando los endpoints reales /usuarios y /login, y devuelve su
// token y su id. Así se prueba, de paso, que esos endpoints de tus
// compañeros siguen funcionando igual que antes.
func registrarYLoguear(t *testing.T, router *gin.Engine, correo string) (token string, id float64) {
	t.Helper()

	respuestaRegistro := hacerPeticion(router, http.MethodPost, "/usuarios", "", map[string]string{
		"nombre":     "Usuario de Prueba",
		"correo":     correo,
		"contrasena": "clave12345",
	})

	if respuestaRegistro.Code != http.StatusCreated {
		t.Fatalf(
			"no se pudo registrar el usuario %s: código %d, cuerpo %s",
			correo, respuestaRegistro.Code, respuestaRegistro.Body.String(),
		)
	}

	usuarioCreado := decodificar(t, respuestaRegistro)["usuario"].(map[string]any)

	respuestaLogin := hacerPeticion(router, http.MethodPost, "/login", "", map[string]string{
		"correo":     correo,
		"contrasena": "clave12345",
	})

	if respuestaLogin.Code != http.StatusOK {
		t.Fatalf(
			"no se pudo iniciar sesión con %s: código %d, cuerpo %s",
			correo, respuestaLogin.Code, respuestaLogin.Body.String(),
		)
	}

	cuerpoLogin := decodificar(t, respuestaLogin)

	return cuerpoLogin["token"].(string), usuarioCreado["id"].(float64)
}

// promoverAAdministrador simula lo que haría un administrador de base de
// datos: convierte un usuario ya registrado en "administrador" editando el
// almacenamiento directamente, y le pide un nuevo login (los tokens creados
// antes de la promoción quedan con el rol viejo, como es esperable).
func promoverAAdministrador(t *testing.T, router *gin.Engine, correo string) (token string) {
	t.Helper()

	datos, err := os.ReadFile(rutaUsuariosJSON)
	if err != nil {
		t.Fatalf("no se pudo leer usuarios.json: %v", err)
	}

	var usuarios []map[string]any
	if err := json.Unmarshal(datos, &usuarios); err != nil {
		t.Fatalf("no se pudo parsear usuarios.json: %v", err)
	}

	encontrado := false
	for i, usuario := range usuarios {
		if usuario["correo"] == correo {
			usuarios[i]["rol"] = "administrador"
			encontrado = true
		}
	}

	if !encontrado {
		t.Fatalf("no se encontró el usuario %s para promoverlo a administrador", correo)
	}

	nuevoContenido, err := json.MarshalIndent(usuarios, "", "    ")
	if err != nil {
		t.Fatalf("no se pudo serializar usuarios.json: %v", err)
	}

	if err := os.WriteFile(rutaUsuariosJSON, nuevoContenido, 0644); err != nil {
		t.Fatalf("no se pudo escribir usuarios.json: %v", err)
	}

	respuestaLogin := hacerPeticion(router, http.MethodPost, "/login", "", map[string]string{
		"correo":     correo,
		"contrasena": "clave12345",
	})

	if respuestaLogin.Code != http.StatusOK {
		t.Fatalf("no se pudo re-loguear como administrador: %s", respuestaLogin.Body.String())
	}

	return decodificar(t, respuestaLogin)["token"].(string)
}

// ---------------------------------------------------------------------
// BLOQUE: test de integración de las funcionalidades de Enrique
// ---------------------------------------------------------------------

func TestFlujoDeNegocioDeVehiculos(t *testing.T) {
	gin.SetMode(gin.TestMode)
	snapshotYRestaurarAlFinal(t)

	router := SetupRouter()

	correoVendedor := fmt.Sprintf("vendedor.%d@automarket.test", os.Getpid())
	correoOtroVendedor := fmt.Sprintf("otrovendedor.%d@automarket.test", os.Getpid())
	correoAdmin := fmt.Sprintf("admin.%d@automarket.test", os.Getpid())

	tokenVendedor, _ := registrarYLoguear(t, router, correoVendedor)
	tokenOtroVendedor, _ := registrarYLoguear(t, router, correoOtroVendedor)

	// El registro público siempre crea usuarios con rol "vendedor", así que
	// para probar las operaciones de administrador primero se registra un
	// usuario normal y luego se "promueve" editando el almacenamiento
	// directamente (como lo haría un administrador de base de datos).
	registrarYLoguear(t, router, correoAdmin)
	tokenAdmin := promoverAAdministrador(t, router, correoAdmin)

	t.Run("visitante puede consultar catálogo sin autenticarse", func(t *testing.T) {
		respuesta := hacerPeticion(router, http.MethodGet, "/catalogo", "", nil)

		if respuesta.Code != http.StatusOK {
			t.Fatalf("se esperaba 200, se obtuvo %d: %s", respuesta.Code, respuesta.Body.String())
		}
	})

	t.Run("sin token no se puede publicar un vehículo", func(t *testing.T) {
		respuesta := hacerPeticion(router, http.MethodPost, "/vehiculos", "", map[string]any{
			"marca": "Mazda", "modelo": "3", "anio": 2020, "precio": 60000000,
		})

		if respuesta.Code != http.StatusUnauthorized {
			t.Fatalf("se esperaba 401, se obtuvo %d: %s", respuesta.Code, respuesta.Body.String())
		}
	})

	var idVehiculo float64

	t.Run("vendedor publica un vehículo y queda pendiente de aprobación", func(t *testing.T) {
		respuesta := hacerPeticion(router, http.MethodPost, "/vehiculos", tokenVendedor, map[string]any{
			"marca":       "Mazda",
			"modelo":      "3",
			"anio":        2020,
			"precio":      60000000,
			"descripcion": "Full equipo, único dueño",
		})

		if respuesta.Code != http.StatusCreated {
			t.Fatalf("se esperaba 201, se obtuvo %d: %s", respuesta.Code, respuesta.Body.String())
		}

		vehiculo := decodificar(t, respuesta)["vehiculo"].(map[string]any)

		if vehiculo["estado"] != "pendiente_aprobacion" {
			t.Errorf("se esperaba estado 'pendiente_aprobacion', se obtuvo %v", vehiculo["estado"])
		}

		idVehiculo = vehiculo["id"].(float64)
	})

	t.Run("el vehículo pendiente NO aparece todavía en el catálogo público", func(t *testing.T) {
		respuesta := hacerPeticion(router, http.MethodGet, "/catalogo", "", nil)
		cuerpo := decodificar(t, respuesta)
		vehiculos := cuerpo["vehiculos"].([]any)

		for _, v := range vehiculos {
			if v.(map[string]any)["id"] == idVehiculo {
				t.Error("un vehículo pendiente de aprobación no debería aparecer en el catálogo")
			}
		}
	})

	t.Run("el vendedor puede ver sus propias publicaciones", func(t *testing.T) {
		respuesta := hacerPeticion(router, http.MethodGet, "/vehiculos/mios", tokenVendedor, nil)

		if respuesta.Code != http.StatusOK {
			t.Fatalf("se esperaba 200, se obtuvo %d: %s", respuesta.Code, respuesta.Body.String())
		}

		vehiculos := decodificar(t, respuesta)["vehiculos"].([]any)
		if len(vehiculos) == 0 {
			t.Error("se esperaba encontrar al menos una publicación propia")
		}
	})

	t.Run("otro vendedor no ve la publicación ajena entre las suyas", func(t *testing.T) {
		respuesta := hacerPeticion(router, http.MethodGet, "/vehiculos/mios", tokenOtroVendedor, nil)
		vehiculos := decodificar(t, respuesta)["vehiculos"].([]any)

		for _, v := range vehiculos {
			if v.(map[string]any)["id"] == idVehiculo {
				t.Error("un vendedor no debería ver publicaciones de otro vendedor como propias")
			}
		}
	})

	t.Run("un vendedor no puede autorizar publicaciones (solo el administrador)", func(t *testing.T) {
		ruta := fmt.Sprintf("/vehiculos/%d/autorizar", int(idVehiculo))
		respuesta := hacerPeticion(router, http.MethodPatch, ruta, tokenVendedor, nil)

		if respuesta.Code != http.StatusForbidden {
			t.Fatalf("se esperaba 403, se obtuvo %d: %s", respuesta.Code, respuesta.Body.String())
		}
	})

	t.Run("el administrador autoriza la publicación", func(t *testing.T) {
		ruta := fmt.Sprintf("/vehiculos/%d/autorizar", int(idVehiculo))
		respuesta := hacerPeticion(router, http.MethodPatch, ruta, tokenAdmin, nil)

		if respuesta.Code != http.StatusOK {
			t.Fatalf("se esperaba 200, se obtuvo %d: %s", respuesta.Code, respuesta.Body.String())
		}

		vehiculo := decodificar(t, respuesta)["vehiculo"].(map[string]any)
		if vehiculo["estado"] != "publicado" {
			t.Errorf("se esperaba estado 'publicado', se obtuvo %v", vehiculo["estado"])
		}
	})

	t.Run("ahora sí aparece en el catálogo público", func(t *testing.T) {
		respuesta := hacerPeticion(router, http.MethodGet, "/catalogo", "", nil)
		vehiculos := decodificar(t, respuesta)["vehiculos"].([]any)

		encontrado := false
		for _, v := range vehiculos {
			if v.(map[string]any)["id"] == idVehiculo {
				encontrado = true
			}
		}

		if !encontrado {
			t.Error("el vehículo autorizado debería aparecer en el catálogo")
		}
	})

	t.Run("otro vendedor no puede reportar como vendido un vehículo ajeno", func(t *testing.T) {
		ruta := fmt.Sprintf("/vehiculos/%d/vender", int(idVehiculo))
		respuesta := hacerPeticion(router, http.MethodPatch, ruta, tokenOtroVendedor, nil)

		if respuesta.Code != http.StatusBadRequest {
			t.Fatalf("se esperaba 400, se obtuvo %d: %s", respuesta.Code, respuesta.Body.String())
		}
	})

	t.Run("el dueño reporta el vehículo como vendido", func(t *testing.T) {
		ruta := fmt.Sprintf("/vehiculos/%d/vender", int(idVehiculo))
		respuesta := hacerPeticion(router, http.MethodPatch, ruta, tokenVendedor, nil)

		if respuesta.Code != http.StatusOK {
			t.Fatalf("se esperaba 200, se obtuvo %d: %s", respuesta.Code, respuesta.Body.String())
		}

		vehiculo := decodificar(t, respuesta)["vehiculo"].(map[string]any)
		if vehiculo["estado"] != "vendido" {
			t.Errorf("se esperaba estado 'vendido', se obtuvo %v", vehiculo["estado"])
		}
	})

	t.Run("un vehículo vendido desaparece del catálogo", func(t *testing.T) {
		respuesta := hacerPeticion(router, http.MethodGet, "/catalogo", "", nil)
		vehiculos := decodificar(t, respuesta)["vehiculos"].([]any)

		for _, v := range vehiculos {
			if v.(map[string]any)["id"] == idVehiculo {
				t.Error("un vehículo vendido no debería seguir en el catálogo")
			}
		}
	})

	var idVehiculoParaEliminar float64

	t.Run("el administrador elimina una publicación", func(t *testing.T) {
		respuestaCrear := hacerPeticion(router, http.MethodPost, "/vehiculos", tokenVendedor, map[string]any{
			"marca": "Renault", "modelo": "Logan", "anio": 2018, "precio": 35000000,
		})
		idVehiculoParaEliminar = decodificar(t, respuestaCrear)["vehiculo"].(map[string]any)["id"].(float64)

		ruta := fmt.Sprintf("/vehiculos/%d", int(idVehiculoParaEliminar))
		respuesta := hacerPeticion(router, http.MethodDelete, ruta, tokenAdmin, nil)

		if respuesta.Code != http.StatusOK {
			t.Fatalf("se esperaba 200, se obtuvo %d: %s", respuesta.Code, respuesta.Body.String())
		}
	})

	t.Run("un vendedor no puede eliminar publicaciones (solo el administrador)", func(t *testing.T) {
		respuestaCrear := hacerPeticion(router, http.MethodPost, "/vehiculos", tokenVendedor, map[string]any{
			"marca": "Chevrolet", "modelo": "Spark", "anio": 2019, "precio": 28000000,
		})
		id := decodificar(t, respuestaCrear)["vehiculo"].(map[string]any)["id"].(float64)

		ruta := fmt.Sprintf("/vehiculos/%d", int(id))
		respuesta := hacerPeticion(router, http.MethodDelete, ruta, tokenVendedor, nil)

		if respuesta.Code != http.StatusForbidden {
			t.Fatalf("se esperaba 403, se obtuvo %d: %s", respuesta.Code, respuesta.Body.String())
		}
	})

	t.Run("un vendedor no puede eliminar usuarios (solo el administrador)", func(t *testing.T) {
		respuesta := hacerPeticion(router, http.MethodDelete, "/usuarios/999999", tokenVendedor, nil)

		if respuesta.Code != http.StatusForbidden {
			t.Fatalf("se esperaba 403, se obtuvo %d: %s", respuesta.Code, respuesta.Body.String())
		}
	})

	t.Run("el administrador elimina un usuario", func(t *testing.T) {
		_, idOtroVendedor := registrarYLoguear(t, router, fmt.Sprintf("borrar.%d@automarket.test", os.Getpid()))

		ruta := fmt.Sprintf("/usuarios/%d", int(idOtroVendedor))
		respuesta := hacerPeticion(router, http.MethodDelete, ruta, tokenAdmin, nil)

		if respuesta.Code != http.StatusOK {
			t.Fatalf("se esperaba 200, se obtuvo %d: %s", respuesta.Code, respuesta.Body.String())
		}
	})

	t.Run("las rutas anteriores de autenticación siguen funcionando igual", func(t *testing.T) {
		// Regresión: /usuarios y /login (hechas por tus compañeros) no deben
		// haberse roto con los cambios de este módulo.
		respuestaLoginMalo := hacerPeticion(router, http.MethodPost, "/login", "", map[string]string{
			"correo":     correoVendedor,
			"contrasena": "clave-incorrecta",
		})

		if respuestaLoginMalo.Code != http.StatusUnauthorized {
			t.Fatalf("se esperaba 401 con clave incorrecta, se obtuvo %d", respuestaLoginMalo.Code)
		}
	})
}

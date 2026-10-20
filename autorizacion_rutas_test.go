package main

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"AutoMarket/services"

	"github.com/gin-gonic/gin"
)

// BLOQUE: test de autorizacion sobre las rutas reales del router
func TestAutorizacionEnRutasReales(t *testing.T) {
	gin.SetMode(gin.TestMode)
	snapshotYRestaurarAlFinal(t)

	router := SetupRouter()

	correoVendedor := fmt.Sprintf("authz.vendedor.%d@automarket.test", os.Getpid())
	correoAdmin := fmt.Sprintf("authz.admin.%d@automarket.test", os.Getpid())

	tokenVendedor, _ := registrarYLoguear(t, router, correoVendedor)
	registrarYLoguear(t, router, correoAdmin)
	tokenAdmin := promoverAAdministrador(t, router, correoAdmin)

	// rolPermitido es el unico rol que puede usar la ruta
	rutasProtegidas := []struct {
		nombre       string
		metodo       string
		ruta         string
		rolPermitido string
	}{
		{"publicar vehículo", http.MethodPost, "/vehiculos", services.RolVendedor},
		{"ver mis publicaciones", http.MethodGet, "/vehiculos/mios", services.RolVendedor},
		{"reportar venta", http.MethodPatch, "/vehiculos/1/vender", services.RolVendedor},
		{"autorizar publicación", http.MethodPatch, "/vehiculos/1/autorizar", services.RolAdministrador},
		{"eliminar publicación", http.MethodDelete, "/vehiculos/1", services.RolAdministrador},
		{"eliminar usuario", http.MethodDelete, "/usuarios/1", services.RolAdministrador},
	}

	for _, ruta := range rutasProtegidas {
		t.Run("sin token devuelve 401: "+ruta.nombre, func(t *testing.T) {
			respuesta := hacerPeticion(router, ruta.metodo, ruta.ruta, "", nil)

			if respuesta.Code != http.StatusUnauthorized {
				t.Fatalf(
					"se esperaba 401, se obtuvo %d: %s",
					respuesta.Code, respuesta.Body.String(),
				)
			}
		})

		// el token del rol contrario al que tiene permiso en esta ruta
		tokenRolEquivocado := tokenVendedor
		if ruta.rolPermitido == services.RolVendedor {
			tokenRolEquivocado = tokenAdmin
		}

		t.Run("rol equivocado devuelve 403: "+ruta.nombre, func(t *testing.T) {
			respuesta := hacerPeticion(router, ruta.metodo, ruta.ruta, tokenRolEquivocado, nil)

			if respuesta.Code != http.StatusForbidden {
				t.Fatalf(
					"se esperaba 403, se obtuvo %d: %s",
					respuesta.Code, respuesta.Body.String(),
				)
			}
		})
	}
}

package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"AutoMarket/middleware"
	"AutoMarket/models"
	"AutoMarket/services"

	"github.com/gin-gonic/gin"
)

// BLOQUE: router de prueba
func construirRouterDePrueba(usuario *models.Usuario, operacion string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.GET(
		"/recurso-protegido",
		func(c *gin.Context) {
			if usuario != nil {
				c.Set(middleware.ClaveUsuarioContexto, *usuario)
			}
			c.Next()
		},
		middleware.RequerirPermiso(operacion),
		func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"mensaje": "acceso permitido"})
		},
	)

	return router
}

// BLOQUE: test vendedor intenta operacion de administrador
func TestVendedorIntentandoOperacionDeAdministrador(t *testing.T) {
	vendedor := models.Usuario{ID: 1, Nombre: "Ana", Rol: services.RolVendedor}
	router := construirRouterDePrueba(&vendedor, services.OpEliminarUsuario)

	peticion := httptest.NewRequest(http.MethodGet, "/recurso-protegido", nil)
	respuesta := httptest.NewRecorder()
	router.ServeHTTP(respuesta, peticion)

	if respuesta.Code != http.StatusForbidden {
		t.Errorf("se esperaba 403 Forbidden, se obtuvo %d", respuesta.Code)
	}
}

// BLOQUE: test administrador hace operacion administrativa
func TestAdministradorHaciendoOperacionAdministrativa(t *testing.T) {
	administrador := models.Usuario{ID: 2, Nombre: "Luis", Rol: services.RolAdministrador}
	router := construirRouterDePrueba(&administrador, services.OpEliminarUsuario)

	peticion := httptest.NewRequest(http.MethodGet, "/recurso-protegido", nil)
	respuesta := httptest.NewRecorder()
	router.ServeHTTP(respuesta, peticion)

	if respuesta.Code != http.StatusOK {
		t.Errorf("se esperaba 200 OK, se obtuvo %d", respuesta.Code)
	}
}

// BLOQUE: test visitante sin autenticar accede a recurso protegido
func TestVisitanteSinAutenticarAccesoRecursoProtegido(t *testing.T) {
	router := construirRouterDePrueba(nil, services.OpEliminarUsuario)

	peticion := httptest.NewRequest(http.MethodGet, "/recurso-protegido", nil)
	respuesta := httptest.NewRecorder()
	router.ServeHTTP(respuesta, peticion)

	if respuesta.Code != http.StatusUnauthorized {
		t.Errorf("se esperaba 401 Unauthorized, se obtuvo %d", respuesta.Code)
	}
}

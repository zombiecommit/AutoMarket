package services

import (
	"errors"
	"testing"
	"time"

	"AutoMarket/models"
)

func TestInvalidarSesionesUsuario(t *testing.T) {
	usuario := models.Usuario{
		ID:     10,
		Nombre: "Usuario de prueba",
		Correo: "prueba@automarket.test",
		Rol:    "vendedor",
	}

	token, err := CrearSesion(usuario)
	if err != nil {
		t.Fatalf("no se pudo crear la sesión: %v", err)
	}

	if _, existe := obtenerSesionConUsuarios(token, func() ([]models.Usuario, error) {
		return []models.Usuario{usuario}, nil
	}); !existe {
		t.Fatal("la sesión debería existir antes de invalidarla")
	}

	InvalidarSesionesUsuario(usuario.ID)

	if _, existe := ObtenerSesion(token); existe {
		t.Fatal("la sesión no debería existir después de invalidarla")
	}
}

func TestObtenerSesionTokenInexistente(t *testing.T) {
	_, existe := ObtenerSesion("token-que-no-existe")

	if existe {
		t.Fatal("un token inexistente no debería tener sesión")
	}
}

func TestObtenerSesionExpirada(t *testing.T) {
	usuario := models.Usuario{
		ID:     20,
		Nombre: "Usuario expirado",
		Correo: "expirado@automarket.test",
		Rol:    "vendedor",
	}

	token, err := CrearSesion(usuario)
	if err != nil {
		t.Fatalf("no se pudo crear la sesión: %v", err)
	}

	sesionesMu.Lock()
	sesion := sesiones[token]
	sesion.ExpiraEn = time.Now().Add(-time.Second)
	sesiones[token] = sesion
	sesionesMu.Unlock()

	if _, existe := ObtenerSesion(token); existe {
		t.Fatal("una sesión expirada no debería ser válida")
	}

	sesionesMu.RLock()
	_, sigueGuardada := sesiones[token]
	sesionesMu.RUnlock()

	if sigueGuardada {
		t.Fatal("la sesión expirada debería eliminarse del mapa")
	}
}

func TestInvalidarSesionIndividual(t *testing.T) {
	usuario := models.Usuario{
		ID:     30,
		Nombre: "Usuario de prueba",
		Correo: "individual@automarket.test",
		Rol:    "vendedor",
	}

	token, err := CrearSesion(usuario)
	if err != nil {
		t.Fatalf("no se pudo crear la sesión: %v", err)
	}

	if _, existe := obtenerSesionConUsuarios(token, func() ([]models.Usuario, error) {
		return []models.Usuario{usuario}, nil
	}); !existe {
		t.Fatal("la sesión debería existir antes de invalidarla")
	}

	InvalidarSesion(token)

	if _, existe := ObtenerSesion(token); existe {
		t.Fatal("la sesión debería dejar de existir después de invalidarla")
	}

}

func TestObtenerSesionActualizaRolUsuario(t *testing.T) {
	usuario := models.Usuario{
		ID:  40,
		Rol: "visitor",
	}
	token, err := CrearSesion(usuario)
	if err != nil {
		t.Fatal(err)
	}
	defer InvalidarSesion(token)

	usuariosActualizados := []models.Usuario{
		{
			ID:  40,
			Rol: "admin",
		},
	}

	sesion, existe := obtenerSesionConUsuarios(token, func() ([]models.Usuario, error) {
		return usuariosActualizados, nil
	})
	if !existe {
		t.Fatal("se esperaba encontrar la sesión")
	}
	if sesion.Usuario.Rol != "admin" {
		t.Fatalf("se esperaba el rol admin, se obtuvo %q", sesion.Usuario.Rol)
	}
}

func TestObtenerSesionFalloAlCargarUsuarios(t *testing.T) {
	usuario := models.Usuario{
		ID:  50,
		Rol: "vendedor",
	}

	token, err := CrearSesion(usuario)
	if err != nil {
		t.Fatalf("no se pudo crear la sesión: %v", err)
	}
	defer InvalidarSesion(token)

	sesion, existe := obtenerSesionConUsuarios(token, func() ([]models.Usuario, error) {
		return nil, errors.New("error de lectura del almacenamiento")
	})

	if existe {
		t.Fatal("la sesión no debería ser válida si falla la carga de usuarios")
	}
	if sesion.Usuario.ID != 0 {
		t.Fatal("no debería devolverse un usuario cuando falla la carga")
	}
}

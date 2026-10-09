package services

import (
	"testing"

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

	_, existe := ObtenerSesion(token)
	if !existe {
		t.Fatal("la sesión debería existir antes de invalidarla")
	}

	InvalidarSesionesUsuario(usuario.ID)

	_, existe = ObtenerSesion(token)
	if existe {
		t.Fatal("la sesión no debería existir después de invalidarla")
	}
}

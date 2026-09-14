package storage

import (
	"os"
	"testing"
)

func TestRegistrarAccounting(t *testing.T) {
	err := RegistrarAccounting(
		1,
		"vendedor",
		"LOGIN",
		"autenticacion",
		true,
	)

	if err != nil {
		t.Fatalf("error al registrar Accounting: %v", err)
	}

	registros, err := CargarRegistrosAccounting()

	if err != nil {
		t.Fatalf("error al cargar Accounting: %v", err)
	}

	if len(registros) != 1 {
		t.Fatalf("se esperaba 1 registro, se encontraron %d", len(registros))
	}

	if registros[0].UsuarioID != 1 {
		t.Errorf("se esperaba UsuarioID 1, se obtuvo %d", registros[0].UsuarioID)
	}

	if registros[0].Operacion != "LOGIN" {
		t.Errorf(
			"se esperaba operación LOGIN, se obtuvo %s",
			registros[0].Operacion,
		)
	}

	if registros[0].Hash == "" {
		t.Error("el registro debería tener un hash")
	}

	if registros[0].HashAnterior != "" {
		t.Error("el primer registro no debería tener hash anterior")
	}

	os.Remove(archivoAccounting)
}

package storage

import (
	"os"
	"testing"
)

// Guarda una copia del accounting.json original y lo restaura
// al terminar la prueba.
func prepararAccountingTest(t *testing.T) {
	t.Helper()

	contenidoOriginal, err := os.ReadFile(archivoAccounting)
	existiaOriginal := err == nil

	t.Cleanup(func() {
		if existiaOriginal {
			_ = os.WriteFile(archivoAccounting, contenidoOriginal, 0644)
		} else {
			_ = os.Remove(archivoAccounting)
		}
	})

	_ = os.Remove(archivoAccounting)
}

// Crea un registro de prueba y devuelve los registros almacenados.
func crearRegistroAccountingTest(t *testing.T) []RegistroAccounting {
	t.Helper()

	err := RegistrarAccounting(
		1,
		"vendedor",
		"LOGIN",
		"autenticacion",
		"",
		true,
	)

	if err != nil {
		t.Fatalf("error al registrar Accounting: %v", err)
	}

	registros, err := CargarRegistrosAccounting()

	if err != nil {
		t.Fatalf("error al cargar Accounting: %v", err)
	}

	return registros
}

// Crea dos registros de prueba y devuelve los registros almacenados.
func crearDosRegistrosAccountingTest(t *testing.T) []RegistroAccounting {
	t.Helper()

	registros := crearRegistroAccountingTest(t)

	err := RegistrarAccounting(
		1,
		"vendedor",
		"PUBLICAR_VEHICULO",
		"vehiculo",
		"10",
		true,
	)

	if err != nil {
		t.Fatalf("error al registrar segundo registro: %v", err)
	}

	registros, err = CargarRegistrosAccounting()

	if err != nil {
		t.Fatalf("error al cargar Accounting: %v", err)
	}

	return registros
}

func TestRegistrarAccounting(t *testing.T) {
	prepararAccountingTest(t)

	registros := crearRegistroAccountingTest(t)

	if len(registros) != 1 {
		t.Fatalf(
			"se esperaba 1 registro, se encontraron %d",
			len(registros),
		)
	}

	if registros[0].UsuarioID != 1 {
		t.Errorf(
			"se esperaba UsuarioID 1, se obtuvo %d",
			registros[0].UsuarioID,
		)
	}

	if registros[0].Operacion != "LOGIN" {
		t.Errorf(
			"se esperaba operación LOGIN, se obtuvo %s",
			registros[0].Operacion,
		)
	}

	if registros[0].RecursoID != "" {
		t.Errorf(
			"se esperaba RecursoID vacío, se obtuvo %s",
			registros[0].RecursoID,
		)
	}

	if registros[0].Hash == "" {
		t.Error("el registro debería tener un hash")
	}

	if registros[0].HashAnterior != "" {
		t.Error("el primer registro no debería tener hash anterior")
	}

	if err := verificarIntegridad(registros); err != nil {
		t.Fatalf(
			"el registro recién creado debería ser íntegro: %v",
			err,
		)
	}
}

func TestAccountingDetectaModificacionDeCampo(t *testing.T) {
	prepararAccountingTest(t)

	registros := crearRegistroAccountingTest(t)

	// Simulamos que alguien modificó un campo del registro.
	registros[0].Operacion = "OPERACION_MODIFICADA"

	err := verificarIntegridad(registros)

	if err == nil {
		t.Fatal(
			"se esperaba que la verificación detectara la modificación del registro",
		)
	}
}

func TestAccountingDetectaModificacionDeHash(t *testing.T) {
	prepararAccountingTest(t)

	registros := crearRegistroAccountingTest(t)

	// Simulamos que alguien modificó el hash.
	registros[0].Hash = "hash_falso"

	err := verificarIntegridad(registros)

	if err == nil {
		t.Fatal(
			"se esperaba que la verificación detectara un hash inválido",
		)
	}
}

func TestAccountingDetectaHashAnteriorIncorrecto(t *testing.T) {
	prepararAccountingTest(t)

	registros := crearDosRegistrosAccountingTest(t)

	// Simulamos que alguien modificó la cadena.
	registros[1].HashAnterior = "hash_anterior_falso"

	err := verificarIntegridad(registros)

	if err == nil {
		t.Fatal(
			"se esperaba que la verificación detectara un HashAnterior incorrecto",
		)
	}
}

func TestAccountingDetectaIDsFueraDeSecuencia(t *testing.T) {
	prepararAccountingTest(t)

	registros := crearDosRegistrosAccountingTest(t)

	// Simulamos una alteración en la secuencia de IDs.
	registros[1].ID = registros[0].ID

	err := verificarIntegridad(registros)

	if err == nil {
		t.Fatal(
			"se esperaba que la verificación detectara IDs fuera de secuencia",
		)
	}
}

func TestRegistrarAccountingDetectaFalloDePersistencia(t *testing.T) {
	prepararAccountingTest(t)

	// Primero creamos un registro válido.
	crearRegistroAccountingTest(t)

	// Convertimos accounting.json en un directorio.
	// De esta manera os.WriteFile() no podrá escribir sobre él.
	if err := os.Remove(archivoAccounting); err != nil {
		t.Fatalf("no se pudo preparar la prueba: %v", err)
	}

	if err := os.Mkdir(archivoAccounting, 0755); err != nil {
		t.Fatalf("no se pudo crear el directorio de prueba: %v", err)
	}

	// Intentamos registrar una nueva operación.
	err := RegistrarAccounting(
		2,
		"vendedor",
		"PUBLICAR_VEHICULO",
		"vehiculo",
		"10",
		true,
	)

	if err == nil {
		t.Fatal(
			"se esperaba un error cuando Accounting no puede guardar el archivo",
		)
	}
}

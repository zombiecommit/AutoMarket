package services

import (
	"encoding/hex"
	"testing"
)

func TestPBKDF2SHA256VectorReferencia(t *testing.T) {
	password := []byte("password")
	salt := []byte("salt")
	iteracionesPrueba := 1
	longitud := 32

	esperado, err := hex.DecodeString(
		"120fb6cffcf8b32c43e7225256c4f837a86548c92ccc35480805987cb70be17b",
	)
	if err != nil {
		t.Fatal(err)
	}

	obtenido := pbkdf2SHA256(password, salt, iteracionesPrueba, longitud)

	if hex.EncodeToString(obtenido) != hex.EncodeToString(esperado) {
		t.Fatalf(
			"PBKDF2 no coincide con el vector de referencia:\nobtenido: %x\nesperado: %x",
			obtenido,
			esperado,
		)
	}
}

func TestGenerarHashUsaSaltDiferente(t *testing.T) {
	contrasena := "ClaveSegura123!"

	hash1, err := generarHash(contrasena)
	if err != nil {
		t.Fatalf("no se pudo generar el primer hash: %v", err)
	}

	hash2, err := generarHash(contrasena)
	if err != nil {
		t.Fatalf("no se pudo generar el segundo hash: %v", err)
	}

	if hash1 == hash2 {
		t.Fatal("dos hashes de la misma contraseña no deberían ser idénticos")
	}

	if !verificarHash(contrasena, hash1) {
		t.Fatal("la contraseña original debería validar contra el primer hash")
	}

	if !verificarHash(contrasena, hash2) {
		t.Fatal("la contraseña original debería validar contra el segundo hash")
	}
}

func TestVerificarHashRechazaHashMalformado(t *testing.T) {
	hashesInvalidos := []string{
		"",
		"sin-separador",
		"salt:hash:extra",
		"%%%:%%% ",
		"c2FsdA==:",
		":aGFzaA==",
	}

	for _, hash := range hashesInvalidos {
		t.Run(hash, func(t *testing.T) {
			if verificarHash("ClaveSegura123!", hash) {
				t.Errorf("se aceptó un hash malformado: %q", hash)
			}
		})
	}
}

func TestVerificarHashRechazaContrasenaIncorrecta(t *testing.T) {
	hash, err := generarHash("ClaveCorrecta123!")
	if err != nil {
		t.Fatalf("no se pudo generar el hash: %v", err)
	}

	if verificarHash("ClaveIncorrecta123!", hash) {
		t.Fatal("se aceptó una contraseña incorrecta")
	}
}

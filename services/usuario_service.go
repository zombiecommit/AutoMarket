package services

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"strings"

	"AutoMarket/models"
	"AutoMarket/storage"
)

const iteraciones = 100000

func RegistrarUsuario(nombre, correo, contrasena string) (models.Usuario, error) {
	nombre = strings.TrimSpace(nombre)
	correo = strings.TrimSpace(correo)

	if nombre == "" || correo == "" || contrasena == "" {
		return models.Usuario{}, errors.New("todos los campos son obligatorios")
	}

	usuarios, err := storage.CargarUsuarios()
	if err != nil {
		return models.Usuario{}, err
	}

	for _, usuario := range usuarios {
		if usuario.Correo == correo {
			return models.Usuario{}, errors.New("el correo ya está registrado")
		}
	}

	contrasenaHash, err := generarHash(contrasena)
	if err != nil {
		return models.Usuario{}, err
	}

	nuevoUsuario := models.Usuario{
		ID:             siguienteID(usuarios),
		Nombre:         nombre,
		Correo:         correo,
		ContrasenaHash: contrasenaHash,
		Rol:            "vendedor",
	}

	usuarios = append(usuarios, nuevoUsuario)

	if err := storage.GuardarUsuarios(usuarios); err != nil {
		return models.Usuario{}, err
	}

	return nuevoUsuario, nil
}

func siguienteID(usuarios []models.Usuario) int {
	mayorID := 0

	for _, usuario := range usuarios {
		if usuario.ID > mayorID {
			mayorID = usuario.ID
		}
	}

	return mayorID + 1
}

func generarHash(contrasena string) (string, error) {
	salt := make([]byte, 16)

	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := pbkdf2SHA256([]byte(contrasena), salt, iteraciones, 32)

	return base64.StdEncoding.EncodeToString(salt) + ":" +
		base64.StdEncoding.EncodeToString(hash), nil
}

func pbkdf2SHA256(password, salt []byte, iteraciones, longitud int) []byte {
	var resultado []byte
	numeroBloque := 1

	for len(resultado) < longitud {
		h := hmac.New(sha256.New, password)
		h.Write(salt)
		h.Write([]byte{
			byte(numeroBloque >> 24),
			byte(numeroBloque >> 16),
			byte(numeroBloque >> 8),
			byte(numeroBloque),
		})

		u := h.Sum(nil)
		bloque := make([]byte, len(u))
		copy(bloque, u)

		for i := 1; i < iteraciones; i++ {
			h = hmac.New(sha256.New, password)
			h.Write(u)
			u = h.Sum(nil)

			for j := range bloque {
				bloque[j] ^= u[j]
			}
		}

		resultado = append(resultado, bloque...)
		numeroBloque++
	}

	return resultado[:longitud]
}

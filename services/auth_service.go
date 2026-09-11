package services

import (
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"strings"

	"AutoMarket/models"
	"AutoMarket/storage"
)

func AutenticarUsuario(correo, contrasena string) (models.Usuario, error) {
	correo = strings.TrimSpace(correo)

	if correo == "" || contrasena == "" {
		return models.Usuario{}, errors.New("correo y contraseña son obligatorios")
	}

	usuarios, err := storage.CargarUsuarios()
	if err != nil {
		return models.Usuario{}, err
	}

	for _, usuario := range usuarios {
		if usuario.Correo == correo {
			if verificarHash(contrasena, usuario.ContrasenaHash) {
				return usuario, nil
			}

			return models.Usuario{}, errors.New("correo o contraseña incorrectos")
		}
	}

	return models.Usuario{}, errors.New("correo o contraseña incorrectos")
}

func verificarHash(contrasena, hashGuardado string) bool {
	partes := strings.Split(hashGuardado, ":")

	if len(partes) != 2 {
		return false
	}

	salt, err := base64.StdEncoding.DecodeString(partes[0])
	if err != nil {
		return false
	}

	hashEsperado, err := base64.StdEncoding.DecodeString(partes[1])
	if err != nil {
		return false
	}

	hashCalculado := pbkdf2SHA256(
		[]byte(contrasena),
		salt,
		iteraciones,
		len(hashEsperado),
	)

	return subtle.ConstantTimeCompare(hashCalculado, hashEsperado) == 1
}

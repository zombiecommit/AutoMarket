package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"

	"AutoMarket/models"
)

type Sesion struct {
	Token   string
	Usuario models.Usuario
}

var (
	sesiones   = make(map[string]Sesion)
	sesionesMu sync.RWMutex
)

func CrearSesion(usuario models.Usuario) (string, error) {
	tokenBytes := make([]byte, 32)

	if _, err := rand.Read(tokenBytes); err != nil {
		return "", errors.New("no se pudo generar la sesión")
	}

	token := hex.EncodeToString(tokenBytes)

	sesion := Sesion{
		Token:   token,
		Usuario: usuario,
	}

	sesionesMu.Lock()
	sesiones[token] = sesion
	sesionesMu.Unlock()

	return token, nil
}

func ObtenerSesion(token string) (Sesion, bool) {
	sesionesMu.RLock()
	sesion, existe := sesiones[token]
	sesionesMu.RUnlock()

	return sesion, existe
}

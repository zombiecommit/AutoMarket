package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"
	"time"

	"AutoMarket/models"
	"AutoMarket/storage"
)

const duracionSesion = 30 * time.Minute

type Sesion struct {
	Token    string
	Usuario  models.Usuario
	ExpiraEn time.Time
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
		Token:    token,
		Usuario:  usuario,
		ExpiraEn: time.Now().Add(duracionSesion),
	}

	sesionesMu.Lock()
	sesiones[token] = sesion
	sesionesMu.Unlock()

	return token, nil
}

func ObtenerSesion(token string) (Sesion, bool) {
	return obtenerSesionConUsuarios(token, storage.CargarUsuarios)
}

func obtenerSesionConUsuarios(
	token string,
	cargarUsuarios func() ([]models.Usuario, error),
) (Sesion, bool) {
	sesionesMu.Lock()
	sesion, existe := sesiones[token]

	if !existe {
		sesionesMu.Unlock()
		return Sesion{}, false
	}

	if !time.Now().Before(sesion.ExpiraEn) {
		delete(sesiones, token)
		sesionesMu.Unlock()
		return Sesion{}, false
	}

	sesionesMu.Unlock()

	usuarios, err := cargarUsuarios()
	if err != nil {
		return Sesion{}, false
	}

	for _, usuario := range usuarios {
		if usuario.ID != sesion.Usuario.ID {
			continue
		}

		sesionesMu.Lock()
		defer sesionesMu.Unlock()

		sesionActual, existe := sesiones[token]
		if !existe || !time.Now().Before(sesionActual.ExpiraEn) {
			delete(sesiones, token)
			return Sesion{}, false
		}

		sesionActual.Usuario = usuario
		sesiones[token] = sesionActual

		return sesionActual, true
	}

	InvalidarSesion(token)
	return Sesion{}, false
}

func InvalidarSesion(token string) {
	sesionesMu.Lock()
	defer sesionesMu.Unlock()

	delete(sesiones, token)
}

func InvalidarSesionesUsuario(usuarioID int) {
	sesionesMu.Lock()
	defer sesionesMu.Unlock()

	for token, sesion := range sesiones {
		if sesion.Usuario.ID == usuarioID {
			delete(sesiones, token)
		}
	}
}

package services

import (
	"AutoMarket/models"
	"AutoMarket/storage"
)

const (
	OpRegistrarUsuario = "REGISTRAR_USUARIO"
	OpIniciarSesion    = "INICIAR_SESION"
)

func RegistrarOperacion(
	usuario models.Usuario,
	operacion string,
	recurso string,
	exito bool,
) error {
	return storage.RegistrarAccounting(
		usuario.ID,
		usuario.Rol,
		operacion,
		recurso,
		exito,
	)
}

func RegistrarOperacionSinAutenticar(
	operacion string,
	recurso string,
	exito bool,
) error {
	return storage.RegistrarAccounting(
		0,
		RolVisitante,
		operacion,
		recurso,
		exito,
	)
}

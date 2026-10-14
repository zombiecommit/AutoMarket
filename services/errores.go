package services

import "errors"

// BLOQUE: errores de negocio del sistema

var (
	ErrCredencialesInvalidas   = errors.New("correo o contraseña incorrectos")
	ErrCorreoDuplicado         = errors.New("el correo ya está registrado")
	ErrUsuarioNoEncontrado     = errors.New("usuario no encontrado")
	ErrPublicacionNoEncontrada = errors.New("publicación no encontrada")
	ErrPublicacionAjena        = errors.New("no puedes reportar la venta de una publicación que no te pertenece")
	ErrPublicacionYaVendida    = errors.New("la publicación ya está marcada como vendida")
	ErrPublicacionNoPendiente  = errors.New("solo se pueden autorizar publicaciones pendientes de aprobación")
	ErrPublicacionNoPublicada  = errors.New("solo se puede reportar como vendida una publicación que esté publicada")
	ErrUsuarioConPublicaciones = errors.New("no se puede eliminar un usuario que tiene publicaciones; elimina primero sus publicaciones")
)

// BLOQUE: error de validación de datos de entrada (corresponde a 400)
type ErrorValidacion struct {
	Mensaje string
}

func (e *ErrorValidacion) Error() string {
	return e.Mensaje
}

func nuevoErrorValidacion(mensaje string) error {
	return &ErrorValidacion{Mensaje: mensaje}
}

func EsErrorDeValidacion(err error) bool {
	var errorValidacion *ErrorValidacion
	return errors.As(err, &errorValidacion)
}

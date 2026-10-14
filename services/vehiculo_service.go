package services

import (
	"strings"

	"AutoMarket/models"
	"AutoMarket/storage"
)

// BLOQUE: estados posibles de una publicacion
//
// Transiciones permitidas (cualquier otra se rechaza con 409):
//
//	pendiente_aprobacion --(administrador autoriza)--> publicado
//	publicado            --(vendedor reporta venta)--> vendido
//
// Una publicacion pendiente NO puede venderse: primero debe publicarse.
const (
	EstadoPendiente = "pendiente_aprobacion"
	EstadoPublicado = "publicado"
	EstadoVendido   = "vendido"
)

// BLOQUE: registrar/publicar vehiculo (vendedor)
// La publicacion queda inicialmente en estado "pendiente_aprobacion".
func RegistrarVehiculo(
	vendedorID int,
	marca string,
	modelo string,
	anio int,
	precio float64,
	descripcion string,
) (models.Vehiculo, error) {
	marca = strings.TrimSpace(marca)
	modelo = strings.TrimSpace(modelo)
	descripcion = strings.TrimSpace(descripcion)

	if marca == "" || modelo == "" {
		return models.Vehiculo{}, nuevoErrorValidacion("marca y modelo son obligatorios")
	}

	if anio <= 0 {
		return models.Vehiculo{}, nuevoErrorValidacion("el año del vehículo es inválido")
	}

	if precio <= 0 {
		return models.Vehiculo{}, nuevoErrorValidacion("el precio debe ser mayor que cero")
	}

	vehiculos, err := storage.CargarVehiculos()
	if err != nil {
		return models.Vehiculo{}, err
	}

	nuevoVehiculo := models.Vehiculo{
		ID:          siguienteIDVehiculo(vehiculos),
		VendedorID:  vendedorID,
		Marca:       marca,
		Modelo:      modelo,
		Anio:        anio,
		Precio:      precio,
		Descripcion: descripcion,
		Estado:      EstadoPendiente,
	}

	vehiculos = append(vehiculos, nuevoVehiculo)

	if err := storage.GuardarVehiculos(vehiculos); err != nil {
		return models.Vehiculo{}, err
	}

	return nuevoVehiculo, nil
}

// BLOQUE: consultar catalogo (visitante / publico)
// Solo se muestran las publicaciones ya autorizadas por el administrador.
// Cada vehiculo se devuelve junto con los datos de contacto del vendedor
// que lo publico.
func ObtenerCatalogo() ([]models.VehiculoPublico, error) {
	vehiculos, err := storage.CargarVehiculos()
	if err != nil {
		return nil, err
	}

	usuarios, err := storage.CargarUsuarios()
	if err != nil {
		return nil, err
	}

	vendedoresPorID := make(map[int]models.Usuario, len(usuarios))
	for _, usuario := range usuarios {
		vendedoresPorID[usuario.ID] = usuario
	}

	catalogo := make([]models.VehiculoPublico, 0)

	for _, vehiculo := range vehiculos {
		if vehiculo.Estado != EstadoPublicado {
			continue
		}

		vendedor := vendedoresPorID[vehiculo.VendedorID]

		catalogo = append(catalogo, models.VehiculoPublico{
			Vehiculo: vehiculo,
			Vendedor: models.ContactoVendedor{
				Nombre: vendedor.Nombre,
				Correo: vendedor.Correo,
			},
		})
	}

	return catalogo, nil
}

// BLOQUE: ver publicaciones propias (vendedor)
func ObtenerPublicacionesDeVendedor(vendedorID int) ([]models.Vehiculo, error) {
	vehiculos, err := storage.CargarVehiculos()
	if err != nil {
		return nil, err
	}

	propias := make([]models.Vehiculo, 0)

	for _, vehiculo := range vehiculos {
		if vehiculo.VendedorID == vendedorID {
			propias = append(propias, vehiculo)
		}
	}

	return propias, nil
}

// BLOQUE: reportar vehiculo vendido (vendedor, solo su propia publicacion)
func ReportarVehiculoVendido(vendedorID int, vehiculoID int) (models.Vehiculo, error) {
	vehiculos, err := storage.CargarVehiculos()
	if err != nil {
		return models.Vehiculo{}, err
	}

	for i, vehiculo := range vehiculos {
		if vehiculo.ID == vehiculoID {
			if vehiculo.VendedorID != vendedorID {
				return models.Vehiculo{}, ErrPublicacionAjena
			}

			if vehiculo.Estado == EstadoVendido {
				return models.Vehiculo{}, ErrPublicacionYaVendida
			}

			if vehiculo.Estado != EstadoPublicado {
				return models.Vehiculo{}, ErrPublicacionNoPublicada
			}

			vehiculos[i].Estado = EstadoVendido

			if err := storage.GuardarVehiculos(vehiculos); err != nil {
				return models.Vehiculo{}, err
			}

			return vehiculos[i], nil
		}
	}

	return models.Vehiculo{}, ErrPublicacionNoEncontrada
}

// BLOQUE: autorizar publicacion (administrador)
// Cambia el estado de "pendiente_aprobacion" a "publicado".
func AutorizarPublicacion(vehiculoID int) (models.Vehiculo, error) {
	vehiculos, err := storage.CargarVehiculos()
	if err != nil {
		return models.Vehiculo{}, err
	}

	for i, vehiculo := range vehiculos {
		if vehiculo.ID == vehiculoID {
			if vehiculo.Estado != EstadoPendiente {
				return models.Vehiculo{}, ErrPublicacionNoPendiente
			}

			vehiculos[i].Estado = EstadoPublicado

			if err := storage.GuardarVehiculos(vehiculos); err != nil {
				return models.Vehiculo{}, err
			}

			return vehiculos[i], nil
		}
	}

	return models.Vehiculo{}, ErrPublicacionNoEncontrada
}

// BLOQUE: eliminar publicacion (administrador)
func EliminarPublicacion(vehiculoID int) error {
	vehiculos, err := storage.CargarVehiculos()
	if err != nil {
		return err
	}

	indice := -1

	for i, vehiculo := range vehiculos {
		if vehiculo.ID == vehiculoID {
			indice = i
			break
		}
	}

	if indice == -1 {
		return ErrPublicacionNoEncontrada
	}

	vehiculos = append(vehiculos[:indice], vehiculos[indice+1:]...)

	return storage.GuardarVehiculos(vehiculos)
}

// BLOQUE: saber si un usuario tiene publicaciones (en cualquier estado)
// Lo usa EliminarUsuario para no dejar publicaciones huerfanas.
func usuarioTienePublicaciones(usuarioID int) (bool, error) {
	vehiculos, err := storage.CargarVehiculos()
	if err != nil {
		return false, err
	}

	for _, vehiculo := range vehiculos {
		if vehiculo.VendedorID == usuarioID {
			return true, nil
		}
	}

	return false, nil
}

// BLOQUE: siguiente ID disponible
func siguienteIDVehiculo(vehiculos []models.Vehiculo) int {
	mayorID := 0

	for _, vehiculo := range vehiculos {
		if vehiculo.ID > mayorID {
			mayorID = vehiculo.ID
		}
	}

	return mayorID + 1
}

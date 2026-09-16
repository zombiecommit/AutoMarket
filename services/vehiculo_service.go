package services

import (
	"errors"
	"strings"

	"AutoMarket/models"
	"AutoMarket/storage"
)

// BLOQUE: estados posibles de una publicacion
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
		return models.Vehiculo{}, errors.New("marca y modelo son obligatorios")
	}

	if anio <= 0 {
		return models.Vehiculo{}, errors.New("el año del vehículo es inválido")
	}

	if precio <= 0 {
		return models.Vehiculo{}, errors.New("el precio debe ser mayor que cero")
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
func ObtenerCatalogo() ([]models.Vehiculo, error) {
	vehiculos, err := storage.CargarVehiculos()
	if err != nil {
		return nil, err
	}

	catalogo := make([]models.Vehiculo, 0)

	for _, vehiculo := range vehiculos {
		if vehiculo.Estado == EstadoPublicado {
			catalogo = append(catalogo, vehiculo)
		}
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
				return models.Vehiculo{}, errors.New("no puedes reportar la venta de una publicación que no te pertenece")
			}

			if vehiculo.Estado == EstadoVendido {
				return models.Vehiculo{}, errors.New("la publicación ya está marcada como vendida")
			}

			vehiculos[i].Estado = EstadoVendido

			if err := storage.GuardarVehiculos(vehiculos); err != nil {
				return models.Vehiculo{}, err
			}

			return vehiculos[i], nil
		}
	}

	return models.Vehiculo{}, errors.New("publicación no encontrada")
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
				return models.Vehiculo{}, errors.New("solo se pueden autorizar publicaciones pendientes de aprobación")
			}

			vehiculos[i].Estado = EstadoPublicado

			if err := storage.GuardarVehiculos(vehiculos); err != nil {
				return models.Vehiculo{}, err
			}

			return vehiculos[i], nil
		}
	}

	return models.Vehiculo{}, errors.New("publicación no encontrada")
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
		return errors.New("publicación no encontrada")
	}

	vehiculos = append(vehiculos[:indice], vehiculos[indice+1:]...)

	return storage.GuardarVehiculos(vehiculos)
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

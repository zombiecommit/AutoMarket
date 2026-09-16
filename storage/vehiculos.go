package storage

import (
	"encoding/json"
	"os"

	"AutoMarket/models"
)

const archivoVehiculos = "storage/vehiculos.json"

// BLOQUE: cargar vehiculos desde el archivo JSON
func CargarVehiculos() ([]models.Vehiculo, error) {
	datos, err := os.ReadFile(archivoVehiculos)

	if err != nil {
		if os.IsNotExist(err) {
			return []models.Vehiculo{}, nil
		}
		return nil, err
	}

	var vehiculos []models.Vehiculo

	if len(datos) == 0 {
		return []models.Vehiculo{}, nil
	}

	err = json.Unmarshal(datos, &vehiculos)
	if err != nil {
		return nil, err
	}

	return vehiculos, nil
}

// BLOQUE: guardar vehiculos en el archivo JSON
func GuardarVehiculos(vehiculos []models.Vehiculo) error {
	datos, err := json.MarshalIndent(vehiculos, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(archivoVehiculos, datos, 0644)
}

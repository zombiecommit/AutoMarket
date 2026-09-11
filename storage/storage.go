package storage

import (
	"encoding/json"
	"os"

	"AutoMarket/models"
)

const archivoUsuarios = "storage/usuarios.json"

func CargarUsuarios() ([]models.Usuario, error) {
	datos, err := os.ReadFile(archivoUsuarios)
	if err != nil {
		return nil, err
	}

	var usuarios []models.Usuario

	err = json.Unmarshal(datos, &usuarios)
	if err != nil {
		return nil, err
	}

	return usuarios, nil
}

func GuardarUsuarios(usuarios []models.Usuario) error {
	datos, err := json.MarshalIndent(usuarios, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(archivoUsuarios, datos, 0644)
}

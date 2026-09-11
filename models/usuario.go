package models

type Usuario struct {
	ID             int    `json:"id"`
	Nombre         string `json:"nombre"`
	Correo         string `json:"correo"`
	ContrasenaHash string `json:"contrasena_hash"`
	Rol            string `json:"rol"`
}

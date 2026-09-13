package handlers

import (
	"net/http"

	"AutoMarket/storage"

	"github.com/gin-gonic/gin"
)

// NOTA: el registro y el login de usuarios los expone el equipo de AAA
// (paquete services/auth_service.go + storage/storage.go), no este archivo.
// Aquí solo vive la acción de administrador de eliminar una cuenta.

// EliminarUsuario permite a un administrador borrar una cuenta existente.
func EliminarUsuario(c *gin.Context) {
	id, err := idDesdeParametro(c, "id")
	if err != nil {
		return
	}

	usuarios, err := storage.CargarUsuarios()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo leer el almacenamiento de usuarios"})
		return
	}

	indice := -1
	for i, u := range usuarios {
		if u.ID == id {
			indice = i
			break
		}
	}

	if indice == -1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "usuario no encontrado"})
		return
	}

	usuarios = append(usuarios[:indice], usuarios[indice+1:]...)

	if err := storage.GuardarUsuarios(usuarios); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo guardar el cambio"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensaje": "usuario eliminado"})
}

package handlers

import (
	"net/http"
	"strconv"

	"AutoMarket/services"

	"github.com/gin-gonic/gin"
)

type registroRequest struct {
	Nombre     string `json:"nombre"`
	Correo     string `json:"correo"`
	Contrasena string `json:"contrasena"`
}

func RegistrarUsuario(c *gin.Context) {
	var request registroRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "el cuerpo de la solicitud debe ser un JSON válido",
		})
		return
	}

	usuario, err := services.RegistrarUsuario(
		request.Nombre,
		request.Correo,
		request.Contrasena,
	)

	if err != nil {
		responderError(c, err, "no se pudo registrar el usuario")
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"mensaje": "usuario registrado correctamente",
		"usuario": gin.H{
			"id":     usuario.ID,
			"nombre": usuario.Nombre,
			"correo": usuario.Correo,
			"rol":    usuario.Rol,
		},
	})
}

// BLOQUE: eliminar usuario (administrador)
func EliminarUsuario(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "el id del usuario es inválido",
		})
		return
	}

	if err := services.EliminarUsuario(id); err != nil {
		responderError(c, err, "no se pudo eliminar el usuario")

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"mensaje": "usuario eliminado correctamente",
	})
}

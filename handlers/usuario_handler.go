package handlers

import (
	"net/http"

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
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
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

package handlers

import (
	"net/http"

	"AutoMarket/services"

	"github.com/gin-gonic/gin"
)

type loginRequest struct {
	Correo     string `json:"correo"`
	Contrasena string `json:"contrasena"`
}

func IniciarSesion(c *gin.Context) {
	var request loginRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "el cuerpo de la solicitud debe ser un JSON válido",
		})
		return
	}

	usuario, err := services.AutenticarUsuario(
		request.Correo,
		request.Contrasena,
	)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	token, err := services.CrearSesion(usuario)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"mensaje": "inicio de sesión exitoso",
		"token":   token,
		"usuario": gin.H{
			"id":     usuario.ID,
			"nombre": usuario.Nombre,
			"correo": usuario.Correo,
			"rol":    usuario.Rol,
		},
	})
}

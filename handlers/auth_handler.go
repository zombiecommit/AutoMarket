package handlers

import (
	"net/http"
	"strings"

	"AutoMarket/middleware"
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

	if strings.TrimSpace(request.Correo) == "" ||
		strings.TrimSpace(request.Contrasena) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "el correo y la contraseña son obligatorios",
		})
		return
	}

	usuario, err := services.AutenticarUsuario(
		request.Correo,
		request.Contrasena,
	)

	if err != nil {
		responderError(c, err, "no se pudo iniciar sesión")
		return
	}
	token, err := services.CrearSesion(usuario)

	if err != nil {
		responderError(c, err, "no se pudo generar la sesión")
		return
	}

	c.Set(middleware.ClaveUsuarioContexto, usuario)

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

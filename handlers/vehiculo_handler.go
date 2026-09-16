package handlers

import (
	"net/http"
	"strconv"

	"AutoMarket/middleware"
	"AutoMarket/services"

	"github.com/gin-gonic/gin"
)

type publicarVehiculoRequest struct {
	Marca       string  `json:"marca"`
	Modelo      string  `json:"modelo"`
	Anio        int     `json:"anio"`
	Precio      float64 `json:"precio"`
	Descripcion string  `json:"descripcion"`
}

// BLOQUE: consultar catalogo (visitante)
func ConsultarCatalogo(c *gin.Context) {
	catalogo, err := services.ObtenerCatalogo()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "no se pudo obtener el catálogo",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"vehiculos": catalogo,
	})
}

// BLOQUE: registrar/publicar vehiculo (vendedor)
func PublicarVehiculo(c *gin.Context) {
	usuario, autenticado := middleware.ObtenerUsuarioDeContexto(c)

	if !autenticado {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "se requiere autenticación",
		})
		return
	}

	var request publicarVehiculoRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "el cuerpo de la solicitud debe ser un JSON válido",
		})
		return
	}

	vehiculo, err := services.RegistrarVehiculo(
		usuario.ID,
		request.Marca,
		request.Modelo,
		request.Anio,
		request.Precio,
		request.Descripcion,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"mensaje":  "publicación registrada, queda pendiente de aprobación",
		"vehiculo": vehiculo,
	})
}

// BLOQUE: ver publicaciones propias (vendedor)
func VerMisPublicaciones(c *gin.Context) {
	usuario, autenticado := middleware.ObtenerUsuarioDeContexto(c)

	if !autenticado {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "se requiere autenticación",
		})
		return
	}

	publicaciones, err := services.ObtenerPublicacionesDeVendedor(usuario.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "no se pudieron obtener las publicaciones",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"vehiculos": publicaciones,
	})
}

// BLOQUE: reportar vehiculo vendido (vendedor)
func ReportarVenta(c *gin.Context) {
	usuario, autenticado := middleware.ObtenerUsuarioDeContexto(c)

	if !autenticado {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "se requiere autenticación",
		})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "el id del vehículo es inválido",
		})
		return
	}

	vehiculo, err := services.ReportarVehiculoVendido(usuario.ID, id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"mensaje":  "publicación marcada como vendida",
		"vehiculo": vehiculo,
	})
}

// BLOQUE: autorizar publicacion (administrador)
func AutorizarPublicacion(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "el id del vehículo es inválido",
		})
		return
	}

	vehiculo, err := services.AutorizarPublicacion(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"mensaje":  "publicación autorizada",
		"vehiculo": vehiculo,
	})
}

// BLOQUE: eliminar publicacion (administrador)
func EliminarPublicacion(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "el id del vehículo es inválido",
		})
		return
	}

	if err := services.EliminarPublicacion(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"mensaje": "publicación eliminada",
	})
}

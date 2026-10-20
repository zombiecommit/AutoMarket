package handlers

import (
	"io"
	"net/http"
	"strconv"
	"strings"

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
		responderError(c, err, "no se pudo obtener el catálogo")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"vehiculos": catalogo,
	})
}

// BLOQUE: registrar/publicar vehiculo (vendedor)
func PublicarVehiculo(c *gin.Context) {
	usuarioID, autenticado := obtenerUsuarioIDAutenticado(c)
	if !autenticado {
		return
	}

	// Con foto: el formulario web envia multipart/form-data.
	if c.ContentType() == "multipart/form-data" {
		publicarVehiculoConFoto(c, usuarioID)
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
		usuarioID,
		request.Marca,
		request.Modelo,
		request.Anio,
		request.Precio,
		request.Descripcion,
	)

	if err != nil {
		responderError(c, err, "no se pudo registrar la publicación")

		return
	}

	c.Set("accounting_recurso_id", strconv.Itoa(vehiculo.ID))

	c.JSON(http.StatusCreated, gin.H{
		"mensaje":  "publicación registrada, queda pendiente de aprobación",
		"vehiculo": vehiculo,
	})
}

// BLOQUE: ver publicaciones propias (vendedor)
func VerMisPublicaciones(c *gin.Context) {
	usuarioID, autenticado := obtenerUsuarioIDAutenticado(c)
	if !autenticado {
		return
	}

	publicaciones, err := services.ObtenerPublicacionesDeVendedor(usuarioID)

	if err != nil {
		responderError(c, err, "no se pudieron obtener las publicaciones")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"vehiculos": publicaciones,
	})
}

// BLOQUE: reportar vehiculo vendido (vendedor)
func ReportarVenta(c *gin.Context) {
	usuarioID, autenticado := obtenerUsuarioIDAutenticado(c)
	if !autenticado {
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "el id del vehículo es inválido",
		})
		return
	}

	vehiculo, err := services.ReportarVehiculoVendido(usuarioID, id)

	if err != nil {
		responderError(c, err, "no se pudo reportar la venta")

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
		responderError(c, err, "no se pudo autorizar la publicación")

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
		responderError(c, err, "no se pudo eliminar la publicación")

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"mensaje": "publicación eliminada",
	})
}

// BLOQUE: obtener id del usuario autenticado o responder 401
func obtenerUsuarioIDAutenticado(c *gin.Context) (int, bool) {
	usuario, autenticado := middleware.ObtenerUsuarioDeContexto(c)

	if !autenticado {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "se requiere autenticación",
		})
		return 0, false
	}

	return usuario.ID, true
}

// BLOQUE: publicar vehiculo con foto (multipart/form-data)
// Campos: marca, modelo, anio, precio, descripcion y el archivo "foto".
func publicarVehiculoConFoto(c *gin.Context, usuarioID int) {
	// Tope un poco mayor que la foto para dejar espacio a los demas campos.
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, services.MaxFotoBytes+(1<<20))

	anio, err := strconv.Atoi(strings.TrimSpace(c.PostForm("anio")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "el año del vehículo es inválido"})
		return
	}

	precio, err := strconv.ParseFloat(strings.TrimSpace(c.PostForm("precio")), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "el precio del vehículo es inválido"})
		return
	}

	archivo, err := c.FormFile("foto")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "la foto del vehículo es obligatoria (JPG, PNG o WEBP, máximo 5 MB)",
		})
		return
	}

	abierto, err := archivo.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no se pudo leer la foto"})
		return
	}
	defer abierto.Close()

	// Se lee un byte de mas para detectar fotos que superen el maximo.
	foto, err := io.ReadAll(io.LimitReader(abierto, services.MaxFotoBytes+1))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no se pudo leer la foto"})
		return
	}

	vehiculo, err := services.RegistrarVehiculoConFoto(
		usuarioID,
		c.PostForm("marca"),
		c.PostForm("modelo"),
		anio,
		precio,
		c.PostForm("descripcion"),
		foto,
	)

	if err != nil {
		responderError(c, err, "no se pudo registrar la publicación")
		return
	}

	c.Set("accounting_recurso_id", strconv.Itoa(vehiculo.ID))

	c.JSON(http.StatusCreated, gin.H{
		"mensaje":  "publicación registrada, queda pendiente de aprobación",
		"vehiculo": vehiculo,
	})
}

package handlers

import (
	"net/http"
	"strconv"

	"AutoMarket/models"
	"AutoMarket/storage"
	"AutoMarket/store"

	"github.com/gin-gonic/gin"
)

// idDesdeParametro extrae y valida un ID entero desde los parámetros de ruta.
// Si el parámetro no es válido, responde con 400 y devuelve err != nil.
func idDesdeParametro(c *gin.Context, nombre string) (int, error) {
	id, err := strconv.Atoi(c.Param(nombre))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "identificador inválido"})
	}
	return id, err
}

// solicitudVehiculo es el cuerpo esperado para publicar un vehículo.
type solicitudVehiculo struct {
	VendedorID  int     `json:"vendedor_id" binding:"required"`
	Marca       string  `json:"marca" binding:"required"`
	Modelo      string  `json:"modelo" binding:"required"`
	Anio        int     `json:"anio" binding:"required"`
	Precio      float64 `json:"precio" binding:"required,gt=0"`
	Descripcion string  `json:"descripcion"`
}

// ---- Visitante ----

// ConsultarCatalogo devuelve únicamente los vehículos ya publicados.
func ConsultarCatalogo(c *gin.Context) {
	vehiculos := store.Datos().ListarVehiculosPublicados()
	c.JSON(http.StatusOK, vehiculos)
}

// ConsultarContacto expone la información de contacto de AutoMarket.
func ConsultarContacto(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"telefono": "+57 300 000 0000",
		"email":    "contacto@automarket.com",
		"horario":  "Lunes a viernes, 8am - 6pm",
	})
}

// ---- Vendedor ----

// PublicarVehiculo registra una nueva publicación en estado "pendiente".
// El administrador debe autorizarla antes de que aparezca en el catálogo.
func PublicarVehiculo(c *gin.Context) {
	var solicitud solicitudVehiculo
	if err := c.ShouldBindJSON(&solicitud); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	usuarios, err := storage.CargarUsuarios()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo leer el almacenamiento de usuarios"})
		return
	}

	vendedorExiste := false
	for _, u := range usuarios {
		if u.ID == solicitud.VendedorID && u.Rol == "vendedor" {
			vendedorExiste = true
			break
		}
	}
	if !vendedorExiste {
		c.JSON(http.StatusNotFound, gin.H{"error": "vendedor no encontrado"})
		return
	}

	vehiculo := &models.Vehiculo{
		VendedorID:  solicitud.VendedorID,
		Marca:       solicitud.Marca,
		Modelo:      solicitud.Modelo,
		Anio:        solicitud.Anio,
		Precio:      solicitud.Precio,
		Descripcion: solicitud.Descripcion,
		Estado:      models.EstadoPendiente,
	}
	vehiculo = store.Datos().CrearVehiculo(vehiculo)

	c.JSON(http.StatusCreated, vehiculo)
}

// MisPublicaciones devuelve todas las publicaciones de un vendedor,
// sin importar su estado (pendiente, publicado o vendido).
func MisPublicaciones(c *gin.Context) {
	vendedorID, err := idDesdeParametro(c, "vendedorId")
	if err != nil {
		return
	}

	vehiculos := store.Datos().ListarVehiculosPorVendedor(vendedorID)
	c.JSON(http.StatusOK, vehiculos)
}

// ReportarVendido marca una publicación del vendedor como vendida.
func ReportarVendido(c *gin.Context) {
	id, err := idDesdeParametro(c, "id")
	if err != nil {
		return
	}

	vehiculo := store.Datos().BuscarVehiculoPorID(id)
	if vehiculo == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "publicación no encontrada"})
		return
	}

	vehiculo.Estado = models.EstadoVendido
	c.JSON(http.StatusOK, vehiculo)
}

// ---- Administrador ----

// AutorizarPublicacion cambia una publicación de "pendiente" a "publicado".
func AutorizarPublicacion(c *gin.Context) {
	id, err := idDesdeParametro(c, "id")
	if err != nil {
		return
	}

	vehiculo := store.Datos().BuscarVehiculoPorID(id)
	if vehiculo == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "publicación no encontrada"})
		return
	}
	if vehiculo.Estado != models.EstadoPendiente {
		c.JSON(http.StatusConflict, gin.H{"error": "la publicación no está pendiente de aprobación"})
		return
	}

	vehiculo.Estado = models.EstadoPublicado
	c.JSON(http.StatusOK, vehiculo)
}

// EliminarPublicacion borra una publicación del catálogo, sin importar su estado.
func EliminarPublicacion(c *gin.Context) {
	id, err := idDesdeParametro(c, "id")
	if err != nil {
		return
	}

	if !store.Datos().EliminarVehiculo(id) {
		c.JSON(http.StatusNotFound, gin.H{"error": "publicación no encontrada"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensaje": "publicación eliminada"})
}

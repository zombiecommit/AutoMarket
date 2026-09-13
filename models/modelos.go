package models

// NOTA: el struct Usuario vive en models/usuario.go (propiedad del equipo de AAA).
// Aquí solo se definen los modelos del catálogo de vehículos.

// EstadoPublicacion representa el ciclo de vida de una publicación de vehículo.
type EstadoPublicacion string

const (
	EstadoPendiente EstadoPublicacion = "pendiente"
	EstadoPublicado EstadoPublicacion = "publicado"
	EstadoVendido   EstadoPublicacion = "vendido"
)

// Vehiculo representa una publicación de un vehículo en el catálogo.
type Vehiculo struct {
	ID          int               `json:"id"`
	VendedorID  int               `json:"vendedor_id"`
	Marca       string            `json:"marca"`
	Modelo      string            `json:"modelo"`
	Anio        int               `json:"anio"`
	Precio      float64           `json:"precio"`
	Descripcion string            `json:"descripcion"`
	Estado      EstadoPublicacion `json:"estado"`
}

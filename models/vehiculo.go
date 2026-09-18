package models

type Vehiculo struct {
	ID          int     `json:"id"`
	VendedorID  int     `json:"vendedor_id"`
	Marca       string  `json:"marca"`
	Modelo      string  `json:"modelo"`
	Anio        int     `json:"anio"`
	Precio      float64 `json:"precio"`
	Descripcion string  `json:"descripcion"`
	Estado      string  `json:"estado"`
}

// ContactoVendedor: datos de contacto del vendedor asociado a una publicacion.
type ContactoVendedor struct {
	Nombre string `json:"nombre"`
	Correo string `json:"correo"`
}

// VehiculoPublico: vehiculo del catalogo junto con los datos de contacto
// del vendedor que lo publico (usado en GET /catalogo).
type VehiculoPublico struct {
	Vehiculo
	Vendedor ContactoVendedor `json:"vendedor"`
}

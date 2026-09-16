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

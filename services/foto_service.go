package services

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"AutoMarket/models"
	"AutoMarket/storage"
)

// BLOQUE: fotos de los vehiculos
// Las fotos se guardan en disco como static/images/vehiculos/<id>.<ext> y la
// ruta publica queda en el campo Foto del vehiculo (se sirve con /static).
const (
	MaxFotoBytes = 5 << 20 // 5 MB
	dirFotos     = "static/images/vehiculos"
	urlFotos     = "/static/images/vehiculos"
)

// BLOQUE: validar la foto y deducir su extension por el contenido real
// (no por el nombre que declare el navegador, que es facil de falsificar).
func extensionDeFoto(datos []byte) (string, error) {
	if len(datos) == 0 {
		return "", nuevoErrorValidacion("la foto del vehículo es obligatoria")
	}

	if len(datos) > MaxFotoBytes {
		return "", nuevoErrorValidacion("la foto no puede superar los 5 MB")
	}

	switch http.DetectContentType(datos) {
	case "image/jpeg":
		return ".jpg", nil
	case "image/png":
		return ".png", nil
	case "image/webp":
		return ".webp", nil
	}

	return "", nuevoErrorValidacion("la foto debe ser una imagen JPG, PNG o WEBP")
}

// BLOQUE: registrar vehiculo junto con su foto (vendedor)
// Primero se valida la foto, luego se crea la publicacion y por ultimo se
// guarda el archivo. Si algo falla despues de crear la publicacion, se
// deshace todo para no dejar publicaciones sin foto ni archivos huerfanos.
func RegistrarVehiculoConFoto(
	vendedorID int,
	marca string,
	modelo string,
	anio int,
	precio float64,
	descripcion string,
	foto []byte,
) (models.Vehiculo, error) {
	extension, err := extensionDeFoto(foto)
	if err != nil {
		return models.Vehiculo{}, err
	}

	vehiculo, err := RegistrarVehiculo(vendedorID, marca, modelo, anio, precio, descripcion)
	if err != nil {
		return models.Vehiculo{}, err
	}

	nombre := fmt.Sprintf("%d%s", vehiculo.ID, extension)
	rutaArchivo := filepath.Join(dirFotos, nombre)

	if err := os.MkdirAll(dirFotos, 0755); err != nil {
		_ = EliminarPublicacion(vehiculo.ID)
		return models.Vehiculo{}, err
	}

	if err := os.WriteFile(rutaArchivo, foto, 0644); err != nil {
		_ = EliminarPublicacion(vehiculo.ID)
		return models.Vehiculo{}, err
	}

	vehiculo.Foto = urlFotos + "/" + nombre

	if err := guardarFotoEnVehiculo(vehiculo.ID, vehiculo.Foto); err != nil {
		_ = os.Remove(rutaArchivo)
		_ = EliminarPublicacion(vehiculo.ID)
		return models.Vehiculo{}, err
	}

	return vehiculo, nil
}

func guardarFotoEnVehiculo(vehiculoID int, rutaFoto string) error {
	vehiculos, err := storage.CargarVehiculos()
	if err != nil {
		return err
	}

	for i := range vehiculos {
		if vehiculos[i].ID == vehiculoID {
			vehiculos[i].Foto = rutaFoto
			break
		}
	}

	return storage.GuardarVehiculos(vehiculos)
}

// BLOQUE: borrar el archivo de la foto al eliminar una publicacion
// Solo borra dentro de static/images/vehiculos para no tocar otros archivos.
func borrarArchivoFoto(rutaFoto string) {
	if rutaFoto == "" {
		return
	}

	ruta := strings.TrimPrefix(rutaFoto, "/")

	if !strings.HasPrefix(ruta, dirFotos+"/") {
		return
	}

	_ = os.Remove(ruta)
}

package store

import (
	"sync"

	"AutoMarket/models"
)

// Store es un almacenamiento en memoria, seguro para uso concurrente,
// exclusivamente para las publicaciones de vehículos.
// Los usuarios NO viven aquí: esos los maneja el paquete "storage"
// (persistencia en JSON) del equipo de AAA.
type Store struct {
	mu                  sync.RWMutex
	vehiculos           map[int]*models.Vehiculo
	siguienteVehiculoID int
}

var datos = &Store{
	vehiculos:           make(map[int]*models.Vehiculo),
	siguienteVehiculoID: 1,
}

// Datos expone la instancia única del almacenamiento en memoria.
func Datos() *Store {
	return datos
}

func (s *Store) CrearVehiculo(v *models.Vehiculo) *models.Vehiculo {
	s.mu.Lock()
	defer s.mu.Unlock()
	v.ID = s.siguienteVehiculoID
	s.siguienteVehiculoID++
	s.vehiculos[v.ID] = v
	return v
}

func (s *Store) BuscarVehiculoPorID(id int) *models.Vehiculo {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.vehiculos[id]
}

func (s *Store) ListarVehiculosPublicados() []*models.Vehiculo {
	s.mu.RLock()
	defer s.mu.RUnlock()
	resultado := []*models.Vehiculo{}
	for _, v := range s.vehiculos {
		if v.Estado == models.EstadoPublicado {
			resultado = append(resultado, v)
		}
	}
	return resultado
}

func (s *Store) ListarVehiculosPorVendedor(vendedorID int) []*models.Vehiculo {
	s.mu.RLock()
	defer s.mu.RUnlock()
	resultado := []*models.Vehiculo{}
	for _, v := range s.vehiculos {
		if v.VendedorID == vendedorID {
			resultado = append(resultado, v)
		}
	}
	return resultado
}

func (s *Store) EliminarVehiculo(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, existe := s.vehiculos[id]; !existe {
		return false
	}
	delete(s.vehiculos, id)
	return true
}

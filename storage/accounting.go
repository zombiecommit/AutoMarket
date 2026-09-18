package storage

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

const archivoAccounting = "accounting.json" //Aquí se guardarán los registros

var accountingMu sync.Mutex //Solo una operación puede modificar Accounting a la vez.

type RegistroAccounting struct {
	ID           int       `json:"id"`
	UsuarioID    int       `json:"usuario_id"`
	Rol          string    `json:"rol"`
	Operacion    string    `json:"operacion"`
	Recurso      string    `json:"recurso"`
	RecursoID    string    `json:"recurso_id,omitempty"`
	FechaHora    time.Time `json:"fecha_hora"`
	Exito        bool      `json:"exito"`
	HashAnterior string    `json:"hash_anterior"`
	Hash         string    `json:"hash"`
}

func CargarRegistrosAccounting() ([]RegistroAccounting, error) {
	accountingMu.Lock()
	defer accountingMu.Unlock()
	return cargarRegistrosAccounting()
}

func cargarRegistrosAccounting() ([]RegistroAccounting, error) {
	datos, err := os.ReadFile(archivoAccounting)
	if os.IsNotExist(err) {
		return []RegistroAccounting{}, nil
	}
	if err != nil {
		return nil, err
	}
	if len(datos) == 0 {
		return []RegistroAccounting{}, nil
	}
	var registros []RegistroAccounting
	err = json.Unmarshal(datos, &registros)
	if err != nil {
		return nil, err
	}
	return registros, nil
}

func GuardarRegistrosAccounting(registros []RegistroAccounting) error {
	accountingMu.Lock()
	defer accountingMu.Unlock()
	return guardarRegistrosAccounting(registros)
}

func guardarRegistrosAccounting(registros []RegistroAccounting) error {
	datos, err := json.MarshalIndent(registros, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile(archivoAccounting, datos, 0644)
}

func calcularHash(registro RegistroAccounting) string {
	contenido := fmt.Sprintf(
		"%d|%d|%s|%s|%s|%s|%s|%t|%s",
		registro.ID,
		registro.UsuarioID,
		registro.Rol,
		registro.Operacion,
		registro.Recurso,
		registro.RecursoID,
		registro.FechaHora.UTC().Format(time.RFC3339Nano),
		registro.Exito,
		registro.HashAnterior,
	)
	hash := sha256.Sum256([]byte(contenido))
	return hex.EncodeToString(hash[:])
}

func verificarIntegridad(registros []RegistroAccounting) error {
	hashAnterior := ""
	for i, registro := range registros {
		if registro.HashAnterior != hashAnterior {
			return fmt.Errorf(
				"integridad comprometida en el registro %d: hash anterior inválido",
				registro.ID,
			)
		}
		hashCalculado := calcularHash(registro)
		if registro.Hash != hashCalculado {
			return fmt.Errorf(
				"integridad comprometida en el registro %d: hash inválido",
				registro.ID,
			)
		}
		if i > 0 && registro.ID <= registros[i-1].ID {
			return fmt.Errorf(
				"integridad comprometida: IDs fuera de secuencia",
			)
		}
		hashAnterior = registro.Hash
	}
	return nil
}

func RegistrarAccounting(
	usuarioID int,
	rol string,
	operacion string,
	recurso string,
	recursoID string,
	exito bool,
) error {
	accountingMu.Lock()
	defer accountingMu.Unlock()

	registros, err := cargarRegistrosAccounting()
	if err != nil {
		return err
	}
	if err := verificarIntegridad(registros); err != nil {
		return err
	}
	id := 1
	hashAnterior := ""
	if len(registros) > 0 {
		ultimo := registros[len(registros)-1]
		id = ultimo.ID + 1
		hashAnterior = ultimo.Hash
	}
	registro := RegistroAccounting{
		ID:           id,
		UsuarioID:    usuarioID,
		Rol:          rol,
		Operacion:    operacion,
		Recurso:      recurso,
		RecursoID:    recursoID,
		FechaHora:    time.Now().UTC(),
		Exito:        exito,
		HashAnterior: hashAnterior,
	}
	registro.Hash = calcularHash(registro)
	registros = append(registros, registro)
	return guardarRegistrosAccounting(registros)
}

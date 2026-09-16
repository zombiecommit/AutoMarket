# Prueba manual: autorización de publicaciones (rol administrador)

Este documento sirve para volver a probar en vivo (demo, sustentación, o debug
rápido) que el middleware de autorización bloquea/permite según el rol.

> Nota: esto mismo ya está automatizado en `main_test.go`
> (`TestFlujoDeNegocioDeVehiculos`). Corre `go test ./...` para verificarlo
> sin pasos manuales. Usa esta guía solo cuando necesites ver las respuestas
> JSON en pantalla (por ejemplo, en una demo en vivo).

## Requisito: servidor corriendo

En una terminal, parado en la raíz del proyecto:

```powershell
go run main.go
```

Déjala corriendo. Todo lo demás va en **otra** terminal.

## 1. Registrar un usuario que será promovido a administrador

```powershell
curl.exe -X POST http://localhost:8080/usuarios -H "Content-Type: application/json" -d '{\"nombre\": \"Admin Prueba\", \"correo\": \"admin@test.com\", \"contrasena\": \"12345678\"}'
```

## 2. Promover ese usuario a administrador (a mano, mientras no exista otra forma)

Abre `storage/usuarios.json`, busca el usuario por su correo, y cambia:

```json
"rol": "vendedor"
```

por:

```json
"rol": "administrador"
```

## 3. Iniciar sesión como administrador

```powershell
curl.exe -X POST http://localhost:8080/login -H "Content-Type: application/json" -d '{\"correo\": \"admin@test.com\", \"contrasena\": \"12345678\"}'
```

Copia el `token` de la respuesta → este es `TOKEN_ADMIN`.

## 4. Registrar e iniciar sesión con un vendedor

```powershell
curl.exe -X POST http://localhost:8080/usuarios -H "Content-Type: application/json" -d '{\"nombre\": \"Vendedor Prueba\", \"correo\": \"vendedor@test.com\", \"contrasena\": \"12345678\"}'

curl.exe -X POST http://localhost:8080/login -H "Content-Type: application/json" -d '{\"correo\": \"vendedor@test.com\", \"contrasena\": \"12345678\"}'
```

Copia el `token` de esa respuesta → este es `TOKEN_VENDEDOR`.

## 5. Publicar un vehículo con el vendedor

```powershell
curl.exe -X POST http://localhost:8080/vehiculos -H "Content-Type: application/json" -H "Authorization: Bearer TOKEN_VENDEDOR" -d '{\"marca\": \"Mazda\", \"modelo\": \"3\", \"anio\": 2020, \"precio\": 60000000, \"descripcion\": \"prueba\"}'
```

Queda en `"estado": "pendiente_aprobacion"`. Copia el `id` devuelto → este es `ID_VEHICULO`.

## 6. Caso permitido: el administrador autoriza la publicación

```powershell
curl.exe -X PATCH http://localhost:8080/vehiculos/ID_VEHICULO/autorizar -H "Authorization: Bearer TOKEN_ADMIN"
```

**Resultado esperado:** `200 OK`, `"estado": "publicado"`.

## 7. Caso bloqueado: el vendedor intenta autorizar (no tiene permiso)

```powershell
curl.exe -X PATCH http://localhost:8080/vehiculos/ID_VEHICULO/autorizar -H "Authorization: Bearer TOKEN_VENDEDOR"
```

**Resultado esperado:** `403 Forbidden`, con un mensaje tipo:
```json
{"error":"el rol 'vendedor' no tiene permiso para realizar esta operación"}
```

## Verificación extra: revisar el registro de accounting

Después de correr los pasos 6 y 7, abre `accounting.json` (se crea en la raíz
del proyecto la primera vez que se registra una operación) y confirma que
quedaron los dos intentos de `AUTORIZAR_PUBLICACION`: uno con `exito: true`
(el del admin) y otro con `exito: false` (el del vendedor).

# AGENTS.md

## Proposito

Este repo es un monolito en Go.

La arquitectura vigente es:

- un solo `ApiService` en protobuf
- un solo servidor gRPC
- transporte concentrado en `api/*`
- logica interna organizada por feature en `internal/features/*`
- webhooks tratados como transporte

Este archivo existe para evitar que nuevos cambios revivan estructuras viejas del repo.

## Mapa rapido

- `cmd/server`
  - bootstrap y wiring
- `api/okgrpc`
  - borde gRPC
  - interceptors
  - metodos del `ApiService`
  - mapping de errores gRPC
- `internal/features/*`
  - dominio, casos de uso, repositorios y errores por feature
- `internal/modules/*`
  - clientes e infraestructura
- `internal/platforms/*`
  - utilidades transversales
- `protobuf/def/v1/*`
  - contrato protobuf
- `protobuf/gen/v1/*`
  - codigo generado

## Regla principal

Si el cambio toca transporte, empieza en `api/*`.

Si el cambio toca negocio o persistencia, empieza en `internal/features/<feature>`.

## Contrato externo

El contrato gRPC se concentra en:

- `protobuf/def/v1/api.proto`
- `protobuf/def/v1/data.proto`
- `protobuf/def/v1/domain.proto`

Regla:

- nuevos metodos van al mismo `ApiService`
- no dividir protobuf en multiples services sin necesidad real

## Responsabilidades

### `api/*`

Debe contener:

- parsing del transporte
- metadata
- auth del transporte
- interceptors
- serializacion de responses
- mapping de errores publicos

No debe contener:

- SQL
- reglas de negocio complejas
- conocimiento detallado del schema

### `internal/features/*`

Debe contener:

- dominio
- services
- repositorios
- DTOs internos
- errores semanticos del feature

No deberia contener codigo nuevo de:

- gRPC
- HTTP handlers
- `status.Error(...)`
- `metadata.FromIncomingContext`

Si algo asi sigue existiendo, es deuda tecnica legacy.

## Auth

La semantica de auth vive en `internal/features/auth`.

El borde gRPC solo debe:

- leer `authorization`
- resolver la sesion
- aplicar autorizacion por metodo
- delegar a `auth.Service`

## Wiring

El punto de composicion principal es `cmd/server/main.go`.

Debe seguir siendo el lugar donde se construyen:

- config
- logger
- postgres
- repositorios
- clientes externos
- `api/okgrpc.Server`
- `grpc.Server`

## Regla para agentes

Haz esto:

- agrega RPCs al `ApiService`
- implementa el borde en `api/okgrpc`
- delega al feature correspondiente
- conserva el monolito simple

Evita esto:

- meter handlers gRPC dentro de features
- crear services protobuf por costumbre
- mover logica de negocio al transporte
- mezclar infraestructura con borde

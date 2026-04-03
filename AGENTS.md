# AGENTS.md

## Proposito

Este repo es un monolito en Go para el backend de `Muy del Campo`.

La arquitectura vigente ya no esta separada por feature en `internal/features/*`.
La decision actual es trabajar con:

- un solo `ApiService` en protobuf
- un solo servidor gRPC
- un solo nucleo de aplicacion en `internal/app`
- un borde de transporte en `api/*`

Este archivo documenta esa decision para evitar regresiones a estructuras viejas del repo.

## Mapa rapido

- `cmd/server`
  - bootstrap y wiring principal
- `api/okgrpc`
  - borde gRPC
  - interceptors
  - metodos del `ApiService`
  - mapping de errores publicos
- `internal/app`
  - dominio
  - inputs y data internos
  - use cases
  - repositorio
  - errores semanticos
- `internal/modules/*`
  - infraestructura y clientes externos
- `internal/platforms/*`
  - config, logger, crypto, validacion y utilidades
- `protobuf/def/v1/*`
  - contrato protobuf
- `protobuf/gen/v1/*`
  - codigo generado
- `infrax/*`
  - artefactos de build/deploy por transporte

## Arquitectura vigente

### 1. Un solo service protobuf

El contrato gRPC se concentra en un unico `ApiService`.

Regla:

- nuevos metodos externos van al mismo `ApiService`
- no dividir protobuf en multiples services salvo necesidad tecnica real

### 2. Un solo nucleo de aplicacion

La logica de negocio vive en `internal/app`.

Eso incluye:

- auth
- sesiones
- usuarios
- validaciones e inputs internos del dominio principal

La decision es intencional:

- el ecommerce no se proyecta como un sistema enorme
- un nucleo unico es mas legible que fragmentar por features sin necesidad
- se prioriza claridad operativa sobre una modularidad excesiva

### 3. Transporte separado del nucleo

El transporte sigue aislado del nucleo.

Hoy existe:

- `api/okgrpc`

Si en el futuro aparece HTTP/webhook como borde activo, debe vivir tambien en `api/*`, no dentro de `internal/app`.

## Regla principal

Si el cambio toca transporte, empieza en `api/*`.

Si el cambio toca negocio, datos internos o persistencia, empieza en `internal/app`.

## Estado real del bootstrap

El entrypoint real es `cmd/server/main.go`.

Secuencia observada hoy:

1. cargar config con `confx.Load()`
2. inicializar logger
3. abrir PostgreSQL
4. construir cliente WhatsApp y `messages.Service`
5. construir `app.Repository`
6. construir `app.UseService`
7. construir `api/okgrpc.Server`
8. construir `grpc.Server` con interceptors
9. registrar `ApiService`
10. abrir listener y servir gRPC

## Responsabilidades por capa

### `api/*`

Debe contener:

- parsing del request del transporte
- metadata
- auth del transporte
- autorizacion por metodo
- interceptors
- serializacion de responses
- mapping de errores publicos

No debe contener:

- SQL
- logica de negocio compleja
- conocimiento detallado del schema

### `internal/app`

Debe contener:

- dominio
- inputs del caso de uso
- data interna
- servicios de aplicacion
- acceso a repositorio
- errores semanticos

Puede contener:

- conversiones a protobuf si el equipo decide tratarlas como conversiones de datos y no como implementacion de gRPC

No debe contener:

- interceptors gRPC
- `metadata.FromIncomingContext`
- registro de services
- wiring del transporte

### `internal/modules/*`

Debe contener:

- postgres
- whatsapp
- otros clientes externos

## Convencion interna del nucleo

Dentro de `internal/app` el patron vigente es:

- `datax.go`
  - input del transporte convertido a input interno
- `data.go`
  - data interna usada por casos de uso o repositorio
- `domain.go`
  - entidades de dominio
- `repository.go`
  - SQL y persistencia
- `useservice.go`
  - orquestacion de casos de uso
- `myerrors.go` / `myerrorx.go`
  - errores semanticos y wrappers

Regla:

- la separacion principal ya no es por feature, sino por responsabilidad dentro del nucleo unico

## Auth

La semantica de auth y sesion vive en `internal/app`.

El borde gRPC en `api/okgrpc` solo debe:

- leer `authorization`
- cargar sesion
- exigir login/staff/root segun el metodo
- delegar al `UseService`

## Wiring

`cmd/server/main.go` sigue siendo el punto de composicion principal.

Debe seguir siendo el lugar donde se construyen:

- config
- logger
- postgres
- clientes externos
- `app.Repository`
- `app.UseService`
- `api/okgrpc.Server`
- `grpc.Server`

## Regla para agentes

Haz esto:

- agrega RPCs al `ApiService`
- implementa el borde en `api/okgrpc`
- delega al `UseService`
- conserva `internal/app` como nucleo unico
- prioriza legibilidad y consistencia sobre separacion artificial

Evita esto:

- reintroducir `internal/features/*` como patron objetivo
- crear services protobuf por costumbre
- mover logica de negocio al transporte
- mezclar SQL o clientes externos dentro de `api/*`

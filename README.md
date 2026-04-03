# apigo

Backend monolitico en Go para `Muy del Campo`.

La arquitectura actual del repo es:

- un solo `ApiService` en protobuf
- un solo servidor gRPC
- un solo nucleo de aplicacion en `internal/app`
- un borde de transporte en `api/*`

## Estructura

- [`cmd/server`](/Users/andrefedev/Documents/Dev/muydelcampo/go/apigo/cmd/server/main.go)
  - bootstrap principal
- [`api/okgrpc`](/Users/andrefedev/Documents/Dev/muydelcampo/go/apigo/api/okgrpc/api.go)
  - borde gRPC
  - interceptors
  - auth del transporte
  - mapping de errores publicos
- [`internal/app`](/Users/andrefedev/Documents/Dev/muydelcampo/go/apigo/internal/app/useservice.go)
  - dominio
  - inputs y data internos
  - use cases
  - repositorio
  - errores semanticos
- [`internal/modules`](/Users/andrefedev/Documents/Dev/muydelcampo/go/apigo/internal/modules/postgres/db.go)
  - infraestructura y clientes externos
- [`internal/platforms`](/Users/andrefedev/Documents/Dev/muydelcampo/go/apigo/internal/platforms/confx/confx.go)
  - config, logger, crypto, validacion y utilidades
- [`protobuf/def/v1`](/Users/andrefedev/Documents/Dev/muydelcampo/go/apigo/protobuf/def/v1/api.proto)
  - contrato protobuf
- [`protobuf/gen/v1`](/Users/andrefedev/Documents/Dev/muydelcampo/go/apigo/protobuf/gen/v1/api_grpc.pb.go)
  - codigo generado

## Contrato

El contrato externo se concentra en un unico service:

- [`ApiService`](/Users/andrefedev/Documents/Dev/muydelcampo/go/apigo/protobuf/def/v1/api.proto)

Regla actual:

- nuevos metodos gRPC se agregan al mismo `ApiService`
- el borde se implementa en `api/okgrpc`

## Flujo interno

El patron actual del proyecto es:

`gRPC -> api/okgrpc -> internal/app/UseService -> internal/app/Repository -> PostgreSQL`

Dentro de `internal/app` el flujo preferido es:

`input del transporte -> input interno -> data interna -> repository`

## Variables de entorno

Requeridas por [`confx.Load()`](/Users/andrefedev/Documents/Dev/muydelcampo/go/apigo/internal/platforms/confx/confx.go):

- `ENV`
- `PORT`
- `WHATSAPP_TOKEN`
- `WHATSAPP_PHONE`
- `PG_DATABASE_URL`

Opcional:

- `GOOGLE_MAPS_API_KEY`

## Desarrollo local

Ejecutar:

```bash
make run
```

Eso carga `.env` y levanta [`cmd/server/main.go`](/Users/andrefedev/Documents/Dev/muydelcampo/go/apigo/cmd/server/main.go).

## Protobuf

Generar Go:

```bash
make _gen_proto_v1
```

Generar Dart:

```bash
make _gen_proto_dart_v2
```

## Deploy

El pipeline gRPC actual vive en:

- [`infrax/grpc/cloudbuild.yaml`](/Users/andrefedev/Documents/Dev/muydelcampo/go/apigo/infrax/grpc/cloudbuild.yaml)

## Regla rapida para cambios

Si el cambio toca transporte:

- empieza en `api/*`

Si el cambio toca negocio, auth, usuarios o persistencia:

- empieza en `internal/app`

Si el cambio toca infraestructura:

- empieza en `internal/modules/*` o `internal/platforms/*`

## Documentacion interna

- [`AGENTS.md`](/Users/andrefedev/Documents/Dev/muydelcampo/go/apigo/AGENTS.md)
- [`plan.md`](/Users/andrefedev/Documents/Dev/muydelcampo/go/apigo/plan.md)
- [`requirements.md`](/Users/andrefedev/Documents/Dev/muydelcampo/go/apigo/requirements.md)

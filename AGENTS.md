# AGENTS.md

## Proposito

Este repositorio implementa un API HTTP en Go para `Muy del Campo`. La base actual ya separa transporte, casos de uso, persistencia e integraciones, pero todavia conviven decisiones nuevas con codigo legacy, drift de documentacion y varias inconsistencias de contrato que generan deuda tecnica real.

Este archivo describe el estado observado en el codigo de hoy. Debe tomarse como referencia operativa para cualquier agente o desarrollador que siga trabajando en este repo.

## Stack real

- Lenguaje: Go `1.26` en [go.mod](/Users/andrefedev/Documents/Dev/muydelcampo/go/apigo/go.mod)
- Transporte activo: HTTP JSON sobre `net/http` + `github.com/go-chi/chi/v5`
- Base de datos activa: PostgreSQL via `github.com/jackc/pgx/v5/pgxpool`
- Runtime objetivo: Google Cloud Run
- Build/deploy: [Dockerfile](/Users/andrefedev/Documents/Dev/muydelcampo/go/apigo/Dockerfile) + [cloudbuild.yaml](/Users/andrefedev/Documents/Dev/muydelcampo/go/apigo/cloudbuild.yaml)
- Logging estructurado: `log/slog` via `internal/platforms/loggerx`
- Integracion externa activa: WhatsApp Cloud API
- Integraciones presentes pero no conectadas al bootstrap actual: Twilio, Google Maps, Storage, Pub/Sub

## Entry point y arranque

El entrypoint real es [cmd/server/main.go](/Users/andrefedev/Documents/Dev/muydelcampo/go/apigo/cmd/server/main.go).

Secuencia observada:

1. Carga configuracion con `configx.Load()`
2. Inicializa logger global con `loggerx.SetupLogger`
3. Abre pool PostgreSQL con `postgres.Open`
4. Construye `Pgdb`
5. Construye repositorios `auth` y `users`
6. Construye cliente WhatsApp y `messages.Service`
7. Construye servicios `auth` y `users`
8. Construye middleware de identidad
9. Construye router HTTP
10. Levanta `http.Server`
11. Hace graceful shutdown con `SIGINT` y `SIGTERM`

Notas relevantes:

- Ya no hay `init()` global para bootstrap.
- El wiring se hace explicitamente desde `main`.
- `auth.Service` depende de `messages.Service`.
- El router tambien monta la feature `app`, no solo `auth` y `users`.

## Estructura del codigo

- `cmd/server`: composicion y bootstrap
- `internal/features/*`: handlers HTTP, servicios y repositorios por feature
- `internal/modules/*`: infraestructura de negocio o integraciones externas
- `internal/platforms/*`: capacidades transversales reutilizables

El patron dominante sigue siendo:

`handler -> service -> repository -> postgres`

Responsabilidad por capa:

- `handler`: parseo HTTP, validacion de entrada y serializacion de salida/error
- `service`: reglas de aplicacion, orquestacion y traduccion de errores
- `repository`: SQL y adaptacion a `aerrx`
- `domain`: modelos del feature
- `data/datax`: DTOs internos y validacion de request
- `myerrors`: errores publicos del feature

## Features activas

### `app`

Vive en `internal/features/app` y hoy solo expone webhook HTTP:

- `GET /v1/app/webhook`
- `POST /v1/app/webhook`

Estado actual:

- `GET` responde el `hub.challenge` si `hub.mode == "subscribe"`
- `POST` lee y loggea headers/body y responde `EVENT_RECEIVED`
- No usa todavia `internal/modules/whatsapp/webhooks`
- No hay verificacion real de firma ni de token del webhook

### `auth`

Expone:

- `POST /v1/auth/code`

Flujo actual:

1. Decodifica JSON con `okhttpx.DecodeJson`
2. Normaliza y valida telefono
3. Genera OTP de 6 digitos
4. Inserta OTP en `users_codes`
5. Envia template de WhatsApp
6. Responde `201` con `{ "id": "<ref>" }`

Importante:

- La limpieza del OTP si falla WhatsApp esta comentada, no activa
- El servicio retorna tambien el codigo OTP, aunque el handler no lo expone

### `users`

Expone:

- `GET /v1/users/me`

Flujo actual:

1. `AttachIdentity` lee `Authorization`
2. Extrae token con prefijo exacto `Bearer `
3. Resuelve identidad con lookup directo en `users.idk`
4. Guarda `*Identity` en `context.Context`
5. `IsAuthenticated` exige identidad valida
6. `me` carga el usuario desde `users`

## Router HTTP

El router principal vive en [internal/platforms/okhttpx/app.go](/Users/andrefedev/Documents/Dev/muydelcampo/go/apigo/internal/platforms/okhttpx/app.go).

Middlewares globales instalados:

- `middleware.Logger`
- `middleware.RealIP`
- `middleware.RequestID`
- `middleware.Recoverer`
- `middleware.URLFormat`
- `middleware.CleanPath`

Rutas base:

- `GET /healthz`
- `GET /readyz`

`/readyz` usa `pool.Ping()`.

## Identidad y autorizacion

La identidad esta encapsulada en `internal/features/auth`.

Modelo actual:

- `UserRef`
- `IdToken`
- `IsSuper`
- `IsStaff`
- `IsActive`

Helpers disponibles:

- `IsAuthenticated()`
- `CanAccessBackoffice()`
- `CanManageUsers()`

Observaciones:

- No hay JWT ni validacion criptografica.
- El bearer token funciona como lookup directo sobre `users.idk`.
- Un token inexistente hoy se traduce a `401`.

## Errores

Convencion actual:

- repository: retorna `apperr` tecnico
- service: envuelve con `apperr.Wrap` o agrega contrato publico con `WithPublic(...)`
- handler/http: usa `okhttpx.Fail`

Paquete activo:

- [internal/platforms/apperr/error.go](/Users/andrefedev/Documents/Dev/muydelcampo/go/apigo/internal/platforms/apperr/error.go): error canonico con `Op`, `Kind`, `Code`, `Message`, `Cause`

Paquetes legacy presentes en el repo pero fuera del flujo activo:

- [internal/platforms/aerr/aerrx/error.go](/Users/andrefedev/Documents/Dev/muydelcampo/go/apigo/internal/platforms/aerr/aerrx/error.go)
- [internal/platforms/aerr/derrx/error.go](/Users/andrefedev/Documents/Dev/muydelcampo/go/apigo/internal/platforms/aerr/derrx/error.go)
- [internal/platforms/aerr/perrx/error.go](/Users/andrefedev/Documents/Dev/muydelcampo/go/apigo/internal/platforms/aerr/perrx/error.go)

Mapping HTTP actual:

- `not_found` -> `404`
- `validation` -> `400`
- `unauthorized` -> `401`
- `forbidden` -> `403`
- `conflict` -> `409`
- default -> `500`

Nota:

- `okhttpx.slogInternalError` usa `slog`.
- Parte del codigo todavia mezcla `log.Printf` y `fmt.Printf` fuera de este esquema.

## Persistencia

El acceso a DB esta en `internal/modules/postgres`.

Piezas relevantes:

- `postgres.Open`: parsea DSN, configura pool y hace `Ping`
- `postgres.Pgdb`: wrapper sobre pool
- `Pgdb.WithTx`: propaga transacciones por `context.Context`

Tablas activamente usadas desde features actuales:

- `users`
- `users_codes`

Advertencia importante:

- [pgdb.sql](/Users/andrefedev/Documents/Dev/muydelcampo/go/apigo/pgdb.sql) no es una fuente canonica del esquema. Mezcla SQL operativo, experimentos, DDL manual y consultas auxiliares.

## WhatsApp

La integracion activa vive en `internal/modules/whatsapp` y `internal/modules/whatsapp/messages`.

Estado actual:

- `whatsapp.Config` soporta `ApiBaseUrl` y `ApiVersion`, pero el bootstrap no los carga desde env
- El cliente usa defaults `https://graph.facebook.com` y `v25.0`
- `messages.Service.SendTemplate` compone la ruta `<phone-id>/messages`
- El cliente hoy loggea request body y URL
- El cliente devuelve errores genericos con el body remoto sin tipado propio

## Variables de entorno reales del bootstrap

Requeridas hoy por `configx.Load()`:

- `ENV`
- `PORT`
- `WHATSAPP_TOKEN`
- `WHATSAPP_PHONE`
- `PG_DATABASE_URL`

Opcionales cargadas hoy:

- `GOOGLE_MAPS_API_KEY`

Campos presentes en `Config` pero no cargados desde env:

- `GoogleGeminiApiKey`

Capacidades soportadas por codigo de WhatsApp pero no expuestas desde `configx.Load()`:

- `WHATSAPP_BASE_URL`
- `WHATSAPP_VERSION`

## Build, deploy y verificacion

- `Dockerfile` compila binario en multi-stage y corre en `scratch`
- `cloudbuild.yaml` despliega a Cloud Run con secretos para `PG_DATABASE_URL`, `WHATSAPP_TOKEN`, `WHATSAPP_PHONE` y `GOOGLE_MAPS_API_KEY`
- El binario `go` no esta en `PATH` del shell, pero si existe en `/Users/andrefedev/sdk/go1.26.0/bin/go`
- `go test ./...` no pudo completarse aqui por falta de acceso de red a `proxy.golang.org`
- No hay archivos `_test.go` en el repo hoy

## Hallazgos y deuda tecnica vigente

Estas son las discrepancias mas importantes detectadas en el estado actual:

- Drift entre documentacion y codigo: el `AGENTS.md` previo ya no describia fielmente rutas, config y verificacion real
- Contrato de auth fragil: token bearer opaco, sin firma, sin expiracion y con semantica `404` cuando falla la identidad
- Filtracion de datos sensibles en logs: bearer token, payloads de WhatsApp y cuerpos completos de webhook
- Inconsistencia de naming de modelo/esquema: aparecen `phone`, `lookups`, `idk` y `ref` sin una convencion unificada
- Repositorio `users` mezcla columnas `lookups` y `phone`, lo que indica drift de esquema o consultas obsoletas
- Webhook en `app` esta acoplado al router final pero desacoplado del modulo `internal/modules/whatsapp/webhooks`, que hoy esta muerto
- Exceso de codigo comentado y placeholders vacios, lo que hace dificil distinguir backlog real de residuos legacy
- No existe suite de tests automatizados para handlers, servicios ni repositorios
- `go.mod` carga dependencias pesadas no visibles en el bootstrap actual, lo que sugiere arrastre historico

## Convenciones para futuros cambios

- Mantener el patron `handler -> service -> repository`
- No crear clientes de infraestructura dentro de handlers
- Dejar el wiring en `cmd/server/main.go` o en un modulo de bootstrap explicito
- Hacer que repository devuelva solo errores tecnicos
- Traducir errores publicos en service o handler
- Centralizar respuestas HTTP con `okhttpx.Json` y `okhttpx.Fail`
- No usar `pgdb.sql` como contrato canonico; si el proyecto sigue creciendo, migrar a migraciones versionadas
- Eliminar logs de secretos o payloads completos antes de seguir ampliando auth/webhooks
- Definir una convencion unica para campos de identidad y telefono antes de sumar mas features

# AGENTS.md

## Propósito

Este repositorio implementa un API HTTP en Go para `Muy del Campo`. El estado actual del código muestra una base modular con foco en autenticación, identidad y usuarios, montada sobre `chi`, `pgxpool` y utilidades internas para configuración, errores y transporte HTTP.

Este archivo sirve como contexto operativo para cualquier agente o desarrollador que continúe el trabajo en este repo. Describe el estado real del código, incluyendo mejoras recientes y las discrepancias de infraestructura que todavía siguen abiertas.

## Stack e infraestructura

- Lenguaje: Go `1.26` en [`go.mod`](/Users/andrefedev/Documents/Dev/muydelcampo/go/apigo/go.mod)
- Transporte activo: HTTP JSON sobre `net/http` + `github.com/go-chi/chi/v5`
- Base de datos activa: PostgreSQL vía `github.com/jackc/pgx/v5/pgxpool`
- Runtime objetivo: Google Cloud Run
- Build/deploy: [`Dockerfile`](/Users/andrefedev/Documents/Dev/muydelcampo/go/apigo/Dockerfile) + [`cloudbuild.yaml`](/Users/andrefedev/Documents/Dev/muydelcampo/go/apigo/cloudbuild.yaml)
- Logging: `log/slog`, texto en `dev`, JSON fuera de `dev`
- Integración externa en curso: WhatsApp Cloud API
- Integraciones auxiliares presentes pero no conectadas al bootstrap actual: Twilio, Google Maps, Google Cloud Storage, Pub/Sub

## Entry point y arranque

El entrypoint es [`cmd/server/main.go`](/Users/andrefedev/Documents/Dev/muydelcampo/go/apigo/cmd/server/main.go).

Flujo esperado de arranque:

1. Cargar configuración con `confx.Load()`
2. Inicializar logger global con `loggex.SetupLogger`
3. Abrir pool PostgreSQL con `postgres.Open`
4. Construir repositorios
5. Construir servicios
6. Construir middlewares
7. Construir router HTTP
8. Levantar `http.Server`
9. Hacer graceful shutdown con `SIGINT` y `SIGTERM`

Observaciones importantes del estado actual:

- `main` ya no usa `init()` global; la carga de config, logger, DB y wiring de dependencias ocurre de forma explícita.
- `auth.Service` ya recibe `MessageService` desde `main`.
- `auth.Code` hace cleanup best-effort del OTP si falla el envío por WhatsApp.
- `go test ./...` no pudo ejecutarse en este entorno porque no existe binario `go` disponible en shell.

## Arquitectura

La estructura sigue un patrón por feature y por plataforma:

- `internal/features/*`: casos de uso y transporte por feature
- `internal/modules/*`: integraciones o módulos de infraestructura de negocio
- `internal/platforms/*`: capacidades transversales reutilizables
- `cmd/server`: composición y bootstrap

### Capas por feature

El patrón dominante es:

`handler -> service -> repository -> postgres`

Responsabilidad por capa:

- `handler`: parsea request HTTP, valida input, invoca service, serializa respuesta o error
- `service`: concentra lógica de aplicación y traducción de errores técnicos a errores públicos
- `repository`: ejecuta SQL y devuelve errores tipificados con `aerrx`
- `domain`: define modelos del feature
- `data/datax`: DTOs internos, validación y mapeo
- `myerrors`: errores públicos del feature

Features activas:

- `auth`: emisión de código OTP y resolución de identidad desde header `Authorization`
- `users`: lectura de usuario autenticado

## Transporte HTTP

El router principal vive en [`internal/platforms/httpx/app.go`](/Users/andrefedev/Documents/Dev/muydelcampo/go/apigo/internal/platforms/httpx/app.go).

### Middlewares globales

Se instalan estos middlewares de `chi`:

- `middleware.Logger`
- `middleware.RealIP`
- `middleware.RequestID`
- `middleware.Recoverer`
- `middleware.URLFormat`
- `middleware.CleanPath`

### Endpoints base

- `GET /healthz`
- `GET /readyz`

`/readyz` usa `pool.Ping()` contra PostgreSQL.

### Rutas versionadas actuales

Bajo `/v1`:

- `/auth`
- `/users`

Rutas concretas observadas:

- `POST /v1/auth/code`
- `GET /v1/users/me`

Notas:

- `/v1/users/me` exige autenticación mediante `auth.Middleware`.
- `users.me` responde con `200 OK`.
- `auth.code` responde con `201 Created` y payload `{ "id": "<ref>" }`.

## Identidad y autorización

La identidad está encapsulada en `internal/features/auth`.

### Flujo

1. `AttachIdentity` lee el header `authorization`
2. Extrae token con prefijo `Bearer `
3. Resuelve identidad en DB con `ResolveIdentityByIdToken`
4. Guarda `*Identity` en `context.Context`
5. `IsAuthenticated` exige una identidad válida

### Modelo de identidad

`Identity` contiene:

- `UserRef`
- `IdToken`
- `IsSuper`
- `IsStaff`
- `IsActive`

Helpers disponibles:

- `IsAuthenticated()`
- `CanAccessBackoffice()`
- `CanManageUsers()`

Actualmente no hay JWT ni validación criptográfica del token; el `Bearer token` funciona como lookup directo del campo `users.idk`.

## Errores

El manejo de errores está bien separado y es una de las convenciones más importantes del repo.

### Paquetes

- [`aerrx`](/Users/andrefedev/Documents/Dev/muydelcampo/go/apigo/internal/platforms/aerr/aerrx/error.go): error técnico, `Kind`, `Oper`, `Cause`
- [`derrx`](/Users/andrefedev/Documents/Dev/muydelcampo/go/apigo/internal/platforms/aerr/derrx/error.go): error público, `Code`, `Body`
- [`perrx`](/Users/andrefedev/Documents/Dev/muydelcampo/go/apigo/internal/platforms/aerr/perrx/error.go): serialización pública final para HTTP

Documento de referencia existente:

- [`docs/aerrx-derrx-perrx.md`](/Users/andrefedev/Documents/Dev/muydelcampo/go/apigo/docs/aerrx-derrx-perrx.md)

### Regla práctica

- repository: retorna `aerrx`
- service: puede envolver con `aerrx.Wrap` o traducir a `derrx`
- handler/http: usa `httpx.Fail`

### Mapping HTTP

`httpx.ParseError` mapea `aerrx.Kind` a status:

- `not_found` -> `404`
- `validation` -> `400`
- `unauthorized` -> `401`
- `forbidden` -> `403`
- `conflict` -> `409`
- default -> `500`

Errores `5xx` se registran con `slog`.

## Validación y normalización

Helpers actuales:

- [`validationx`](/Users/andrefedev/Documents/Dev/muydelcampo/go/apigo/internal/platforms/validator/validationx/validator.go)
- [`normalizex`](/Users/andrefedev/Documents/Dev/muydelcampo/go/apigo/internal/platforms/validator/normalizex/normalize.go)

Puntos relevantes:

- `IsPhoneNumber` valida números móviles colombianos
- `IsOneTimeCode` valida OTP de 6 dígitos
- `ClearString` colapsa espacios
- `NormalizeTitle` y `NormalizarStreet` normalizan texto
- `httpx.DecodeJson` rechaza campos desconocidos y múltiples objetos JSON en el body

## Persistencia y acceso a datos

El acceso a DB está en `internal/modules/postgres`.

### Componentes

- `postgres.Open`: crea `pgxpool.Pool`, configura límites y hace `Ping`
- `postgres.Pgdb`: wrapper sobre pool con soporte para `Exec`, `Query`, `QueryRow`
- `Pgdb.WithTx`: propaga transacciones usando `context.Context`

### Tablas observadas en código

Las features activas interactúan con:

- `users`
- `users_codes`

Consultas actuales:

- `auth.Repository.CodeInsert`
- `auth.Repository.CodeDelete`
- `auth.Repository.CodeSelect`
- `users.Repository.Select`
- `users.Repository.SelectByPhone`
- `auth.Repository.SelectIdentityByIdToken`

Nota sobre esquema:

- [`pgdb.sql`](/Users/andrefedev/Documents/Dev/muydelcampo/go/apigo/pgdb.sql) contiene SQL mixto, migraciones manuales y consultas de apoyo; no debe asumirse como fuente canónica completa del esquema actual.

## Integraciones externas

### WhatsApp

La integración activa en código apunta a Meta WhatsApp Cloud API dentro de `internal/modules/whatsapp`.

Piezas disponibles:

- `Client`
- `MessageService.SendTemplate`

Estado:

- hay cliente HTTP y envío de template
- el wiring ya se hace desde `main`
- `MessageService.SendTemplate` usa `context.Context`
- el cliente compone la URL con `baseURL + version + phoneID`
- si el envío falla, `auth.Service.Code` intenta borrar el OTP recién creado

### Twilio

Existe soporte en [`internal/platforms/confx/conf.go`](/Users/andrefedev/Documents/Dev/muydelcampo/go/apigo/internal/platforms/confx/conf.go), pero no está conectado al arranque actual.

### Google Maps, Storage y Pub/Sub

También existen helpers en `confx/conf.go`, pero hoy no forman parte del bootstrap efectivo.

## Variables de entorno

### Requeridas por el bootstrap actual

- `ENV`
- `PORT`
- `WHATSAPP_TOKEN`
- `WHATSAPP_PHONE`
- `WHATSAPP_BASE_URL`
- `WHATSAPP_VERSION`
- `PG_DATABASE_URL`

Notas:

- `PORT` se normaliza a formato `:8080` si llega como `8080`.
- `WHATSAPP_BASE_URL` y `WHATSAPP_VERSION` tienen defaults en `confx.Load()`.

### Declaradas en `Config` pero no completamente usadas

- `GOOGLE_MAPS_API_KEY`
- `GoogleGeminiApiKey` existe en `Config`, pero no se carga desde env en `Load()`

### Variables legacy o usadas por helpers no conectados

- `PGDB_URI`
- `TWILIO_SID`
- `TWILIO_TOKEN`
- `TWILIO_ACCOUNT_SID`
- `GOOGLE_APPLICATION_CREDENTIALS`

### Regla de seguridad

No copiar secretos reales a documentación, commits ni respuestas. El repo contiene un `.env` local, pero este archivo no debe replicarse ni volcar sus valores en otros artefactos.

## Deploy y operación

### Docker

El binario se compila en multi-stage build y termina en imagen `scratch`.

Observación:

- [`Dockerfile`](/Users/andrefedev/Documents/Dev/muydelcampo/go/apigo/Dockerfile) usa una imagen `golang` alineada con `go.mod`.

### Cloud Build / Cloud Run

`cloudbuild.yaml` construye imagen, hace push y despliega a Cloud Run con:

- `--use-http2`
- `--min-instances=1`
- `--max-instances=10`
- `--session-affinity`
- `--no-cpu-throttling`
- `--set-cloudsql-instances`

Observación importante:

- `cloudbuild.yaml` ahora publica `PG_DATABASE_URL`, `WHATSAPP_TOKEN` y `WHATSAPP_PHONE`, y fija `WHATSAPP_BASE_URL`/`WHATSAPP_VERSION` como env vars.

## Estado funcional actual

Lo que sí está más cerca de funcionar conceptualmente:

- health checks
- conexión PostgreSQL
- resolución de identidad por bearer token contra DB
- lectura de `/v1/users/me`
- generación de OTP en `auth`

Lo que está incompleto o inconsistente:

- alineación entre `.env` local y secretos reales configurados en Cloud Run

## Convenciones para futuros cambios

- Mantener el patrón `handler -> service -> repository`
- No devolver mensajes públicos desde repository
- Traducir errores de dominio en service usando `derrx`
- Centralizar respuestas HTTP con `httpx.Json` y `httpx.Fail`
- Si se agrega autenticación real, evolucionar `auth` sin romper el contrato de contexto `WithIdentity/IdentityFromContext`
- Mantener el wiring de integraciones en `main` y no crear clientes dentro de handlers o services
- Si se formalizan migraciones, separar `pgdb.sql` en migraciones versionadas

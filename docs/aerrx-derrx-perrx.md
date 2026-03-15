# aerrx, derrx y perrx

Este documento describe cómo se manejan los errores en la estructura actual del proyecto.

Paquetes involucrados:

- `internal/platforms/aerr/aerrx`
- `internal/platforms/aerr/derrx`
- `internal/platforms/aerr/perrx`
- `internal/platforms/httpx`

La idea central es separar tres responsabilidades:

1. clasificar técnicamente el error
2. definir el contrato público del error
3. serializar ese error en el transporte HTTP

## Resumen corto

### `aerrx`

Responsabilidad:

- clasificar el error por `Kind`
- guardar la operación donde ocurrió (`Oper`)
- preservar la causa real (`Cause`)

No debe decidir el mensaje que verá el cliente.

Se usa principalmente en:

- repositories
- módulos de infraestructura
- services cuando solo agregan contexto técnico
- helpers HTTP para etiquetar errores de transporte como `validation`

### `derrx`

Responsabilidad:

- exponer un `Code` estable
- exponer un `Body` seguro para el cliente
- envolver la causa real

Sí representa el contrato público que frontend o mobile pueden consumir.

Se usa principalmente en:

- validaciones de request
- services cuando traducen un error técnico a uno de dominio
- middlewares o helpers HTTP cuando el error nace en transporte

### `perrx`

Responsabilidad:

- leer la cadena de errores
- construir el payload público final
- aplicar fallbacks cuando no existe `derrx`

No toma decisiones de negocio. Solo presenta el error.

Se usa desde `httpx.ParseError`.

## Estructura actual del flujo

La estructura dominante del repo es:

```text
handler -> service -> repository -> postgres
          middleware
```

La cadena de error esperada es:

```text
error real -> aerrx -> derrx -> perrx -> JSON HTTP
```

No siempre aparecen las tres capas. Por ejemplo:

- un repository puede retornar solo `aerrx`
- un handler o middleware puede crear directamente `derrx` si el error nace en el transporte
- `perrx` siempre se aplica al final en HTTP

## Qué hace cada paquete

## `aerrx`

Archivo:

- `internal/platforms/aerr/aerrx/error.go`

Modelo:

```go
type Error struct {
    Kind  Kind
    Oper  string
    Cause error
}
```

`Kind` hoy soporta:

- `not_found`
- `internal`
- `validation`
- `unauthorized`
- `forbidden`
- `conflict`

Helpers importantes:

- `aerrx.New(kind, oper, cause)`
- `aerrx.Wrap(oper, cause)`
- `aerrx.KindOf(err)`
- `aerrx.IsKind(err, kind)`
- `aerrx.OperOf(err)`

Uso correcto:

- en `repository`, para clasificar errores de DB
- en `service`, para agregar contexto técnico sin volver público el detalle
- en transporte, para etiquetar errores nacidos en HTTP como `validation` o `unauthorized`

## `derrx`

Archivo:

- `internal/platforms/aerr/derrx/error.go`

Modelo:

```go
type Error struct {
    Code  string
    Body  string
    Cause error
}
```

Helpers importantes:

- `derrx.New(code, body, cause)`
- `derrx.NewKind(kind, oper, code, body, cause)`
- `derrx.Validation(oper, code, body, cause)`
- `derrx.Unauthorized(oper, code, body, cause)`
- `derrx.Forbidden(oper, code, body, cause)`
- `derrx.NotFound(oper, code, body, cause)`
- `derrx.Conflict(oper, code, body, cause)`
- `derrx.Internal(oper, code, body, cause)`
- `derrx.CodeOf(err)`
- `derrx.BodyOf(err)`

Uso correcto:

- cuando el cliente necesita un `code` estable
- cuando el cliente necesita un mensaje seguro
- cuando un `service` traduce un error técnico a uno de dominio
- cuando el handler o middleware detecta un error de transporte corregible por el cliente

Regla importante:

- `derrx` por sí solo no define el `HTTP status`
- el `HTTP status` siempre sale del `aerrx.Kind`
- si devuelves un `derrx` sin ningún `aerrx` en la cadena, `aerrx.KindOf(err)` caerá en `internal`
- en la práctica eso puede terminar en `500`

Por eso, para errores públicos nuevos, prefiere `derrx.Validation`, `derrx.Unauthorized`, `derrx.NotFound`, etc. en lugar de `derrx.New(...)` directo.

## Sobre `cause`

`cause` es opcional.

No necesitas inventar un `error` solo para completar el constructor.

Casos donde `cause = nil` está bien:

- falta identidad en el contexto
- la ruta requiere autenticación
- un request incumple una regla propia del handler
- construyes un error de dominio sin una causa externa real

Casos donde conviene preservar `cause`:

- `json.Decode`
- errores de DB
- errores de red
- errores de librerías externas

La mejora introducida en `derrx` justamente busca eso: poder construir errores públicos tipados aun cuando no tengas una causa real.

## `perrx`

Archivo:

- `internal/platforms/aerr/perrx/error.go`

Modelo HTTP final:

```go
type PublicError struct {
    Code string `json:"code"`
    Body string `json:"body"`
}
```

`perrx.FromError(err)` hace esto:

1. intenta leer `Code` y `Body` desde `derrx`
2. obtiene el `Kind` desde `aerrx`
3. si falta `Code` o `Body`, aplica un fallback por `Kind`

Fallbacks actuales:

- `not_found` -> `No se encontró el recurso solicitado`
- `validation` -> `Los datos enviados no son válidos`
- `unauthorized` -> `Debes iniciar sesión`
- `forbidden` -> `No tienes permisos para realizar esta acción`
- `conflict` -> `La operación entra en conflicto con el estado actual`
- default -> `Ha ocurrido un error inesperado`

## Cómo se convierte en HTTP

El transporte HTTP vive en `internal/platforms/httpx`.

Flujo actual:

1. `httpx.Fail(w, r, err)` recibe cualquier `error`
2. `httpx.ParseError(err)` resuelve `status` y `perrx.PublicError`
3. si el status es `5xx`, se registra el error con `slog`
4. `httpx.Json` serializa el payload

Mapping actual de `Kind -> status`:

- `not_found` -> `404`
- `validation` -> `400`
- `unauthorized` -> `401`
- `forbidden` -> `403`
- `conflict` -> `409`
- default -> `500`

La consecuencia práctica es esta:

- si la cadena contiene `aerrx.KindUnauthorized`, el status será `401`
- si la cadena solo contiene `derrx` y ningún `aerrx`, el status caerá en `500`

## Flujo recomendado por capa

### Repository

Debe devolver `aerrx` y nada más.

Ejemplo alineado al repo:

```go
if errors.Is(err, pgx.ErrNoRows) {
    return nil, aerrx.New(aerrx.KindNotFound, "User.Repository.Select", err)
}

return nil, aerrx.New(aerrx.KindInternal, "User.Repository.Select", err)
```

No debería devolver:

- `derrx`
- mensajes públicos
- `http status`

### Service

Puede hacer dos cosas:

1. propagar el error técnico con más contexto
2. traducir un error técnico a un error público de dominio

Ejemplo real del patrón:

```go
user, err := s.deps.UserRepository.Select(ctx, ref)
if err != nil {
    if aerrx.IsKind(err, aerrx.KindNotFound) {
        return nil, ErrUserNotFound(err)
    }

    return nil, aerrx.Wrap("Users.Service.GetByRef", err)
}
```

Regla útil:

- si el error sigue siendo interno al sistema, usa `aerrx.Wrap`
- si ya tiene significado para el cliente, tradúcelo a `derrx`
- si el mensaje genérico por `Kind` es suficiente, no necesitas `derrx`

### Handler

El handler no debería decidir el `status` manualmente para cada error.

Debe:

1. decodificar request
2. validar request
3. invocar el service
4. delegar la respuesta de error a `httpx.Fail`

Ejemplo real de `auth`:

```go
var req CodeRequest
if err := httpx.DecodeJson(r, &req, oper); err != nil {
    httpx.Fail(w, r, err)
    return
}

req.Normalize()
if err := req.Validate(); err != nil {
    httpx.Fail(w, r, err)
    return
}

id, code, err := h.deps.Service.Code(ctx, req.Phone)
if err != nil {
    httpx.Fail(w, r, err)
    return
}
```

### Middleware

En este proyecto los middlewares también participan del mismo contrato de errores.

Ejemplo real de `auth.IsAuthenticated`:

```go
httpx.Fail(
    w, r,
    derrx.Unauthorized(
        "Auth.Middlex.IsAuthenticated",
        "auth.authentication_required",
        "Debes iniciar sesión para continuar",
        nil,
    ),
)
```

Eso es correcto porque:

- el error nace en transporte/autorización
- el cliente necesita un código público estable
- el `Kind` sigue permitiendo mapear a `401`

## Errores que nacen en HTTP

Hay errores que no pertenecen al dominio sino al transporte HTTP:

- body vacío
- JSON inválido
- tipo de dato incorrecto
- campo desconocido
- múltiples objetos JSON en el body

Esos errores deben nacer en el handler o en helpers HTTP, no en repository ni service.

Hoy eso vive en `httpx.DecodeJson`.

Ejemplo:

```go
return derrx.Validation(
    op,
    "http.invalid_json",
    "El cuerpo JSON no es válido",
    err,
)
```

## Cuándo usar solo `aerrx`

No todo error necesita `derrx`.

Usa solo `aerrx` cuando el fallback público por `Kind` sea suficiente.

Ejemplos:

- `not_found` genérico
- `unauthorized` genérico
- `forbidden` genérico
- `internal` donde no quieres exponer detalle

Ejemplo:

```go
return aerrx.New(aerrx.KindUnauthorized, "Users.Handlex.Me", nil)
```

Eso permite que HTTP responda correctamente con `401` y que `perrx` construya:

```json
{
  "code": "unauthorized",
  "body": "Debes iniciar sesión"
}
```

Usa `derrx` solo cuando realmente necesitas un `code` o `body` más específico para cliente.

## Ejemplos de uso reales

### Validación de request

`auth.CodeRequest.Validate` retorna un `derrx` con `KindValidation`:

```go
return derrx.Validation(
    "Auth.CodeRequest.Validate",
    "auth.invalid_phone",
    "El número de teléfono no es válido",
    nil,
)
```

### Error técnico desde repository

`users.Repository.Select` clasifica `not_found` e `internal` usando `aerrx`.

### Traducción a error público en service

`users.Service.GetByRef` traduce `not_found` a `users.user_not_found`.

### Serialización final

`httpx.Fail` termina devolviendo un JSON como:

```json
{
  "code": "users.user_not_found",
  "body": "Usuario no encontrado"
}
```

Si no existe `derrx`, `perrx` genera un fallback:

```json
{
  "code": "internal",
  "body": "Ha ocurrido un error inesperado"
}
```

## Qué no hacer

No uses `derrx` en repositories.

No uses `derrx.New(code, body, nil)` para errores HTTP o de dominio que necesitan un status distinto de `500`.

No uses `perrx` fuera del transporte HTTP.

No devuelvas `err.Error()` directo al cliente.

No mezcles `status code` como parte del contrato de `service`.

No hagas que el handler conozca todos los detalles internos de la causa real.

## Regla práctica para el repo

Si estás trabajando en este proyecto, piensa así:

- `aerrx` clasifica
- `derrx` comunica
- `perrx` presenta
- `httpx` conecta todo eso con HTTP

Y toma esta decisión rápida:

- si solo necesitas clasificar, usa `aerrx`
- si además necesitas `code/body` específico, usa `derrx` tipado
- si no tienes causa real, pasa `nil`

Si una sola capa intenta hacer las cuatro cosas al mismo tiempo, el diseño se vuelve confuso y difícil de mantener.

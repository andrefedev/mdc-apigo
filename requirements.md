# requirements.md

## Reglas de arquitectura

Estas reglas son obligatorias para cambios nuevos.

## 1. Transporte

- Todo codigo de transporte nuevo debe vivir en `api/*`.
- Esto incluye gRPC, webhooks HTTP y futuros adapters externos.
- `api/*` puede parsear requests, leer metadata, aplicar auth del transporte y mapear errores publicos.

## 2. Protobuf

- El contrato gRPC del proyecto se concentra en un unico `ApiService`.
- No se deben crear multiples services protobuf sin una necesidad tecnica explicita.
- Los mensajes se separan por responsabilidad en `api.proto`, `data.proto` y `domain.proto`.

## 3. Features

- La logica interna debe vivir en `internal/features/*`.
- Cada feature puede contener dominio, services, repositorios, DTOs internos y errores semanticos.
- Las features no deben recibir codigo nuevo de transporte.

## 4. Auth

- La logica de sesion y autenticacion pertenece a `internal/features/auth`.
- El transporte solo resuelve metadata, bearer token y autorizacion por metodo.
- La expiracion, revocacion y semantica de auth no deben duplicarse en el borde.

## 5. Wiring

- `cmd/server/main.go` es el punto principal de composicion.
- El bootstrap debe seguir siendo explicito.

## 6. Errores

- Los errores publicos del transporte se resuelven en el borde.
- Los services internos retornan errores del dominio o tecnicos con contexto.
- No se debe introducir `status.Error(...)` como convencion por defecto dentro de features nuevas.

## 7. Webhooks

- Los webhooks se consideran transporte.
- No deben usarse para reintroducir handlers dentro de features.

## 8. No regresion

No se debe volver a una estructura donde:

- cada feature define su propio server gRPC
- el contrato protobuf se fragmenta por inercia
- handlers y adapters vuelven a vivir dentro de `internal/features/*`

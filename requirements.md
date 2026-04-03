# requirements.md

## Reglas de arquitectura

Estas reglas son obligatorias para cambios nuevos.

## 1. Nucleo unico

- La logica interna del sistema vive en `internal/app`.
- No se deben introducir nuevas carpetas `internal/features/*` como direccion arquitectonica principal.
- Si un dominio nuevo aparece, primero debe evaluarse si realmente necesita salir de `internal/app`.

## 2. Transporte

- Todo codigo de transporte nuevo debe vivir en `api/*`.
- Hoy el borde activo es `api/okgrpc`.
- Si aparece webhook/HTTP, debe vivir tambien en `api/*`.

## 3. Protobuf

- El contrato gRPC del proyecto se concentra en un unico `ApiService`.
- No se deben crear multiples services protobuf sin necesidad tecnica explicita.
- Los mensajes siguen separados por responsabilidad en `api.proto`, `data.proto` y `domain.proto`.

## 4. Use service

- La orquestacion principal de casos de uso vive en `internal/app/useservice.go`.
- El borde no debe absorber logica de negocio compleja.
- El repositorio no debe absorber reglas de caso de uso.

## 5. Data flow interno

- El flujo recomendado es `input del transporte -> input interno -> data interna -> repository`.
- `datax.go` representa inputs de entrada.
- `data.go` representa data interna del caso de uso o persistencia.
- La duplicacion entre `input` y `data` es aceptable cuando mejora claridad o desacopla el transporte del repositorio.

## 6. Auth

- La semantica de sesion y auth pertenece a `internal/app`.
- El transporte solo resuelve bearer token, sesion y autorizacion por metodo.
- La expiracion, revocacion y errores semanticos de auth no deben duplicarse fuera del nucleo.

## 7. Errores

- Los errores publicos del transporte se resuelven en `api/*`.
- `internal/app` retorna errores del dominio o tecnicos con contexto.
- No se debe usar `status.Error(...)` como convencion por defecto dentro de `internal/app`.

## 8. Wiring

- `cmd/server/main.go` es el punto principal de composicion.
- El bootstrap debe seguir siendo explicito y legible.

## 9. Infraestructura

- Clientes externos y acceso a DB compartido viven en `internal/modules/*`.
- Utilidades transversales viven en `internal/platforms/*`.

## 10. No regresion

No se debe volver a una estructura donde:

- el nucleo se fragmenta por feature sin una necesidad real
- cada dominio define su propio server gRPC
- el contrato protobuf se fragmenta por inercia
- handlers y adapters vuelven a dispersarse dentro del nucleo

# plan.md

## Objetivo

Reducir riesgo operativo y deuda tecnica antes de seguir expandiendo features sobre una base que hoy funciona, pero tiene contratos inestables entre auth, webhooks, modelo de datos y documentacion.

## Estado resumido

- El bootstrap principal esta ordenado y explicitamente cableado
- La separacion `handler -> service -> repository` ya existe
- El mayor problema no es de forma, sino de consistencia y seguridad

## Prioridad P0

### 1. Cerrar fugas de seguridad y contratos equivocados

- Eliminar logs de `Authorization`, payloads de WhatsApp y cuerpos completos de webhook
- Cambiar la semantica de token invalido para que auth falle con `401`, no con `404`
- Definir verificacion real del webhook de WhatsApp: token de verificacion y firma

Por que:

- Hoy hay fuga potencial de secretos y comportamiento HTTP incorrecto para autenticacion

## Prioridad P1

### 2. Unificar el modelo de datos

- Elegir una convencion canonica para telefono: `phone` o `lookups`, pero no ambas
- Elegir una convencion canonica para identidad: `idk` vs nombre mas expresivo
- Revisar `users.Repository.SelectByPhone` y todas las estructuras `User`, `UserRaw`, `Code`
- Alinear JSON, tags `db`, SQL y documentacion

Por que:

- Hay drift entre consultas, modelos y schema informal; eso ya es deuda estructural

### 3. Formalizar persistencia

- Sacar el esquema real de `pgdb.sql`
- Crear migraciones versionadas
- Documentar tablas canonicas activas y restricciones para `users` y `users_codes`

Por que:

- Hoy no existe una fuente confiable del schema

## Prioridad P2

### 4. Endurecer el flujo de OTP

- Reactivar cleanup del OTP si falla WhatsApp
- Definir TTL, verificacion y consumo del OTP como flujo completo
- Evaluar transaccion u outbox si el envio y la persistencia deben quedar consistentes

Por que:

- El flujo actual crea OTP y envia mensaje, pero el contrato de consistencia esta incompleto

### 5. Separar feature HTTP de webhook del modulo de integracion

- Decidir si `internal/features/app` sera realmente la feature de webhook
- Si no, mover webhook a una feature o adaptar `internal/modules/whatsapp/webhooks`
- Eliminar paquetes muertos o conectarlos al bootstrap

Por que:

- Hoy hay codigo activo y codigo muerto resolviendo el mismo dominio

## Prioridad P3

### 6. Mejorar testabilidad y mantenimiento

- Introducir tests de handler para `/v1/auth/code`, `/v1/users/me`, `/healthz`, `/readyz`
- Introducir tests de servicio para auth y users
- Reemplazar dependencias concretas por interfaces donde haga falta para pruebas
- Limpiar codigo comentado, placeholders y dependencias no usadas en `go.mod`

Por que:

- La arquitectura es razonable, pero todavia no esta defendida por tests ni por contratos estables

## Verificacion recomendada

1. Hacer que `go` quede disponible en `PATH` del entorno local o documentar el wrapper oficial
2. Ejecutar `go test ./...` con red o con mod cache precargado
3. Agregar CI minima de compilacion y tests
4. Validar manualmente:
   `/healthz`
   `/readyz`
   `POST /v1/auth/code`
   `GET /v1/users/me`
   `GET/POST /v1/app/webhook`

## Criterio de exito

- Auth sin fuga de secretos
- Contratos HTTP coherentes
- Modelo de datos consistente
- Schema versionado
- Webhooks validados
- Tests basicos cubriendo rutas criticas

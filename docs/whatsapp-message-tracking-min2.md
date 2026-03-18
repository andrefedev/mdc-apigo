# WhatsApp Message Tracking Min 2

## Objetivo

Guardar solo lo necesario para:

- listar mensajes en UI
- saber su estado actual
- conservar el historial de estados importantes

## Qué guardar

### Tabla principal

`whatsapp_messages_min`

Una fila por mensaje aceptado por la API.

Campos mínimos:

- `id`
- `user_ref`
- `wa_id`
- `to_phone`
- `message_id`
- `category`
- `template_name`
- `status`
- `created_at`
- `updated_at`

Uso:

- listar mensajes
- filtrar por categoría
- ver estado actual

### Tabla de eventos

`whatsapp_message_events_min`

Una fila por cambio de estado.

Campos mínimos:

- `id`
- `message_id`
- `event_type`
- `status`
- `created_at`

Uso:

- saber si pasó por `sent`
- saber si fue `delivered`
- saber si fue `read`
- saber si falló

## Estados recomendados

- `accepted`
- `sent`
- `delivered`
- `read`
- `failed`

## Flujo

### POST a Meta exitoso

Insertar en `whatsapp_messages_min`:

- `message_id`
- `wa_id`
- `to_phone`
- `category`
- `template_name`
- `status = accepted`

Opcionalmente insertar evento `accepted`.

### Webhook de estado

Por cada webhook:

1. resolver por `message_id`
2. actualizar `whatsapp_messages_min.status`
3. insertar evento en `whatsapp_message_events_min`

## Integración con el repo actual

### En `messages.Service` o `Service2`

Después de un `POST` exitoso:

- guardar el mensaje principal

### En webhook

Cuando llegue `sent`, `delivered`, `read` o `failed`:

- actualizar el estado agregado
- insertar evento

## Cuándo crecer esto

Solo si luego necesitas:

- payloads crudos
- auditoría completa
- reintentos
- análisis de campañas

Por ahora no hace falta.

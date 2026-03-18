# WhatsApp Message Ledger 2

## Objetivo

Diseñar un registro confiable de mensajes salientes de WhatsApp que permita:

- Saber qué se intentó enviar, aunque el POST falle
- Guardar `wamid`, contacto resuelto y metadata útil
- Consolidar el estado más reciente para una UI de listado
- Conservar el historial completo de eventos vía webhook
- Relacionar mensajes con `auth`, campañas y otros flujos futuros

## Principio

La respuesta síncrona del POST y los webhooks son dos fuentes parciales de verdad.

El sistema local debe unificarlas así:

1. Crear un registro local antes del envío
2. Actualizarlo con el resultado del POST
3. Completarlo y corregirlo con los webhooks de estado

## Modelo propuesto

### Tabla agregada

`whatsapp_messages`

Sirve para:

- listados
- filtros
- detalle resumido
- estado actual

Debe guardar:

- referencia local
- categoría de negocio
- template
- destinatario
- `wamid`
- estado actual
- timestamps relevantes
- errores relevantes
- contexto útil para UI y auditoría

### Tabla de eventos

`whatsapp_message_events`

Sirve para:

- timeline del mensaje
- debugging
- auditoría
- re-procesamiento

Debe guardar cada transición o evento observado:

- creación local
- aceptación por API
- error de envío
- cambios de estado por webhook

## Estados

Estados locales recomendados:

- `queued`
- `accepted`
- `sent`
- `delivered`
- `read`
- `failed`

Notas:

- `queued`: el sistema decidió enviar y creó el registro
- `accepted`: Meta aceptó la solicitud y devolvió `wamid`
- `sent`/`delivered`/`read`: vienen de webhooks
- `failed`: puede venir del POST o del webhook

## Categorías

Usar una categoría explícita de negocio:

- `auth`
- `utility`
- `marketing`

No depender solo del template name para la UI.

## Flujo

### 1. Antes del POST

Insertar `whatsapp_messages` con:

- `status = queued`
- `category`
- `template_name`
- `template_language`
- `to_phone`
- `user_ref` si aplica
- `context_json`
- `request_payload_json`

Insertar evento `queued`.

### 2. POST a Meta

Si falla:

- marcar `failed`
- guardar error code/message/body
- insertar evento `send_error`

Si Meta responde bien:

- guardar `provider_message_id`
- guardar `provider_contact_wa_id`
- marcar `accepted`
- guardar `response_payload_json`
- insertar evento `api_accept`

### 3. Webhooks

Por cada webhook de estado:

- resolver por `provider_message_id`
- insertar evento crudo
- actualizar estado agregado
- actualizar timestamps:
  - `sent_at`
  - `delivered_at`
  - `read_at`
  - `failed_at`

## Campos útiles para UI

En la tabla agregada:

- `category`
- `status`
- `template_name`
- `template_language`
- `to_phone`
- `provider_contact_wa_id`
- `provider_message_id`
- `provider_error_code`
- `provider_error_message`
- `created_at`
- `accepted_at`
- `delivered_at`
- `read_at`
- `failed_at`
- `context_json`

En la UI puedes listar:

- destinatario
- categoría
- template
- estado actual
- fecha de creación
- fecha de entrega/lectura
- error

Y en el detalle mostrar:

- timeline de eventos
- payload de request sanitizado
- payload/respuesta de webhook

## Contexto recomendado

Guardar `context_json` con estructura libre pero consistente.

Ejemplos:

```json
{"feature":"auth","code_ref":"...","user_ref":"..."}
```

```json
{"feature":"campaigns","campaign_ref":"...","segment_ref":"..."}
```

## Reglas importantes

- `provider_message_id` debe ser `UNIQUE` cuando exista
- Los webhooks nunca deben crear un segundo mensaje si ya existe uno por `wamid`
- El POST nunca debe ser la única fuente de estado final
- El historial crudo nunca debe reemplazar la tabla agregada

## Integración sugerida con el repo actual

### `auth.Service.Code`

1. Crea OTP
2. Inserta ledger `queued`
3. Llama a `messages.Service`
4. Si acepta, guarda `accepted`
5. Si falla, marca `failed`

### webhook HTTP

1. Recibe webhook
2. Normaliza eventos de estado
3. Resuelve por `provider_message_id`
4. Inserta evento
5. Actualiza tabla agregada

## Próximo paso de implementación

1. Crear tablas
2. Crear repositorio de ledger
3. Crear servicio de tracking
4. Integrarlo a `auth`
5. Integrarlo a webhook

    # Muy del Campo - Backend (apigo)

    Este es el backend oficial de **Muy del Campo**, un e-commerce especializado en la venta de huevos y pollo campesino. El sistema está diseñado para ser escalable, seguro y altamente integrado con servicios de mensajería (WhatsApp) y geolocalización.

    ## 🥚 Dominio de Negocio
    - **Productos Principales:** Huevos de campo, Pollo campesino (fresco/procesado).
    - **Logística:** Gestión de rutas de entrega basada en Google Maps. Con limitaciones,
    - **Comunicación:** Notificaciones y flujos de compra vía WhatsApp Business API.
    - **IA:** Integración con Gemini para soporte, creación de pedidos, si un usuario envía un texto para comprar o un audio debemos poder crear y/o consultar pedidos y optimización de catálogo.

    ## 🏗️ Arquitectura y Estructura
    El proyecto sigue una arquitectura limpia (Clean Architecture) simplificada, organizada por **features** y **modules**.

    ### Directorios Clave:
    - `/cmd/server`: Punto de entrada de la aplicación.
    - `/internal/features`: Contiene la lógica de negocio agrupada por dominio (auth, user, orders, products). Cada feature debe seguir el patrón: `handler -> service -> repository`.
    - `/internal/modules`: Adaptadores para servicios externos (Postgres, WhatsApp, Google Cloud Storage/PubSub, Gemini).
    - `/internal/platform`: Código base reutilizable (HTTP router, validación, errores personalizados, criptografía).

    ## 🛠️ Stack Tecnológico
    - **Lenguaje:** Go 1.26+
    - **Base de Datos:** PostgreSQL (usando `pgx/v5` para pooling y transacciones).
    - **Routing:** `go-chi/chi/v5`.
    - **Autenticación:** Token opaque.
    - **Nube:** Google Cloud (Storage, Pub/Sub, Cloud Run). CloudBuild
    - **Integraciones:** WhatsApp Business API, Google Maps API, Twilio.

    ## 📏 Reglas de Oro para Desarrollo (AI Instructions)

    1. **Gestión de Errores:** Nunca ignores un error. Usa el paquete `internal/platform/aerr` para envolver errores con contexto y códigos de estado HTTP apropiados.
    2. **Inyección de Dependencias:** Usa structs de `Deps` (ej. `ServiceDeps`, `HandlerDeps`) para pasar dependencias. No uses estados globales para repositorios o servicios.
    3. **Validación:** Toda entrada desde el `Handler` debe ser validada usando el paquete `internal/platform/validator` antes de pasar al `Service`.
    4. **Nomenclatura:**
       - Repositorios: `repository.go`
       - Servicios: `useservice.go` (Use Case Service)
       - Handlers: `handler.go`
    5. **Base de Datos:** Prefiere el uso de `pgxpool.Pool`. Todas las consultas SQL deben estar en la capa de `repository.go`. Tener en cuenta el Wrapper de WithTx para transacciones.
    6. **Consistencia de Idioma:** El código (variables, funciones, comentarios) debe estar en **Inglés**, pero los mensajes de error orientados al usuario final pueden soportar localización (predeterminado: Español).
    7. **Manejo de errores:** Tengo los errores en aerr.aerrx, aerr.derrx, aerr.perrx, la idea es tener errores internos, por capas y luego exponerlos de manera estructurada a la hora de exponer la API.

    ## 🚀 Comandos Útiles
    - `make build`: Compila el binario del servidor.
    - `go run cmd/server/main.go`: Levanta el servidor localmente.
    - `docker build -t apigo .`: Construye la imagen de contenedor para despliegue.

    ---
    *Documento generado para guiar la evolución de Muy del Campo. Mantén este archivo actualizado ante cambios arquitectónicos.*

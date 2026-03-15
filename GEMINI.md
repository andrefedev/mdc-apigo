# Muy del Campo - Backend (apigo)

Este es el backend oficial de **Muy del Campo**, un e-commerce especializado en la venta de huevos y pollo campesino. El sistema está diseñado para ser escalable, seguro y altamente integrado con servicios de mensajería (WhatsApp) y geolocalización.

## 🥚 Dominio de Negocio
- **Productos Principales:** Huevos de campo, Pollo campesino (fresco/procesado).
- **Logística:** Gestión de rutas de entrega basada en Google Maps.
- **Comunicación:** Notificaciones y flujos de compra vía WhatsApp Business API y Twilio.
- **IA:** Integración con Google Gemini para soporte, creación de pedidos conversacional (texto y audio), y optimización de catálogo.

## 🏗️ Arquitectura y Estructura
El proyecto sigue una arquitectura limpia (Clean Architecture) simplificada, organizada por **features**, **modules** y **platforms**.

### Directorios Clave:
- `/cmd/server`: Punto de entrada de la aplicación (`main.go`). Se encarga de la inyección de dependencias y el ciclo de vida del servidor.
- `/internal/features`: Contiene la lógica de negocio agrupada por dominio (ej. `auth`, `users`). Cada feature implementa el flujo: `handler -> service -> repository`. Incluye definiciones de dominio (`domain.go`), estructuras de datos (`data.go`), y mapeos (`mappers.go`).
- `/internal/modules`: Adaptadores concretos para servicios externos (`gemini`, `gmaps`, `gpubsub`, `gstorage`, `postgres`, `whatsapp`).
- `/internal/platforms`: Código base transversal y reutilizable (`aerr` para errores, `confx` para configuración, `crypx` para criptografía, `httpx` para ruteo/servidor, `loggex` para logs, `validator` para normalización y validación).
- `/docs`: Documentación técnica del proyecto.

## 🛠️ Stack Tecnológico
- **Lenguaje:** Go 1.26
- **Base de Datos:** PostgreSQL (usando `pgx/v5` para pooling y transacciones).
- **Routing:** `go-chi/chi/v5` envolviendo el servidor HTTP en `/internal/platforms/httpx`.
- **Autenticación:** Token opaque.
- **Nube:** Google Cloud (Cloud Run, Cloud SQL, Storage, Pub/Sub). CI/CD con Google Cloud Build (`cloudbuild.yaml`).
- **Integraciones Externas:** WhatsApp Business API, Google Maps API, Twilio.
- **Logging:** Estándar `log/slog` (configurado en `loggex`).

## 📏 Reglas de Oro para Desarrollo (AI Instructions)

1. **Gestión de Errores:** **NUNCA ignores un error**. Usa el paquete `internal/platforms/aerr` para envolver errores con contexto y códigos de estado HTTP apropiados (`aerrx`, `derrx`, `perrx`). Documentación adicional en `docs/aerrx-derrx-perrx.md`.
2. **Inyección de Dependencias:** Usa structs estandarizados de tipo `Deps` (ej. `ServiceDeps`, `HandlerDeps`, `AppRouterDeps`) para la inyección explícita, como se ve en `main.go`. Evita a toda costa los estados globales.
3. **Validación:** Toda entrada o request capturado en el `Handler` debe ser validado y/o normalizado usando `internal/platforms/validator` antes de invocar los métodos del `Service`.
4. **Nomenclatura de Archivos en Features:**
   - **Repositorios:** `repository.go`
   - **Servicios:** `useservice.go` (Use Case Service)
   - **Handlers:** `handler.go`
   - **Modelos:** `domain.go` (entidades de negocio), `data.go` (modelos DTO/DB), `datax.go` (extensiones).
   - **Middlewares/Contexto:** `middlex.go`, `ctx.go`.
5. **Base de Datos:** Usa siempre el `pgx/v5` pool inyectado. Las consultas y la lógica SQL pertenecen estrictamente a la capa de `repository.go`. Utiliza transacciones de manera segura si la plataforma provee un wrapper o se pasa el contexto de la transacción.
6. **Consistencia de Idioma:** El código fuente (nombres de variables, funciones, paquetes, comentarios internos) debe estar **en Inglés**. Los mensajes de error o comunicación orientada al cliente/usuario final (ej. respuestas de WhatsApp) pueden soportar localización (predeterminado: Español).
7. **Logging Estructurado:** Utiliza `slog.Info`, `slog.Error`, etc., inyectando atributos clave como `err` u otros datos de contexto, para facilitar la observabilidad en GCP.

## 🚀 Construcción y Ejecución
- **Ejecución Local:** Asegúrate de tener el archivo `.env` configurado. Usa `make run` para cargar las variables del `.env` y ejecutar el servidor, o alternativamente `go run cmd/server/main.go`.
- **Compilación:** `go build -mod=readonly -o server ./cmd/server`
- **Docker:** Usa el comando `docker build -t mdc-apigo .` para construir la imagen multi-stage definida en el `Dockerfile`.
- **Despliegue:** Se realiza de manera automatizada hacia Google Cloud Run utilizando Google Cloud Build (`cloudbuild.yaml`).

---
*Documento generado y mantenido por Gemini CLI para guiar la evolución arquitectónica de Muy del Campo.*
# CRUD Clientes - Arquitectura Hexagonal (Go)

Sistema CRUD de clientes con:

- **Backend**: Go (net/http)
- **Frontend**: HTML + CSS + JavaScript vanilla
- **Base de datos**: PostgreSQL (Docker)
- **Arquitectura**: Hexagonal (Ports & Adapters)

## Estructura del proyecto

```
CRUDCLIENTESHEX/
├── cmd/api/main.go                 # Entry point
├── internal/
│   ├── domain/                     # Núcleo (entidades + puertos)
│   │   ├── client.go
│   │   └── ports.go
│   ├── application/                # Casos de uso
│   │   └── client_service.go
│   ├── adapters/
│   │   ├── http/                   # Driving adapter (HTTP)
│   │   │   └── client_handler.go
│   │   └── repository/             # Driven adapter (Postgres)
│   │       └── postgres_client_repo.go
│   └── infrastructure/
│       └── database.go
├── web/                            # Frontend estático
│   ├── index.html
│   ├── styles.css
│   └── app.js
├── docker-compose.yml
├── go.mod
└── README.md
```

## Requisitos

- Go 1.22+
- Docker + Docker Compose

## 1. Levantar PostgreSQL

```bash
docker compose up -d
```

Esto crea el contenedor `crudclientesgrok-db` con:

- Usuario: `postgres`
- Password: `norimorienair4614`
- Base de datos: `crudclientesgrok`
- Puerto: `5432`

> **Nota**: El contenedor anterior que tenías (`db-1`) está en estado `exited` por un problema de red. Este `docker-compose` crea uno limpio y dedicado para este proyecto.

## 2. Configurar variables de entorno (opcional)

```bash
cp .env.example .env
```

## 3. Instalar dependencias y ejecutar

```bash
go mod tidy
go run ./cmd/api
```

El servidor arranca en: **http://localhost:8080**

- Frontend: http://localhost:8080/
- API:     http://localhost:8080/api/clients

## Endpoints API

| Método | Ruta                  | Descripción          |
|--------|-----------------------|----------------------|
| GET    | /api/clients          | Listar todos         |
| POST   | /api/clients          | Crear cliente        |
| GET    | /api/clients/{id}     | Obtener por ID       |
| PUT    | /api/clients/{id}     | Actualizar           |
| DELETE | /api/clients/{id}     | Eliminar             |

### Ejemplo de body (POST / PUT)

```json
{
  "nombre": "María González",
  "email": "maria@ejemplo.com",
  "telefono": "+54 11 5555-1234",
  "direccion": "Calle Falsa 123"
}
```

## Cómo copiar a tu ruta de Windows

Como este entorno es Linux, el proyecto se generó en:

```
/home/workdir/artifacts/CRUDCLIENTESHEX
```

Copia todo el contenido a tu ruta deseada:

```
c:\code\GROKPROJECTS\GO\CRUDCLIENTESHEX
```

## Arquitectura Hexagonal - Resumen

| Capa              | Responsabilidad                              |
|-------------------|----------------------------------------------|
| **Domain**        | Entidad `Client` + interfaz `ClientRepository` (puerto) |
| **Application**   | Casos de uso (`ClientService`)               |
| **Adapters**      | HTTP handlers + Postgres repository          |
| **Infrastructure**| Conexión a base de datos                     |

El dominio **no depende** de HTTP ni de Postgres. Solo conoce las interfaces (puertos).

## Campos del Cliente

- `id` (UUID)
- `nombre` (obligatorio)
- `email` (obligatorio)
- `telefono`
- `direccion`
- `created_at` / `updated_at`

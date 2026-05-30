# Lista de tareas (Go + SQLite)

App web de tareas (to-do) que **guarda los datos en una base de datos SQLite**,
así que persisten al cerrar y reabrir. Hecha en **Go**. Demuestra el **CRUD**
completo (Crear, Leer, Actualizar, Borrar).

## Uso

```bash
go run .
```

Abre **http://localhost:8080**. Añade tareas, márcalas como hechas (⬜ → ✅) y
bórralas (🗑️). Cierra y vuelve a abrir: **tus tareas siguen ahí**.

Los datos se guardan en `tasks.db` (configurable con la variable `DB_PATH`).

## Cómo funciona

- `internal/store`: usa `database/sql` con el driver **modernc.org/sqlite**
  (SQLite en Go puro, sin instalar nada). Crea la tabla `tasks` y ofrece
  `Add`, `List`, `Toggle` y `Delete`. Las consultas usan parámetros (`?`) para
  evitar inyección SQL.
- `internal/web`: la lista (`GET /`) y las acciones (`POST /add`,
  `POST /toggle/{id}`, `POST /delete/{id}`).

## Estructura

```
main.go                 arranque (abre la BD, sirve)
internal/store/         base de datos SQLite (CRUD) + pruebas
internal/web/           lista, añadir, marcar y borrar (handlers + plantilla)
```

## Pruebas

```bash
go test ./...
```

La prueba hace el **CRUD completo** con una base de datos temporal: añadir,
listar, marcar como hecha y borrar.

## Stack

Go (net/http, database/sql, html/template, go:embed) ·
SQLite (modernc.org/sqlite, 100% Go).

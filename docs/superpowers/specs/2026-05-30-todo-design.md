# Diseño: Lista de tareas con base de datos (Go + SQLite)

**Fecha:** 2026-05-30 · **Estado:** Aprobado · **Autor:** Samuel (15º proyecto)

## Objetivo

Una app web de tareas (to-do) que **guarda los datos en una base de datos
(SQLite)**, así que persisten al cerrar y reabrir. Objetivo de aprendizaje:
persistencia y **CRUD** (crear, leer, actualizar, borrar) con `database/sql`.

## Pantalla y rutas

- **Inicio** — `GET /`: formulario para añadir + lista de tareas (pendientes y
  hechas), cada una con botón de marcar/desmarcar y borrar.
- **Añadir** — `POST /add`: crea una tarea.
- **Marcar** — `POST /toggle/{id}`: cambia hecha ↔ no hecha.
- **Borrar** — `POST /delete/{id}`: elimina la tarea.

## Arquitectura

```
todo-go/
  main.go                 arranca el servidor; ruta de la BD (DB_PATH)
  internal/store/
    store.go              SQLite: Task; New (crea tabla); Add, List, Toggle, Delete
    store_test.go         CRUD con una BD temporal: añadir → listar → marcar → borrar
  internal/web/
    handler.go            GET / + POST /add + POST /toggle/{id} + POST /delete/{id}
    templates/index.html
    static/style.css
  README.md
```

- **store.go:** `Task{ ID int; Title string; Done bool }`. `Store{ db *sql.DB }`.
  `New(path)` abre SQLite (driver `modernc.org/sqlite`, Go puro) y crea la tabla
  `tasks(id, title, done)` si no existe. Métodos `Add(title)`, `List()` (pendientes
  primero, luego por id desc), `Toggle(id)` (`done = 1 - done`), `Delete(id)`, `Close()`.
- **handler.go:** los handlers usan el store y redirigen a `/` (patrón PRG).

## Pruebas

- **store_test.go:** con una BD en `t.TempDir()`: `Add` → `List` da 1 tarea
  pendiente; `Toggle` la marca hecha; `Delete` la quita.
- `go build/vet/test` limpios.

## Seguridad

Consultas con **parámetros** (`?`) → previene inyección SQL. HTML escapado con
`html/template`.

## Fuera de alcance (YAGNI)

Cuentas de usuario, fechas/recordatorios, categorías, editar el texto.

## Criterios de éxito

1. `go run .` sirve en http://localhost:8080.
2. Añadir, marcar y borrar tareas funciona; **persisten** al reiniciar.
3. La prueba del CRUD pasa.

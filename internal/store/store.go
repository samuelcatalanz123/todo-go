// Package store guarda las tareas en una base de datos SQLite.
package store

import (
	"database/sql"

	_ "modernc.org/sqlite" // driver de SQLite en Go puro (sin instalar nada)
)

// Task es una tarea.
type Task struct {
	ID    int
	Title string
	Done  bool
}

// Store guarda las tareas en SQLite.
type Store struct {
	db *sql.DB
}

// New abre (o crea) la base de datos y prepara la tabla.
func New(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS tasks (
		id    INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		done  INTEGER NOT NULL DEFAULT 0
	)`); err != nil {
		return nil, err
	}
	return &Store{db: db}, nil
}

// Close cierra la base de datos.
func (s *Store) Close() error { return s.db.Close() }

// Add crea una tarea nueva (sin marcar).
func (s *Store) Add(title string) error {
	_, err := s.db.Exec("INSERT INTO tasks (title) VALUES (?)", title)
	return err
}

// List devuelve las tareas: primero las pendientes, luego las hechas; las más
// nuevas arriba.
func (s *Store) List() ([]Task, error) {
	rows, err := s.db.Query("SELECT id, title, done FROM tasks ORDER BY done ASC, id DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		var done int
		if err := rows.Scan(&t.ID, &t.Title, &done); err != nil {
			return nil, err
		}
		t.Done = done == 1
		tasks = append(tasks, t)
	}
	return tasks, rows.Err()
}

// Toggle cambia una tarea de hecha a no hecha (y viceversa).
func (s *Store) Toggle(id int) error {
	_, err := s.db.Exec("UPDATE tasks SET done = 1 - done WHERE id = ?", id)
	return err
}

// Delete borra una tarea.
func (s *Store) Delete(id int) error {
	_, err := s.db.Exec("DELETE FROM tasks WHERE id = ?", id)
	return err
}

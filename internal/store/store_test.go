package store

import (
	"path/filepath"
	"testing"
)

func newTestStore(t *testing.T) *Store {
	t.Helper()
	s, err := New(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	t.Cleanup(func() { _ = s.Close() })
	return s
}

func TestCRUD(t *testing.T) {
	s := newTestStore(t)

	// Crear
	if err := s.Add("Comprar pan"); err != nil {
		t.Fatalf("Add: %v", err)
	}

	// Leer
	tasks, err := s.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(tasks) != 1 {
		t.Fatalf("esperaba 1 tarea, obtuve %d", len(tasks))
	}
	if tasks[0].Title != "Comprar pan" || tasks[0].Done {
		t.Fatalf("tarea mal: %+v", tasks[0])
	}
	id := tasks[0].ID

	// Actualizar (marcar hecha)
	if err := s.Toggle(id); err != nil {
		t.Fatalf("Toggle: %v", err)
	}
	tasks, _ = s.List()
	if !tasks[0].Done {
		t.Error("la tarea debería estar marcada como hecha")
	}

	// Borrar
	if err := s.Delete(id); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	tasks, _ = s.List()
	if len(tasks) != 0 {
		t.Errorf("esperaba 0 tareas tras borrar, obtuve %d", len(tasks))
	}
}

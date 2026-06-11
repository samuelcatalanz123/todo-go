package web

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/samuelcatalanz123/todo-go/internal/store"
)

// nuevoHandler crea un Handler con una base de datos SQLite temporal
// (en t.TempDir(), que se borra solo al terminar el test).
func nuevoHandler(t *testing.T) *Handler {
	t.Helper()
	s, err := store.New(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	t.Cleanup(func() { _ = s.Close() })
	h, err := New(s)
	if err != nil {
		t.Fatalf("web.New: %v", err)
	}
	return h
}

func post(h *Handler, path string, form url.Values) {
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	h.Routes().ServeHTTP(httptest.NewRecorder(), req)
}

// TestCRUDPorLaWeb prueba el ciclo completo por la web: añadir una tarea,
// marcarla como hecha y borrarla, usando una base de datos temporal.
func TestCRUDPorLaWeb(t *testing.T) {
	h := nuevoHandler(t)

	// Añadir
	post(h, "/add", url.Values{"title": {"Comprar pan"}})
	tasks, _ := h.store.List()
	if len(tasks) != 1 || tasks[0].Title != "Comprar pan" {
		t.Fatalf("la tarea no se añadió: %v", tasks)
	}
	id := tasks[0].ID

	// Marcar como hecha
	post(h, "/toggle/"+strconv.Itoa(id), nil)
	tasks, _ = h.store.List()
	if !tasks[0].Done {
		t.Errorf("la tarea debería estar marcada como hecha")
	}

	// Borrar
	post(h, "/delete/"+strconv.Itoa(id), nil)
	tasks, _ = h.store.List()
	if len(tasks) != 0 {
		t.Errorf("la tarea debería haberse borrado, quedan %d", len(tasks))
	}
}

// TestHome comprueba que la página principal responde 200.
func TestHome(t *testing.T) {
	h := nuevoHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	h.Routes().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("home: código %d, esperaba 200", rec.Code)
	}
}

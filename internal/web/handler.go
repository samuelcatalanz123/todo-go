// Package web sirve la lista de tareas.
package web

import (
	"embed"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/samuelcatalanz123/todo-go/internal/store"
)

//go:embed templates/*.html static/*
var files embed.FS

// Handler sirve la web.
type Handler struct {
	tmpl  *template.Template
	store *store.Store
}

// New crea el Handler con el store dado.
func New(s *store.Store) (*Handler, error) {
	tmpl, err := template.ParseFS(files, "templates/*.html")
	if err != nil {
		return nil, err
	}
	return &Handler{tmpl: tmpl, store: s}, nil
}

// Routes monta las rutas.
func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.FileServerFS(files))
	mux.HandleFunc("GET /{$}", h.home)
	mux.HandleFunc("POST /add", h.add)
	mux.HandleFunc("POST /toggle/{id}", h.toggle)
	mux.HandleFunc("POST /delete/{id}", h.delete)
	return mux
}

func (h *Handler) home(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.store.List()
	if err != nil {
		http.Error(w, "error del servidor", http.StatusInternalServerError)
		return
	}
	if err := h.tmpl.ExecuteTemplate(w, "index.html", tasks); err != nil {
		http.Error(w, "error del servidor", http.StatusInternalServerError)
	}
}

func (h *Handler) add(w http.ResponseWriter, r *http.Request) {
	title := strings.TrimSpace(r.FormValue("title"))
	if title != "" {
		_ = h.store.Add(title)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) toggle(w http.ResponseWriter, r *http.Request) {
	if id, err := strconv.Atoi(r.PathValue("id")); err == nil {
		_ = h.store.Toggle(id)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	if id, err := strconv.Atoi(r.PathValue("id")); err == nil {
		_ = h.store.Delete(id)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

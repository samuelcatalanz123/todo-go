// todo-go: una lista de tareas que se guarda en una base de datos SQLite.
package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/samuelcatalanz123/todo-go/internal/store"
	"github.com/samuelcatalanz123/todo-go/internal/web"
)

func main() {
	dbPath := envOr("DB_PATH", "tasks.db")
	s, err := store.New(dbPath)
	if err != nil {
		log.Fatalf("no se pudo abrir la base de datos: %v", err)
	}
	defer s.Close()

	h, err := web.New(s)
	if err != nil {
		log.Fatalf("no se pudo crear el handler: %v", err)
	}

	addr := ":8080"
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}
	slog.Info("servidor iniciado", "abre", "http://localhost"+addr, "bd", dbPath)
	if err := http.ListenAndServe(addr, h.Routes()); err != nil {
		log.Fatal(err)
	}
}

func envOr(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}

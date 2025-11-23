package main

import (
    "net/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"example.com/todo-backend/objects"
)

func main() {

    r := chi.NewRouter()

    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)

    // CORS middleware
    r.Use(func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Access-Control-Allow-Origin", "*")
            w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
            w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

            if r.Method == http.MethodOptions {
                w.WriteHeader(204)
                return
            }

            next.ServeHTTP(w, r)
        })
    })

    // Routes (same structure as before)
    r.Route("/api", func(api chi.Router) {
        api.Get("/todos", tasks.Get)
        api.Post("/todos", tasks.Create)
        api.Put("/todos/{id}", tasks.Update)     
        api.Delete("/todos/{id}", tasks.Delete)
    })

    http.ListenAndServe(":8080", r)
}

package tasks

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "os"
    "strconv"
    "example.com/todo-backend/objects/users"

    "github.com/go-chi/chi/v5"
    "github.com/redis/go-redis/v9"
)

var (
    rdb *redis.Client
    ctx = context.Background()
)

func init() {
    addr := os.Getenv("REDIS_ADDR")
    if addr == "" {
        addr = "localhost:6379"
    }

    password := os.Getenv("REDIS_PASSWORD") // can be empty

    rdb = redis.NewClient(&redis.Options{
        Addr:     addr,
        Password: password,
        DB:       0,
    })
}

type Todo struct {
    ID    int64  `json:"id"`
    Title string `json:"title"`
    Done  bool   `json:"done"`
}

// Redis keys:
// todo:id        -> counter (INCR)
// todo:<id>      -> JSON of Todo
// todos          -> set of all IDs

// GET /api/todos
func Get(w http.ResponseWriter, r *http.Request) {
    
    log.Println("Current user: %v",users.CurrentUser)

    ids, err := rdb.SMembers(ctx, "todos").Result()
    if err != nil {
        http.Error(w, "redis error: "+err.Error(), http.StatusInternalServerError)
        return
    }

    todos := make([]Todo, 0, len(ids))

    for _, idStr := range ids {
        key := "todo:" + idStr
        data, err := rdb.Get(ctx, key).Bytes()
        if err == redis.Nil {
            continue // key vanished, ignore
        } else if err != nil {
            http.Error(w, "redis error: "+err.Error(), http.StatusInternalServerError)
            return
        }

        var t Todo
        if err := json.Unmarshal(data, &t); err != nil {
            http.Error(w, "json error: "+err.Error(), http.StatusInternalServerError)
            return
        }
        todos = append(todos, t)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(todos)
}

// POST /api/todos
func Create(w http.ResponseWriter, r *http.Request) {
    var input struct {
        Title string `json:"title"`
    }

    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
        return
    }
    if input.Title == "" {
        http.Error(w, "title required", http.StatusBadRequest)
        return
    }

    // generate new ID
    id, err := rdb.Incr(ctx, "todo:id").Result()
    if err != nil {
        http.Error(w, "redis error: "+err.Error(), http.StatusInternalServerError)
        return
    }

    todo := Todo{
        ID:    id,
        Title: input.Title,
        Done:  false,
    }

    data, err := json.Marshal(todo)
    if err != nil {
        http.Error(w, "json error: "+err.Error(), http.StatusInternalServerError)
        return
    }

    key := "todo:" + strconv.FormatInt(id, 10)

    // store todo and add to set
    if err := rdb.Set(ctx, key, data, 0).Err(); err != nil {
        http.Error(w, "redis error: "+err.Error(), http.StatusInternalServerError)
        return
    }
    if err := rdb.SAdd(ctx, "todos", id).Err(); err != nil {
        http.Error(w, "redis error: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(todo)
}

// PUT /api/todos/{id}
func Update(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    if idStr == "" {
        http.Error(w, "missing id", http.StatusBadRequest)
        return
    }

    key := "todo:" + idStr

    data, err := rdb.Get(ctx, key).Bytes()
    if err == redis.Nil {
        http.Error(w, "not found", http.StatusNotFound)
        return
    } else if err != nil {
        http.Error(w, "redis error: "+err.Error(), http.StatusInternalServerError)
        return
    }

    var todo Todo
    if err := json.Unmarshal(data, &todo); err != nil {
        http.Error(w, "json error: "+err.Error(), http.StatusInternalServerError)
        return
    }

    var patch struct {
        Title *string `json:"title"`
        Done  *bool   `json:"done"`
    }

    if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
        http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
        return
    }

    if patch.Title != nil {
        todo.Title = *patch.Title
    }
    if patch.Done != nil {
        todo.Done = *patch.Done
    }

    data, err = json.Marshal(todo)
    if err != nil {
        http.Error(w, "json error: "+err.Error(), http.StatusInternalServerError)
        return
    }

    if err := rdb.Set(ctx, key, data, 0).Err(); err != nil {
        http.Error(w, "redis error: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(todo)
}

// DELETE /api/todos/{id}
func Delete(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    if idStr == "" {
        http.Error(w, "missing id", http.StatusBadRequest)
        return
    }

    key := "todo:" + idStr

    // delete key and remove from set
    if err := rdb.Del(ctx, key).Err(); err != nil {
        http.Error(w, "redis error: "+err.Error(), http.StatusInternalServerError)
        return
    }
    if err := rdb.SRem(ctx, "todos", idStr).Err(); err != nil {
        http.Error(w, "redis error: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

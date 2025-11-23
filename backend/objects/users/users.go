package users

import (
    "context"
    "encoding/json"
    "net/http"
    "os"
    "strconv"
    "strings"

    "github.com/redis/go-redis/v9"
    "golang.org/x/crypto/bcrypt"
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

type User struct {
    ID           int64  `json:"id"`
    Email        string `json:"email"`
    Name         string `json:"name"`
    PasswordHash string `json:"password_hash"`
}

// ---------- helpers ----------

func normalizeEmail(email string) string {
    return strings.TrimSpace(strings.ToLower(email))
}

func userKey(id int64) string {
    return "user:" + strconv.FormatInt(id, 10)
}

func userEmailKey(email string) string {
    return "user:email:" + email
}

// ---------- handlers ----------

// POST /api/register
// body: { "email": "...", "name": "...", "password": "..." }
func Register(w http.ResponseWriter, r *http.Request) {
    var input struct {
        Email    string `json:"email"`
        Name     string `json:"name"`
        Password string `json:"password"`
    }

    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, "invalid JSON", http.StatusBadRequest)
        return
    }

    input.Email = normalizeEmail(input.Email)
    input.Name = strings.TrimSpace(input.Name)

    if input.Email == "" || input.Name == "" || input.Password == "" {
        http.Error(w, "email, name and password required", http.StatusBadRequest)
        return
    }

    // check if email already exists
    if _, err := rdb.Get(ctx, userEmailKey(input.Email)).Result(); err == nil {
        http.Error(w, "email already registered", http.StatusConflict)
        return
    } else if err != redis.Nil && err != nil {
        http.Error(w, "redis error: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // hash password
    hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, "hash error: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // generate new user id
    id, err := rdb.Incr(ctx, "user:id").Result()
    if err != nil {
        http.Error(w, "redis error: "+err.Error(), http.StatusInternalServerError)
        return
    }

    user := User{
        ID:           id,
        Email:        input.Email,
        Name:         input.Name,
        PasswordHash: string(hash),
    }

    data, err := json.Marshal(user)
    if err != nil {
        http.Error(w, "json error: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // store: user:<id> -> JSON, user:email:<email> -> id
    pipe := rdb.TxPipeline()
    pipe.Set(ctx, userKey(id), data, 0)
    pipe.Set(ctx, userEmailKey(user.Email), id, 0)
    if _, err := pipe.Exec(ctx); err != nil {
        http.Error(w, "redis error: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    // do NOT send password_hash to client
    json.NewEncoder(w).Encode(map[string]interface{}{
        "id":    user.ID,
        "email": user.Email,
        "name":  user.Name,
    })
}

// POST /api/login
// body: { "email": "...", "password": "..." }
func Login(w http.ResponseWriter, r *http.Request) {
    var input struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, "invalid JSON", http.StatusBadRequest)
        return
    }

    input.Email = normalizeEmail(input.Email)
    if input.Email == "" || input.Password == "" {
        http.Error(w, "email and password required", http.StatusBadRequest)
        return
    }

    // lookup user id by email
    idStr, err := rdb.Get(ctx, userEmailKey(input.Email)).Result()
    if err == redis.Nil {
        http.Error(w, "invalid credentials", http.StatusUnauthorized)
        return
    } else if err != nil {
        http.Error(w, "redis error: "+err.Error(), http.StatusInternalServerError)
        return
    }

    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        http.Error(w, "data error", http.StatusInternalServerError)
        return
    }

    // load user
    data, err := rdb.Get(ctx, userKey(id)).Bytes()
    if err == redis.Nil {
        http.Error(w, "invalid credentials", http.StatusUnauthorized)
        return
    } else if err != nil {
        http.Error(w, "redis error: "+err.Error(), http.StatusInternalServerError)
        return
    }

    var user User
    if err := json.Unmarshal(data, &user); err != nil {
        http.Error(w, "json error: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // check password
    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
        http.Error(w, "invalid credentials", http.StatusUnauthorized)
        return
    }

    // success – for now just return basic user info
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "status": "ok",
        "user": map[string]interface{}{
            "id":    user.ID,
            "email": user.Email,
            "name":  user.Name,
        },
    })
}

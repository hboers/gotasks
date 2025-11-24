package users

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

// -----------------------------------------------------------------------------
// Redis Setup
// -----------------------------------------------------------------------------

var (
	rdb *redis.Client
	ctx = context.Background()
)

func init() {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	password := os.Getenv("REDIS_PASSWORD") // kann leer sein

	rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})
}

// -----------------------------------------------------------------------------
// Types
// -----------------------------------------------------------------------------

type User struct {
	ID           int64  `json:"id"`
	Email        string `json:"email"`
	Name         string `json:"name"`
	PasswordHash string `json:"password_hash"`
}

// Context-Key für User
type ctxKey int

const userKeyCtx ctxKey = 1

// -----------------------------------------------------------------------------
// Key-Helper
// -----------------------------------------------------------------------------

func normalizeEmail(email string) string {
	return strings.TrimSpace(strings.ToLower(email))
}

func userKey(id int64) string {
	return "user:" + strconv.FormatInt(id, 10)
}

func userEmailKey(email string) string {
	return "user:email:" + email
}

func sessionKey(sid string) string {
	return "session:" + sid
}

// generate random session ID
func newSessionID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// -----------------------------------------------------------------------------
// HTTP Handler: Register
// POST /api/register
// body: { "email": "...", "name": "...", "password": "..." }
// -----------------------------------------------------------------------------

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

	// E-Mail schon vorhanden?
	if _, err := rdb.Get(ctx, userEmailKey(input.Email)).Result(); err == nil {
		http.Error(w, "email already registered", http.StatusConflict)
		return
	} else if err != redis.Nil && err != nil {
		http.Error(w, "redis error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Passwort hashen
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "hash error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// neue User-ID
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

	// user:<id> -> JSON, user:email:<email> -> id
	pipe := rdb.TxPipeline()
	pipe.Set(ctx, userKey(id), data, 0)
	pipe.Set(ctx, userEmailKey(user.Email), id, 0)
	if _, err := pipe.Exec(ctx); err != nil {
		http.Error(w, "redis error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"id":    user.ID,
		"email": user.Email,
		"name":  user.Name,
	})
}

// -----------------------------------------------------------------------------
// HTTP Handler: Login
// POST /api/login
// body: { "email": "...", "password": "..." }
// Legt Session in Redis an und setzt session_id-Cookie
// -----------------------------------------------------------------------------

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

	// user-id via email
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

	// user laden
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

	// Passwort prüfen
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	// Session anlegen
	sid, err := newSessionID()
	if err != nil {
		http.Error(w, "session error", http.StatusInternalServerError)
		return
	}

	if err := rdb.Set(ctx, sessionKey(sid), user.ID, 24*time.Hour).Err(); err != nil {
		http.Error(w, "redis error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Cookie setzen – für localhost:5173 ↔ 8080: SameSite=None, Secure=false
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sid,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,                 // in Produktion: true + HTTPS
		//SameSite: http.SameSiteNoneMode, // wichtig für Vite-Frontend auf anderem Origin
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status": "ok",
		"user": map[string]any{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		},
	})
}

// -----------------------------------------------------------------------------
// CurrentUser: liest session_id-Cookie, lädt User aus Redis
// -----------------------------------------------------------------------------

func CurrentUser(r *http.Request) (*User, error) {
	c, err := r.Cookie("session_id")
	if err != nil {
		return nil, err
	}

	sid := c.Value
	if sid == "" {
		return nil, http.ErrNoCookie
	}

	idStr, err := rdb.Get(ctx, sessionKey(sid)).Result()
	if err != nil {
		return nil, err
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, err
	}

	data, err := rdb.Get(ctx, userKey(id)).Bytes()
	if err != nil {
		return nil, err
	}

	var u User
	if err := json.Unmarshal(data, &u); err != nil {
		return nil, err
	}
	return &u, nil
}

// -----------------------------------------------------------------------------
// Middleware: RequireAuth – schützt Routen
// -----------------------------------------------------------------------------

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, err := CurrentUser(r)
		if err != nil || u == nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), userKeyCtx, u)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// User aus Context auslesen (z.B. in Handlern)

func GetUserFromContext(r *http.Request) *User {
	if v := r.Context().Value(userKeyCtx); v != nil {
		if u, ok := v.(*User); ok {
			return u
		}
	}
	return nil
}

// -----------------------------------------------------------------------------
// Logout: Session löschen + Cookie invalidieren
// -----------------------------------------------------------------------------

func Logout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_id")
	if err == nil && c.Value != "" {
		_ = rdb.Del(ctx, sessionKey(c.Value)).Err()
	}

	// Cookie beim Client invalidieren
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
		// SameSite: http.SameSiteNoneMode,
	})

	w.WriteHeader(http.StatusNoContent)
}

func Me(w http.ResponseWriter, r *http.Request) {
    u, err := CurrentUser(r)
    if err != nil || u == nil {
        http.Error(w, "unauthorized", http.StatusUnauthorized)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]any{
        "id":    u.ID,
        "email": u.Email,
        "name":  u.Name,
    })
}
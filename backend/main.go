package main

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
)

type Todo struct {
    ID    int    `json:"id"`
    Title string `json:"title"`
    Done  bool   `json:"done"`
}

// in-memory data (replace with DB later)
var todos = []Todo{
    {ID: 1, Title: "Learn Go", Done: false},
    {ID: 2, Title: "Build REST API", Done: false},
}

func main() {
    r := gin.Default()

    // CORS – simple/naive version for local dev
    r.Use(func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        if c.Request.Method == http.MethodOptions {
            c.AbortWithStatus(204)
            return
        }
        c.Next()
    })

    api := r.Group("/api")
    {
        api.GET("/todos", getTodos)
        api.POST("/todos", createTodo)
        api.PUT("/todos/:id", updateTodo)
        api.DELETE("/todos/:id", deleteTodo)
    }

    r.Run(":8080") // http://localhost:8080
}

func getTodos(c *gin.Context) {
    c.JSON(http.StatusOK, todos)
}

func createTodo(c *gin.Context) {
    var input struct {
        Title string `json:"title"`
    }
    if err := c.ShouldBindJSON(&input); err != nil || input.Title == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "title required"})
        return
    }

    // simple ID generation
    newID := 1
    if len(todos) > 0 {
        newID = todos[len(todos)-1].ID + 1
    }

    todo := Todo{ID: newID, Title: input.Title, Done: false}
    todos = append(todos, todo)
    c.JSON(http.StatusCreated, todo)
}

func updateTodo(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }

    var input struct {
        Title *string `json:"title"`
        Done  *bool   `json:"done"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
        return
    }

    for i, t := range todos {
        if t.ID == id {
            if input.Title != nil {
                t.Title = *input.Title
            }
            if input.Done != nil {
                t.Done = *input.Done
            }
            todos[i] = t
            c.JSON(http.StatusOK, t)
            return
        }
    }

    c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
}

func deleteTodo(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }

    for i, t := range todos {
        if t.ID == id {
            todos = append(todos[:i], todos[i+1:]...)
            c.Status(http.StatusNoContent)
            return
        }
    }

    c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
}

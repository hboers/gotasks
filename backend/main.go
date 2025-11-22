package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
	"example.com/todo-backend/objects/tasks"
)


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
        api.GET("/todos", tasks.Get)
        api.POST("/todos", tasks.Create)
        api.PUT("/todos/:id", tasks.Update)
        api.DELETE("/todos/:id", tasks.Delete)
    }

    r.Run(":8080") // http://localhost:8080
}


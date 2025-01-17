package main

import (
	"log"

	"1uciuszzz/todo-app-golang/internal/config"
	"1uciuszzz/todo-app-golang/internal/handler"
	"1uciuszzz/todo-app-golang/internal/model"
	"1uciuszzz/todo-app-golang/internal/repository"
	"1uciuszzz/todo-app-golang/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 连接数据库
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 自动迁移数据库表结构
	if err := db.AutoMigrate(&model.Todo{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 初始化依赖
	todoRepo := repository.NewTodoRepository(db)
	todoService := service.NewTodoService(todoRepo)
	todoHandler := handler.NewTodoHandler(todoService)

	// 设置路由
	r := gin.Default()

	// Todo 路由组
	v1 := r.Group("/api/v1")
	{
		todos := v1.Group("/todos")
		{
			todos.POST("/", todoHandler.Create)
			todos.GET("/", todoHandler.List)
			todos.GET("/:id", todoHandler.Get)
			todos.PUT("/:id", todoHandler.Update)
			todos.DELETE("/:id", todoHandler.Delete)
		}
	}

	// 启动服务器
	if err := r.Run(cfg.ServerAddress); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

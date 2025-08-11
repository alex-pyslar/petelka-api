package main

import (
	"net/http"
	"time"

	"github.com/alex-pyslar/online-store/internal/config"
	"github.com/alex-pyslar/online-store/internal/handler"
	"github.com/alex-pyslar/online-store/internal/logger"
	"github.com/alex-pyslar/online-store/internal/repository"
	"github.com/alex-pyslar/online-store/internal/service"

	_ "github.com/alex-pyslar/online-store/docs"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Online Store API
// @version 1.0
// @description REST API for an ecommerce platform
// @host localhost:8080
// @BasePath /api
func main() {
	// Инициализация логгера
	log, err := logger.NewLogger()
	if err != nil {
		panic(err)
	}

	// Инициализация конфигурации и зависимостей
	cfg, err := config.NewConfig(log)
	if err != nil {
		log.Fatalf("Error initializing config: %v", err)
	}
	defer cfg.DB.Close()
	defer cfg.Redis.Close()
	log.Info("Connected to PostgreSQL and Redis successfully")

	// Инициализация репозиториев
	userRepo := repository.NewUserRepository(cfg.DB, cfg.Redis)
	productRepo := repository.NewProductRepository(cfg.DB, cfg.Redis)
	categoryRepo := repository.NewCategoryRepository(cfg.DB, cfg.Redis)
	orderRepo := repository.NewOrderRepository(cfg.DB, cfg.Redis)
	commentRepo := repository.NewCommentRepository(cfg.DB, cfg.Redis)

	// Инициализация сервисов с логгером
	userService := service.NewUserService(userRepo, log)
	productService := service.NewProductService(productRepo, log)
	categoryService := service.NewCategoryService(categoryRepo, log)
	orderService := service.NewOrderService(orderRepo, log)
	commentService := service.NewCommentService(commentRepo, log)

	// Инициализация хендлеров без логгера
	userHandler := handler.NewUserHandler(userService)
	productHandler := handler.NewProductHandler(productService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	orderHandler := handler.NewOrderHandler(orderService)
	commentHandler := handler.NewCommentHandler(commentService)
	authHandler := handler.NewAuthHandler(userService)

	// Настройка роутера
	router := mux.NewRouter()

	// Публичные маршруты (без аутентификации)
	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/auth/register", authHandler.Register).Methods("POST")
	api.HandleFunc("/auth/login", authHandler.Login).Methods("POST")

	// Защищенные маршруты (требуют JWT токен)
	protected := api.PathPrefix("").Subrouter()
	protected.Use(handler.JWTMiddleware(log))

	// Пользователи (защищенные)
	protected.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	protected.HandleFunc("/users/{id}", userHandler.GetUser).Methods("GET")
	protected.HandleFunc("/users", userHandler.ListUsers).Methods("GET")
	protected.HandleFunc("/users/{id}", userHandler.UpdateUser).Methods("PUT")
	protected.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods("DELETE")

	// Товары (защищенные)
	protected.HandleFunc("/products", productHandler.CreateProduct).Methods("POST")
	protected.HandleFunc("/products/{id}", productHandler.GetProduct).Methods("GET")
	protected.HandleFunc("/products", productHandler.ListProducts).Methods("GET")
	protected.HandleFunc("/products/{id}", productHandler.UpdateProduct).Methods("PUT")
	protected.HandleFunc("/products/{id}", productHandler.DeleteProduct).Methods("DELETE")

	// Категории (защищенные)
	protected.HandleFunc("/categories", categoryHandler.CreateCategory).Methods("POST")
	protected.HandleFunc("/categories/{id}", categoryHandler.GetCategory).Methods("GET")
	protected.HandleFunc("/categories", categoryHandler.ListCategories).Methods("GET")
	protected.HandleFunc("/categories/{id}", categoryHandler.UpdateCategory).Methods("PUT")
	protected.HandleFunc("/categories/{id}", categoryHandler.DeleteCategory).Methods("DELETE")

	// Заказы (защищенные)
	protected.HandleFunc("/orders", orderHandler.CreateOrder).Methods("POST")
	protected.HandleFunc("/orders/{id}", orderHandler.GetOrder).Methods("GET")
	protected.HandleFunc("/orders", orderHandler.ListOrders).Methods("GET")
	protected.HandleFunc("/orders/{id}", orderHandler.UpdateOrder).Methods("PUT")
	protected.HandleFunc("/orders/{id}", orderHandler.DeleteOrder).Methods("DELETE")

	// Комментарии (защищенные)
	protected.HandleFunc("/comments", commentHandler.CreateComment).Methods("POST")
	protected.HandleFunc("/comments/{id}", commentHandler.GetComment).Methods("GET")
	protected.HandleFunc("/comments", commentHandler.ListComments).Methods("GET")
	protected.HandleFunc("/comments/{id}", commentHandler.UpdateComment).Methods("PUT")
	protected.HandleFunc("/comments/{id}", commentHandler.DeleteComment).Methods("DELETE")

	// Технические маршруты
	router.Handle("/metrics", promhttp.Handler())
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Логирование старта сервера
	log.Infof("Server starting on :8080 at %s", time.Now().Format("2006-01-02 15:04:05"))
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

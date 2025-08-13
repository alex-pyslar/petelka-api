package main

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/alex-pyslar/online-store/internal/config"
	"github.com/alex-pyslar/online-store/internal/handler"
	"github.com/alex-pyslar/online-store/internal/logger"
	"github.com/alex-pyslar/online-store/internal/repository"
	"github.com/alex-pyslar/online-store/internal/service"

	_ "github.com/alex-pyslar/online-store/docs"
)

func main() {
	log, err := logger.NewLogger()
	if err != nil {
		panic(err)
	}

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

	// Инициализация сервисов
	userService := service.NewUserService(userRepo, log)
	productService := service.NewProductService(productRepo, log)
	categoryService := service.NewCategoryService(categoryRepo, log)
	orderService := service.NewOrderService(orderRepo, log)
	commentService := service.NewCommentService(commentRepo, log)

	// Инициализация хендлеров
	userHandler := handler.NewUserHandler(userService)
	productHandler := handler.NewProductHandler(productService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	orderHandler := handler.NewOrderHandler(orderService)
	commentHandler := handler.NewCommentHandler(commentService)
	authHandler := handler.NewAuthHandler(userService)

	router := mux.NewRouter()
	router.Use(handler.CorsMiddleware) // Применяем CORS ко всем запросам
	api := router.PathPrefix("/api").Subrouter()

	// --- Public маршруты (доступны всем) ---
	public := api.PathPrefix("").Subrouter()
	public.HandleFunc("/auth/register", authHandler.Register).Methods("POST")
	public.HandleFunc("/auth/login", authHandler.Login).Methods("POST")
	public.HandleFunc("/products", productHandler.ListProducts).Methods("GET")
	public.HandleFunc("/products/{id}", productHandler.GetProduct).Methods("GET")
	public.HandleFunc("/products/search", productHandler.SearchProducts).Methods("GET")
	public.HandleFunc("/categories", categoryHandler.ListCategories).Methods("GET")
	public.HandleFunc("/categories/{id}", categoryHandler.GetCategory).Methods("GET")

	// --- Protected маршруты (доступны авторизованным пользователям) ---
	protected := api.PathPrefix("").Subrouter()
	protected.Use(handler.AuthMiddleware(log))                                      // Проверяем авторизацию
	protected.HandleFunc("/comments", commentHandler.CreateComment).Methods("POST") // Создание комментария
	protected.HandleFunc("/orders", orderHandler.CreateOrder).Methods("POST")       // Создание заказа

	// --- Admin маршруты (доступны только администраторам) ---
	admin := api.PathPrefix("").Subrouter()
	admin.Use(handler.AuthMiddleware(log), handler.AdminMiddleware(log)) // Сначала авторизация, потом проверка роли
	admin.HandleFunc("/products", productHandler.CreateProduct).Methods("POST")
	admin.HandleFunc("/products/{id}", productHandler.UpdateProduct).Methods("PUT")
	admin.HandleFunc("/products/{id}", productHandler.DeleteProduct).Methods("DELETE")
	admin.HandleFunc("/categories", categoryHandler.CreateCategory).Methods("POST")
	admin.HandleFunc("/categories/{id}", categoryHandler.UpdateCategory).Methods("PUT")
	admin.HandleFunc("/categories/{id}", categoryHandler.DeleteCategory).Methods("DELETE")
	admin.HandleFunc("/users", userHandler.ListUsers).Methods("GET")
	admin.HandleFunc("/users/{id}", userHandler.UpdateUser).Methods("PUT")
	admin.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods("DELETE")

	// Технические маршруты
	router.Handle("/metrics", promhttp.Handler())
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	log.Infof("Server starting on :8080 at %s", time.Now().Format("2006-01-02 15:04:05"))
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

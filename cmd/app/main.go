package main

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/alex-pyslar/petelka-api/internal/config"
	"github.com/alex-pyslar/petelka-api/internal/handler"
	"github.com/alex-pyslar/petelka-api/internal/logger"
	"github.com/alex-pyslar/petelka-api/internal/repository"
	"github.com/alex-pyslar/petelka-api/internal/service"

	_ "github.com/alex-pyslar/petelka-api/docs"
)

func main() {
	// Инициализация логгера
	log, err := logger.NewLogger()
	if err != nil {
		panic(err)
	}

	// Инициализация конфигурации (DB, Redis, MinIO)
	cfg, err := config.NewConfig(log)
	if err != nil {
		log.Fatalf("Error initializing config: %v", err)
	}
	defer cfg.DB.Close()
	defer cfg.Redis.Close()
	log.Info("Connected to PostgreSQL and Redis successfully")

	// === Репозитории ===
	userRepo := repository.NewUserRepository(cfg.DB, cfg.Redis)
	productRepo := repository.NewProductRepository(cfg.DB, cfg.Redis)
	categoryRepo := repository.NewCategoryRepository(cfg.DB, cfg.Redis)
	orderRepo := repository.NewOrderRepository(cfg.DB, cfg.Redis)
	commentRepo := repository.NewCommentRepository(cfg.DB, cfg.Redis)

	// MinIO как репозиторий
	photoRepo, err := repository.NewPhotoRepository(
		cfg.MinioEndpoint,
		cfg.MinioAccessKey,
		cfg.MinioSecretKey,
		cfg.MinioBucket,
		cfg.MinioUseSSL,
		cfg.Redis, // только Redis для кэша
	)
	if err != nil {
		log.Fatalf("Failed to initialize PhotoRepository: %v", err)
	}

	// === Сервисы ===
	userService := service.NewUserService(userRepo, log)
	productService := service.NewProductService(productRepo, log)
	categoryService := service.NewCategoryService(categoryRepo, log)
	orderService := service.NewOrderService(orderRepo, log)
	commentService := service.NewCommentService(commentRepo, log)
	photoService := service.NewPhotoService(photoRepo, log)

	// === Хендлеры ===
	userHandler := handler.NewUserHandler(userService)
	productHandler := handler.NewProductHandler(productService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	orderHandler := handler.NewOrderHandler(orderService)
	commentHandler := handler.NewCommentHandler(commentService)
	authHandler := handler.NewAuthHandler(userService)
	photoHandler := handler.NewPhotoHandler(photoService)

	// === Роутинг ===
	router := mux.NewRouter()
	router.Use(handler.CorsMiddleware)
	api := router.PathPrefix("/api").Subrouter()

	// --- Публичные маршруты ---
	public := api.PathPrefix("").Subrouter()
	public.HandleFunc("/auth/register", authHandler.Register).Methods("POST")
	public.HandleFunc("/auth/login", authHandler.Login).Methods("POST")
	public.HandleFunc("/products", productHandler.ListProducts).Methods("GET")
	public.HandleFunc("/products/search", productHandler.SearchProducts).Methods("GET")
	public.HandleFunc("/products/{id}", productHandler.GetProduct).Methods("GET")
	public.HandleFunc("/categories", categoryHandler.ListCategories).Methods("GET")
	public.HandleFunc("/categories/{id}", categoryHandler.GetCategory).Methods("GET")
	public.HandleFunc("/photos/{objectName}", photoHandler.Download).Methods("GET")

	// --- Защищённые маршруты ---
	protected := api.PathPrefix("").Subrouter()
	protected.Use(handler.AuthMiddleware(log))
	protected.HandleFunc("/comments", commentHandler.CreateComment).Methods("POST")
	protected.HandleFunc("/orders", orderHandler.CreateOrder).Methods("POST")

	// --- Админские маршруты ---
	admin := api.PathPrefix("").Subrouter()
	admin.Use(handler.AuthMiddleware(log), handler.AdminMiddleware(log))
	admin.HandleFunc("/products", productHandler.CreateProduct).Methods("POST")
	admin.HandleFunc("/products/{id}", productHandler.UpdateProduct).Methods("PUT")
	admin.HandleFunc("/products/{id}", productHandler.DeleteProduct).Methods("DELETE")
	admin.HandleFunc("/categories", categoryHandler.CreateCategory).Methods("POST")
	admin.HandleFunc("/categories/{id}", categoryHandler.UpdateCategory).Methods("PUT")
	admin.HandleFunc("/categories/{id}", categoryHandler.DeleteCategory).Methods("DELETE")
	admin.HandleFunc("/users", userHandler.ListUsers).Methods("GET")
	admin.HandleFunc("/users/{id}", userHandler.UpdateUser).Methods("PUT")
	admin.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods("DELETE")
	admin.HandleFunc("/photos", photoHandler.Upload).Methods("POST")

	// --- Технические ---
	router.Handle("/metrics", promhttp.Handler())
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	log.Infof("Server starting on :8080 at %s", time.Now().Format("2006-01-02 15:04:05"))
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

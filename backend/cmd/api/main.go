package main

import (
	"expense-management-system/internal/handler"
	"expense-management-system/internal/middleware"
	"expense-management-system/internal/repository"
	"expense-management-system/internal/usecase"
	"expense-management-system/internal/worker"
	"expense-management-system/pkg/config"
	"expense-management-system/pkg/database"
	"expense-management-system/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	godotenv.Load()
	logger.Init()

	cfg := config.Load()

	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		logger.ErrorLogger.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	logger.InfoLogger.Println("Connected to database successfully")

	paymentChan := make(chan usecase.PaymentJob, 100)

	userRepo := repository.NewUserRepository(db)
	expenseRepo := repository.NewExpenseRepository(db)
	approvalRepo := repository.NewApprovalRepository(db)
	auditRepo := repository.NewAuditLogRepository(db)

	authUsecase := usecase.NewAuthUsecase(userRepo, cfg)
	expenseUsecase := usecase.NewExpenseUsecase(expenseRepo, approvalRepo, auditRepo, userRepo, paymentChan)

	paymentService := worker.NewPaymentService(cfg, expenseRepo, auditRepo)
	workerPool := worker.NewWorkerPool(paymentChan, paymentService, cfg.WorkerPoolSize, cfg.WorkerMaxRetries)
	workerPool.Start()

	authHandler := handler.NewAuthHandler(authUsecase)
	expenseHandler := handler.NewExpenseHandler(expenseUsecase)
	healthHandler := handler.NewHealthHandler()
	docsHandler := handler.NewDocsHandler()

	router := mux.NewRouter()

	rateLimiter := middleware.NewRateLimiter(100, 20)
	go rateLimiter.CleanupVisitors()

	router.Use(middleware.LoggingMiddleware)
	router.Use(rateLimiter.Middleware)

	// Documentation endpoints (no auth required)
	router.HandleFunc("/docs", docsHandler.SwaggerUI).Methods("GET")
	router.HandleFunc("/docs/openapi.yaml", docsHandler.ServeOpenAPISpec).Methods("GET")

	router.HandleFunc("/api/health", healthHandler.Health).Methods("GET")
	router.HandleFunc("/api/auth/login", authHandler.Login).Methods("POST")

	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.Use(middleware.AuthMiddleware(authUsecase))

	apiRouter.HandleFunc("/expenses", expenseHandler.Submit).Methods("POST")
	apiRouter.HandleFunc("/expenses", expenseHandler.List).Methods("GET")

	// Manager-only routes - MUST be before /{id} route to avoid conflicts
	apiRouter.Handle("/expenses/pending", middleware.ManagerOnly(http.HandlerFunc(expenseHandler.GetPendingApprovals))).Methods("GET")
	apiRouter.Handle("/expenses/{id}/approve", middleware.ManagerOnly(http.HandlerFunc(expenseHandler.Approve))).Methods("PUT")
	apiRouter.Handle("/expenses/{id}/reject", middleware.ManagerOnly(http.HandlerFunc(expenseHandler.Reject))).Methods("PUT")

	// Generic /{id} route must be last
	apiRouter.HandleFunc("/expenses/{id}", expenseHandler.GetByID).Methods("GET")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://frontend:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	server := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: handler,
	}

	go func() {
		logger.InfoLogger.Printf("Server starting on port %s", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.ErrorLogger.Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.InfoLogger.Println("Shutting down server...")
	workerPool.Stop()

	if err := server.Close(); err != nil {
		logger.ErrorLogger.Printf("Error closing server: %v", err)
	}

	logger.InfoLogger.Println("Server stopped")
}

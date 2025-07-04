package main

import (
    "net/http"
    "transfer-service/api/handler"
    "transfer-service/repository"
    "transfer-service/service"
    "transfer-service/middleware"
    "github.com/gorilla/mux"
    "go.uber.org/zap"
)

func main() {
    // Initialize logger
    middleware.InitLogger()
    log := middleware.GetLogger()
    
    log.Info("Starting transfer service")
    
    // Initialize database middleware
    dbMiddleware, err := middleware.NewDatabaseMiddleware()
    if err != nil {
        log.Fatal("Failed to initialize database", zap.Error(err))
    }
    defer dbMiddleware.Close()

    accountRepo := repository.NewAccountRepository(dbMiddleware.GetDB())
    transactionRepo := repository.NewTransactionRepository(dbMiddleware.GetDB())
    
    accountSvc := service.NewAccountService(accountRepo)
    transactionSvc := service.NewTransactionService(accountRepo, transactionRepo)

    accountHandler := handler.NewAccountHandler(accountSvc)
    txHandler := handler.NewTransactionHandler(transactionSvc)

    r := mux.NewRouter()
    
    // Apply logging middleware to all routes
    r.Use(middleware.LoggingMiddleware)
    
    r.HandleFunc("/accounts", accountHandler.CreateAccount).Methods("POST")
    r.HandleFunc("/accounts/{id}", accountHandler.GetAccount).Methods("GET")
    r.HandleFunc("/transactions", txHandler.Transfer).Methods("POST")

    // Not in the scope of the project...
    r.HandleFunc("/transactions", txHandler.GetTransactionHistory).Methods("GET")
    r.HandleFunc("/accounts/{id}/transactions", txHandler.GetAccountTransactionHistory).Methods("GET")

    log.Info("Server listening on :8080")
    if err := http.ListenAndServe(":8080", r); err != nil {
        log.Fatal("Server failed to start", zap.Error(err))
    }
}

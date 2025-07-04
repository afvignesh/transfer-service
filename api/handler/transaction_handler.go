package handler

import (
    "encoding/json"
    "net/http"
    "strconv"
    "transfer-service/model"
    "transfer-service/service"
    "github.com/gorilla/mux"
)

type TransactionHandler struct {
    svc *service.TransactionService
}

func NewTransactionHandler(s *service.TransactionService) *TransactionHandler {
    return &TransactionHandler{svc: s}
}

func (h *TransactionHandler) Transfer(w http.ResponseWriter, r *http.Request) {
    var tx model.Transaction
    if err := json.NewDecoder(r.Body).Decode(&tx); err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(model.APIResponse{
            Success: false,
            Message: "Invalid request body",
            Error:   err.Error(),
        })
        return
    }

    // Get result from service
    result := h.svc.Transfer(r.Context(), tx)
    
    // Pass through the service response
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(result.Status)
    if result.Success {
        json.NewEncoder(w).Encode(model.APIResponse{
            Success: result.Success,
            Message: result.Message,
            Data:    result.Data,
        })
    } else {
        json.NewEncoder(w).Encode(model.APIResponse{
            Success: result.Success,
            Message: result.Message,
            Error:   result.Error,
        })
    }
}

func (h *TransactionHandler) GetTransactionHistory(w http.ResponseWriter, r *http.Request) {
    // Get result from service
    result := h.svc.GetTransactionHistory(r.Context())
    
    // Pass through the service response
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(result.Status)
    if result.Success {
        json.NewEncoder(w).Encode(model.APIResponse{
            Success: result.Success,
            Message: result.Message,
            Data:    result.Data,
        })
    } else {
        json.NewEncoder(w).Encode(model.APIResponse{
            Success: result.Success,
            Message: result.Message,
            Error:   result.Error,
        })
    }
}

func (h *TransactionHandler) GetAccountTransactionHistory(w http.ResponseWriter, r *http.Request) {
    idStr := mux.Vars(r)["id"]
    id, err := strconv.Atoi(idStr)
    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(model.APIResponse{
            Success: false,
            Message: "Invalid account ID",
            Error:   err.Error(),
        })
        return
    }

    // Get result from service
    result := h.svc.GetAccountTransactionHistory(r.Context(), id)
    
    // Pass through the service response
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(result.Status)
    if result.Success {
        json.NewEncoder(w).Encode(model.APIResponse{
            Success: result.Success,
            Message: result.Message,
            Data:    result.Data,
        })
    } else {
        json.NewEncoder(w).Encode(model.APIResponse{
            Success: result.Success,
            Message: result.Message,
            Error:   result.Error,
        })
    }
}

package handler

import (
    "encoding/json"
    "net/http"
    "strconv"
    "transfer-service/model"
    "transfer-service/service"
    "github.com/gorilla/mux"
)

type AccountHandler struct {
    svc *service.AccountService
}

func NewAccountHandler(s *service.AccountService) *AccountHandler {
    return &AccountHandler{svc: s}
}

func (h *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
    var acc model.Account
    if err := json.NewDecoder(r.Body).Decode(&acc); err != nil {
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
    result := h.svc.CreateAccount(r.Context(), acc)
    
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

func (h *AccountHandler) GetAccount(w http.ResponseWriter, r *http.Request) {
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
    result := h.svc.GetAccount(r.Context(), id)
    
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

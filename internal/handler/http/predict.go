package http

import (
	"log"
    "encoding/json"
    "net/http"
    
    "github.com/kevo-1/model-serving-platform/internal/domain"
)

func (h *Handler) handlePredict(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    var req domain.PredictionRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        log.Printf("JSON decode error: %v", err)
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    
    res, err := h.predictionService.Predict(r.Context(), req)
    
    if err != nil {
        log.Printf("Prediction error: %v", err)
        switch e := err.(type) {
            case *domain.ValidationError:
                http.Error(w, e.Error(), http.StatusBadRequest)
                return
            case *domain.ModelNotFoundError:
                http.Error(w, e.Error(), http.StatusNotFound)
                return
            case *domain.InvalidInputError:
                http.Error(w, e.Error(), http.StatusBadRequest)
                return
            default:
                log.Printf("Unhandled error type: %T", err)
                http.Error(w, "Internal server error", http.StatusInternalServerError)
                return
        }
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(res)
}
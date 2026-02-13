package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/kevo-1/model-serving-platform/internal/domain"
	"github.com/kevo-1/model-serving-platform/internal/metrics"
	"github.com/kevo-1/model-serving-platform/internal/repository"
)

type PredictionService struct {
	registry *repository.ModelRegistry
}

func NewPredictionService(registry *repository.ModelRegistry) *PredictionService {
    return &PredictionService{
        registry: registry,
    }
}

func (s *PredictionService) Predict(ctx context.Context, req domain.PredictionRequest) (domain.PredictionResponse, error) {
    //validate request
	if err := req.Validate();err != nil {
		return domain.PredictionResponse{}, err
	}
	
    //generate RequestID if empty
	if req.RequestID == "" {
		req.RequestID = uuid.New().String()
	}
    
    //get model from registry
	model, err := s.registry.Get(req.ModelID)
	if err != nil {
		return domain.PredictionResponse{}, err
	}

    inferenceStart := time.Now()
    
    prediction, err := model.Predict(ctx, req.Features)
    
    inferenceDuration := time.Since(inferenceStart).Seconds()
    
    // Record metrics
    success := err == nil
    metrics.RecordPrediction(req.ModelID, success, inferenceDuration)
    
    if err != nil {
        return domain.PredictionResponse{}, err
    }

    //Build response with timing
	totalLatency := float64(time.Since(inferenceStart).Microseconds()) / 1000
	
	response := domain.PredictionResponse{
		ModelID: req.ModelID,
		RequestID: req.RequestID,
		LatencyMs: totalLatency,
		Prediction: prediction,
		Timestamp: time.Now(),
	}

	return response, nil
}
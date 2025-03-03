package pkg

type ConfidenceResult struct {
	Confidence int `json:"confidence"`
}

type ConfidenceProvider interface {
	PredictConfidence(news string) (ConfidenceResult, error)
}

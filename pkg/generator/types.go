package generator

type ConfidenceResult struct {
	Confidence int `json:"confidence"`
}

type ConfidenceGenerator interface {
	GenerateConfidence(news string) (ConfidenceResult, error)
}

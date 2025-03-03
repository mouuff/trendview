package pkg

type ProviderConfiguration struct {
	Model  string
	Prompt string
}

type ConfidenceResult struct {
	Confidence int `json:"confidence"`
}

type ConfidenceProvider interface {
	Predict(config ProviderConfiguration, news string) (ConfidenceResult, error)
}

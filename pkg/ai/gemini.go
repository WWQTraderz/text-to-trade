package ai

import (
	"context"

	"cloud.google.com/go/vertexai/genai"
)

func Gemini(ctx context.Context) (*genai.Client, error) {
	return genai.NewClient(ctx, "miami-ai-hack24mia-901", "us-central1")
}

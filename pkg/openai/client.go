package openai

import (
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"github.com/llamacto/llama-gin-kit/config"
	"io"
)

var client *openai.Client

// Init initializes the OpenAI client
func Init(cfg *config.Config) error {
	client = openai.NewClient(cfg.OpenAI.APIKey)
	return nil
}

// GetClient returns the OpenAI client instance
func GetClient() *openai.Client {
	return client
}

// GenerateAudio generates audio from text using OpenAI's TTS API
func GenerateAudio(ctx context.Context, text string) ([]byte, error) {
	// Create audio file
	req := openai.CreateSpeechRequest{
		Model: openai.TTSModel1,
		Input: text,
		Voice: openai.VoiceAlloy,
	}

	// Get audio data from OpenAI
	resp, err := client.CreateSpeech(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create speech: %v", err)
	}

	// Read audio data
	data, err := io.ReadAll(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to read audio data: %v", err)
	}

	return data, nil
}

package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/llamacto/llama-gin-kit/internal/modules/tts"
	"github.com/llamacto/llama-gin-kit/middleware"
	"github.com/llamacto/llama-gin-kit/pkg/database"
)

// RegisterTTSRoutes registers all TTS-related routes
func RegisterTTSRoutes(v1 *gin.RouterGroup) {
	// Initialize TTS repository and service
	ttsRepo := tts.NewTTSRepository(database.GetDB())
	ttsService := tts.NewTTSService(ttsRepo)
	ttsHandler := tts.NewTTSHandler(ttsService)

	// Public TTS routes
	v1.GET("/tts/audio/:id", ttsHandler.GetAudio)
	v1.GET("/tts/audio/:id/download", ttsHandler.DownloadAudio)
	v1.GET("/tts/voices", ttsHandler.GetVoices)

	// Protected TTS routes
	ttsGroup := v1.Group("/tts")
	{
		ttsGroup.Use(middleware.JWT())
		ttsGroup.POST("/generate", ttsHandler.Generate)
		ttsGroup.POST("/translate", ttsHandler.Translate)
		ttsGroup.GET("/history", ttsHandler.GetHistory)
		ttsGroup.DELETE("/history/:id", ttsHandler.DeleteHistory)
	}
}

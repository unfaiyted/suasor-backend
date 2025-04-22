package bundles

import (
	"suasor/client/types"
	"suasor/handlers"
)

// AIClientHandlersImpl implements the AIClientHandlers interface
type AIClientHandlersImpl struct {
	claudeHandler *handlers.AIHandler[*types.ClaudeConfig]
	openaiHandler *handlers.AIHandler[*types.OpenAIConfig]
	ollamaHandler *handlers.AIHandler[*types.OllamaConfig]
}

func NewAIClientHandlers(
	claudeHandler *handlers.AIHandler[*types.ClaudeConfig],
	openaiHandler *handlers.AIHandler[*types.OpenAIConfig],
	ollamaHandler *handlers.AIHandler[*types.OllamaConfig],
) AIClientHandlers {
	return &AIClientHandlersImpl{
		claudeHandler: claudeHandler,
		openaiHandler: openaiHandler,
		ollamaHandler: ollamaHandler,
	}
}

func (h *AIClientHandlersImpl) ClaudeAIHandler() *handlers.AIHandler[*types.ClaudeConfig] {
	return h.claudeHandler
}

func (h *AIClientHandlersImpl) OpenAIHandler() *handlers.AIHandler[*types.OpenAIConfig] {
	return h.openaiHandler
}

func (h *AIClientHandlersImpl) OllamaHandler() *handlers.AIHandler[*types.OllamaConfig] {
	return h.ollamaHandler
}

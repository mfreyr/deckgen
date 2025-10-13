package llm

import (
	"fmt"

	"github.com/mfreyr/deckgen/internal/config"
	"github.com/mfreyr/deckgen/internal/service"
)

type LLMFactory struct {
	providers map[service.LLMProviderName]service.LLMProvider
}

func NewLLMFactory(cfg map[string]config.LLMProviderConfig) (*LLMFactory, error) {
	providers := make(map[service.LLMProviderName]service.LLMProvider)
	for name, providerCfg := range cfg {
		if !providerCfg.Enabled {
			continue
		}

		switch service.LLMProviderName(name) {
		case "openai":
			provider, err := NewOpenAIProvider(providerCfg)
			if err != nil {
				return nil, fmt.Errorf("failed to initialize %s provider: %w", name, err)
			}
			providers[service.LLMProviderName(name)] = provider
		default:
			return nil, fmt.Errorf("failed to initialize %s provider: %s", name, "provider is not supported")
		}

	}
	return &LLMFactory{providers: providers}, nil
}

func (f *LLMFactory) GetProvider(providerType service.LLMProviderName) (service.LLMProvider, error) {
	provider, ok := f.providers[providerType]
	if !ok {
		return nil, fmt.Errorf("provider '%s' is not supported or not enabled in config", providerType)
	}
	return provider, nil
}

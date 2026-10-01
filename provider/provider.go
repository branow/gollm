// Package provider maps an gollm.Config onto one of the provider packages. It is
// separate because the providers import gollm, so gollm cannot import them back.
package provider

import (
	"sort"

	"github.com/branow/gollm"
	"github.com/branow/gollm/anthropic"
	"github.com/branow/gollm/claudecode"
	"github.com/branow/gollm/openai"
)

// constructors is the whole registry; nothing registers globally.
var constructors = map[string]func(gollm.Config) (gollm.Client, error){
	anthropic.Name: func(cfg gollm.Config) (gollm.Client, error) {
		return anthropic.New(anthropic.Config{
			APIKey:    cfg.APIKey,
			BaseURL:   cfg.BaseURL,
			Model:     cfg.Model,
			MaxTokens: cfg.MaxTokens,
		})
	},
	openai.Name: func(cfg gollm.Config) (gollm.Client, error) {
		return openai.New(openai.Config{
			APIKey:    cfg.APIKey,
			BaseURL:   cfg.BaseURL,
			Model:     cfg.Model,
			MaxTokens: cfg.MaxTokens,
		})
	},
	claudecode.Name: func(cfg gollm.Config) (gollm.Client, error) {
		return claudecode.New(claudecode.Config{
			Command: cfg.Command,
			Model:   cfg.Model,
		})
	},
}

// New builds the client a config names.
func New(cfg gollm.Config) (gollm.Client, error) {
	construct, ok := constructors[cfg.Provider]
	if !ok {
		return nil, &gollm.Error{
			Class: gollm.ClassBadRequest, Provider: cfg.Provider, Detail: "unknown provider",
		}
	}
	return construct(cfg)
}

// Names lists the supported providers, sorted.
func Names() []string {
	names := make([]string, 0, len(constructors))
	for name := range constructors {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

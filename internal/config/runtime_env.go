package config

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// ApplyRuntimeEnv applies cloud/host overrides (Render, Docker, etc.).
func ApplyRuntimeEnv(cfg *Config) {
	if cfg == nil {
		return
	}
	if v := strings.TrimSpace(os.Getenv("PORT")); v != "" {
		if p, err := strconv.Atoi(v); err == nil && p > 0 {
			cfg.Server.Port = p
		}
	}
	if v := strings.TrimSpace(os.Getenv("AUTH_PASSWORD")); v != "" {
		cfg.Auth.Password = v
	}
	if v := strings.TrimSpace(os.Getenv("OPENAI_API_KEY")); v != "" {
		cfg.OpenAI.APIKey = v
	}
	if strings.TrimSpace(os.Getenv("MEOWCYBER_HTTPS")) == "0" {
		cfg.Server.TLSEnabled = false
		cfg.Server.TLSAutoSelfSign = false
	}
	if dataDir := strings.TrimSpace(os.Getenv("MEOWCYBER_DATA_DIR")); dataDir != "" {
		cfg.Database.Path = filepath.Join(dataDir, "conversations.db")
		if strings.TrimSpace(cfg.Database.KnowledgeDBPath) != "" {
			cfg.Database.KnowledgeDBPath = filepath.Join(dataDir, "knowledge.db")
		}
	}
}

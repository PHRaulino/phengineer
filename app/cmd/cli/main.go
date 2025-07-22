package main

import (
	"context"
	"log"

	"github.com/PHRaulino/phengineer/internal/infrastructure/config"
	"github.com/PHRaulino/phengineer/internal/infrastructure/utils/logger"
)

func main() {
	logger.SetupLogger()

	ctx := context.Background()

	// Adiciona config ao context com verbose
	ctx, err := config.WithConfig(ctx, "app/.phengineer")
	if err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	// Mostra diagnósticos
	config.PrintDiagnostics(ctx)
}

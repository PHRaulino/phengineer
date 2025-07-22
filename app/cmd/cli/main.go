package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/PHRaulino/phengineer/internal/domain/discovery"
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

	service := discovery.NewService()
	// Usa discovery com lock para tracking de mudanças
	result, changes, err := service.DiscoverFilesWithLock(ctx)
	if err != nil {
		fmt.Printf("Erro: %v\n", err)
		return
	}

	fmt.Printf("=== Resultado da Descoberta ===\n")
	fmt.Printf("Total encontrados: %d\n", result.TotalFound)
	fmt.Printf("Total filtrados: %d\n", result.TotalFiltered)
	fmt.Printf("Arquivos válidos: %d\n", len(result.Files))
	fmt.Printf("Arquivos grandes: %d\n", len(result.OversizedFiles))
	fmt.Printf("Commit atual: %s\n", result.GitCommit)
	fmt.Printf("Timestamp: %v\n", time.Unix(result.Timestamp, 0))

	fmt.Printf("\n=== Análise de Mudanças ===\n")
	fmt.Printf("Houve mudanças: %t\n", changes.HasChanges)
	fmt.Printf("Arquivos novos: %d\n", len(changes.NewFiles))
	fmt.Printf("Arquivos alterados: %d\n", len(changes.ChangedFiles))
	fmt.Printf("Arquivos deletados: %d\n", len(changes.DeletedFiles))
	fmt.Printf("Arquivos inalterados: %d\n", len(changes.UnchangedFiles))

	if len(changes.NewFiles) > 0 {
		fmt.Println("\n=== Arquivos Novos ===")
		for i, file := range changes.NewFiles {
			if i >= 5 {
				fmt.Printf("... e mais %d arquivos novos\n", len(changes.NewFiles)-5)
				break
			}
			fmt.Printf("+ %s/%s (%s) [%s]\n", file.Path, file.Name, file.Type, file.PatternType)
		}
	}

	if len(changes.ChangedFiles) > 0 {
		fmt.Println("\n=== Arquivos Alterados ===")
		for i, file := range changes.ChangedFiles {
			if i >= 5 {
				fmt.Printf("... e mais %d arquivos alterados\n", len(changes.ChangedFiles)-5)
				break
			}
			fmt.Printf("~ %s/%s (%s) [commit: %s]\n", file.Path, file.Name, file.Type, file.CommitHash[:8])
		}
	}

	if len(changes.DeletedFiles) > 0 {
		fmt.Println("\n=== Arquivos Deletados ===")
		for _, file := range changes.DeletedFiles {
			fmt.Printf("- %s/%s (%s)\n", file.Path, file.Name, file.Type)
		}
	}

	// Separa por tipo de pattern
	snippetFiles := make([]discovery.File, 0)
	customFiles := make([]discovery.File, 0)

	for _, file := range result.Files {
		if file.PatternType == discovery.PatternTypeSnippet {
			snippetFiles = append(snippetFiles, file)
		} else {
			customFiles = append(customFiles, file)
		}
	}

	fmt.Printf("\n=== Por Tipo de Pattern ===\n")
	fmt.Printf("Arquivos de código (snippet): %d\n", len(snippetFiles))
	fmt.Printf("Arquivos de config/docs (custom): %d\n", len(customFiles))
}

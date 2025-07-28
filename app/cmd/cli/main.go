package main

import (
	"context"
	"fmt"
	"os"
	"time"

	cli "github.com/PHRaulino/phengineer/cmd/cli/commands"
	"github.com/PHRaulino/phengineer/internal/domain/discovery"
	"github.com/PHRaulino/phengineer/internal/infrastructure/auth"
	"github.com/PHRaulino/phengineer/internal/infrastructure/config"
	"github.com/PHRaulino/phengineer/internal/infrastructure/utils/logger"
	"github.com/spf13/cobra"
)

func showWelcomeScreen() {
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║                     🛠️  PHEngineer CLI                       ║")
	fmt.Println("║                                                              ║")
	fmt.Println("║        Professional File Discovery & Analysis Tool          ║")
	fmt.Println("║                                                              ║")
	fmt.Println("║  Discover, analyze and track changes in your codebase       ║")
	fmt.Println("║  with intelligent pattern matching and Git integration      ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	fmt.Println()
}

var rootCmd = &cobra.Command{
	Use:   "phengineer",
	Short: "🛠️ PHEngineer CLI - Professional File Discovery & Analysis Tool",
	Long: `PHEngineer CLI é uma ferramenta profissional para descoberta e análise de arquivos.
	
Descubra, analise e acompanhe mudanças em sua base de código
com correspondência inteligente de padrões e integração Git.`,
	RunE: runDiscovery,
}

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Gerenciar autenticação",
	Long:  "Comandos para configurar e gerenciar autenticação com diferentes provedores",
}

func init() {
	logger.SetupLogger()
	auth.SetupGenerators()
	
	// Adicionar comando auth
	rootCmd.AddCommand(authCmd)
	authCmd.AddCommand(cli.GetAuthSetupCmd())
}

func runDiscovery(cmd *cobra.Command, args []string) error {
	showWelcomeScreen()

	ctx := context.Background()

	// Adiciona config ao context com verbose
	ctx, err := config.WithConfig(ctx, ".phengineer")
	if err != nil {
		return fmt.Errorf("failed to initialize config: %w", err)
	}

	service := discovery.NewService()
	// Usa discovery com lock para tracking de mudanças
	result, changes, err := service.DiscoverFilesWithLock(ctx)
	if err != nil {
		return fmt.Errorf("erro na descoberta: %w", err)
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
	
	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Erro: %v\n", err)
		os.Exit(1)
	}
}

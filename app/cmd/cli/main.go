package main

import (
	"context"
	"fmt"
	"log"

	"github.com/PHRaulino/phengineer/internal/infrastructure/config"
)

func main() {
	// Exemplo básico
	basicExample()

	// Exemplo com verbose
	verboseExample()

	// Exemplo de uso em diferentes funções
	contextFlowExample()
}

// basicExample mostra uso básico
func basicExample() {
	fmt.Println("=== Basic Context Config Example ===")

	ctx := context.Background()

	// Adiciona config ao context
	ctx, err := config.WithConfig(ctx, ".phengineer")
	if err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	// Usa a config
	cfg := config.FromContext(ctx)
	fmt.Printf("✅ App: %s\n", cfg.Auto.AppName)
	fmt.Printf("✅ Remote: %s\n", cfg.Auto.RemoteURL)

	// Ou acessa diretamente as settings
	settings := config.GetSettings(ctx)
	fmt.Printf("✅ Language: %s %s\n", settings.Project.Language.Name, settings.Project.Language.Version)
}

// verboseExample mostra inicialização com detalhes
func verboseExample() {
	fmt.Println("\n=== Verbose Context Config Example ===")

	ctx := context.Background()

	// Adiciona config ao context com verbose
	ctx, err := config.WithConfigVerbose(ctx, ".phengineer")
	if err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	// Mostra diagnósticos
	config.PrintDiagnostics(ctx)
}

// contextFlowExample mostra como passar context entre funções
func contextFlowExample() {
	fmt.Println("\n=== Context Flow Example ===")

	ctx := context.Background()

	// Inicializa config
	ctx, err := config.WithConfig(ctx, ".phengineer")
	if err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	// Passa context para outras funções
	processProject(ctx)
	analyzeFiles(ctx)
	generateReport(ctx)
	optionalConfigFunction(ctx)
	middlewareExample()
	cliExample()
}

// processProject simula processamento usando config
func processProject(ctx context.Context) {
	cfg := config.FromContext(ctx)
	settings := cfg.Settings

	fmt.Printf("🔄 Processing project: %s\n", cfg.Auto.AppName)
	fmt.Printf("   Type: %s\n", settings.Project.Type)
	fmt.Printf("   Language: %s %s\n", settings.Project.Language.Name, settings.Project.Language.Version)
	fmt.Printf("   Remote: %s\n", cfg.Auto.RemoteURL)
}

// analyzeFiles simula análise de arquivos usando config
func analyzeFiles(ctx context.Context) {
	settings := config.GetSettings(ctx)
	auto := config.GetAutoConfig(ctx)

	fmt.Printf("📊 Analyzing files for: %s\n", auto.AppName)
	fmt.Printf("   Include pattern: %s\n", settings.Analysis.FilesIncludePath)
	fmt.Printf("   Exclude pattern: %s\n", settings.Analysis.FilesExcludePath)
	fmt.Printf("   Max file size: %s\n", settings.Analysis.FileLimits.MaxFileSize)
	fmt.Printf("   Max files: %d\n", settings.Analysis.FileLimits.MaxFiles)

	// Simula lógica de análise baseada nas configurações
	if settings.Project.Language.Name == "go" {
		fmt.Println("   Using Go-specific analysis rules")
	}
}

// generateReport simula geração de relatório usando config
func generateReport(ctx context.Context) {
	cfg := config.FromContext(ctx)

	fmt.Printf("📄 Generating report for: %s\n", cfg.Auto.AppName)
	fmt.Printf("   Project type: %s\n", cfg.Settings.Project.Type)
	fmt.Printf("   Config dir: %s\n", cfg.Auto.ConfigDirPath)

	// Aqui você usaria as configurações para gerar o relatório
	// com base no tipo de projeto, linguagem, etc.
}

// Exemplo de uma função que pode receber context com ou sem config
func optionalConfigFunction(ctx context.Context) {
	if config.HasConfig(ctx) {
		cfg := config.FromContext(ctx)
		fmt.Printf("Using configured app: %s\n", cfg.Auto.AppName)
	} else {
		fmt.Println("No config found, using defaults")
	}
}

// Exemplo de middleware/handler pattern
func withConfigMiddleware(next func(context.Context)) func(context.Context) {
	return func(ctx context.Context) {
		// Adiciona config se não estiver presente
		if !config.HasConfig(ctx) {
			var err error
			ctx, err = config.WithConfig(ctx, ".phengineer")
			if err != nil {
				log.Printf("Failed to load config: %v", err)
				return
			}
		}

		// Chama próxima função
		next(ctx)
	}
}

// Exemplo de uso do middleware
func middlewareExample() {
	fmt.Println("\n=== Middleware Example ===")

	// Função que precisa de config
	handler := func(ctx context.Context) {
		cfg := config.FromContext(ctx)
		fmt.Printf("Handler called for app: %s\n", cfg.Auto.AppName)
	}

	// Envolve com middleware
	wrappedHandler := withConfigMiddleware(handler)

	// Chama sem config - middleware adiciona automaticamente
	ctx := context.Background()
	wrappedHandler(ctx)
}

// Exemplo de CLI completa
func cliExample() {
	fmt.Println("\n=== CLI Example ===")

	ctx := context.Background()

	// Inicializa config com verbose para CLI
	ctx, err := config.WithConfigVerbose(ctx, ".phengineer")
	if err != nil {
		fmt.Printf("❌ Configuration failed: %v\n", err)
		showConfigHelp()
		return
	}

	// CLI está pronta para usar
	runCLI(ctx)
}

// runCLI simula execução da CLI
func runCLI(ctx context.Context) {
	cfg := config.FromContext(ctx)

	fmt.Printf("🚀 %s CLI started successfully!\n", cfg.Auto.AppName)
	fmt.Printf("📁 Working directory: %s\n", cfg.Auto.ConfigDirPath)

	// Aqui você implementaria os comandos da CLI
	// Todos receberiam o context com a config
}

// showConfigHelp mostra ajuda de configuração
func showConfigHelp() {
	fmt.Printf(`
Configuration Setup Help:

1. Make sure you're in a Git repository
2. Create a 'config' folder in the repository root
3. The settings.yml file will be created automatically

Requirements:
• Git installed and available in PATH
• Inside a Git repository
• Git repository has at least one remote
• Config folder exists in repository root
`)
}

# PHEngineer - Estrutura Core/UseCase

## 📁 Estrutura Completa de Diretórios

```
phengineer/
├── cmd/
│   └── cli/
│       └── main.go                    # CLI entry point
├── internal/
│   ├── settings/
│   │   ├── core/
│   │   │   ├── models.go              # ProjectSettings, AnalysisConfig, FileLimits
│   │   │   ├── ports.go               # SettingsReader, PatternLoader, Validator
│   │   │   └── validator.go           # Settings validation logic
│   │   └── usecase/
│   │       ├── load_settings.go       # LoadSettingsUseCase
│   │       └── load_patterns.go       # LoadPatternsUseCase
│   ├── discovery/
│   │   ├── core/
│   │   │   ├── models.go              # FileFilter, DiscoveryResult, FileInfo
│   │   │   ├── ports.go               # FileWalker, PatternMatcher, FileSizeChecker
│   │   │   ├── discoverer.go          # File discovery logic
│   │   │   └── filter.go              # File filtering logic
│   │   └── usecase/
│   │       └── discover_files.go      # DiscoverFilesUseCase
│   ├── analyzer/
│   │   ├── core/
│   │   │   ├── models.go              # AnalysisResult, FileAnalysis, CodeElement
│   │   │   ├── ports.go               # FileReader, CodeParser, StackDetector
│   │   │   ├── file_analyzer.go       # File structure analysis logic
│   │   │   ├── code_analyzer.go       # Code parsing and function extraction
│   │   │   └── stack_analyzer.go      # Technology stack detection
│   │   └── usecase/
│   │       └── analyze_project.go     # AnalyzeProjectUseCase
│   ├── generator/
│   │   ├── core/
│   │   │   ├── models.go              # Context, Metadata, ContextData
│   │   │   ├── ports.go               # ContextWriter, SchemaValidator, JsonSerializer
│   │   │   ├── filetree_generator.go  # FileTree context generation
│   │   │   ├── functions_generator.go # Functions context generation
│   │   │   ├── stack_generator.go     # Stack context generation
│   │   │   └── context_generator.go   # ProjectContext generation
│   │   └── usecase/
│   │       └── generate_contexts.go   # GenerateContextsUseCase
│   ├── pipeline/
│   │   ├── core/
│   │   │   ├── models.go              # PipelineConfig, ExecutionResult, StepResult
│   │   │   ├── ports.go               # StepExecutor, ResultAggregator
│   │   │   └── orchestrator.go        # Pipeline orchestration logic
│   │   └── usecase/
│   │       └── execute_pipeline.go    # ExecutePipelineUseCase
│   └── adapters/
│       ├── filesystem/
│       │   ├── settings_reader.go     # Implements SettingsReader
│       │   ├── pattern_loader.go      # Implements PatternLoader
│       │   ├── file_walker.go         # Implements FileWalker
│       │   └── file_reader.go         # Implements FileReader
│       ├── parsers/
│       │   ├── code_parser.go         # Implements CodeParser
│       │   └── stack_detector.go      # Implements StackDetector
│       └── storage/
│           ├── json_writer.go         # Implements ContextWriter
│           └── schema_validator.go    # Implements SchemaValidator
├── pkg/
│   ├── utils/
│   │   ├── patterns.go                # Pattern matching utilities
│   │   ├── filesize.go                # File size parsing utilities
│   │   └── pathutils.go               # Path manipulation utilities
│   └── errors/
│       └── errors.go                  # Custom error types
├── schemas/
│   ├── file-tree.schema.json
│   ├── functions.schema.json
│   ├── stack.schema.json
│   └── project-context.schema.json
├── testdata/
│   ├── sample-project/               # Projeto exemplo para testes
│   │   ├── .phengineer/
│   │   │   ├── settings.yml
│   │   │   ├── .ignorefiles
│   │   │   └── .analyzefiles
│   │   ├── src/
│   │   │   ├── main.py
│   │   │   ├── utils.py
│   │   │   └── models.py
│   │   ├── tests/
│   │   │   └── test_main.py
│   │   ├── requirements.txt
│   │   └── README.md
│   └── expected/                     # Outputs esperados
│       ├── file-tree.json
│       ├── functions.json
│       ├── stack.json
│       └── project-context.json
├── scripts/
│   ├── build.sh                      # Build script
│   └── test.sh                       # Test script
├── docs/
│   ├── architecture.md               # Documentação da arquitetura
│   └── examples.md                   # Exemplos de uso
├── .gitignore
├── Makefile
├── go.mod
├── go.sum
└── README.md
```

## 📦 Estrutura Detalhada dos Packages

### internal/settings/

#### core/models.go

```go
// Tipos centrais para configuração
type ProjectSettings struct {
    Project  ProjectConfig  `yaml:"project"`
    Analysis AnalysisConfig `yaml:"analysis"`
    Contexts ContextsConfig `yaml:"contexts"`
}

type ProjectConfig struct {
    Name     string `yaml:"name"`
    Type     string `yaml:"type"`
    Language string `yaml:"language"`
    Version  string `yaml:"version"`
}

type AnalysisConfig struct {
    FileLimits FileLimits `yaml:"file_limits"`
}

type FileLimits struct {
    MaxFileSize string `yaml:"max_file_size"`
    MaxFiles    int    `yaml:"max_files"`
}

type ContextsConfig struct {
    Enabled []string `yaml:"enabled"`
}

type PatternConfig struct {
    IgnorePatterns  []string
    AnalyzePatterns []string
}
```

#### core/ports.go

```go
// Interfaces para infraestrutura
type SettingsReader interface {
    ReadSettings(path string) (*ProjectSettings, error)
}

type PatternLoader interface {
    LoadIgnorePatterns(path string) ([]string, error)
    LoadAnalyzePatterns(path string) ([]string, error)
}

type SettingsValidator interface {
    Validate(settings *ProjectSettings) error
}
```

#### core/validator.go

```go
// Lógica de validação pura
type Validator struct{}

func (v *Validator) Validate(settings *ProjectSettings) error
func (v *Validator) validateProject(config ProjectConfig) error
func (v *Validator) validateAnalysis(config AnalysisConfig) error
func (v *Validator) validateContexts(config ContextsConfig) error
```

#### usecase/load_settings.go

```go
// Orquestração de carregamento
type LoadSettingsUseCase struct {
    reader    SettingsReader
    validator SettingsValidator
}

func (uc *LoadSettingsUseCase) Execute(settingsPath string) (*ProjectSettings, error)
```

#### usecase/load_patterns.go

```go
// Orquestração de carregamento de patterns
type LoadPatternsUseCase struct {
    loader PatternLoader
}

func (uc *LoadPatternsUseCase) Execute(basePath string) (*PatternConfig, error)
```

### internal/discovery/

#### core/models.go

```go
// Tipos para descoberta de arquivos
type FileInfo struct {
    Path     string
    Size     int64
    ModTime  time.Time
    IsDir    bool
    Priority int
}

type DiscoveryResult struct {
    Files         []FileInfo
    TotalFiles    int
    FilteredFiles int
    Errors        []error
}

type FileFilter struct {
    IgnorePatterns  []string
    AnalyzePatterns []string
    MaxFileSize     int64
    MaxFiles        int
}
```

#### core/ports.go

```go
// Interfaces para descoberta
type FileWalker interface {
    Walk(rootPath string, walkFn func(path string, info os.FileInfo) error) error
}

type PatternMatcher interface {
    MatchesAny(path string, patterns []string) bool
    Matches(path string, pattern string) bool
}

type FileSizeChecker interface {
    GetSize(path string) (int64, error)
    ExceedsLimit(path string, limit int64) (bool, error)
}
```

#### core/discoverer.go

```go
// Lógica de descoberta pura
type Discoverer struct {
    walker  FileWalker
    matcher PatternMatcher
    checker FileSizeChecker
}

func (d *Discoverer) DiscoverFiles(rootPath string, filter FileFilter) (*DiscoveryResult, error)
func (d *Discoverer) shouldIgnore(path string, patterns []string) bool
func (d *Discoverer) calculatePriority(path string, patterns []string) int
```

#### core/filter.go

```go
// Lógica de filtragem pura
type Filter struct {
    matcher PatternMatcher
}

func (f *Filter) ApplyIgnoreFilters(files []FileInfo, patterns []string) []FileInfo
func (f *Filter) PrioritizeFiles(files []FileInfo, patterns []string) []FileInfo
func (f *Filter) ApplyLimits(files []FileInfo, limits FileLimits) []FileInfo
```

#### usecase/discover_files.go

```go
// Orquestração de descoberta
type DiscoverFilesUseCase struct {
    discoverer *Discoverer
    filter     *Filter
}

func (uc *DiscoverFilesUseCase) Execute(rootPath string, filter FileFilter) (*DiscoveryResult, error)
```

### internal/analyzer/

#### core/models.go

```go
// Tipos para análise
type AnalysisResult struct {
    Type      string                 `json:"type"`
    Data      map[string]interface{} `json:"data"`
    FileCount int                    `json:"file_count"`
    Timestamp time.Time              `json:"timestamp"`
    Errors    []string               `json:"errors,omitempty"`
}

type FileAnalysis struct {
    Path        string
    Extension   string
    Size        int64
    LineCount   int
    Importance  ImportanceLevel
    Category    FileCategory
}

type CodeElement struct {
    Name       string
    Type       ElementType // function, class, method, etc.
    Signature  string
    Line       int
    File       string
    Visibility string
    Params     []Parameter
    Returns    []Return
}

type ImportanceLevel int
type FileCategory string
type ElementType string
```

#### core/ports.go

```go
// Interfaces para análise
type FileReader interface {
    ReadFile(path string) ([]byte, error)
    GetFileInfo(path string) (os.FileInfo, error)
}

type CodeParser interface {
    ParseFile(path string, content []byte) ([]CodeElement, error)
    SupportedExtensions() []string
}

type StackDetector interface {
    DetectFromFile(path string, content []byte) ([]Technology, error)
    DetectFromProject(files []FileInfo) (TechStack, error)
}
```

#### core/file_analyzer.go

```go
// Lógica de análise de arquivos pura
type FileAnalyzer struct {
    reader FileReader
}

func (fa *FileAnalyzer) Analyze(files []FileInfo) (*AnalysisResult, error)
func (fa *FileAnalyzer) buildFileTree(files []FileInfo) map[string]interface{}
func (fa *FileAnalyzer) analyzeExtensions(files []FileInfo) map[string]int
func (fa *FileAnalyzer) classifyFiles(files []FileInfo) []FileAnalysis
func (fa *FileAnalyzer) calculateImportance(file FileInfo) ImportanceLevel
```

#### core/code_analyzer.go

```go
// Lógica de análise de código pura
type CodeAnalyzer struct {
    reader FileReader
    parser CodeParser
}

func (ca *CodeAnalyzer) Analyze(files []FileInfo) (*AnalysisResult, error)
func (ca *CodeAnalyzer) extractFunctions(files []FileInfo) ([]CodeElement, error)
func (ca *CodeAnalyzer) groupByFile(elements []CodeElement) map[string][]CodeElement
func (ca *CodeAnalyzer) calculateComplexity(element CodeElement) int
```

#### core/stack_analyzer.go

```go
// Lógica de detecção de stack pura
type StackAnalyzer struct {
    reader   FileReader
    detector StackDetector
}

func (sa *StackAnalyzer) Analyze(files []FileInfo) (*AnalysisResult, error)
func (sa *StackAnalyzer) detectLanguages(files []FileInfo) []Language
func (sa *StackAnalyzer) detectFrameworks(files []FileInfo) []Framework
func (sa *StackAnalyzer) detectDependencies(files []FileInfo) []Dependency
```

#### usecase/analyze_project.go

```go
// Orquestração de análises
type AnalyzeProjectUseCase struct {
    fileAnalyzer  *FileAnalyzer
    codeAnalyzer  *CodeAnalyzer
    stackAnalyzer *StackAnalyzer
}

func (uc *AnalyzeProjectUseCase) Execute(files []FileInfo) (map[string]*AnalysisResult, error)
func (uc *AnalyzeProjectUseCase) runAnalysis(analyzer Analyzer, files []FileInfo) (*AnalysisResult, error)
```

### internal/generator/

#### core/models.go

```go
// Tipos para geração de contextos
type Context struct {
    Type     string      `json:"type"`
    Metadata Metadata    `json:"metadata"`
    Data     interface{} `json:"data"`
}

type Metadata struct {
    GeneratedAt time.Time `json:"generated_at"`
    ProjectName string    `json:"project_name"`
    ProjectType string    `json:"project_type"`
    Version     string    `json:"version"`
    FileCount   int       `json:"file_count"`
}

type ContextData interface {
    GetType() string
    Validate() error
}

type FileTreeData struct {
    Structure   map[string]interface{} `json:"structure"`
    Extensions  map[string]int         `json:"extensions"`
    Important   []string               `json:"important_files"`
    Conventions map[string]string      `json:"conventions"`
}

type FunctionsData struct {
    Files     []FileFunction `json:"files"`
    Summary   FunctionSummary `json:"summary"`
    Languages []string       `json:"languages"`
}

type StackData struct {
    Languages    []Language    `json:"languages"`
    Frameworks   []Framework   `json:"frameworks"`
    Dependencies []Dependency  `json:"dependencies"`
    DevTools     []DevTool     `json:"dev_tools"`
}

type ProjectContextData struct {
    Documentation []DocumentFile `json:"documentation"`
    Configuration []ConfigFile   `json:"configuration"`
    Contracts     []ContractFile `json:"contracts"`
    Summary       ProjectSummary `json:"summary"`
}
```

#### core/ports.go

```go
// Interfaces para geração
type ContextWriter interface {
    WriteContext(context *Context, path string) error
}

type SchemaValidator interface {
    ValidateAgainstSchema(context *Context) error
    LoadSchema(contextType string) ([]byte, error)
}

type JsonSerializer interface {
    Serialize(data interface{}) ([]byte, error)
    Deserialize(data []byte, target interface{}) error
}
```

#### core/filetree_generator.go

```go
// Lógica de geração de file-tree pura
type FileTreeGenerator struct {
    settings *ProjectSettings
}

func (fg *FileTreeGenerator) Generate(analysisResults map[string]*AnalysisResult) (*Context, error)
func (fg *FileTreeGenerator) buildTreeStructure(fileAnalysis *AnalysisResult) map[string]interface{}
func (fg *FileTreeGenerator) detectConventions(fileAnalysis *AnalysisResult) map[string]string
func (fg *FileTreeGenerator) identifyImportantFiles(fileAnalysis *AnalysisResult) []string
```

#### core/functions_generator.go

```go
// Lógica de geração de functions pura
type FunctionsGenerator struct {
    settings *ProjectSettings
}

func (fg *FunctionsGenerator) Generate(analysisResults map[string]*AnalysisResult) (*Context, error)
func (fg *FunctionsGenerator) extractFunctions(codeAnalysis *AnalysisResult) []FileFunction
func (fg *FunctionsGenerator) buildSummary(functions []FileFunction) FunctionSummary
func (fg *FunctionsGenerator) groupByLanguage(functions []FileFunction) map[string][]FileFunction
```

#### core/stack_generator.go

```go
// Lógica de geração de stack pura
type StackGenerator struct {
    settings *ProjectSettings
}

func (sg *StackGenerator) Generate(analysisResults map[string]*AnalysisResult) (*Context, error)
func (sg *StackGenerator) processStackAnalysis(stackAnalysis *AnalysisResult) *StackData
func (sg *StackGenerator) prioritizeByImportance(items []Technology) []Technology
func (sg *StackGenerator) detectDeploymentStrategy(stack *StackData) string
```

#### core/context_generator.go

```go
// Lógica de geração de project-context pura
type ProjectContextGenerator struct {
    settings *ProjectSettings
    reader   FileReader
}

func (pcg *ProjectContextGenerator) Generate(analysisResults map[string]*AnalysisResult) (*Context, error)
func (pcg *ProjectContextGenerator) findDocumentation(fileAnalysis *AnalysisResult) []DocumentFile
func (pcg *ProjectContextGenerator) findConfiguration(fileAnalysis *AnalysisResult) []ConfigFile
func (pcg *ProjectContextGenerator) findContracts(fileAnalysis *AnalysisResult) []ContractFile
func (pcg *ProjectContextGenerator) buildProjectSummary(docs []DocumentFile, configs []ConfigFile) ProjectSummary
```

#### usecase/generate_contexts.go

```go
// Orquestração de geração
type GenerateContextsUseCase struct {
    fileTreeGen     *FileTreeGenerator
    functionsGen    *FunctionsGenerator
    stackGen        *StackGenerator
    contextGen      *ProjectContextGenerator
    writer          ContextWriter
    validator       SchemaValidator
    enabledContexts []string
}

func (uc *GenerateContextsUseCase) Execute(analysisResults map[string]*AnalysisResult) error
func (uc *GenerateContextsUseCase) generateContext(generator ContextGenerator, analysisResults map[string]*AnalysisResult) (*Context, error)
func (uc *GenerateContextsUseCase) saveAndValidate(context *Context) error
```

### internal/pipeline/

#### core/models.go

```go
// Tipos para pipeline
type PipelineConfig struct {
    ProjectPath     string
    SettingsPath    string
    EnabledContexts []string
    Parallel        bool
}

type ExecutionResult struct {
    Success      bool
    Duration     time.Duration
    StepResults  []StepResult
    Errors       []error
    GeneratedContexts []string
}

type StepResult struct {
    StepName  string
    Success   bool
    Duration  time.Duration
    Error     error
    Metadata  map[string]interface{}
}
```

#### core/ports.go

```go
// Interfaces para pipeline
type StepExecutor interface {
    Execute(ctx context.Context) (*StepResult, error)
    GetName() string
}

type ResultAggregator interface {
    Aggregate(results []StepResult) *ExecutionResult
}
```

#### core/orchestrator.go

```go
// Lógica de orquestração pura
type Orchestrator struct {
    config     PipelineConfig
    aggregator ResultAggregator
}

func (o *Orchestrator) Execute(ctx context.Context, steps []StepExecutor) (*ExecutionResult, error)
func (o *Orchestrator) executeSequential(ctx context.Context, steps []StepExecutor) []StepResult
func (o *Orchestrator) executeParallel(ctx context.Context, steps []StepExecutor) []StepResult
func (o *Orchestrator) handleStepFailure(step StepExecutor, err error) *StepResult
```

#### usecase/execute_pipeline.go

```go
// Orquestração completa do pipeline
type ExecutePipelineUseCase struct {
    settingsUC   *LoadSettingsUseCase
    patternsUC   *LoadPatternsUseCase
    discoveryUC  *DiscoverFilesUseCase
    analyzerUC   *AnalyzeProjectUseCase
    generatorUC  *GenerateContextsUseCase
    orchestrator *Orchestrator
}

func (uc *ExecutePipelineUseCase) Execute(ctx context.Context, config PipelineConfig) (*ExecutionResult, error)
func (uc *ExecutePipelineUseCase) buildSteps(config PipelineConfig) []StepExecutor
```

## 🔧 Adapters (Implementações de Infraestrutura)

### internal/adapters/filesystem/

```go
// settings_reader.go - Implementa SettingsReader
type YamlSettingsReader struct{}
func (r *YamlSettingsReader) ReadSettings(path string) (*ProjectSettings, error)

// pattern_loader.go - Implementa PatternLoader  
type FilePatternLoader struct{}
func (l *FilePatternLoader) LoadIgnorePatterns(path string) ([]string, error)
func (l *FilePatternLoader) LoadAnalyzePatterns(path string) ([]string, error)

// file_walker.go - Implementa FileWalker
type OSFileWalker struct{}
func (w *OSFileWalker) Walk(rootPath string, walkFn func(path string, info os.FileInfo) error) error

// file_reader.go - Implementa FileReader
type OSFileReader struct{}
func (r *OSFileReader) ReadFile(path string) ([]byte, error)
func (r *OSFileReader) GetFileInfo(path string) (os.FileInfo, error)
```

### internal/adapters/parsers/

```go
// code_parser.go - Implementa CodeParser
type MultiLanguageParser struct {
    parsers map[string]LanguageParser
}
func (p *MultiLanguageParser) ParseFile(path string, content []byte) ([]CodeElement, error)
func (p *MultiLanguageParser) SupportedExtensions() []string

// stack_detector.go - Implementa StackDetector
type FileBasedStackDetector struct{}
func (d *FileBasedStackDetector) DetectFromFile(path string, content []byte) ([]Technology, error)
func (d *FileBasedStackDetector) DetectFromProject(files []FileInfo) (TechStack, error)
```

### internal/adapters/storage/

```go
// json_writer.go - Implementa ContextWriter
type JsonFileWriter struct {
    basePath string
}
func (w *JsonFileWriter) WriteContext(context *Context, path string) error

// schema_validator.go - Implementa SchemaValidator
type JsonSchemaValidator struct {
    schemaPath string
}
func (v *JsonSchemaValidator) ValidateAgainstSchema(context *Context) error
func (v *JsonSchemaValidator) LoadSchema(contextType string) ([]byte, error)
```

## 🎯 CLI Integration

### cmd/cli/main.go

```go
package main

import (
    "context"
    "fmt"
    "os"
    
    "github.com/spf13/cobra"
    
    // Use cases
    "github.com/user/phengineer/internal/pipeline/usecase"
    "github.com/user/phengineer/internal/pipeline/core"
    
    // Adapters
    "github.com/user/phengineer/internal/adapters/filesystem"
    "github.com/user/phengineer/internal/adapters/parsers"
    "github.com/user/phengineer/internal/adapters/storage"
)

var (
    projectPath string
    contextType string
)

var rootCmd = &cobra.Command{
    Use:   "phengineer",
    Short: "PHEngineer - AI-powered development context generator",
}

var generateCmd = &cobra.Command{
    Use:   "generate",
    Short: "Generate project contexts",
    Run:   runGenerate,
}

func runGenerate(cmd *cobra.Command, args []string) {
    // Wire up dependencies (manual DI for MVP)
    config := core.PipelineConfig{
        ProjectPath:  projectPath,
        SettingsPath: projectPath + "/.phengineer/settings.yml",
        Parallel:     false,
    }
    
    // Create adapters
    settingsReader := &filesystem.YamlSettingsReader{}
    patternLoader := &filesystem.FilePatternLoader{}
    fileWalker := &filesystem.OSFileWalker{}
    fileReader := &filesystem.OSFileReader{}
    codeParser := &parsers.MultiLanguageParser{}
    stackDetector := &parsers.FileBasedStackDetector{}
    contextWriter := &storage.JsonFileWriter{BasePath: projectPath + "/.phengineer/context"}
    schemaValidator := &storage.JsonSchemaValidator{SchemaPath: "schemas"}
    
    // Build use case with dependencies
    pipelineUC := usecase.NewExecutePipelineUseCase(
        settingsReader,
        patternLoader,
        fileWalker,
        fileReader,
        codeParser,
        stackDetector,
        contextWriter,
        schemaValidator,
    )
    
    // Execute pipeline
    result, err := pipelineUC.Execute(context.Background(), config)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }
    
    fmt.Printf("✅ Pipeline completed in %v\n", result.Duration)
    fmt.Printf("Generated contexts: %v\n", result.GeneratedContexts)
}

func init() {
    generateCmd.Flags().StringVarP(&projectPath, "path", "p", ".", "Project path")
    generateCmd.Flags().StringVarP(&contextType, "context", "c", "", "Specific context to generate")
    
    rootCmd.AddCommand(generateCmd)
}

func main() {
    if err := rootCmd.Execute(); err != nil {
        os.Exit(1)
    }
}
```

## 🧪 Testing Strategy

### Core Testing (Unit Tests)

```go
// internal/analyzer/core/file_analyzer_test.go
func TestFileAnalyzer_Analyze(t *testing.T) {
    // Test pure business logic
    // No I/O dependencies
    // Fast execution
}

// internal/generator/core/filetree_generator_test.go  
func TestFileTreeGenerator_Generate(t *testing.T) {
    // Test generation logic
    // Mock analysis results
    // Validate JSON structure
}
```

### UseCase Testing (Integration Tests)

```go
// internal/analyzer/usecase/analyze_project_test.go
func TestAnalyzeProjectUseCase_Execute(t *testing.T) {
    // Test orchestration
    // Use mock adapters
    // Validate integration flow
}

// internal/pipeline/usecase/execute_pipeline_test.go
func TestExecutePipelineUseCase_Execute(t *testing.T) {
    // End-to-end pipeline test
    // Use testdata project
    // Validate complete flow
}
```

### Adapter Testing (Infrastructure Tests)

```go
// internal/adapters/filesystem/settings_reader_test.go
func TestYamlSettingsReader_ReadSettings(t *testing.T) {
    // Test actual file I/O
    // Use testdata files
    // Validate YAML parsing
}
```

## 🚀 Benefícios da Estrutura Core/UseCase

### Separação Clara de Responsabilidades

- **Core**: Lógica de negócio pura, testável sem I/O
- **UseCase**: Orquestração, coordenação entre cores
- **Adapters**: Implementações de infraestrutura

### Testabilidade Excepcional

- **Unit Tests**: Core isolado, execução rápida
- **Integration Tests**: UseCase com mocks
- **E2E Tests**: Pipeline completo com dados reais

### Evolução Natural para Lambda

- **Core + UseCase**: Reutilizáveis entre CLI e Lambda
- **Adapters**: Diferentes para cada ambiente
- **Zero reescrita** da lógica de negócio

### Manutenibilidade

- **Baixo acoplamento**: Core independe de infraestrutura
- **Alta coesão**: Cada módulo tem responsabilidade única
- **Substituição fácil**: Adapters são intercambiáveis

-----

> **Nota**: Esta estrutura segue o padrão core/usecase mantendo simplicidade e clareza. Cada camada tem responsabilidade bem definida e a arquitetura permite evolução natural conforme necessidade.
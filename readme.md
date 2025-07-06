# PHEngineer 🚀

Automatização Inteligente de Geração de Código com IA para Workflows no GitHub

---

## 📌 Visão Geral

O **PHEngineer** é um orquestrador desenvolvido em **Go** que integra **GitHub Issues** com **Inteligência Artificial (StackSpot AI)** para automatizar a geração de código, testes, documentação e validação seguindo conceitos como **Clean Architecture** e **Design Patterns**.

Através da abertura de Issues em um repositório, o sistema interpreta os pedidos dos usuários, executa múltiplos agentes de IA especializados e entrega um **Pull Request pronto para revisão** — tudo isso com execução paralela para máxima eficiência.

---

## 🎯 Objetivos do Projeto

- **Automatizar a geração de código** a partir de pedidos em linguagem natural ou estruturada
- **Padronizar entregas** seguindo arquiteturas e práticas previamente definidas
- **Reduzir o tempo de desenvolvimento** através de múltiplos agentes IA executando em paralelo
- **Permitir fácil integração** com qualquer repositório GitHub através de GitHub Actions
- **Introduzir Go** como linguagem backend na organização de forma prática e segura

---

## 🏗️ Arquitetura

### **Padrão Arquitetural**

- **Arquitetura Hexagonal (Ports & Adapters)**
  - Separação clara entre domínio e infraestrutura
  - Facilita testes e manutenção
  - Permite troca de implementações (HTTP ↔ MCP)

### **Camadas da Aplicação**

```
┌─────────────────────────────────────────────────────────────┐
│                    GitHub Issues (Trigger)                  │
└─────────────────────────┬───────────────────────────────────┘
                         │
┌─────────────────────────▼───────────────────────────────────┐
│                   Interface Layer                           │
│  ├─ GitHub Actions Handler                                  │
│  ├─ Issue Parser                                           │
│  └─ Response Formatter                                     │
└─────────────────────────┬───────────────────────────────────┘
                         │
┌─────────────────────────▼───────────────────────────────────┐
│                  Application Layer                          │
│  ├─ Orchestrator (Core Business Logic)                     │
│  ├─ Pipeline Manager                                       │
│  ├─ Agent Coordinator                                      │
│  └─ Result Consolidator                                    │
└─────────────────────────┬───────────────────────────────────┘
                         │
┌─────────────────────────▼───────────────────────────────────┐
│                    Domain Layer                             │
│  ├─ Feature (Entity)                                       │
│  ├─ Agent (Value Object)                                   │
│  ├─ Pipeline (Aggregate)                                   │
│  └─ Repository Interfaces (Ports)                          │
└─────────────────────────┬───────────────────────────────────┘
                         │
┌─────────────────────────▼───────────────────────────────────┐
│                 Infrastructure Layer                        │
│  ├─ AI Provider Adapters (HTTP/MCP)                        │
│  ├─ GitHub API Client                                      │
│  ├─ File System Handler                                    │
│  └─ Configuration Manager                                  │
└─────────────────────────────────────────────────────────────┘
```

---

## ⚙️ Tecnologias Utilizadas

### **Core Stack**

- **Go (Golang)** — Linguagem principal do projeto
- **Arquitetura Hexagonal** — Padrão arquitetural para desacoplamento
- **GitHub API** — Integração para Issues, Pull Requests e repositórios
- **GitHub Actions** — Automação e execução do pipeline CI/CD

### **Comunicação com IA**

- **StackSpot AI** — Geração de código, testes e documentação
- **mark3labs/mcp-go** — Implementação MCP para prototipagem
- **HTTP REST** — Fallback e compatibilidade
- **YAML** — Formato de comunicação estruturada

### **Ferramentas de Desenvolvimento**

- **Goroutines + Channels** — Execução paralela dos agentes
- **Testify** — Framework de testes
- **Golangci-lint** — Análise estática de código

---

## 🔗 Fluxo de Funcionamento

### **1. Trigger & Interpretação**

```
GitHub Issue → Issue Parser → Interpretation Agent → YAML Structure
```

### **2. Orquestração & Execução**

```
YAML → Pipeline Manager → Agent Coordinator → Parallel Execution
                                          ├─ Code Generation Agent
                                          ├─ Test Generation Agent
                                          └─ Documentation Agent
```

### **3. Consolidação & Entrega**

```
Parallel Results → Result Consolidator → GitHub PR Creation
```

---

## 🧩 Estrutura de Pastas

```
ia-first/
├── cmd/
│   └── orchestrator/
│       └── main.go                    # Ponto de entrada da aplicação
├── internal/
│   ├── domain/                        # Camada de Domínio
│   │   ├── entities/
│   │   │   ├── feature.go
│   │   │   └── pipeline.go
│   │   ├── valueobjects/
│   │   │   └── agent.go
│   │   └── repositories/              # Interfaces (Ports)
│   │       ├── ai_provider.go
│   │       └── github_client.go
│   ├── application/                   # Camada de Aplicação
│   │   ├── orchestrator.go
│   │   ├── pipeline_manager.go
│   │   ├── agent_coordinator.go
│   │   └── result_consolidator.go
│   ├── infrastructure/                # Camada de Infraestrutura
│   │   ├── adapters/
│   │   │   ├── stackspot_mcp.go      # Adapter MCP
│   │   │   ├── stackspot_http.go     # Adapter HTTP
│   │   │   └── github_api.go
│   │   ├── config/
│   │   │   └── config.go
│   │   └── parsers/
│   │       ├── yaml_parser.go
│   │       └── issue_parser.go
│   └── interfaces/                    # Camada de Interface
│       ├── handlers/
│       │   └── github_actions.go
│       └── formatters/
│           └── response_formatter.go
├── pkg/                              # Packages públicos
│   ├── github/
│   │   └── client.go
│   └── agents/
│       ├── code_generator.go
│       ├── test_generator.go
│       └── doc_generator.go
├── configs/
│   ├── config.yaml
│   └── agents.yaml
├── scripts/
│   └── setup.sh
├── tests/
│   ├── unit/
│   ├── integration/
│   └── fixtures/
├── docs/
│   ├── architecture.md
│   ├── deployment.md
│   └── contributing.md
├── .github/
│   └── workflows/
│       └── ia-first.yml
├── go.mod
├── go.sum
├── Dockerfile
├── docker-compose.yml
└── README.md
```

---

## 🔌 Arquitetura de Adaptadores

### **AI Provider Interface (Port)**

```go
type AIProvider interface {
    GenerateCode(ctx context.Context, req CodeGenerationRequest) (*CodeGenerationResponse, error)
    GenerateTests(ctx context.Context, req TestGenerationRequest) (*TestGenerationResponse, error)
    GenerateDocumentation(ctx context.Context, req DocumentationRequest) (*DocumentationResponse, error)
    InterpretIssue(ctx context.Context, req IssueInterpretationRequest) (*IssueInterpretationResponse, error)
}
```

### **Implementações (Adapters)**

- **StackSpotMCPAdapter** — Usa MCP para comunicação
- **StackSpotHTTPAdapter** — Usa HTTP REST para comunicação
- **MockAIAdapter** — Para testes unitários

### **Configuração Dinâmica**

```yaml
ai:
  provider: "stackspot"
  protocol: "mcp" # ou "http"
  config:
    endpoint: "..."
    credentials: "..."
```

---

## 🔑 Agentes Especializados

| Agente                      | Função                                       | Input               | Output        |
| --------------------------- | -------------------------------------------- | ------------------- | ------------- |
| **Issue Interpreter**       | Interpreta descrição e gera YAML estruturado | Issue Description   | YAML Config   |
| **Code Generator**          | Produz código-fonte seguindo padrões         | YAML + Context      | Source Code   |
| **Test Generator**          | Cria testes unitários/integração             | Code + Requirements | Test Files    |
| **Documentation Generator** | Produz documentação técnica                  | Code + Context      | Documentation |

---

## 📝 Exemplo de YAML Estruturado

```yaml
feature:
  name: "Adicionar endpoint de criação de usuário"
  description: "Endpoint REST para criação de usuários com validação"

architecture:
  pattern: "Clean Architecture"
  language: "Go"
  framework: "Gin"

design_patterns:
  - "Repository Pattern"
  - "DTO Pattern"
  - "Builder Pattern"

requirements:
  testing: true
  documentation: true
  validation: true

agents:
  - name: "code_generator"
    priority: 1
    config:
      template: "rest_endpoint"
      validation: true
  - name: "test_generator"
    priority: 2
    dependencies: ["code_generator"]
  - name: "doc_generator"
    priority: 3
    dependencies: ["code_generator"]
```

---

## 🚀 Estratégia de Implementação

### **Fase 1: Foundation (2-3 semanas)**

- [ ] Estrutura hexagonal básica
- [ ] Interfaces e contratos
- [ ] Configuração e parsers
- [ ] Testes unitários base

### **Fase 2: Core Features (3-4 semanas)**

- [ ] Orquestrador principal
- [ ] Pipeline manager
- [ ] Agent coordinator
- [ ] Adapter MCP (prototipagem)

### **Fase 3: Integration (2-3 semanas)**

- [ ] GitHub API integration
- [ ] StackSpot AI integration
- [ ] GitHub Actions workflow
- [ ] Testes de integração

### **Fase 4: Production Ready (1-2 semanas)**

- [ ] Adapter HTTP (fallback)
- [ ] Monitoring e logs
- [ ] Documentação completa
- [ ] Deploy e CI/CD

---

## 🔄 Estratégia de Migração

### **Flexibilidade de Protocolos**

```go
// Factory pattern para criação de adapters
func NewAIProvider(config Config) AIProvider {
    switch config.Protocol {
    case "mcp":
        return NewStackSpotMCPAdapter(config)
    case "http":
        return NewStackSpotHTTPAdapter(config)
    default:
        return NewMockAIAdapter()
    }
}
```

### **Plano de Migração**

1. **Atual**: `mark3labs/mcp-go` para prototipagem
2. **Fallback**: HTTP adapter como backup
3. **Futuro**: Migração para SDK oficial MCP quando disponível
4. **Flexibilidade**: Troca de protocolo via configuração

---

## 🛡️ Segurança e Confiabilidade

### **Validação de Entrada**

- Sanitização de Issues do GitHub
- Validação de YAML estruturado
- Rate limiting para APIs

### **Tratamento de Erros**

- Circuit breaker para APIs externas
- Retry com backoff exponencial
- Fallback para HTTP em caso de falha MCP

### **Monitoramento**

- Métricas de performance
- Logs estruturados
- Health checks

---

## 🧠 Próximos Passos

### **Sprint 1**

1. [ ] Setup do projeto Go com estrutura hexagonal
2. [ ] Implementação das interfaces core
3. [ ] Parser de Issues e YAML
4. [ ] Testes unitários básicos

### **Sprint 2**

1. [ ] Orquestrador e pipeline manager
2. [ ] Agent coordinator com goroutines
3. [ ] Adapter MCP básico
4. [ ] Mock do StackSpot AI

### **Sprint 3**

1. [ ] Integração real com StackSpot AI
2. [ ] GitHub API client
3. [ ] Testes de integração
4. [ ] Configuração dinâmica

---

## 🎯 Benefícios Esperados

### **Técnicos**

- **Performance**: Execução paralela de agentes
- **Escalabilidade**: Arquitetura modular e desacoplada
- **Manutenibilidade**: Código limpo e testável
- **Flexibilidade**: Múltiplos protocolos de comunicação

### **Organizacionais**

- **Introdução do Go**: Showcase da linguagem
- **Automação**: Redução de trabalho manual
- **Padronização**: Código consistente
- **Inovação**: Uso de tecnologias emergentes (MCP)

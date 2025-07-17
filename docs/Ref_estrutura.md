# Estrutura Final do Projeto Go CLI - Wire Ready

## Estrutura de Diretórios Completa

```
myapp/
├── cmd/                                        # Múltiplos entrypoints
│   ├── cli/
│   │   └── main.go                            # Entrypoint CLI
│   ├── api/
│   │   └── main.go                            # Entrypoint API REST
│   ├── lambda/
│   │   └── main.go                            # Entrypoint AWS Lambda
│   └── worker/
│       └── main.go                            # Entrypoint Background Worker
│
├── internal/                                   # Código privado da aplicação
│   ├── domain/                                # ═══ DOMAIN LAYER ═══
│   │   │                                      # Regras de negócio puras
│   │   ├── user/
│   │   │   ├── user.go                        # Entidade User com regras de negócio
│   │   │   ├── repository.go                  # Interface Repository (definida no domain)
│   │   │   ├── service.go                     # Domain Service (validações complexas)
│   │   │   └── errors.go                      # Erros específicos do domínio
│   │   │
│   │   ├── order/
│   │   │   ├── order.go                       # Entidade Order
│   │   │   ├── item.go                        # Entidade Item
│   │   │   ├── repository.go                  # Interface Repository
│   │   │   ├── pricing_service.go             # Domain Service para cálculos
│   │   │   └── status.go                      # Value Objects (enum status)
│   │   │
│   │   └── shared/                            # Value Objects compartilhados
│   │       ├── email.go                       # Value Object Email
│   │       ├── money.go                       # Value Object Money
│   │       └── pagination.go                  # Value Objects de paginação
│   │
│   ├── application/                           # ═══ APPLICATION LAYER ═══
│   │   │                                      # Use Cases / Application Services
│   │   ├── user/                              # Use cases de usuário
│   │   │   ├── register_user_use_case.go      # UC: Registrar novo usuário
│   │   │   ├── activate_user_use_case.go      # UC: Ativar usuário
│   │   │   ├── get_user_use_case.go           # UC: Buscar usuário
│   │   │   ├── update_user_use_case.go        # UC: Atualizar dados
│   │   │   └── dto.go                         # DTOs de entrada/saída
│   │   │
│   │   ├── order/                             # Use cases de pedidos
│   │   │   ├── create_order_use_case.go       # UC: Criar pedido
│   │   │   ├── process_order_use_case.go      # UC: Processar pagamento
│   │   │   ├── cancel_order_use_case.go       # UC: Cancelar pedido
│   │   │   ├── list_orders_use_case.go        # UC: Listar pedidos
│   │   │   └── dto.go                         # DTOs específicos
│   │   │
│   │   └── interfaces/                        # Interfaces definidas pela application
│   │       ├── email_sender.go                # Interface para envio de email
│   │       ├── payment_gateway.go             # Interface para gateway de pagamento
│   │       ├── notification_service.go        # Interface para notificações
│   │       └── file_storage.go                # Interface para storage
│   │
│   ├── infrastructure/                        # ═══ INFRASTRUCTURE LAYER ═══
│   │   │                                      # Implementações técnicas
│   │   ├── adapters/                          # Adapters para sistemas externos
│   │   │   ├── database/                      # Adapters de banco de dados
│   │   │   │   ├── sqlc/                      # ┌─ Código GERADO pelo SQLC
│   │   │   │   │   ├── db.go                  # │ Interface Queries gerada
│   │   │   │   │   ├── models.go              # │ Structs baseadas no schema
│   │   │   │   │   ├── user.sql.go            # │ Funções de query de usuário
│   │   │   │   │   ├── order.sql.go           # │ Funções de query de pedidos
│   │   │   │   │   └── queries/               # │ Arquivos SQL fonte
│   │   │   │   │       ├── schema.sql         # │ Schema do banco (DDL)
│   │   │   │   │       ├── user.sql           # │ Queries SQL de usuário
│   │   │   │   │       └── order.sql          # │ Queries SQL de pedidos
│   │   │   │   │                              # └─
│   │   │   │   ├── sqlc_user_repository.go    # Adapter: implementa domain user.Repository
│   │   │   │   ├── sqlc_order_repository.go   # Adapter: implementa domain order.Repository
│   │   │   │   ├── memory_user_repository.go  # Adapter: implementação em memória (testes)
│   │   │   │   ├── connection.go              # Configuração e pool de conexões
│   │   │   │   └── migrations/                # Migrações do banco
│   │   │   │       ├── 001_create_users.up.sql
│   │   │   │       ├── 001_create_users.down.sql
│   │   │   │       ├── 002_create_orders.up.sql
│   │   │   │       └── 002_create_orders.down.sql
│   │   │   │
│   │   │   ├── email/                         # Adapters de email
│   │   │   │   ├── smtp_sender.go             # Implementação SMTP
│   │   │   │   ├── sendgrid_sender.go         # Implementação SendGrid
│   │   │   │   └── mock_sender.go             # Mock para testes
│   │   │   │
│   │   │   ├── payment/                       # Adapters de pagamento
│   │   │   │   ├── stripe_gateway.go          # Implementação Stripe
│   │   │   │   ├── paypal_gateway.go          # Implementação PayPal
│   │   │   │   └── mock_gateway.go            # Mock para testes
│   │   │   │
│   │   │   ├── notification/                  # Adapters de notificação
│   │   │   │   ├── slack_notifier.go          # Implementação Slack
│   │   │   │   ├── webhook_notifier.go        # Implementação Webhook
│   │   │   │   └── console_notifier.go        # Implementação console (dev)
│   │   │   │
│   │   │   └── storage/                       # Adapters de storage
│   │   │       ├── s3_storage.go              # Implementação AWS S3
│   │   │       ├── local_storage.go           # Implementação local
│   │   │       └── memory_storage.go          # Implementação em memória
│   │   │
│   │   ├── config/                            # Configuração da aplicação
│   │   │   ├── config.go                      # Structs de configuração
│   │   │   ├── loader.go                      # Carregamento de arquivos YAML
│   │   │   └── validator.go                   # Validação de configurações
│   │   │
│   │   ├── cli/                               # Interface CLI
│   │   │   ├── root.go                        # Comando raiz
│   │   │   ├── user_commands.go               # Comandos de usuário
│   │   │   ├── order_commands.go              # Comandos de pedidos
│   │   │   └── migration_commands.go          # Comandos de migração
│   │   │
│   │   ├── http/                              # Interface HTTP/REST
│   │   │   ├── server.go                      # Configuração do servidor
│   │   │   ├── handlers/
│   │   │   │   ├── user_handler.go            # Handlers de usuário
│   │   │   │   └── order_handler.go           # Handlers de pedidos
│   │   │   ├── middleware/
│   │   │   │   ├── auth.go                    # Middleware de autenticação
│   │   │   │   ├── logging.go                 # Middleware de logging
│   │   │   │   └── cors.go                    # Middleware CORS
│   │   │   └── routes/
│   │   │       └── routes.go                  # Definição das rotas
│   │   │
│   │   ├── lambda/                            # Interface AWS Lambda
│   │   │   ├── handlers/
│   │   │   │   ├── user_handler.go            # Lambda handlers de usuário
│   │   │   │   └── order_handler.go           # Lambda handlers de pedidos
│   │   │   └── response.go                    # Formatação de respostas Lambda
│   │   │
│   │   ├── worker/                            # Interface Background Worker
│   │   │   ├── worker.go                      # Worker principal
│   │   │   ├── jobs/
│   │   │   │   ├── email_job.go               # Job de envio de email
│   │   │   │   └── cleanup_job.go             # Job de limpeza
│   │   │   └── queue/
│   │   │       └── queue.go                   # Abstração de fila
│   │   │
│   │   └── container/                         # ═══ DEPENDENCY INJECTION ═══
│   │       ├── providers.go                   # Provider functions (Wire-ready)
│   │       ├── container.go                   # Container manual (atual)
│   │       ├── wire.go                        # Configuração Wire (futuro)
│   │       └── wire_gen.go                    # Código gerado pelo Wire
│   │
│   └── pkg/                                   # Pacotes exportáveis
│       ├── logger/
│       │   ├── logger.go                      # Interface de logging
│       │   └── levels.go                      # Níveis de log
│       ├── validator/
│       │   ├── validator.go                   # Validações genéricas
│       │   └── rules.go                       # Regras customizadas
│       └── httpclient/
│           ├── client.go                      # Cliente HTTP reutilizável
│           └── retry.go                       # Lógica de retry
│
├── configs/                                   # Templates de configuração
│   ├── config.example.yaml                   # Template principal
│   ├── cli.yaml                              # Config específica CLI
│   ├── api.yaml                              # Config específica API
│   ├── lambda.yaml                           # Config específica Lambda
│   ├── worker.yaml                           # Config específica Worker
│   ├── development.yaml                      # Config desenvolvimento
│   ├── staging.yaml                          # Config staging
│   └── production.yaml                       # Template produção
│
├── deployments/                              # Configurações de deploy
│   ├── docker/
│   │   ├── Dockerfile.api                    # Docker para API
│   │   ├── Dockerfile.worker                 # Docker para Worker
│   │   ├── Dockerfile.cli                    # Docker para CLI
│   │   └── docker-compose.yml                # Orquestração local
│   ├── kubernetes/
│   │   ├── api-deployment.yaml               # Deployment Kubernetes API
│   │   ├── worker-deployment.yaml            # Deployment Kubernetes Worker
│   │   └── configmap.yaml                    # ConfigMap
│   └── terraform/
│       ├── lambda.tf                         # Terraform para Lambda
│       ├── api.tf                            # Terraform para API
│       └── variables.tf                      # Variáveis Terraform
│
├── scripts/                                  # Scripts auxiliares
│   ├── build/
│   │   ├── build-cli.sh                      # Build CLI
│   │   ├── build-api.sh                      # Build API
│   │   ├── build-lambda.sh                   # Build Lambda
│   │   ├── build-worker.sh                   # Build Worker
│   │   └── build-all.sh                      # Build todos
│   ├── dev/
│   │   ├── setup.sh                          # Setup ambiente desenvolvimento
│   │   ├── migrate.sh                        # Rodar migrações
│   │   └── generate.sh                       # Gerar código (SQLC, Wire)
│   └── deploy/
│       ├── deploy-api.sh                     # Deploy API
│       ├── deploy-lambda.sh                  # Deploy Lambda
│       └── rollback.sh                       # Rollback
│
├── docs/                                     # Documentação
│   ├── README.md                             # Documentação principal
│   ├── architecture.md                      # Documentação da arquitetura
│   ├── api.md                               # Documentação da API
│   ├── cli.md                               # Documentação do CLI
│   └── deployment.md                        # Guia de deployment
│
├── .github/                                  # GitHub Actions
│   └── workflows/
│       ├── ci.yml                           # Pipeline de CI
│       ├── cd-api.yml                       # Pipeline CD API
│       ├── cd-lambda.yml                    # Pipeline CD Lambda
│       └── release.yml                      # Pipeline de release
│
├── sqlc.yaml                                # Configuração do SQLC
├── .gitignore                               # Arquivos ignorados pelo Git
├── .env.example                             # Template de variáveis de ambiente
├── Makefile                                 # Comandos principais do projeto
├── Dockerfile                               # Dockerfile padrão (se necessário)
├── docker-compose.yml                       # Desenvolvimento local
├── go.mod                                   # Dependências do Go
└── go.sum                                   # Checksums das dependências
```

## Arquivos de Exemplo Principais

### **Makefile**
```makefile
# Build commands
.PHONY: build-all build-cli build-api build-lambda build-worker
build-all: build-cli build-api build-lambda build-worker

build-cli:
	CGO_ENABLED=0 go build -o bin/myapp-cli cmd/cli/main.go

build-api:
	CGO_ENABLED=0 go build -o bin/myapp-api cmd/api/main.go

build-lambda:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/myapp-lambda cmd/lambda/main.go
	zip bin/myapp-lambda.zip bin/myapp-lambda

build-worker:
	CGO_ENABLED=0 go build -o bin/myapp-worker cmd/worker/main.go

# Development commands
.PHONY: dev setup generate migrate test
dev:
	go run cmd/cli/main.go --config configs/development.yaml

setup:
	./scripts/dev/setup.sh

generate:
	./scripts/dev/generate.sh

migrate:
	./scripts/dev/migrate.sh

test:
	go test ./...

# Docker commands
.PHONY: docker-api docker-worker docker-cli
docker-api:
	docker build -f deployments/docker/Dockerfile.api -t myapp-api .

docker-worker:
	docker build -f deployments/docker/Dockerfile.worker -t myapp-worker .

docker-cli:
	docker build -f deployments/docker/Dockerfile.cli -t myapp-cli .

# Wire generation
.PHONY: wire
wire:
	cd internal/container && go generate

# Clean
.PHONY: clean
clean:
	rm -rf bin/
	rm -f internal/container/wire_gen.go
```

### **SQLC Configuration** (sqlc.yaml)
```yaml
version: "2"
sql:
  - engine: "mysql"
    queries: "internal/infrastructure/adapters/database/sqlc/queries"
    schema: "internal/infrastructure/adapters/database/sqlc/queries/schema.sql"
    gen:
      go:
        package: "sqlc"
        out: "internal/infrastructure/adapters/database/sqlc"
        sql_package: "database/sql"
        emit_json_tags: true
        emit_db_tags: true
        emit_prepared_queries: false
        emit_interface: true
        emit_exact_table_names: false
```

### **.gitignore**
```gitignore
# Binários
/bin/
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with `go test -c`
*.test

# Output of the go coverage tool
*.out

# Go workspace file
go.work

# Dependency directories
vendor/

# IDEs
.vscode/
.idea/
*.swp
*.swo
*~

# OS generated files
.DS_Store
.DS_Store?
._*
.Spotlight-V100
.Trashes
ehthumbs.db
Thumbs.db

# Environment variables
.env
.env.local
.env.*.local

# Config files with secrets
config.yaml
*.local.yaml

# Logs
*.log

# Coverage
coverage.html
coverage.out

# Wire generated (opcional - alguns preferem versionar)
# internal/container/wire_gen.go

# Build artifacts
dist/
release/
```

## Principais Características

### **🎯 Multi-Entrypoint Ready**
- CLI, API, Lambda, Worker compartilham mesma lógica
- Builds independentes e otimizados
- Configurações específicas por entrypoint

### **🏗️ Clean Architecture**
- Domain independente de infraestrutura
- Application orquestra domain + infrastructure
- Infrastructure implementa interfaces

### **🔧 SQLC Integration**
- Código type-safe gerado automaticamente
- Adapters fazem bridge SQLC ↔ Domain
- Queries organizadas por contexto

### **📦 Wire Ready**
- Providers preparados para migração
- Container manual funcional agora
- Migração simples quando necessário

### **🧪 Testing Friendly**
- Mocks em cada adapter
- Interfaces bem definidas
- Containers de teste isolados

### **🚀 Production Ready**
- Docker multi-stage builds
- CI/CD pipelines preparados
- Configurações por ambiente
- Monitoring e observabilidade

Esta estrutura te dá a base sólida para começar e escalar conforme necessário! �


```md
myapp/
├── cmd/
│   ├── cli/main.go           # ✅ Entrypoint principal
│   └── api/main.go           # ✅ Para demonstrar versatilidade
├── internal/
│   ├── domain/
│   │   └── user/
│   │       ├── user.go       # ✅ Entidade simples
│   │       └── repository.go # ✅ Interface
│   ├── application/
│   │   └── user/
│   │       ├── register_use_case.go # ✅ Use case simples
│   │       └── dto.go        # ✅ DTOs básicos
│   └── infrastructure/
│       ├── adapters/
│       │   └── database/
│       │       └── memory_user_repository.go # ✅ Simples para começar
│       ├── config/
│       │   ├── config.go     # ✅ Estrutura básica
│       │   └── loader.go     # ✅ YAML loading
│       ├── cli/
│       │   └── commands.go   # ✅ Comandos CLI
│       ├── http/             # ✅ Para API
│       │   └── handlers/
│       └── container/
│           ├── providers.go  # ✅ Wire-ready
│           └── container.go  # ✅ Manual por agora
├── configs/
│   ├── config.example.yaml  # ✅ Template
│   ├── cli.yaml             # ✅ Config CLI
│   └── api.yaml             # ✅ Config API
├── Makefile                 # ✅ Comandos básicos
├── go.mod
└── README.md               # ✅ Documentação clara
```
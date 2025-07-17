
```text
myapp/
├── cmd/
│   ├── cli/main.go           # ✅ Entrypoint principal
│   └── api/main.go           # ✅ Para demonstrar versatilidade
├── internal/
│   ├── domain/
│   │   └── agent/
│   │       ├── agent.go       # ✅ Entidade simples
│   │       └── repository.go # ✅ Interface
│   ├── application/
│   │   └── agent/
│   │       ├── comunication_use_case.go # ✅ Use case simples
│   │       └── dto.go        # ✅ DTOs básicos
│   └── infrastructure/
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

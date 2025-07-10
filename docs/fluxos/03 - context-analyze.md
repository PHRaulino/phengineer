# PHEngineer - Context Generation Flow

## 📋 Overview

O fluxo de geração de contexto é responsável por analisar profundamente um projeto e criar contextos estruturados em JSON que servem como base de conhecimento para os agentes IA. O processo combina análise determinística via código com interpretação inteligente via IA, utilizando atualizações incrementais baseadas em Git para máxima eficiência.

## 🎯 Objetivos

- **Compreender o projeto** - Estrutura, tecnologias, padrões e código
- **Gerar contextos JSON** - Estrutura universal para qualquer LLM
- **Atualização inteligente** - Incremental baseado em mudanças Git
- **Filtrar relevância** - Focar em informações úteis, remover ruído
- **Máxima compatibilidade** - JSON funciona com GPT, Claude, Gemini, OSS

## 📊 Tipos de Projeto e Contextos

### Contextos Base (Todos os Projetos)
Sempre gerados independente do tipo:

| Contexto | Método | Atualização | Schema |
|----------|--------|-------------|--------|
| `file-tree.json` | 🔧 Code | Completa | FileTree |
| `statistics.json` | 🔧 Code | Completa | Statistics |
| `stack.json` | 🤖 IA | Completa | Stack |
| `architecture.json` | 🤖 IA | Completa | Architecture |
| `functions.json` | 🤖 IA | **Incremental** | Functions |
| `project-context.json` | 🤖 IA | **Incremental** | ProjectContext |

### Contextos Específicos por Tipo

| Tipo de Projeto | Contextos Específicos | Total |
|-----------------|----------------------|-------|
| **Frontend** | `components.json`, `styling.json` | **8** |
| **Backend** | Nenhum específico | **6** |
| **Infrastructure** | Nenhum específico | **6** |
| **Library** | Nenhum específico | **6** |
| **Generic** | Nenhum específico | **6** |

## 🔄 Modo de Análise

### Estratégia Baseada em Git Commit
```yaml
# .phengineer/config.yml
analysis:
  last_commit_hash: ""              # Vazio = análise completa
  incremental_contexts: ["functions", "project-context"]
```

### Lógica de Determinação
- **Análise Completa**: Campo `last_commit_hash` vazio ou inexistente
- **Análise Incremental**: Campo preenchido + Git diff para detectar mudanças

## 📋 JSON Schemas dos Contextos

### 1. file-tree.json
```json
{
  "metadata": {
    "generated_at": "2025-01-15T10:30:00Z",
    "total_files": 156,
    "total_directories": 23,
    "max_depth": 5
  },
  "structure": {
    "root": ".",
    "directories": [
      {
        "path": "src/",
        "type": "source",
        "files_count": 45,
        "subdirectories": ["components/", "services/", "utils/"],
        "importance": "high"
      },
      {
        "path": "tests/",
        "type": "test",
        "files_count": 23,
        "subdirectories": ["unit/", "integration/"],
        "importance": "medium"
      }
    ],
    "important_files": [
      {
        "path": "src/main.py",
        "type": "entry_point",
        "size": "2.3KB",
        "importance": "critical"
      },
      {
        "path": "requirements.txt",
        "type": "dependency",
        "size": "1.1KB", 
        "importance": "high"
      }
    ]
  },
  "conventions": {
    "naming": "snake_case",
    "organization": "by_feature",
    "test_pattern": "test_*.py"
  }
}
```

### 2. statistics.json
```json
{
  "metadata": {
    "generated_at": "2025-01-15T10:30:00Z",
    "analysis_scope": "all_files",
    "languages_detected": ["python", "typescript"]
  },
  "summary": {
    "total_files": 89,
    "total_lines": 12450,
    "total_functions": 156,
    "total_classes": 23,
    "total_modules": 34
  },
  "by_language": [
    {
      "language": "python",
      "files": 67,
      "lines": 9200,
      "functions": 123,
      "classes": 18,
      "modules": 25,
      "complexity_avg": 3.2
    },
    {
      "language": "typescript", 
      "files": 22,
      "lines": 3250,
      "functions": 33,
      "classes": 5,
      "modules": 9,
      "complexity_avg": 2.8
    }
  ],
  "by_directory": [
    {
      "path": "src/services/",
      "functions": 45,
      "classes": 8,
      "avg_file_size": 180
    },
    {
      "path": "src/utils/",
      "functions": 28,
      "classes": 2, 
      "avg_file_size": 95
    }
  ],
  "metrics": {
    "avg_functions_per_file": 1.75,
    "avg_lines_per_function": 15.2,
    "largest_file": {
      "path": "src/services/user_service.py",
      "lines": 245,
      "functions": 12
    }
  }
}
```

### 3. stack.json
```json
{
  "metadata": {
    "generated_at": "2025-01-15T10:30:00Z",
    "detection_method": "ai_analysis",
    "confidence": "high"
  },
  "languages": [
    {
      "name": "Python",
      "version": "3.11",
      "primary": true,
      "usage_percentage": 85
    },
    {
      "name": "TypeScript",
      "version": "4.9",
      "primary": false,
      "usage_percentage": 15
    }
  ],
  "frameworks": [
    {
      "name": "FastAPI",
      "version": "0.68.0",
      "category": "web_framework",
      "language": "python",
      "purpose": "REST API development",
      "essential": true
    },
    {
      "name": "React",
      "version": "18.2.0",
      "category": "frontend_framework", 
      "language": "typescript",
      "purpose": "User interface",
      "essential": true
    }
  ],
  "databases": [
    {
      "name": "PostgreSQL",
      "version": "13",
      "type": "relational",
      "orm": "SQLAlchemy",
      "purpose": "primary_storage"
    }
  ],
  "key_dependencies": [
    {
      "name": "pydantic",
      "version": "1.10.0",
      "category": "validation",
      "purpose": "Data validation and serialization",
      "essential": true,
      "language": "python"
    },
    {
      "name": "tailwindcss",
      "version": "3.3.0",
      "category": "styling",
      "purpose": "CSS utility framework",
      "essential": true,
      "language": "typescript"
    }
  ],
  "development_tools": [
    {
      "name": "pytest",
      "category": "testing",
      "purpose": "Python testing framework"
    },
    {
      "name": "black",
      "category": "formatting", 
      "purpose": "Code formatting"
    }
  ],
  "deployment": {
    "containerization": "Docker",
    "orchestration": null,
    "cloud_services": ["AWS"],
    "ci_cd": "GitHub Actions"
  }
}
```

### 4. architecture.json
```json
{
  "metadata": {
    "generated_at": "2025-01-15T10:30:00Z",
    "detection_method": "ai_analysis",
    "analysis_depth": "comprehensive"
  },
  "pattern": {
    "primary": "Clean Architecture",
    "secondary": ["Repository Pattern", "Dependency Injection"],
    "confidence": "high"
  },
  "layers": [
    {
      "name": "Domain",
      "path": "src/domain/",
      "purpose": "Business entities and rules",
      "dependencies": []
    },
    {
      "name": "Use Cases",
      "path": "src/usecases/",
      "purpose": "Application business rules",
      "dependencies": ["domain"]
    },
    {
      "name": "Interface Adapters",
      "path": "src/adapters/",
      "purpose": "Controllers, presenters, gateways",
      "dependencies": ["usecases", "domain"]
    },
    {
      "name": "Infrastructure",
      "path": "src/infrastructure/",
      "purpose": "External interfaces",
      "dependencies": ["adapters", "usecases", "domain"]
    }
  ],
  "data_models": [
    {
      "name": "User",
      "type": "entity",
      "file": "src/domain/entities/user.py",
      "fields": [
        {"name": "id", "type": "UUID", "required": true},
        {"name": "email", "type": "str", "required": true},
        {"name": "password_hash", "type": "str", "required": true},
        {"name": "created_at", "type": "datetime", "required": true}
      ],
      "relationships": [
        {"type": "one_to_many", "target": "Order", "field": "orders"}
      ]
    },
    {
      "name": "Order",
      "type": "entity", 
      "file": "src/domain/entities/order.py",
      "fields": [
        {"name": "id", "type": "UUID", "required": true},
        {"name": "user_id", "type": "UUID", "required": true},
        {"name": "total", "type": "Decimal", "required": true},
        {"name": "status", "type": "OrderStatus", "required": true}
      ],
      "relationships": [
        {"type": "many_to_one", "target": "User", "field": "user"}
      ]
    }
  ],
  "principles": [
    {
      "name": "Dependency Inversion",
      "applied": true,
      "evidence": "Interfaces defined in domain, implementations in infrastructure"
    },
    {
      "name": "Single Responsibility",
      "applied": true,
      "evidence": "Each class has single, well-defined purpose"
    }
  ],
  "design_patterns": [
    {
      "name": "Repository Pattern",
      "location": "src/adapters/repositories/",
      "purpose": "Data access abstraction"
    },
    {
      "name": "Factory Pattern",
      "location": "src/infrastructure/factories/",
      "purpose": "Object creation"
    }
  ],
  "conventions": {
    "naming": "snake_case for files, PascalCase for classes",
    "organization": "by_layer_then_feature",
    "dependency_direction": "inward_to_domain"
  }
}
```

### 5. functions.json (Incremental)
```json
{
  "metadata": {
    "generated_at": "2025-01-15T10:30:00Z",
    "total_files": 47,
    "total_functions": 156,
    "last_commit_hash": "a1b2c3d4e5f6",
    "incremental": true
  },
  "files": [
    {
      "file_path": "src/services/user_service.py",
      "file_hash": "abc123def456",
      "language": "python",
      "last_modified": "2025-01-15T09:15:00Z",
      "functions": [
        {
          "name": "create_user",
          "signature": "create_user(user_data: CreateUserRequest) -> User",
          "line_number": 15,
          "visibility": "public",
          "async": false,
          "purpose": "Creates a new user account with validation and password hashing",
          "params": [
            {
              "name": "user_data",
              "type": "CreateUserRequest",
              "description": "User registration data including email, password, and profile information"
            }
          ],
          "returns": {
            "type": "User",
            "description": "Created user object with generated ID and timestamps"
          },
          "dependencies": [
            "bcrypt.hashpw",
            "UserRepository.save",
            "EmailValidator.validate"
          ],
          "side_effects": [
            "database_write",
            "password_hashing",
            "email_validation"
          ],
          "raises": [
            "ValidationError",
            "DuplicateEmailError"
          ]
        },
        {
          "name": "get_user_by_id", 
          "signature": "get_user_by_id(user_id: UUID) -> Optional[User]",
          "line_number": 45,
          "visibility": "public",
          "async": false,
          "purpose": "Retrieves user by ID with error handling",
          "params": [
            {
              "name": "user_id",
              "type": "UUID",
              "description": "Unique user identifier"
            }
          ],
          "returns": {
            "type": "Optional[User]",
            "description": "User object if found, None otherwise"
          },
          "dependencies": [
            "UserRepository.get_by_id"
          ],
          "side_effects": [
            "database_read"
          ],
          "raises": [
            "DatabaseError"
          ]
        }
      ]
    }
  ]
}
```

### 6. project-context.json (Incremental + Relevância)
```json
{
  "metadata": {
    "generated_at": "2025-01-15T10:30:00Z",
    "total_files_analyzed": 23,
    "relevant_files": 15,
    "irrelevant_files": 8,
    "last_commit_hash": "a1b2c3d4e5f6",
    "incremental": true
  },
  "files": [
    {
      "file_path": "README.md",
      "file_hash": "abc123def456",
      "file_type": "documentation",
      "last_modified": "2025-01-10T15:20:00Z",
      "size": "3.2KB",
      "relevant": true,
      "relevance_score": 0.95,
      "relevance_reason": "Main project documentation with comprehensive setup instructions and architecture overview",
      "content_summary": "Complete project setup guide covering installation, configuration, API usage, and deployment procedures",
      "key_sections": [
        "Installation", 
        "Configuration",
        "API Documentation",
        "Deployment"
      ]
    },
    {
      "file_path": "docker-compose.yml",
      "file_hash": "def789ghi012",
      "file_type": "configuration",
      "last_modified": "2025-01-12T09:30:00Z",
      "size": "1.8KB",
      "relevant": true,
      "relevance_score": 0.90,
      "relevance_reason": "Active configuration for local development environment with all required services",
      "content_summary": "Development environment setup with PostgreSQL, Redis, and API service configuration",
      "services": [
        "postgresql",
        "redis", 
        "api",
        "frontend"
      ]
    },
    {
      "file_path": "docs/TEMPLATE.md",
      "file_hash": "ghi345jkl678",
      "file_type": "documentation",
      "last_modified": "2024-08-15T10:00:00Z",
      "size": "0.8KB",
      "relevant": false,
      "relevance_score": 0.15,
      "relevance_reason": "Generic template file with placeholder content, not customized for this project",
      "content_summary": "Boilerplate template with standard sections and placeholder text"
    },
    {
      "file_path": "openapi.yml",
      "file_hash": "jkl901mno234",
      "file_type": "contract",
      "last_modified": "2025-01-14T14:45:00Z",
      "size": "12.5KB",
      "relevant": true,
      "relevance_score": 0.98,
      "relevance_reason": "Current API specification defining all 15 endpoints with detailed schemas",
      "content_summary": "REST API contract with authentication, CRUD operations for users and orders",
      "endpoints_count": 15,
      "schemas_count": 8,
      "security_schemes": ["BearerAuth"]
    }
  ],
  "summary": {
    "project_purpose": "E-commerce API platform with microservices architecture for handling product catalog, user management, and order processing",
    "key_configurations": [
      {
        "file": "docker-compose.yml",
        "purpose": "Local development environment"
      },
      {
        "file": ".env.example",
        "purpose": "Environment variables template"
      }
    ],
    "active_contracts": [
      {
        "file": "openapi.yml",
        "type": "REST API",
        "endpoints": 15
      },
      {
        "file": "schema.graphql", 
        "type": "GraphQL",
        "operations": 8
      }
    ],
    "relevant_documentation": [
      {
        "file": "README.md",
        "focus": "Setup and usage"
      },
      {
        "file": "docs/architecture.md",
        "focus": "System design"
      }
    ],
    "deployment_strategy": "Docker containers with PostgreSQL database and Redis caching layer",
    "external_integrations": [
      "Payment gateway API",
      "Email service",
      "Analytics platform"
    ]
  }
}
```

### 7. components.json (Frontend)
```json
{
  "metadata": {
    "generated_at": "2025-01-15T10:30:00Z",
    "framework": "React",
    "typescript": true,
    "total_components": 34
  },
  "hierarchy": {
    "layout": [
      {
        "name": "App",
        "path": "src/App.tsx",
        "type": "root",
        "children": ["Layout", "Router"],
        "props": [],
        "state_management": "none"
      },
      {
        "name": "Layout",
        "path": "src/components/Layout/Layout.tsx",
        "type": "layout",
        "children": ["Header", "Sidebar", "Main", "Footer"],
        "props": [
          {"name": "children", "type": "ReactNode", "required": true}
        ],
        "state_management": "local"
      }
    ],
    "pages": [
      {
        "name": "UserDashboard",
        "path": "src/pages/UserDashboard.tsx",
        "type": "page",
        "components_used": ["UserProfile", "OrderList", "ActivityFeed"],
        "props": [],
        "state_management": "react-query",
        "routes": ["/dashboard"]
      }
    ],
    "shared": [
      {
        "name": "Button",
        "path": "src/components/Button/Button.tsx",
        "type": "ui_component",
        "reusable": true,
        "props": [
          {"name": "variant", "type": "'primary' | 'secondary'", "required": false, "default": "primary"},
          {"name": "size", "type": "'sm' | 'md' | 'lg'", "required": false, "default": "md"},
          {"name": "onClick", "type": "() => void", "required": false},
          {"name": "disabled", "type": "boolean", "required": false, "default": false}
        ],
        "state_management": "none"
      }
    ]
  },
  "patterns": {
    "composition": "Higher-order components and render props",
    "state_management": {
      "global": "React Query for server state",
      "local": "useState and useReducer",
      "forms": "React Hook Form"
    },
    "props": {
      "naming": "camelCase",
      "typing": "TypeScript interfaces",
      "defaults": "defaultProps or default parameters"
    },
    "organization": "by_feature_with_shared_components"
  },
  "hooks": [
    {
      "name": "useAuth",
      "path": "src/hooks/useAuth.ts", 
      "purpose": "Authentication state and methods",
      "returns": "AuthContext with user, login, logout methods"
    },
    {
      "name": "useApi",
      "path": "src/hooks/useApi.ts",
      "purpose": "API calls with loading and error states",
      "returns": "Query state and mutation methods"
    }
  ],
  "styling_integration": {
    "method": "Tailwind CSS with CSS modules",
    "component_styling": "className props with Tailwind utilities",
    "theme_system": "CSS custom properties for design tokens"
  }
}
```

### 8. styling.json (Frontend)
```json
{
  "metadata": {
    "generated_at": "2025-01-15T10:30:00Z",
    "framework": "React",
    "primary_approach": "Tailwind CSS",
    "secondary_approaches": ["CSS Modules"]
  },
  "architecture": {
    "approach": "Utility-first with component abstractions",
    "organization": "Tailwind utilities + custom CSS modules for complex components",
    "naming_convention": "Tailwind classes + BEM for custom CSS",
    "responsive_strategy": "Mobile-first with Tailwind breakpoints"
  },
  "design_system": {
    "tokens": {
      "colors": {
        "primary": {
          "50": "#eff6ff",
          "500": "#3b82f6", 
          "900": "#1e3a8a"
        },
        "secondary": {
          "50": "#f9fafb",
          "500": "#6b7280",
          "900": "#111827"
        }
      },
      "spacing": {
        "unit": "0.25rem",
        "scale": [0, 1, 2, 3, 4, 5, 6, 8, 10, 12, 16, 20, 24, 32, 40, 48, 56, 64]
      },
      "typography": {
        "font_families": {
          "sans": ["Inter", "system-ui", "sans-serif"],
          "mono": ["JetBrains Mono", "monospace"]
        },
        "scales": {
          "xs": "0.75rem",
          "sm": "0.875rem", 
          "base": "1rem",
          "lg": "1.125rem",
          "xl": "1.25rem"
        }
      },
      "breakpoints": {
        "sm": "640px",
        "md": "768px",
        "lg": "1024px",
        "xl": "1280px"
      }
    }
  },
  "components_styling": [
    {
      "name": "Button",
      "approach": "Tailwind utilities with variants",
      "base_classes": "px-4 py-2 rounded-md font-medium transition-colors",
      "variants": {
        "primary": "bg-blue-500 text-white hover:bg-blue-600",
        "secondary": "bg-gray-200 text-gray-900 hover:bg-gray-300"
      },
      "responsive": "responsive padding and text sizing"
    },
    {
      "name": "Card",
      "approach": "CSS Module + Tailwind",
      "file": "Card.module.css",
      "reason": "Complex shadows and animations not easily achievable with Tailwind"
    }
  ],
  "responsive_design": {
    "strategy": "Mobile-first design with progressive enhancement",
    "breakpoints_usage": "Tailwind responsive prefixes (sm:, md:, lg:, xl:)",
    "layout_approach": "CSS Grid and Flexbox through Tailwind utilities",
    "image_strategy": "Responsive images with Next.js Image component"
  },
  "theming": {
    "dark_mode": {
      "enabled": true,
      "strategy": "CSS custom properties with Tailwind dark: prefix",
      "toggle": "useTheme hook with localStorage persistence"
    },
    "customization": {
      "method": "Tailwind config + CSS custom properties",
      "file": "tailwind.config.js + globals.css"
    }
  },
  "performance": {
    "optimization": "PurgeCSS automatic with Tailwind build process",
    "bundle_size": "Only used utilities included in production",
    "critical_css": "Inlined for above-the-fold content"
  }
}
```

## 🔄 Fluxo de Geração Atualizado

### Fase 1: Preparação e Validação
```
1.1. Validar .phengineer/ configurado
1.2. Carregar config.yml 
1.3. Carregar .analyzefiles e .ignorefiles
1.4. Git diff analysis para modo incremental
1.5. Detectar linguagem e tipo de projeto
```

### Fase 2: Code Analysis → JSON
```
2.1. Gerar file-tree.json (estrutura hierárquica)
2.2. Gerar statistics.json (métricas quantitativas)
```

### Fase 3: IA Analysis → JSON
```
3.1. Gerar stack.json (tecnologias e dependências)
3.2. Gerar architecture.json (padrões e data models)
3.3. Atualizar functions.json (incremental por arquivo)
3.4. Atualizar project-context.json (incremental + relevância)
3.5. Frontend: components.json + styling.json
```

### Fase 4: Finalização
```
4.1. Validar schemas JSON gerados
4.2. Atualizar last_commit_hash
4.3. Salvar contextos em .phengineer/context/
```

## ⚙️ Configuração

### config.yml
```yaml
project:
  name: "my-project"
  type: "frontend"  # auto-detected
  language: "typescript"  # auto-detected

analysis:
  last_commit_hash: ""
  incremental_contexts: ["functions", "project-context"]
  validate_schemas: true
  
contexts:
  format: "json"
  generate_human_docs: false  # Opcional
  compress_output: false
```

## 📊 Estrutura Final

### Contextos JSON (Para IA)
```
.phengineer/context/
├── file-tree.json        # 🤖 Estrutura hierárquica
├── statistics.json       # 🤖 Métricas quantitativas  
├── stack.json           # 🤖 Stack tecnológico
├── architecture.json    # 🤖 Padrões arquiteturais + data models
├── functions.json       # 🤖 Inventário funções (incremental)
├── project-context.json # 🤖 Docs + configs (incremental + relevância)
├── components.json      # 🤖 Frontend: componentes
└── styling.json         # 🤖 Frontend: styling
```

### Schemas de Validação
```
.phengineer/schemas/
├── file-tree.schema.json
├── statistics.schema.json
├── stack.schema.json
├── architecture.schema.json
├── functions.schema.json
├── project-context.schema.json
├── components.schema.json
└── styling.schema.json
```

## 🎯 Benefícios dos Schemas JSON

### Para IA
- **Estrutura garantida** - Parsing consistente
- **Acesso direto** - `context.functions[0].name`
- **Menos tokens** - Sem formatação markdown
- **Universal** - Funciona com qualquer LLM

### Para Desenvolvimento
- **Validação automática** - Schemas garantem estrutura
- **Versionamento** - Evolução controlada dos contextos
- **Tooling** - Processamento programático
- **Debug** - Estrutura clara e inspecionável

### Para Manutenção
- **Consistência** - Mesmo formato entre projetos
- **Incremental preciso** - Updates granulares
- **Relevância inteligente** - Filtro automático de ruído
- **Eficiência** - Apenas mudanças reais processadas

---

> **Nota**: Todos os contextos são gerados em JSON estruturado para máxima compatibilidade com LLMs e eficiência de processamento. A estrutura incremental e filtros de relevância garantem atualizações precisas e contextos limpos.
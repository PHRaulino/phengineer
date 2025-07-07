# PHEngineer - Fluxo de Entrada e Validação

## 🔄 Fluxo Definido

### 1. Trigger Inicial
- **Entrada:** Issue criada/editada no GitHub
- **Ação:** GitHub Actions detecta mudança (`issues: [opened, edited]`)

### 2. Interpretação de Requisitos
- **Agente:** Requirements Interpreter
- **Input:** Texto da Issue
- **Output:** JSON estruturado com especificação técnica

### 3. Geração de Especificação
- **Responsável:** Orquestrador
- **Ação:** Usar template Go para renderizar MD da especificação
- **Output:** Issue atualizada com especificação estruturada

### 4. Validação do Usuário
- **Mecanismo:** Seção "STATUS: PENDING/APPROVED" no MD
- **Trigger:** Edição da Issue pelo usuário
- **Estados:** PENDING → APPROVED/REJECTED

### 5. Detecção de Aprovação
- **Método:** Parser detecta "STATUS: APPROVED" no body da Issue
- **Ação:** Inicia pipeline de desenvolvimento

## 📋 JSON de Especificação

```json
{
  "generation_type": "feature|test|fix|doc|refactor",
  "summary": "Descrição resumida da funcionalidade",
  "architecture": {
    "pattern": "Clean Architecture",
    "framework": "Gin",
    "database": "PostgreSQL"
  },
  "files_changes": [
    {
      "file_path": "caminho/do/arquivo.go",
      "change": "Descrição da mudança",
      "type": "new_file|modify|delete"
    }
  ],
  "tests": [
    {
      "type": "unit|integration|e2e",
      "description": "Descrição dos testes"
    }
  ],
  "complexity": "low|medium|high",
  "dor": [
    "Lista de critérios de prontidão"
  ],
  "dod": [
    "Lista de critérios de conclusão"
  ]
}
```

## 📄 Template de Issue

A Issue é atualizada com:
- **Solicitação Original:** Preserva prompt inicial
- **Especificação Técnica:** Dados do JSON formatados
- **DOR/DOD:** Critérios de qualidade
- **Seção de Correções:** Para iterações do usuário
- **Status de Aprovação:** Controle de fluxo

## 🎯 Próximo Fluxo: Pipeline de Desenvolvimento

**Entrada:** Issue com STATUS: APPROVED
**Saída:** Pull Request com código gerado

### Componentes Necessários:
1. **Seletor de Pipeline:** Identifica tipo de geração
2. **Agent Coordinator:** Orquestra agentes especializados
3. **Agentes Paralelos:**
   - Code Generator
   - Test Generator
   - Documentation Generator
4. **Result Consolidator:** Agrupa resultados
5. **PR Creator:** Gera Pull Request final

### Decisões Pendentes:
- Estratégia de execução paralela vs sequencial
- Tratamento de dependências entre agentes
- Validação de código gerado antes do PR
- Estrutura de branches (feature/fix/doc)
- Formato de commit messages    
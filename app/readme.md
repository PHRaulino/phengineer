# PHEngineer MVP - Context Generation System

## 🎯 Visão Geral

O PHEngineer MVP é um sistema de geração de contextos estruturados que analisa projetos de software e cria representações JSON otimizadas para agentes de IA. O objetivo é fornecer contexto rico e relevante para automação de desenvolvimento, começando com o Requirements Interpreter Agent.

## 🏗️ Arquitetura do Sistema

### Componentes Principais

```
Projeto → Configuração → Análise → Contextos JSON → Agentes IA
```

### Estrutura de Arquivos

```
.phengineer/
├── settings.yml        # Configurações do projeto
├── .ignorefiles       # Padrões de arquivos a ignorar
├── .analyzefiles      # Padrões de arquivos prioritários
└── context/           # Contextos gerados
    ├── file-tree.json
    ├── functions.json
    ├── stack.json
    └── project-context.json
```

## ⚙️ Configuração

### settings.yml

Arquivo central de configuração com informações do projeto e comportamento da análise:

```yaml
project:
  name: "meu-projeto"
  type: "lambda"           # Tipo arquitetural do projeto
  language: "python"       # Linguagem principal
  version: "3.13"         # Versão da linguagem

analysis:
  file_limits:
    max_file_size: "1MB"   # Limite de tamanho por arquivo
    max_files: 1000        # Limite total de arquivos

contexts:
  enabled: ["file-tree", "functions", "stack", "project-context"]
```

### .ignorefiles

Lista de padrões de arquivos que devem ser ignorados na análise (sintaxe similar ao .gitignore):

```
# Dependencies
node_modules/
__pycache__/
.venv/

# Build artifacts
dist/
build/
*.pyc

# System files
.git/
.DS_Store
*.log
```

### .analyzefiles

Lista de padrões de arquivos que devem ser priorizados na análise:

```
# Source code
**/*.py
**/*.js
**/*.ts

# Configuration
package.json
requirements.txt
Dockerfile

# Documentation
README.md
docs/**/*.md
```

## 📋 Contextos Gerados

### 1. file-tree.json

Estrutura hierárquica do projeto com metadados organizacionais.

**Propósito**: Fornecer navegação e compreensão da organização do projeto.

**Conteúdo**:
- Árvore de diretórios e arquivos
- Convenções de nomenclatura detectadas
- Arquivos importantes identificados
- Métricas básicas (contagem de arquivos, profundidade)

### 2. functions.json

Inventário completo de funções, métodos e classes do projeto.

**Propósito**: Mapeamento de funcionalidades existentes para geração de código complementar.

**Conteúdo**:
- Assinaturas de funções com tipos
- Documentação e propósito de cada função
- Dependências entre funções
- Localização no código (arquivo e linha)

### 3. stack.json

Tecnologias, frameworks e dependências utilizadas no projeto.

**Propósito**: Compreensão do ecossistema tecnológico para decisões arquiteturais.

**Conteúdo**:
- Linguagens e versões
- Frameworks principais
- Dependências críticas
- Ferramentas de desenvolvimento
- Estratégia de deployment

### 4. project-context.json

Documentação, configurações e contratos relevantes do projeto.

**Propósito**: Contexto de negócio e especificações técnicas para alinhamento de requisitos.

**Conteúdo**:
- Documentação principal (README, docs/)
- Arquivos de configuração ativos
- Contratos de API (OpenAPI, GraphQL)
- Configurações de ambiente

## 🔄 Fluxo de Operação

### Fase 1: Configuração
1. Criar estrutura `.phengineer/`
2. Configurar `settings.yml` baseado no tipo de projeto
3. Definir padrões em `.ignorefiles` e `.analyzefiles`

### Fase 2: Descoberta de Arquivos
1. Escanear projeto respeitando limites configurados
2. Aplicar filtros de ignore para remover arquivos irrelevantes
3. Priorizar arquivos definidos em `.analyzefiles`
4. Validar limites de tamanho e quantidade

### Fase 3: Geração de Contextos
1. Analisar estrutura de arquivos → `file-tree.json`
2. Extrair funções e classes → `functions.json`
3. Detectar stack tecnológico → `stack.json`
4. Processar documentação → `project-context.json`

### Fase 4: Validação
1. Verificar schemas JSON
2. Confirmar completude dos contextos
3. Registrar métricas de geração

## 🎯 Tipos de Projeto Suportados

### Lambda Functions
- **Características**: Projetos pequenos, foco em função específica
- **Contextos recomendados**: file-tree, functions, stack
- **Limites típicos**: <100 arquivos, <1MB por arquivo

### APIs
- **Características**: Endpoints REST/GraphQL, documentação de contratos
- **Contextos recomendados**: Todos os contextos
- **Limites típicos**: <1000 arquivos, <1MB por arquivo

### Frontend
- **Características**: Componentes, styling, estado de aplicação
- **Contextos recomendados**: file-tree, functions, stack, project-context
- **Limites típicos**: <1000 arquivos, assets limitados

### Backend
- **Características**: Lógica de negócio, integração com bancos, arquitetura
- **Contextos recomendados**: Todos os contextos
- **Limites típicos**: <1000 arquivos, <1MB por arquivo

## 📊 Benefícios

### Para Agentes IA
- **Contexto estruturado**: JSON parsing direto sem ambiguidade
- **Informação relevante**: Filtros eliminam ruído desnecessário
- **Compreensão profunda**: Múltiplas perspectivas do mesmo projeto
- **Eficiência**: Menos tokens, mais precisão

### Para Desenvolvedores
- **Setup simples**: Configuração mínima necessária
- **Controle fino**: Filtros customizáveis por projeto
- **Transparência**: Contextos inspecionáveis e versionáveis
- **Performance**: Limites evitam processamento excessivo

### Para Equipes
- **Consistência**: Mesma estrutura entre projetos
- **Colaboração**: Configuração compartilhável
- **Documentação**: Contextos servem como documentação viva
- **Onboarding**: Novos membros compreendem projeto rapidamente

## 🚀 Casos de Uso

### Requirements Interpreter Agent
1. Desenvolvedor cria Issue com solicitação
2. Sistema carrega contextos do projeto
3. Agente analisa solicitação + contextos
4. Gera especificação técnica estruturada
5. Desenvolvedor valida antes da implementação

### Code Generation
1. Especificação técnica aprovada
2. Agente usa contextos para entender padrões existentes
3. Gera código alinhado com arquitetura atual
4. Mantém consistência de naming e estrutura

### Documentation Generation
1. Contextos fornecem inventário completo
2. Agente gera documentação baseada em código atual
3. Mantém sincronização automática
4. Identifica gaps de documentação

## ⚠️ Limitações do MVP

### Escopo Reduzido
- **Setup manual**: Configuração inicial requer intervenção humana
- **Contextos básicos**: Apenas 4 contextos essenciais
- **Análise estática**: Sem execução de código ou testes
- **Single-shot**: Sem updates incrementais automáticos

### Suporte Limitado
- **Linguagens**: Foco inicial em Python/JavaScript/TypeScript
- **Frameworks**: Detecção básica sem análise profunda
- **Projetos grandes**: Limites podem ser restritivos para monorepos

### Validação Mínima
- **Schemas básicos**: Validação estrutural simples
- **Error handling**: Recuperação limitada de falhas
- **Performance**: Sem otimizações para projetos complexos

## 🎯 Critérios de Sucesso

### Funcionalidade
- [ ] Geração completa dos 4 contextos em <5 minutos
- [ ] Contextos contêm informações relevantes e precisas
- [ ] Requirements Interpreter Agent produz especificações coerentes
- [ ] Sistema funciona em projetos Python e JavaScript

### Qualidade
- [ ] Contextos passam validação de schema
- [ ] Filtros removem >90% de arquivos irrelevantes
- [ ] Informações extraídas estão corretas e atualizadas
- [ ] Performance aceitável em projetos de até 1000 arquivos

### Usabilidade
- [ ] Setup em <10 minutos para desenvolvedor experiente
- [ ] Configuração intuitiva e bem documentada
- [ ] Logs e feedback claros durante processamento
- [ ] Contextos são humanamente legíveis para debug

## 📈 Evolução Futura

### Automação
- Setup automático com detecção de projeto
- Updates incrementais baseados em Git
- Integração com CI/CD pipelines
- CLI robusta com múltiplas opções

### Contextos Avançados
- Análise de componentes para frontends
- Detecção de padrões arquiteturais
- Mapeamento de dependências complexas
- Métricas de qualidade e complexidade

### Inteligência
- Machine learning para detecção de padrões
- Sugestões automáticas de configuração
- Otimização baseada em histórico de uso
- Previsão de necessidades de contexto

---

> **Nota**: Este MVP foca na validação do conceito central - contextos estruturados geram valor real para agentes IA. A automação e sofisticação serão adicionadas após comprovação do valor fundamental.
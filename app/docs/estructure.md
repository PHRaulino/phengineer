# Implementação Completa - Cadastro de Credenciais

## 📁 Estrutura do Fluxo

```
├── cmd/cli/
│   └── auth.go                    # Comando CLI
├── internal/
│   ├── domain/
│   │   ├── auth/
│   │   │   ├── entity.go         # Entidades
│   │   │   ├── repository.go     # Interface do repositório
│   │   │   └── service.go        # Lógica de negócio
│   ├── infrastructure/
│   │   ├── auth/
│   │   │   └── keyring_repository.go  # Implementação do storage
│   │   └── stackspot/
│   │       └── client.go          # Cliente da API
│   └── presentation/
│       └── tui/
│           ├── screens/
│           │   ├── auth_setup.go  # Tela principal
│           │   ├── auth_user.go   # Tela de user auth
│           │   └── auth_system.go # Tela de system auth
│           └── messages/
│               └── auth.go        # Mensagens específicas
```

## 1. Comando CLI (cmd/cli/auth.go)

```go
package cli

import (
    "fmt"
    "github.com/spf13/cobra"
    "github.com/PHRaulino/phengineer/internal/domain/auth"
    "github.com/PHRaulino/phengineer/internal/infrastructure/auth"
    "github.com/PHRaulino/phengineer/internal/infrastructure/stackspot"
    "github.com/PHRaulino/phengineer/internal/presentation/tui"
    "github.com/PHRaulino/phengineer/internal/presentation/tui/screens"
)

var authCmd = &cobra.Command{
    Use:   "auth",
    Short: "Manage authentication",
}

var authSetupCmd = &cobra.Command{
    Use:   "setup",
    Short: "Configure authentication interactively",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Inicializar dependências
        keyringRepo := authinfra.NewKeyringRepository()
        stackspotClient := stackspot.NewClient()
        authService := auth.NewService(keyringRepo, stackspotClient)

        // Criar e iniciar TUI
        app := tui.NewApp()
        setupScreen := screens.NewAuthSetupScreen(authService)
        app.SetInitialScreen(setupScreen)

        return app.Start()
    },
}

func init() {
    authCmd.AddCommand(authSetupCmd)
    rootCmd.AddCommand(authCmd)
}
```

## 2. Entidades de Domínio (internal/domain/auth/entity.go)

```go
package auth

import (
    "time"
    "errors"
)

// Tipos de autenticação
type AuthMode string

const (
    AuthModeUser    AuthMode = "stackspot_user"
    AuthModeService AuthMode = "stackspot_service"
)

// Credenciais base
type Credentials struct {
    Mode      AuthMode  `json:"mode"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// Credenciais de usuário
type UserCredentials struct {
    Credentials
    ClientID     string `json:"client_id"`
    ClientSecret string `json:"client_secret"`
}

// Credenciais de sistema (Vault)
type SystemCredentials struct {
    Credentials
    VaultURL   string `json:"vault_url"`
    VaultPath  string `json:"vault_path"`
    VaultToken string `json:"vault_token"`
}

// Token armazenado
type Token struct {
    AccessToken  string    `json:"access_token"`
    RefreshToken string    `json:"refresh_token"`
    TokenType    string    `json:"token_type"`
    ExpiresAt    time.Time `json:"expires_at"`
    Scope        string    `json:"scope"`
}

// Validações
func (uc *UserCredentials) Validate() error {
    if uc.ClientID == "" {
        return errors.New("client ID é obrigatório")
    }
    if uc.ClientSecret == "" {
        return errors.New("client secret é obrigatório")
    }
    return nil
}

func (sc *SystemCredentials) Validate() error {
    if sc.VaultURL == "" {
        return errors.New("vault URL é obrigatório")
    }
    if sc.VaultPath == "" {
        return errors.New("vault path é obrigatório")
    }
    if sc.VaultToken == "" {
        return errors.New("vault token é obrigatório")
    }
    return nil
}
```

## 3. Interface do Repositório (internal/domain/auth/repository.go)

```go
package auth

import "context"

// Repository define as operações de persistência
type Repository interface {
    // Credenciais
    SaveUserCredentials(ctx context.Context, creds *UserCredentials) error
    SaveSystemCredentials(ctx context.Context, creds *SystemCredentials) error
    GetCredentials(ctx context.Context) (interface{}, error)
    DeleteCredentials(ctx context.Context) error

    // Tokens
    SaveToken(ctx context.Context, scope string, token *Token) error
    GetToken(ctx context.Context, scope string) (*Token, error)
    DeleteToken(ctx context.Context, scope string) error
    DeleteAllTokens(ctx context.Context) error
}
```

## 4. Serviço de Domínio (internal/domain/auth/service.go)

```go
package auth

import (
    "context"
    "fmt"
    "time"
)

// StackspotClient interface para o cliente da API
type StackspotClient interface {
    AuthenticateUser(ctx context.Context, clientID, clientSecret string) (*Token, error)
    AuthenticateSystem(ctx context.Context, vaultCreds *SystemCredentials) (*Token, error)
    ValidateToken(ctx context.Context, token string) error
}

// Service implementa a lógica de negócio
type Service struct {
    repo   Repository
    client StackspotClient
}

func NewService(repo Repository, client StackspotClient) *Service {
    return &Service{
        repo:   repo,
        client: client,
    }
}

// SetupUserAuth configura autenticação de usuário
func (s *Service) SetupUserAuth(ctx context.Context, clientID, clientSecret string) error {
    // Validar credenciais
    creds := &UserCredentials{
        Credentials: Credentials{
            Mode:      AuthModeUser,
            CreatedAt: time.Now(),
            UpdatedAt: time.Now(),
        },
        ClientID:     clientID,
        ClientSecret: clientSecret,
    }

    if err := creds.Validate(); err != nil {
        return fmt.Errorf("credenciais inválidas: %w", err)
    }

    // Testar autenticação
    token, err := s.client.AuthenticateUser(ctx, clientID, clientSecret)
    if err != nil {
        return fmt.Errorf("falha na autenticação: %w", err)
    }

    // Salvar credenciais
    if err := s.repo.SaveUserCredentials(ctx, creds); err != nil {
        return fmt.Errorf("erro ao salvar credenciais: %w", err)
    }

    // Salvar token
    if err := s.repo.SaveToken(ctx, "execution", token); err != nil {
        return fmt.Errorf("erro ao salvar token: %w", err)
    }

    return nil
}

// SetupSystemAuth configura autenticação de sistema
func (s *Service) SetupSystemAuth(ctx context.Context, vaultURL, vaultPath, vaultToken string) error {
    creds := &SystemCredentials{
        Credentials: Credentials{
            Mode:      AuthModeService,
            CreatedAt: time.Now(),
            UpdatedAt: time.Now(),
        },
        VaultURL:   vaultURL,
        VaultPath:  vaultPath,
        VaultToken: vaultToken,
    }

    if err := creds.Validate(); err != nil {
        return fmt.Errorf("credenciais inválidas: %w", err)
    }

    // Testar conexão com Vault e autenticação
    token, err := s.client.AuthenticateSystem(ctx, creds)
    if err != nil {
        return fmt.Errorf("falha na autenticação: %w", err)
    }

    // Salvar credenciais
    if err := s.repo.SaveSystemCredentials(ctx, creds); err != nil {
        return fmt.Errorf("erro ao salvar credenciais: %w", err)
    }

    // Salvar token
    if err := s.repo.SaveToken(ctx, "execution", token); err != nil {
        return fmt.Errorf("erro ao salvar token: %w", err)
    }

    return nil
}

// GetCurrentAuth retorna as credenciais atuais
func (s *Service) GetCurrentAuth(ctx context.Context) (interface{}, error) {
    return s.repo.GetCredentials(ctx)
}

// ClearAuth remove todas as credenciais e tokens
func (s *Service) ClearAuth(ctx context.Context) error {
    if err := s.repo.DeleteAllTokens(ctx); err != nil {
        return err
    }
    return s.repo.DeleteCredentials(ctx)
}
```

## 5. Tela Principal de Setup (internal/presentation/tui/screens/auth_setup.go)

```go
package screens

import (
    "context"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/PHRaulino/phengineer/internal/domain/auth"
    "github.com/PHRaulino/phengineer/internal/presentation/tui/components/navigation"
    "github.com/PHRaulino/phengineer/internal/presentation/tui/messages"
    "github.com/PHRaulino/phengineer/internal/presentation/tui/models"
)

type AuthSetupScreen struct {
    models.BaseModel
    authService *auth.Service
    menu        *navigation.Menu
    ctx         context.Context
}

func NewAuthSetupScreen(authService *auth.Service) *AuthSetupScreen {
    s := &AuthSetupScreen{
        BaseModel: models.BaseModel{
            Theme: styles.DefaultTheme,
        },
        authService: authService,
        ctx:         context.Background(),
    }

    s.initMenu()
    return s
}

func (s *AuthSetupScreen) initMenu() {
    s.menu = navigation.NewMenu(s.Theme).
        AddItem(navigation.MenuItem{
            Title:       "Stackspot User",
            Description: "Autenticação pessoal para desenvolvedores",
            Icon:        "👤",
            Action: func() tea.Msg {
                return messages.ChangeScreenMsg{
                    Screen: NewAuthUserScreen(s.authService),
                }
            },
        }).
        AddItem(navigation.MenuItem{
            Title:       "Stackspot Service",
            Description: "Autenticação de sistema via Hashicorp Vault",
            Icon:        "🏢",
            Action: func() tea.Msg {
                return messages.ChangeScreenMsg{
                    Screen: NewAuthSystemScreen(s.authService),
                }
            },
        }).
        AddItem(navigation.MenuItem{
            Title:       "Verificar Configuração",
            Description: "Visualizar autenticação atual",
            Icon:        "🔍",
            Action: func() tea.Msg {
                return s.checkCurrentAuth()
            },
        }).
        AddItem(navigation.MenuItem{
            Title:       "Limpar Autenticação",
            Description: "Remover todas as credenciais",
            Icon:        "🗑️",
            Action: func() tea.Msg {
                return messages.ConfirmActionMsg{
                    Message: "Tem certeza que deseja remover todas as credenciais?",
                    Action:  s.clearAuth,
                }
            },
        })

    s.menu.SetWidth(s.Width - 4)
}

func (s *AuthSetupScreen) Init() tea.Cmd {
    return tea.Batch(
        s.Header.Init(),
        s.Footer.Init(),
    )
}

func (s *AuthSetupScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "q", "ctrl+c":
            return s, tea.Quit
        }

    case messages.AuthSuccessMsg:
        // Mostrar mensagem de sucesso e voltar ao menu principal
        return s, tea.Batch(
            messages.ShowNotification("✅ Autenticação configurada com sucesso!", "success"),
            tea.Tick(2*time.Second, func(time.Time) tea.Msg {
                return messages.PopScreenMsg{}
            }),
        )

    case messages.AuthErrorMsg:
        s.Error = msg.Error
        return s, nil
    }

    // Atualizar menu
    newMenu, cmd := s.menu.Update(msg)
    s.menu = newMenu.(*navigation.Menu)

    return s, cmd
}

func (s *AuthSetupScreen) View() string {
    // Header
    header := s.Header.View()

    // Título
    title := s.Theme.GetStyles().Title.Render("🔐 Configuração de Autenticação")
    subtitle := s.Theme.GetStyles().Subtitle.Render("Escolha o tipo de autenticação")

    // Menu
    menuView := s.menu.View()

    // Error
    var errorView string
    if s.Error != nil {
        errorView = s.Theme.GetStyles().Error.Render("❌ " + s.Error.Error())
    }

    // Footer
    footer := s.Footer.View()

    // Compor tudo
    content := lipgloss.JoinVertical(
        lipgloss.Left,
        title,
        subtitle,
        "",
        menuView,
        errorView,
    )

    // Centralizar verticalmente
    availableHeight := s.Height - lipgloss.Height(header) - lipgloss.Height(footer) - 2
    content = lipgloss.Place(
        s.Width,
        availableHeight,
        lipgloss.Center,
        lipgloss.Center,
        content,
    )

    return lipgloss.JoinVertical(
        lipgloss.Left,
        header,
        content,
        footer,
    )
}

func (s *AuthSetupScreen) checkCurrentAuth() tea.Msg {
    creds, err := s.authService.GetCurrentAuth(s.ctx)
    if err != nil {
        return messages.ShowNotification("❌ Nenhuma autenticação configurada", "error")
    }

    // Mostrar informações das credenciais
    switch c := creds.(type) {
    case *auth.UserCredentials:
        return messages.ShowNotification(
            fmt.Sprintf("✅ Autenticado como usuário: %s", c.ClientID),
            "success",
        )
    case *auth.SystemCredentials:
        return messages.ShowNotification(
            fmt.Sprintf("✅ Autenticado via Vault: %s", c.VaultURL),
            "success",
        )
    }

    return nil
}

func (s *AuthSetupScreen) clearAuth() tea.Msg {
    if err := s.authService.ClearAuth(s.ctx); err != nil {
        return messages.AuthErrorMsg{Error: err}
    }
    return messages.ShowNotification("✅ Autenticação removida com sucesso", "success")
}

func (s *AuthSetupScreen) GetTitle() string {
    return "Configuração de Autenticação"
}
```

## 6. Tela de Autenticação de Usuário (internal/presentation/tui/screens/auth_user.go)

```go
package screens

import (
    "context"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/PHRaulino/phengineer/internal/domain/auth"
    "github.com/PHRaulino/phengineer/internal/presentation/tui/components/forms"
    "github.com/PHRaulino/phengineer/internal/presentation/tui/messages"
    "github.com/PHRaulino/phengineer/internal/presentation/tui/models"
)

type AuthUserScreen struct {
    models.BaseModel
    authService  *auth.Service
    form         *forms.Form
    processing   bool
    clientID     string
    clientSecret string
}

func NewAuthUserScreen(authService *auth.Service) *AuthUserScreen {
    s := &AuthUserScreen{
        BaseModel: models.BaseModel{
            Theme: styles.DefaultTheme,
        },
        authService: authService,
    }

    s.initForm()
    return s
}

func (s *AuthUserScreen) initForm() {
    s.form = forms.NewForm(
        "👤 Credenciais de Usuário Stackspot",
        "Digite suas credenciais para autenticação pessoal",
    ).
        AddField("Client ID", forms.NewInput().
            WithPlaceholder("seu-client-id-aqui").
            WithValidation(forms.Required).
            WithHelp("ID fornecido pela Stackspot")).
        AddField("Client Secret", forms.NewPassword().
            WithPlaceholder("seu-client-secret-aqui").
            WithValidation(forms.Required).
            WithHelp("Secret fornecido pela Stackspot")).
        OnSubmit(s.handleSubmit)

    s.form.SetTheme(s.Theme)
    s.form.SetWidth(60)
}

func (s *AuthUserScreen) Init() tea.Cmd {
    return tea.Batch(
        s.form.Init(),
        textinput.Blink,
    )
}

func (s *AuthUserScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "esc":
            if !s.processing {
                return s, messages.PopScreen()
            }
        case "ctrl+c":
            return s, tea.Quit
        }

    case forms.SubmitMsg:
        values := msg.Values
        s.clientID = values["Client ID"]
        s.clientSecret = values["Client Secret"]
        s.processing = true
        return s, s.authenticate()

    case messages.AuthSuccessMsg:
        return s, messages.ChangeScreen(NewAuthSuccessScreen())

    case messages.AuthErrorMsg:
        s.processing = false
        s.Error = msg.Error
        return s, nil
    }

    if !s.processing {
        newForm, cmd := s.form.Update(msg)
        s.form = newForm.(*forms.Form)
        return s, cmd
    }

    return s, nil
}

func (s *AuthUserScreen) View() string {
    if s.processing {
        spinner := s.Theme.GetStyles().Spinner.Render("⠋")
        message := s.Theme.GetStyles().Muted.Render("Autenticando com Stackspot...")

        return lipgloss.Place(
            s.Width,
            s.Height,
            lipgloss.Center,
            lipgloss.Center,
            lipgloss.JoinVertical(
                lipgloss.Center,
                spinner,
                message,
            ),
        )
    }

    // Form view
    formView := s.form.View()

    // Error
    var errorView string
    if s.Error != nil {
        errorStyle := s.Theme.GetStyles().Error.
            Width(60).
            Align(lipgloss.Center).
            MarginTop(1)
        errorView = errorStyle.Render("❌ " + s.Error.Error())
    }

    // Help
    helpStyle := s.Theme.GetStyles().Help
    help := helpStyle.Render("ESC voltar • TAB próximo campo • ENTER confirmar")

    content := lipgloss.JoinVertical(
        lipgloss.Center,
        formView,
        errorView,
        "",
        help,
    )

    return lipgloss.Place(
        s.Width,
        s.Height,
        lipgloss.Center,
        lipgloss.Center,
        content,
    )
}

func (s *AuthUserScreen) authenticate() tea.Cmd {
    return func() tea.Msg {
        ctx := context.Background()

        err := s.authService.SetupUserAuth(ctx, s.clientID, s.clientSecret)
        if err != nil {
            return messages.AuthErrorMsg{Error: err}
        }

        return messages.AuthSuccessMsg{}
    }
}

func (s *AuthUserScreen) handleSubmit(values map[string]string) tea.Cmd {
    return func() tea.Msg {
        return forms.SubmitMsg{Values: values}
    }
}

func (s *AuthUserScreen) GetTitle() string {
    return "Autenticação de Usuário"
}
```

## 7. Mensagens Customizadas (internal/presentation/tui/messages/auth.go)

```go
package messages

import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/PHRaulino/phengineer/internal/presentation/tui/models"
)

// AuthSuccessMsg indica sucesso na autenticação
type AuthSuccessMsg struct{}

// AuthErrorMsg indica erro na autenticação
type AuthErrorMsg struct {
    Error error
}

// ConfirmActionMsg solicita confirmação do usuário
type ConfirmActionMsg struct {
    Message string
    Action  func() tea.Msg
}

// ShowNotification exibe uma notificação temporária
type ShowNotificationMsg struct {
    Message string
    Type    string // "success", "error", "info", "warning"
}

// Helpers para criar comandos
func ShowNotification(message, notificationType string) tea.Cmd {
    return func() tea.Msg {
        return ShowNotificationMsg{
            Message: message,
            Type:    notificationType,
        }
    }
}
```

## 8. Implementação do Repositório Keyring (internal/infrastructure/auth/keyring_repository.go)

```go
package authinfra

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/zalando/go-keyring"
    "github.com/PHRaulino/phengineer/internal/domain/auth"
)

const (
    serviceName = "stackspot-cli"
    credsKey    = "credentials"
    tokenPrefix = "token_"
)

type KeyringRepository struct{}

func NewKeyringRepository() *KeyringRepository {
    return &KeyringRepository{}
}

// SaveUserCredentials salva credenciais de usuário
func (r *KeyringRepository) SaveUserCredentials(ctx context.Context, creds *auth.UserCredentials) error {
    data, err := json.Marshal(creds)
    if err != nil {
        return fmt.Errorf("erro ao serializar credenciais: %w", err)
    }

    return keyring.Set(serviceName, credsKey, string(data))
}

// SaveSystemCredentials salva credenciais de sistema
func (r *KeyringRepository) SaveSystemCredentials(ctx context.Context, creds *auth.SystemCredentials) error {
    data, err := json.Marshal(creds)
    if err != nil {
        return fmt.Errorf("erro ao serializar credenciais: %w", err)
    }

    return keyring.Set(serviceName, credsKey, string(data))
}

// GetCredentials recupera as credenciais armazenadas
func (r *KeyringRepository) GetCredentials(ctx context.Context) (interface{}, error) {
    data, err := keyring.Get(serviceName, credsKey)
    if err != nil {
        return nil, fmt.Errorf("credenciais não encontradas: %w", err)
    }

    // Tentar decodificar como UserCredentials primeiro
    var userCreds auth.UserCredentials
    if err := json.Unmarshal([]byte(data), &userCreds); err == nil && userCreds.Mode == auth.AuthModeUser {
        return &userCreds, nil
    }

    // Tentar como SystemCredentials
    var sysCreds auth.SystemCredentials
    if err := json.Unmarshal([]byte(data), &sysCreds); err == nil && sysCreds.Mode == auth.AuthModeService {
        return &sysCreds, nil
    }

    return nil, fmt.Errorf("formato de credenciais inválido")
}

// DeleteCredentials remove as credenciais
func (r *KeyringRepository) DeleteCredentials(ctx context.Context) error {
    return keyring.Delete(serviceName, credsKey)
}

// SaveToken salva um token para um escopo específico
func (r *KeyringRepository) SaveToken(ctx context.Context, scope string, token *auth.Token) error {
    data, err := json.Marshal(token)
    if err != nil {
        return fmt.Errorf("erro ao serializar token: %w", err)
    }

    key := tokenPrefix + scope
    return keyring.Set(serviceName, key, string(data))
}

// GetToken recupera um token por escopo
func (r *KeyringRepository) GetToken(ctx context.Context, scope string) (*auth.Token, error) {
    key := tokenPrefix + scope
    data, err := keyring.Get(serviceName, key)
    if err != nil {
        return nil, fmt.Errorf("token não encontrado: %w", err)
    }

    var token auth.Token
    if err := json.Unmarshal([]byte(data), &token); err != nil {
        return nil, fmt.Errorf("erro ao decodificar token: %w", err)
    }

    return &token, nil
}

// DeleteToken remove um token específico
func (r *KeyringRepository) DeleteToken(ctx context.Context, scope string) error {
    key := tokenPrefix + scope
    return keyring.Delete(serviceName, key)
}

// DeleteAllTokens remove todos os tokens
func (r *KeyringRepository) DeleteAllTokens(ctx context.Context) error {
    // Lista de escopos conhecidos
    scopes := []string{"execution", "creation", "read", "write"}

    for _, scope := range scopes {
        key := tokenPrefix + scope
        _ = keyring.Delete(serviceName, key) // Ignora erro se não existir
    }

    return nil
}
```

## 🎯 Fluxo Completo

1. **Usuário executa**: `stackspot auth setup`
2. **CLI inicializa**: Cria o serviço com suas dependências
3. **TUI abre**: Mostra menu de opções de autenticação
4. **Usuário seleciona**: "Stackspot User"
5. **Formulário aparece**: Campos para Client ID e Secret
6. **Usuário preenche**: Dados validados em tempo real
7. **Submit**: Chama o serviço de autenticação
8. **Serviço**:
   - Valida credenciais
   - Testa autenticação com API
   - Salva no Keyring
   - Armazena token
9. **TUI**: Mostra sucesso e volta ao menu
10. **Credenciais salvas**: Disponíveis para uso futuro

## 🔧 Pontos de Extensão

- **Novos tipos de auth**: Adicionar nova tela e método no serviço
- **Validações customizadas**: Extender o sistema de validação de forms
- **Storage alternativo**: Implementar nova interface Repository
- **Notificações**: Sistema de toast notifications
- **Temas**: Adicionar novos temas visuais

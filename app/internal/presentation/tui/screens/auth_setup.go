package screens

import (
	"github.com/PHRaulino/phengineer/internal/infrastructure/auth"
	"github.com/PHRaulino/phengineer/internal/presentation/tui/components/forms"
	"github.com/PHRaulino/phengineer/internal/presentation/tui/messages"
	"github.com/PHRaulino/phengineer/internal/presentation/tui/models"
	"github.com/PHRaulino/phengineer/internal/presentation/tui/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type AuthSetupScreen struct {
	models.BaseModel
	form        *forms.Form
	authType    string
	step        int
	credentials map[string]string
}

func NewAuthSetupScreen() *AuthSetupScreen {
	s := &AuthSetupScreen{
		BaseModel: models.BaseModel{
			Theme: styles.DefaultTheme,
		},
		credentials: make(map[string]string),
		step:        0,
	}

	s.initForm()
	return s
}

func (s *AuthSetupScreen) initForm() {
	switch s.step {
	case 0: // Seleção do tipo de auth
		s.form = forms.NewForm(
			"🔐 Configuração de Autenticação",
			"Escolha o tipo de autenticação que deseja configurar",
		).AddField("Tipo", forms.NewSelect([]string{
			"Stackspot User (Desenvolvedor)",
			"Stackspot Service (Sistema via Vault)", 
			"GitHub (Personal Access Token)",
		}))

	case 1: // Configuração específica
		if s.authType == "user" {
			s.form = forms.NewForm(
				"👤 Credenciais de Usuário",
				"Insira suas credenciais do Stackspot",
			).
				AddField("Client ID", forms.NewInput().
					WithPlaceholder("seu-client-id").
					WithValidation(forms.Required)).
				AddField("Client Secret", forms.NewPassword().
					WithPlaceholder("seu-client-secret").
					WithValidation(forms.Required))
		} else if s.authType == "service" {
			s.form = forms.NewForm(
				"🏢 Configuração do Vault",
				"Configure a integração com Hashicorp Vault",
			).
				AddField("Vault URL", forms.NewInput().
					WithPlaceholder("https://vault.empresa.com").
					WithValidation(forms.ValidateURL)).
				AddField("AWS Role", forms.NewInput().
					WithPlaceholder("stackspot-role").
					WithValidation(forms.Required)).
				AddField("StackSpot Secret Path", forms.NewInput().
					WithPlaceholder("secret/data/stackspot").
					WithValidation(forms.Required))
		} else if s.authType == "github" {
			s.form = forms.NewForm(
				"🐙 Configuração GitHub",
				"Configure seu Personal Access Token do GitHub",
			).
				AddField("Personal Access Token", forms.NewPassword().
					WithPlaceholder("ghp_xxxxxxxxxxxxxxxxxxxx").
					WithValidation(forms.Required))
		}
	}
}

func (s *AuthSetupScreen) Init() tea.Cmd {
	return s.form.Init()
}

func (s *AuthSetupScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return s, func() tea.Msg {
				return messages.PopScreenMsg{}
			}
		}
	case forms.SubmitMsg:
		// Processar valores do formulário
		s.credentials = msg.Values
		if s.step == 0 {
			// Avançar para próximo step
			switch msg.Values["Tipo"] {
			case "Stackspot User (Desenvolvedor)":
				s.authType = "user"
			case "Stackspot Service (Sistema via Vault)":
				s.authType = "service"
			case "GitHub (Personal Access Token)":
				s.authType = "github"
			}
			s.step = 1
			s.initForm()
			return s, s.form.Init()
		} else {
			// Finalizar configuração - salvar credenciais
			return s, s.saveCredentials()
		}
	}

	// Atualizar formulário
	newForm, cmd := s.form.Update(msg)
	s.form = newForm.(*forms.Form)
	return s, cmd
}

func (s *AuthSetupScreen) View() string {
	containerStyle := lipgloss.NewStyle().
		Width(s.Width).
		Height(s.Height).
		Align(lipgloss.Center, lipgloss.Center).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(s.Theme.Border).
		Padding(1, 2)

	content := s.form.View()

	// Adicionar navegação
	helpStyle := s.Theme.GetStyles().Info.
		MarginTop(2).
		Align(lipgloss.Center)

	help := helpStyle.Render("ESC para voltar • Tab/Enter para navegar")

	return containerStyle.Render(content + "\n" + help)
}

func (s *AuthSetupScreen) SetSize(width, height int) {
	s.BaseModel.SetSize(width, height)
	if s.form != nil {
		s.form.SetWidth(width - 8)
	}
}

func (s *AuthSetupScreen) SetTheme(theme *styles.Theme) {
	s.BaseModel.Theme = theme
	if s.form != nil {
		s.form.SetTheme(theme)
	}
}

func (s *AuthSetupScreen) GetTitle() string {
	return "Auth Setup"
}

func (s *AuthSetupScreen) HandleError(err error) tea.Cmd {
	s.BaseModel.Error = err
	return nil
}

func (s *AuthSetupScreen) saveCredentials() tea.Cmd {
	return func() tea.Msg {
		var err error
		
		switch s.authType {
		case "user":
			// Salvar credenciais StackSpot diretas
			provider := auth.GetStackSpotProvider()
			err = provider.SaveCredentials(
				s.credentials["Client ID"],
				s.credentials["Client Secret"],
			)
			
		case "service":
			// Salvar configuração do Vault
			provider := auth.GetVaultProvider()
			err = provider.SaveConfig(
				s.credentials["Vault URL"],
				s.credentials["AWS Role"],
				s.credentials["StackSpot Secret Path"],
			)
			
		case "github":
			// Salvar token GitHub
			provider := auth.GetGitHubProvider()
			err = provider.SaveToken(s.credentials["Personal Access Token"])
		}
		
		if err != nil {
			// TODO: Retornar erro para a UI
			return messages.PopScreenMsg{}
		}
		
		return messages.PopScreenMsg{}
	}
}

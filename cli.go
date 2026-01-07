package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Boca Juniors colors 💙💛
var (
	bocaBlue   = lipgloss.Color("#0038A8")
	bocaGold   = lipgloss.Color("#FFCC00")
	greenColor = lipgloss.Color("#00FF00")
	whiteColor = lipgloss.Color("#FFFFFF")
)

// Styles
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(bocaGold).
			Background(bocaBlue).
			Padding(1, 4).
			MarginBottom(1).
			Align(lipgloss.Center)

	disclaimerStyle = lipgloss.NewStyle().
			Italic(true).
			Foreground(greenColor).
			MarginBottom(1)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(bocaGold).
			Padding(1, 2).
			MarginBottom(1)

	selectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(bocaBlue).
			Background(bocaGold).
			Padding(0, 2)

	normalStyle = lipgloss.NewStyle().
			Foreground(whiteColor).
			Padding(0, 2)

	promptStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(bocaGold).
			MarginTop(1)

	successStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(greenColor).
			MarginTop(1)

	savedPathStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Italic(true)

	warningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFA500")).
			Italic(true)

	pixelArt = `
 ███████╗███╗   ██╗██╗   ██╗    ███████╗██╗    ██╗██╗████████╗ ██████╗██╗  ██╗
 ██╔════╝████╗  ██║██║   ██║    ██╔════╝██║    ██║██║╚══██╔══╝██╔════╝██║  ██║
 █████╗  ██╔██╗ ██║██║   ██║    ███████╗██║ █╗ ██║██║   ██║   ██║     ███████║
 ██╔══╝  ██║╚██╗██║╚██╗ ██╔╝    ╚════██║██║███╗██║██║   ██║   ██║     ██╔══██║
 ███████╗██║ ╚████║ ╚████╔╝     ███████║╚███╔███╔╝██║   ██║   ╚██████╗██║  ██║
 ╚══════╝╚═╝  ╚═══╝  ╚═══╝      ╚══════╝ ╚══╝╚══╝ ╚═╝   ╚═╝    ╚═════╝╚═╝  ╚═╝
`
)

// PersistentConfig stores saved paths per app
type PersistentConfig struct {
	Apps map[string]AppConfig `json:"apps"`
}

// AppConfig stores the saved paths for an app
type AppConfig struct {
	ConfigDir  string `json:"configDir"`
	TargetPath string `json:"targetPath"`
	LastEnv    string `json:"lastEnv"`
	UseJS      bool   `json:"useJS"`
}

// CLI states
type state int

const (
	stateSelectApp state = iota
	stateAppMenu
	stateInputConfigDir
	stateInputTargetPath
	stateInputEnv
	stateConfirm
	stateDone
	stateAddAppName
	stateAddAppConfigDir
	stateAddAppTargetPath
	stateAddAppUseJS
)

// Menu options for app menu
const (
	menuOptionSwitch = iota
	menuOptionEditPaths
	menuOptionDelete
)

// Model represents the application state
type model struct {
	state            state
	selectedApp      int
	apps             []string // App names from config
	configDir        string
	targetPath       string
	env              string
	useJS            bool
	textInput        textinput.Model
	err              error
	result           string
	quitting         bool
	persistentConfig PersistentConfig
	hasSavedConfig   bool
	menuOption       int
	newAppName       string
}

// getConfigPath returns the path to the persistent config file
func getConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ".envswitch-config.json"
	}
	return filepath.Join(homeDir, ".envswitch-config.json")
}

// loadPersistentConfig loads the saved configuration
func loadPersistentConfig() PersistentConfig {
	config := PersistentConfig{
		Apps: make(map[string]AppConfig),
	}

	data, err := os.ReadFile(getConfigPath())
	if err != nil {
		// Initialize with default app - The Vault uses JS
		config.Apps["The Vault"] = AppConfig{UseJS: true}
		return config
	}

	json.Unmarshal(data, &config)
	if config.Apps == nil {
		config.Apps = make(map[string]AppConfig)
		config.Apps["The Vault"] = AppConfig{UseJS: true}
	}

	// Ensure "The Vault" always has UseJS = true (legacy fix)
	if vault, exists := config.Apps["The Vault"]; exists {
		vault.UseJS = true
		config.Apps["The Vault"] = vault
	}

	return config
}

// savePersistentConfig saves the configuration to disk
func savePersistentConfig(config PersistentConfig) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(getConfigPath(), data, 0644)
}

// getAppNames returns sorted list of app names
func getAppNames(config PersistentConfig) []string {
	names := make([]string, 0, len(config.Apps)+1)
	for name := range config.Apps {
		names = append(names, name)
	}
	names = append(names, "➕ Add New App...")
	return names
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Type here..."
	ti.Focus()
	ti.CharLimit = 500
	ti.Width = 60

	persistentConfig := loadPersistentConfig()
	apps := getAppNames(persistentConfig)

	return model{
		state:            stateSelectApp,
		selectedApp:      0,
		apps:             apps,
		textInput:        ti,
		persistentConfig: persistentConfig,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "q":
			// Only quit if not in an input state
			if m.state == stateSelectApp || m.state == stateAppMenu {
				m.quitting = true
				return m, tea.Quit
			}
			// Allow 'q' in input fields
			var cmd tea.Cmd
			m.textInput, cmd = m.textInput.Update(msg)
			return m, cmd
		case "up", "k":
			switch m.state {
			case stateSelectApp:
				if m.selectedApp > 0 {
					m.selectedApp--
				}
			case stateAppMenu:
				if m.menuOption > 0 {
					m.menuOption--
				}
			case stateAddAppUseJS:
				m.useJS = !m.useJS
			}
		case "down", "j":
			switch m.state {
			case stateSelectApp:
				if m.selectedApp < len(m.apps)-1 {
					m.selectedApp++
				}
			case stateAppMenu:
				if m.menuOption < 2 {
					m.menuOption++
				}
			case stateAddAppUseJS:
				m.useJS = !m.useJS
			}
		case "enter":
			return m.handleEnter()
		case "esc":
			return m.handleEsc()
		}
	}

	// Handle text input for input states
	if m.isInputState() {
		var cmd tea.Cmd
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) isInputState() bool {
	return m.state == stateInputConfigDir ||
		m.state == stateInputTargetPath ||
		m.state == stateInputEnv ||
		m.state == stateAddAppName ||
		m.state == stateAddAppConfigDir ||
		m.state == stateAddAppTargetPath
}

func (m model) handleEsc() (tea.Model, tea.Cmd) {
	switch m.state {
	case stateAppMenu:
		m.state = stateSelectApp
		m.menuOption = 0
	case stateInputConfigDir:
		m.state = stateAppMenu
		m.menuOption = 0
	case stateInputTargetPath:
		m.state = stateInputConfigDir
		m.textInput.SetValue(m.configDir)
		m.textInput.Placeholder = "Config directory path..."
	case stateInputEnv:
		// If we have saved paths, go back to app menu, otherwise to target path
		if m.hasSavedConfig {
			m.state = stateAppMenu
			m.menuOption = 0
		} else {
			m.state = stateInputTargetPath
			m.textInput.SetValue(m.targetPath)
			m.textInput.Placeholder = "Target file path..."
		}
	case stateConfirm:
		m.state = stateInputEnv
		m.textInput.SetValue(m.env)
		m.textInput.Placeholder = "Environment name..."
	case stateAddAppName:
		m.state = stateSelectApp
	case stateAddAppConfigDir:
		m.state = stateAddAppName
		m.textInput.SetValue(m.newAppName)
		m.textInput.Placeholder = "App name..."
	case stateAddAppTargetPath:
		m.state = stateAddAppConfigDir
		m.textInput.SetValue(m.configDir)
		m.textInput.Placeholder = "Config directory path..."
	case stateAddAppUseJS:
		m.state = stateAddAppTargetPath
		m.textInput.SetValue(m.targetPath)
		m.textInput.Placeholder = "Target file path..."
	}
	return m, nil
}

func (m model) handleEnter() (tea.Model, tea.Cmd) {
	switch m.state {
	case stateSelectApp:
		appName := m.apps[m.selectedApp]

		// Check if "Add New App" was selected
		if appName == "➕ Add New App..." {
			m.state = stateAddAppName
			m.textInput.SetValue("")
			m.textInput.Placeholder = "Enter new app name..."
			m.textInput.Focus()
			m.newAppName = ""
			m.configDir = ""
			m.targetPath = ""
			m.useJS = false
			return m, textinput.Blink
		}

		// Load saved config for this app
		if savedConfig, exists := m.persistentConfig.Apps[appName]; exists {
			m.configDir = savedConfig.ConfigDir
			m.targetPath = savedConfig.TargetPath
			m.env = savedConfig.LastEnv
			m.useJS = savedConfig.UseJS
			m.hasSavedConfig = savedConfig.ConfigDir != "" && savedConfig.TargetPath != ""
		}

		// If we have saved paths, show menu. Otherwise go to config input
		if m.hasSavedConfig {
			m.state = stateAppMenu
			m.menuOption = 0
		} else {
			m.state = stateInputConfigDir
			m.textInput.SetValue("")
			m.textInput.Placeholder = "Paste the absolute path to config directory..."
			m.textInput.Focus()
			return m, textinput.Blink
		}
		return m, nil

	case stateAppMenu:
		switch m.menuOption {
		case menuOptionSwitch:
			// Quick switch - go directly to env input
			m.state = stateInputEnv
			m.textInput.SetValue(m.env)
			if m.env != "" {
				m.textInput.Placeholder = fmt.Sprintf("Enter to use '%s', or type new...", m.env)
			} else {
				m.textInput.Placeholder = "Environment name (e.g., test, stress, prod)..."
			}
			m.textInput.Focus()
			return m, textinput.Blink
		case menuOptionEditPaths:
			// Edit paths
			m.state = stateInputConfigDir
			m.textInput.SetValue(m.configDir)
			m.textInput.Placeholder = "Edit config directory path..."
			m.textInput.Focus()
			return m, textinput.Blink
		case menuOptionDelete:
			// Delete app
			appName := m.apps[m.selectedApp]
			delete(m.persistentConfig.Apps, appName)
			savePersistentConfig(m.persistentConfig)
			m.apps = getAppNames(m.persistentConfig)
			if m.selectedApp >= len(m.apps) {
				m.selectedApp = len(m.apps) - 1
			}
			m.state = stateSelectApp
			m.hasSavedConfig = false
			return m, nil
		}

	case stateInputConfigDir:
		value := strings.TrimSpace(m.textInput.Value())
		if value == "" {
			return m, nil
		}
		m.configDir = value
		m.state = stateInputTargetPath
		m.textInput.SetValue(m.targetPath)
		if m.targetPath != "" {
			m.textInput.Placeholder = "Edit target file path..."
		} else {
			m.textInput.Placeholder = "Paste the absolute path to target file..."
		}
		return m, textinput.Blink

	case stateInputTargetPath:
		value := strings.TrimSpace(m.textInput.Value())
		if value == "" {
			return m, nil
		}
		m.targetPath = value
		m.state = stateInputEnv
		m.textInput.SetValue(m.env)
		if m.env != "" {
			m.textInput.Placeholder = fmt.Sprintf("Enter to use '%s', or type new...", m.env)
		} else {
			m.textInput.Placeholder = "Environment name (e.g., test, stress, prod)..."
		}
		return m, textinput.Blink

	case stateInputEnv:
		value := strings.TrimSpace(m.textInput.Value())
		if value == "" && m.env == "" {
			return m, nil
		}
		if value != "" {
			m.env = value
		}
		m.state = stateConfirm
		return m, nil

	case stateConfirm:
		// Save the config before executing
		appName := m.apps[m.selectedApp]
		m.persistentConfig.Apps[appName] = AppConfig{
			ConfigDir:  m.configDir,
			TargetPath: m.targetPath,
			LastEnv:    m.env,
			UseJS:      m.useJS,
		}
		savePersistentConfig(m.persistentConfig)

		// Execute the switch
		err := executeSwitchNew(
			m.configDir,
			m.targetPath,
			m.env,
			m.useJS,
		)
		if err != nil {
			m.err = err
			m.result = fmt.Sprintf("❌ Error: %v", err)
		} else {
			m.result = fmt.Sprintf("✅ Successfully switched to %s environment!", m.env)
		}
		m.state = stateDone
		return m, nil

	case stateDone:
		m.quitting = true
		return m, tea.Quit

	// Add new app flow
	case stateAddAppName:
		value := strings.TrimSpace(m.textInput.Value())
		if value == "" {
			return m, nil
		}
		// Check if app already exists
		if _, exists := m.persistentConfig.Apps[value]; exists {
			m.err = fmt.Errorf("app '%s' already exists", value)
			return m, nil
		}
		m.newAppName = value
		m.state = stateAddAppConfigDir
		m.textInput.SetValue("")
		m.textInput.Placeholder = "Paste the absolute path to config directory..."
		return m, textinput.Blink

	case stateAddAppConfigDir:
		value := strings.TrimSpace(m.textInput.Value())
		if value == "" {
			return m, nil
		}
		m.configDir = value
		m.state = stateAddAppTargetPath
		m.textInput.SetValue("")
		m.textInput.Placeholder = "Paste the absolute path to target file..."
		return m, textinput.Blink

	case stateAddAppTargetPath:
		value := strings.TrimSpace(m.textInput.Value())
		if value == "" {
			return m, nil
		}
		m.targetPath = value
		m.state = stateAddAppUseJS
		m.useJS = false
		return m, nil

	case stateAddAppUseJS:
		// Save the new app
		m.persistentConfig.Apps[m.newAppName] = AppConfig{
			ConfigDir:  m.configDir,
			TargetPath: m.targetPath,
			UseJS:      m.useJS,
		}
		savePersistentConfig(m.persistentConfig)
		m.apps = getAppNames(m.persistentConfig)

		// Select the new app and go to env input
		for i, name := range m.apps {
			if name == m.newAppName {
				m.selectedApp = i
				break
			}
		}
		m.hasSavedConfig = true
		m.state = stateInputEnv
		m.textInput.SetValue("")
		m.textInput.Placeholder = "Environment name (e.g., test, stress, prod)..."
		m.textInput.Focus()
		return m, textinput.Blink
	}

	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return "\n  👋 See you later! VAMOS BOCA! 💙💛\n\n"
	}

	var s strings.Builder

	// Title with pixel art
	titleArt := lipgloss.NewStyle().
		Foreground(bocaGold).
		Bold(true).
		Render(pixelArt)
	s.WriteString(titleArt)
	s.WriteString("\n")

	// Disclaimer
	disclaimer := disclaimerStyle.Render("  ✨ Full vibe - it will not make any harm (i think) ✨")
	s.WriteString(disclaimer)
	s.WriteString("\n\n")

	// Main content based on state
	switch m.state {
	case stateSelectApp:
		s.WriteString(m.viewSelectApp())
	case stateAppMenu:
		s.WriteString(m.viewAppMenu())
	case stateInputConfigDir:
		s.WriteString(m.viewInputConfigDir())
	case stateInputTargetPath:
		s.WriteString(m.viewInputTargetPath())
	case stateInputEnv:
		s.WriteString(m.viewInputEnv())
	case stateConfirm:
		s.WriteString(m.viewConfirm())
	case stateDone:
		s.WriteString(m.viewDone())
	case stateAddAppName:
		s.WriteString(m.viewAddAppName())
	case stateAddAppConfigDir:
		s.WriteString(m.viewAddAppConfigDir())
	case stateAddAppTargetPath:
		s.WriteString(m.viewAddAppTargetPath())
	case stateAddAppUseJS:
		s.WriteString(m.viewAddAppUseJS())
	}

	return s.String()
}

func (m model) viewSelectApp() string {
	var s strings.Builder

	prompt := promptStyle.Render("  🏗️  In which app do you want to work?")
	s.WriteString(prompt)
	s.WriteString("\n\n")

	for i, appName := range m.apps {
		cursor := "  "
		style := normalStyle
		if i == m.selectedApp {
			cursor = "▸ "
			style = selectedStyle
		}

		line := fmt.Sprintf("%s%s", cursor, style.Render(appName))

		// Show saved info for configured apps
		if appName != "➕ Add New App..." {
			if savedConfig, exists := m.persistentConfig.Apps[appName]; exists && savedConfig.ConfigDir != "" {
				line += savedPathStyle.Render(fmt.Sprintf(" [%s]", savedConfig.LastEnv))
			} else {
				line += warningStyle.Render(" (not configured)")
			}
		}
		s.WriteString(line)
		s.WriteString("\n")
	}

	s.WriteString("\n")
	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#666666"))
	s.WriteString(helpStyle.Render("  ↑/↓: navigate • enter: select • q: quit"))
	s.WriteString("\n")

	return s.String()
}

func (m model) viewAppMenu() string {
	var s strings.Builder

	appName := m.apps[m.selectedApp]
	header := promptStyle.Render(fmt.Sprintf("  � %s", appName))
	s.WriteString(header)
	s.WriteString("\n\n")

	// Show current paths
	pathInfo := savedPathStyle.Render(fmt.Sprintf("  📁 %s\n  🎯 %s", m.configDir, m.targetPath))
	s.WriteString(pathInfo)
	s.WriteString("\n\n")

	menuOptions := []string{
		"🚀 Quick Switch (enter environment)",
		"✏️  Edit Paths",
		"🗑️  Delete App",
	}

	for i, opt := range menuOptions {
		cursor := "  "
		style := normalStyle
		if i == m.menuOption {
			cursor = "▸ "
			style = selectedStyle
		}
		s.WriteString(fmt.Sprintf("%s%s\n", cursor, style.Render(opt)))
	}

	s.WriteString("\n")
	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#666666"))
	s.WriteString(helpStyle.Render("  ↑/↓: navigate • enter: select • esc: back"))
	s.WriteString("\n")

	return s.String()
}

func (m model) viewInputConfigDir() string {
	var s strings.Builder

	appName := m.apps[m.selectedApp]
	header := promptStyle.Render(fmt.Sprintf("  � Config for: %s", appName))
	s.WriteString(header)
	s.WriteString("\n\n")

	prompt := lipgloss.NewStyle().Foreground(whiteColor).Render("  Config directory path:")
	s.WriteString(prompt)
	s.WriteString("\n\n  ")
	s.WriteString(m.textInput.View())
	s.WriteString("\n\n")

	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#666666"))
	s.WriteString(helpStyle.Render("  enter: confirm • esc: back"))
	s.WriteString("\n")

	return s.String()
}

func (m model) viewInputTargetPath() string {
	var s strings.Builder

	appName := m.apps[m.selectedApp]
	header := promptStyle.Render(fmt.Sprintf("  📁 Config for: %s", appName))
	s.WriteString(header)
	s.WriteString("\n\n")

	configInfo := lipgloss.NewStyle().Foreground(greenColor).Render(fmt.Sprintf("  ✓ Config dir: %s", m.configDir))
	s.WriteString(configInfo)
	s.WriteString("\n\n")

	prompt := lipgloss.NewStyle().Foreground(whiteColor).Render("  Target file path:")
	s.WriteString(prompt)
	s.WriteString("\n\n  ")
	s.WriteString(m.textInput.View())
	s.WriteString("\n\n")

	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#666666"))
	s.WriteString(helpStyle.Render("  enter: confirm • esc: back"))
	s.WriteString("\n")

	return s.String()
}

func (m model) viewInputEnv() string {
	var s strings.Builder

	appName := m.apps[m.selectedApp]
	header := promptStyle.Render(fmt.Sprintf("  🌍 Switch environment: %s", appName))
	s.WriteString(header)
	s.WriteString("\n\n")

	// Show paths
	pathInfo := savedPathStyle.Render(fmt.Sprintf("  � %s\n  🎯 %s", m.configDir, m.targetPath))
	s.WriteString(pathInfo)
	s.WriteString("\n\n")

	prompt := lipgloss.NewStyle().Foreground(whiteColor).Render("  Enter environment name:")
	s.WriteString(prompt)
	s.WriteString("\n\n  ")
	s.WriteString(m.textInput.View())
	s.WriteString("\n\n")

	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#666666"))
	s.WriteString(helpStyle.Render("  enter: confirm • esc: back"))
	s.WriteString("\n")

	return s.String()
}

func (m model) viewConfirm() string {
	var s strings.Builder

	header := promptStyle.Render("  📋 Confirm switch:")
	s.WriteString(header)
	s.WriteString("\n\n")

	appName := m.apps[m.selectedApp]

	info := boxStyle.Render(fmt.Sprintf(
		"  App:         %s\n"+
			"  Environment: %s\n"+
			"  Config Dir:  %s\n"+
			"  Target:      %s\n"+
			"  Use JS:      %t",
		appName, m.env, m.configDir, m.targetPath, m.useJS,
	))
	s.WriteString(info)
	s.WriteString("\n\n")

	s.WriteString(lipgloss.NewStyle().Bold(true).Foreground(bocaGold).Render("  Press ENTER to execute!"))
	s.WriteString("\n\n")

	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#666666"))
	s.WriteString(helpStyle.Render("  enter: execute • esc: go back"))
	s.WriteString("\n")

	return s.String()
}

func (m model) viewDone() string {
	var s strings.Builder

	if m.err != nil {
		errorStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF0000"))
		s.WriteString(errorStyle.Render(m.result))
	} else {
		s.WriteString(successStyle.Render(m.result))
	}
	s.WriteString("\n\n")

	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#666666"))
	s.WriteString(helpStyle.Render("  Press ENTER to exit"))
	s.WriteString("\n")

	return s.String()
}

func (m model) viewAddAppName() string {
	var s strings.Builder

	header := promptStyle.Render("  ➕ Add New App")
	s.WriteString(header)
	s.WriteString("\n\n")

	prompt := lipgloss.NewStyle().Foreground(whiteColor).Render("  Enter app name:")
	s.WriteString(prompt)
	s.WriteString("\n\n  ")
	s.WriteString(m.textInput.View())

	if m.err != nil {
		s.WriteString("\n\n")
		errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
		s.WriteString(errorStyle.Render(fmt.Sprintf("  ⚠️  %v", m.err)))
	}

	s.WriteString("\n\n")
	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#666666"))
	s.WriteString(helpStyle.Render("  enter: confirm • esc: cancel"))
	s.WriteString("\n")

	return s.String()
}

func (m model) viewAddAppConfigDir() string {
	var s strings.Builder

	header := promptStyle.Render(fmt.Sprintf("  ➕ Add New App: %s", m.newAppName))
	s.WriteString(header)
	s.WriteString("\n\n")

	prompt := lipgloss.NewStyle().Foreground(whiteColor).Render("  Config directory path:")
	s.WriteString(prompt)
	s.WriteString("\n\n  ")
	s.WriteString(m.textInput.View())
	s.WriteString("\n\n")

	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#666666"))
	s.WriteString(helpStyle.Render("  enter: confirm • esc: back"))
	s.WriteString("\n")

	return s.String()
}

func (m model) viewAddAppTargetPath() string {
	var s strings.Builder

	header := promptStyle.Render(fmt.Sprintf("  ➕ Add New App: %s", m.newAppName))
	s.WriteString(header)
	s.WriteString("\n\n")

	configInfo := lipgloss.NewStyle().Foreground(greenColor).Render(fmt.Sprintf("  ✓ Config dir: %s", m.configDir))
	s.WriteString(configInfo)
	s.WriteString("\n\n")

	prompt := lipgloss.NewStyle().Foreground(whiteColor).Render("  Target file path:")
	s.WriteString(prompt)
	s.WriteString("\n\n  ")
	s.WriteString(m.textInput.View())
	s.WriteString("\n\n")

	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#666666"))
	s.WriteString(helpStyle.Render("  enter: confirm • esc: back"))
	s.WriteString("\n")

	return s.String()
}

func (m model) viewAddAppUseJS() string {
	var s strings.Builder

	header := promptStyle.Render(fmt.Sprintf("  ➕ Add New App: %s", m.newAppName))
	s.WriteString(header)
	s.WriteString("\n\n")

	configInfo := lipgloss.NewStyle().Foreground(greenColor).Render(fmt.Sprintf("  ✓ Config dir: %s", m.configDir))
	s.WriteString(configInfo)
	s.WriteString("\n")
	targetInfo := lipgloss.NewStyle().Foreground(greenColor).Render(fmt.Sprintf("  ✓ Target: %s", m.targetPath))
	s.WriteString(targetInfo)
	s.WriteString("\n\n")

	prompt := lipgloss.NewStyle().Foreground(whiteColor).Render("  Use JavaScript config files?")
	s.WriteString(prompt)
	s.WriteString("\n\n")

	options := []struct {
		label    string
		selected bool
	}{
		{"Yes (.js files)", m.useJS},
		{"No (.json files)", !m.useJS},
	}

	for i, opt := range options {
		cursor := "  "
		style := normalStyle
		if (i == 0 && m.useJS) || (i == 1 && !m.useJS) {
			cursor = "▸ "
			style = selectedStyle
		}
		s.WriteString(fmt.Sprintf("%s%s\n", cursor, style.Render(opt.label)))
	}

	s.WriteString("\n")
	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#666666"))
	s.WriteString(helpStyle.Render("  ↑/↓: toggle • enter: confirm • esc: back"))
	s.WriteString("\n")

	return s.String()
}

// executeSwitchNew runs the actual environment switch
func executeSwitchNew(configDir, targetPath, env string, useJS bool) error {
	var config *Config
	var err error
	var configPath string

	if useJS {
		configPath = filepath.Join(configDir, fmt.Sprintf("config.%s.js", env))
	} else {
		configPath = filepath.Join(configDir, fmt.Sprintf("config.%s.json", env))
	}

	// Check if config file exists first
	if _, statErr := os.Stat(configPath); os.IsNotExist(statErr) {
		return fmt.Errorf("config file not found: %s", configPath)
	}

	if useJS {
		config, err = LoadConfigFromJS(configPath)
	} else {
		config, err = loadConfig(configPath)
	}

	if err != nil {
		return fmt.Errorf("loading config %s: %v", configPath, err)
	}

	// Check if target file exists
	if _, statErr := os.Stat(targetPath); os.IsNotExist(statErr) {
		return fmt.Errorf("target file not found: %s", targetPath)
	}

	content, err := os.ReadFile(targetPath)
	if err != nil {
		return fmt.Errorf("reading target file %s: %v", targetPath, err)
	}

	result := applyReplacements(string(content), config, false)

	err = os.WriteFile(targetPath, []byte(result), 0644)
	if err != nil {
		return fmt.Errorf("writing target file %s: %v", targetPath, err)
	}

	return nil
}

// RunInteractiveCLI starts the interactive CLI
func RunInteractiveCLI() error {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	_, err := p.Run()
	return err
}

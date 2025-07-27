// version 0.3.1
package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	appStyle       = lipgloss.NewStyle().Margin(1, 2)
	headerStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("42"))
	successStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true)
	errorStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("197")).Bold(true)
	statusStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	commandStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("86")).Bold(true)
	directoryStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Italic(true)
)

var (
	debugMode  bool
	targetDir  string
	projectDir string
)

var initCmd = &cobra.Command{
	Use:   "init [project-name]",
	Short: "Create a new DocuBook project",
	Long:  "Create a new DocuBook documentation site with a modern and clean design.",
	Args:  cobra.MaximumNArgs(1),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if debug, _ := cmd.Flags().GetBool("debug"); debug {
			debugMode = true
		}
		projectName := "docs"
		if len(args) > 0 {
			projectName = args[0]
		}
		if targetDir == "" {
			targetDir = "."
		}
		absPath, err := filepath.Abs(targetDir)
		if err != nil {
			fmt.Printf("❌ Error resolving path: %v\n", err)
			os.Exit(1)
		}
		projectDir = filepath.Join(absPath, projectName)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		p := tea.NewProgram(initialModel())
		finalModel, err := p.Run()
		if err != nil {
			return fmt.Errorf("error running program: %w", err)
		}
		if m, ok := finalModel.(model); ok && m.err != nil {
			return m.err
		}
		return nil
	},
}

type model struct {
	width       int
	height      int
	progress    float64
	status      string
	err         error
	done        bool
	steps       []step
	currentStep int
}

type step struct {
	title   string
	command tea.Cmd
	done    bool
	success bool
	status  string
}

func initialModel() model {
	return model{
		status:      "Starting setup...",
		currentStep: 0,
		steps: []step{
			{title: "🚀 Setting up project directory...", status: "Creating directory...", command: createProjectDir},
			{title: "📦 Initializing project...", status: "Running npm init...", command: initProject},
			{title: "✨ Installing dependencies...", status: "Running npm install...", command: installDepsAndSetup},
		},
	}
}

func (m model) Init() tea.Cmd {
	return m.steps[0].command
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, nil
	case errMsg:
		m.err = msg
		m.steps[m.currentStep].success = false
		m.steps[m.currentStep].done = true
		return m, tea.Quit
	case stepCompleteMsg:
		m.steps[m.currentStep].done = true
		m.steps[m.currentStep].success = true
		m.currentStep++
		m.progress = float64(m.currentStep) / float64(len(m.steps))
		if m.currentStep >= len(m.steps) {
			m.done = true
			m.status = "✅ Setup complete!"
			return m, tea.Quit
		}
		m.status = m.steps[m.currentStep].status
		return m, m.steps[m.currentStep].command
	}
	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("\n%s\n\n", errorStyle.Render("Error: "+m.err.Error()))
	}
	var sb strings.Builder
	sb.WriteString(headerStyle.Render("✨ DocuBook CLI Setup") + "\n\n")
	for i, s := range m.steps {
		icon, style := "○", statusStyle
		if s.done {
			if s.success {
				icon, style = "✓", successStyle
			} else {
				icon, style = "✗", errorStyle
			}
		}
		sb.WriteString(fmt.Sprintf("  %s %s\n", style.Render(icon), s.title))
		if !s.done && i == m.currentStep {
			sb.WriteString(fmt.Sprintf("    %s\n", statusStyle.Render(m.status)))
		}
	}
	sb.WriteString("\n")
	if !m.done {
		progressBar := progress.New(progress.WithWidth(min(60, m.width-4)), progress.WithGradient("#5A56E0", "#EE6FF8"))
		sb.WriteString(progressBar.ViewAs(m.progress) + "\n\n")
		sb.WriteString(statusStyle.Render("Press 'q' to quit."))
	} else {
		sb.WriteString(successStyle.Render("✓ All done!") + "\n\n")
		sb.WriteString("Next steps:\n")
		sb.WriteString(fmt.Sprintf("  %s %s\n", commandStyle.Render("cd"), directoryStyle.Render(filepath.Base(projectDir))))
		sb.WriteString(fmt.Sprintf("  %s\n", commandStyle.Render("create")))
	}
	return appStyle.Render(sb.String())
}

type (
	errMsg          error
	stepCompleteMsg struct{}
)

func createProjectDir() tea.Msg {
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return errMsg(fmt.Errorf("failed to create project directory: %w", err))
	}
	if err := os.Chdir(projectDir); err != nil {
		return errMsg(fmt.Errorf("failed to change to project directory: %w", err))
	}
	return stepCompleteMsg{}
}

func initProject() tea.Msg {
	if err := runCommand(context.Background(), "npm", "init", "-y"); err != nil {
		return errMsg(fmt.Errorf("failed to initialize project (is npm installed?): %w", err))
	}
	return stepCompleteMsg{}
}

func installDepsAndSetup() tea.Msg {
	if err := runCommand(context.Background(), "npm", "install", "--silent", "--no-progress", "@docubook/create"); err != nil {
		return errMsg(fmt.Errorf("failed to install dependencies: %w", err))
	}
	if err := runCommand(context.Background(), "npx", "--yes", "--quiet", "@docubook/create@latest", "."); err != nil {
		return errMsg(fmt.Errorf("failed to setup project: %w", err))
	}
	return stepCompleteMsg{}
}

func runCommand(ctx context.Context, name string, args ...string) error {
	cmd := exec.CommandContext(ctx, name, args...)
	var stdout, stderr bytes.Buffer
	var writers []io.Writer
	writers = append(writers, &stdout, &stderr)
	if debugMode {
		writers = append(writers, os.Stdout)
	}
	cmd.Stdout = io.MultiWriter(writers...)
	cmd.Stderr = io.MultiWriter(writers...)
	err := cmd.Run()
	if err != nil {
		output := strings.TrimSpace(stdout.String() + "\n" + stderr.String())
		if output != "" {
			return fmt.Errorf("command '%s %s' failed: %w\n--- Output ---\n%s", name, strings.Join(args, " "), err, output)
		}
		return fmt.Errorf("command '%s %s' failed: %w", name, strings.Join(args, " "), err)
	}
	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func init() {
	initCmd.Flags().BoolVarP(&debugMode, "debug", "d", false, "Enable debug mode")
	initCmd.Flags().StringVarP(&targetDir, "dir", "D", ".", "Directory to create the project in")
	rootCmd.AddCommand(initCmd)
}

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

// Styles
var (
	appStyle = lipgloss.NewStyle().Margin(1, 2)

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("42")) // Green

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")). // Green
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("197")). // Red
			Bold(true)

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")) // Gray

	commandStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")). // Purple
			Bold(true)

	directoryStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")). // Green
			Italic(true)
)

// Global flags, set by Cobra
var (
	debugMode  bool
	targetDir  string
	projectDir string
)

// createCmd defines the "create" command
var createCmd = &cobra.Command{
	Use:   "create [project-name]",
	Short: "Create a new DocuBook project",
	Long:  "Create a new DocuBook documentation site with a modern and clean design.",
	Args:  cobra.MaximumNArgs(1),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Initialize debug mode if flag is set
		if debug, _ := cmd.Flags().GetBool("debug"); debug {
			debugMode = true
		}

		// Set up project directory name, defaulting to "docs"
		projectName := "docs"
		if len(args) > 0 {
			projectName = args[0]
		}

		// Resolve the full, absolute path for the project
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
		// Initialize and run the Bubble Tea program
		p := tea.NewProgram(initialModel())
		model, err := p.Run()
		if err != nil {
			return fmt.Errorf("error running program: %w", err)
		}

		// Check if the final model has an error state
		if finalModel, ok := model.(model); ok && finalModel.err != nil {
			return finalModel.err
		}
		return nil
	},
}

// TUI state model
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

// A single step in the setup process
type step struct {
	title    string
	command  tea.Cmd
	done     bool
	success  bool
	status   string
	complete string
}

// initialModel sets up the initial state of the TUI
func initialModel() model {
	return model{
		status:      "Starting setup...",
		currentStep: 0,
		steps: []step{
			{
				title:    "🚀 Setting up project directory...",
				status:   "Creating directory...",
				complete: "Project directory created",
				command:  createProjectDir,
			},
			{
				title:    "📦 Initializing project...",
				status:   "Running npm init...",
				complete: "Project initialized",
				command:  initProject,
			},
			{
				title:    "✨ Installing dependencies and scaffolding project...",
				status:   "Running npm install and npx...",
				complete: "Dependencies installed and project scaffolded",
				command:  installDepsAndSetup,
			},
		},
	}
}

// Init is the first command that is run when the program starts
func (m model) Init() tea.Cmd {
	// Start the first step
	return m.steps[0].command
}

// Update handles messages and updates the model
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case errMsg:
		m.err = msg
		m.steps[m.currentStep].success = false
		m.steps[m.currentStep].done = true
		return m, tea.Quit

	case stepCompleteMsg:
		// Mark the current step as done and successful
		m.steps[m.currentStep].done = true
		m.steps[m.currentStep].success = true

		// Move to the next step
		m.currentStep++
		m.progress = float64(m.currentStep) / float64(len(m.steps))

		// If all steps are done, we're finished
		if m.currentStep >= len(m.steps) {
			m.done = true
			m.status = "✅ Setup complete!"
			return m, tea.Quit
		}

		// Start the next step
		m.status = m.steps[m.currentStep].status
		return m, m.steps[m.currentStep].command
	}

	return m, nil
}

// View renders the UI
func (m model) View() string {
	if m.err != nil {
		// On error, show the error message and quit.
		return fmt.Sprintf("\n%s\n\n", errorStyle.Render("Error: "+m.err.Error()))
	}

	var sb strings.Builder

	// Header
	sb.WriteString(headerStyle.Render("✨ DocuBook CLI Setup"))
	sb.WriteString("\n\n")

	// Steps
	for i, s := range m.steps {
		var icon string
		var style lipgloss.Style

		if s.done {
			if s.success {
				icon = "✓"
				style = successStyle
			} else {
				icon = "✗"
				style = errorStyle
			}
			sb.WriteString(fmt.Sprintf("  %s %s\n", style.Render(icon), s.title))
		} else if i == m.currentStep {
			icon = "○"
			style = statusStyle
			sb.WriteString(fmt.Sprintf("  %s %s\n", style.Render(icon), s.title))
			sb.WriteString(fmt.Sprintf("    %s\n", statusStyle.Render(m.status)))
		} else {
			icon = "○"
			style = statusStyle
			sb.WriteString(fmt.Sprintf("  %s %s\n", style.Render(icon), s.title))
		}
	}

	sb.WriteString("\n")

	// Progress bar or final message
	if !m.done {
		progressBar := progress.New(
			progress.WithWidth(min(60, m.width-4)),
			progress.WithGradient("#5A56E0", "#EE6FF8"),
		)
		sb.WriteString(progressBar.ViewAs(m.progress) + "\n\n")
		sb.WriteString(statusStyle.Render("Press 'q' to quit."))
	} else {
		sb.WriteString(successStyle.Render("✓ All done!"))
		sb.WriteString("\n\n")

		// Show next steps
		sb.WriteString("Next steps:\n")
		sb.WriteString(fmt.Sprintf("  %s %s\n",
			commandStyle.Render("cd"),
			directoryStyle.Render(filepath.Base(projectDir))))
		sb.WriteString(fmt.Sprintf("  %s\n", commandStyle.Render("docubook create")))
	}

	return appStyle.Render(sb.String())
}

// Messages for TUI communication
type (
	errMsg          error
	stepCompleteMsg struct{}
)

// Commands for each setup step
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

// runCommand executes a shell command with the given context and arguments
func runCommand(ctx context.Context, name string, args ...string) error {
	cmd := exec.CommandContext(ctx, name, args...)

	// Always capture output for better error reporting
	var stdout, stderr bytes.Buffer
	var writers []io.Writer
	writers = append(writers, &stdout, &stderr)

	// Also show live command output in debug mode
	if debugMode {
		writers = append(writers, os.Stdout)
	}
	cmd.Stdout = io.MultiWriter(writers...)
	cmd.Stderr = io.MultiWriter(writers...)

	err := cmd.Run()
	if err != nil {
		// Include command output in error message for better debugging
		output := strings.TrimSpace(stdout.String() + "\n" + stderr.String())
		if output != "" {
			return fmt.Errorf("command '%s %s' failed: %w\n--- Output ---\n%s", name, strings.Join(args, " "), err, output)
		}
		return fmt.Errorf("command '%s %s' failed: %w", name, strings.Join(args, " "), err)
	}

	return nil
}

// min returns the smaller of two integers.
// Added for clarity and compatibility with Go versions < 1.21.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func init() {
	createCmd.Flags().BoolVarP(&debugMode, "debug", "d", false, "Enable debug mode (shows detailed output)")
	createCmd.Flags().StringVarP(&targetDir, "dir", "D", ".", "Directory to create the project in")
	rootCmd.AddCommand(createCmd)
}

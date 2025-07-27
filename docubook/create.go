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
	"time"

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
			Foreground(lipgloss.Color("42"))

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("197")).
			Bold(true)

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			MarginTop(1)

	commandStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")).
			Bold(true)

	directoryStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")).
			Italic(true)
)

var (
	debugMode  bool
	targetDir  string
	projectDir string
)

var createCmd = &cobra.Command{
	Use:   "cli [project-name]",
	Short: "Create a new DocuBook project",
	Long:  "Create a new DocuBook documentation site with a modern and clean design.",
	Args:  cobra.MaximumNArgs(1),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Initialize debug mode if flag is set
		if debug, _ := cmd.Flags().GetBool("debug"); debug {
			debugMode = true
		}

		// Set up project directory
		projectName := "cli-docubook"
		if len(args) > 0 {
			projectName = args[0]
		}

		// Resolve full path
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
		if _, err := p.Run(); err != nil {
			return fmt.Errorf("error running program: %w", err)
		}
		return nil
	},
}

type model struct {
	progress float64
	status   string
	err      error
	width    int
	height   int
	steps    []step
}

type step struct {
	text    string
	done    bool
	success bool
}

func initialModel() model {
	return model{
		progress: 0,
		status:   "Starting setup...",
		steps: []step{
			{text: "🚀 Setting up project directory..."},
			{text: "📦 Initializing project..."},
			{text: "✨ Installing dependencies..."},
		},
	}
}

func (m model) Init() tea.Cmd {
	// No need to modify m here as it's a value receiver
	// The steps are already initialized in initialModel()
	return tea.Batch(
		setupProject(),
	)
}

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

	case progressMsg:
		m.progress = float64(msg)
		// Only update progress if we're not done yet
		if m.progress < 1.0 {
			return m, tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
				return progressMsg(m.progress + 0.01) // Small increment for smooth animation
			})
		}
		return m, nil

	case statusMsg:
		m.status = string(msg)
		return m, nil

	case errMsg:
		m.err = msg
		return m, tea.Quit

	case stepCompleteMsg:
		if msg.step < len(m.steps) {
			m.steps[msg.step].done = true
			m.steps[msg.step].success = msg.success
			// Update progress based on completed steps
			m.progress = float64(msg.step+1) / float64(len(m.steps))

			// Update status based on current step
			switch msg.step {
			case 0:
				m.status = "Project directory created"
			case 1:
				m.status = "Project initialized"
			case 2:
				m.status = "Dependencies installed"
			}

			// If all steps are done, we're finished
			if msg.step == len(m.steps)-1 {
				m.status = "✅ Setup complete!"
				m.progress = 1.0
				return m, tea.Quit
			}

			// Return a command to update the progress
			return m, tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
				return progressMsg(m.progress)
			})
		}
		return m, nil

	}

	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("\n%s\n\n", errorStyle.Render("Error: "+m.err.Error()))
	}

	var sb strings.Builder

	// Header
	sb.WriteString(headerStyle.Render("✨ DocuBook CLI"))
	sb.WriteString("\n\n")

	// Steps
	for i, step := range m.steps {
		icon := "○"
		if step.done {
			icon = "✓"
		}

		style := statusStyle
		if step.done && step.success {
			style = style.Foreground(lipgloss.Color("42"))
		} else if step.done && !step.success {
			style = style.Foreground(lipgloss.Color("197"))
		}

		sb.WriteString(fmt.Sprintf("  %s %s\n", style.Render(icon), step.text))

		// Show progress for current step
		if !step.done && i == int(float64(len(m.steps)-1)*m.progress) {
			sb.WriteString(fmt.Sprintf("    %s\n", statusStyle.Render(m.status)))
		}
	}

	sb.WriteString("\n")

	// Progress bar
	if m.progress < 1.0 {
		progressBar := progress.New(
			progress.WithWidth(min(50, m.width-4)),
			progress.WithGradient("#5A56E0", "#EE6FF8"),
		)
		sb.WriteString(progressBar.ViewAs(m.progress) + "\n\n")
	} else {
		sb.WriteString(successStyle.Render("✓ All done!"))
		sb.WriteString("\n\n")

		// Show next steps
		sb.WriteString("Next steps:\n")
		sb.WriteString(fmt.Sprintf("  %s %s\n",
			commandStyle.Render("cd"),
			directoryStyle.Render(filepath.Base(projectDir))))
		sb.WriteString(fmt.Sprintf("  %s\n", commandStyle.Render("create")))
	}

	return appStyle.Render(sb.String())
}

// Messages
type (
	progressMsg     float64
	statusMsg       string
	errMsg          error
	stepCompleteMsg struct {
		step    int
		success bool
	}
)

// Commands

func setupProject() tea.Cmd {
	return tea.Sequence(
		func() tea.Msg {
			return statusMsg("Starting DocuBook setup...")
		},
		// Step 1: Create project directory
		func() tea.Msg {
			if err := os.MkdirAll(projectDir, 0755); err != nil {
				return errMsg(fmt.Errorf("failed to create project directory: %w", err))
			}
			if err := os.Chdir(projectDir); err != nil {
				return errMsg(fmt.Errorf("failed to change to project directory: %w", err))
			}
			return stepCompleteMsg{step: 0, success: true}
		},
		// Step 2: Initialize project
		func() tea.Msg {
			if err := runCommand(context.Background(), "npm", "init", "-y"); err != nil {
				return errMsg(fmt.Errorf("failed to initialize project: %w", err))
			}
			return stepCompleteMsg{step: 1, success: true}
		},
		// Step 3: Install dependencies and setup project
		func() tea.Msg {
			if err := runCommand(context.Background(), "npm", "install", "--silent", "--no-progress", "@docubook/create"); err != nil {
				return errMsg(fmt.Errorf("failed to install dependencies: %w", err))
			}

			if err := runCommand(context.Background(), "npx", "--yes", "--quiet", "@docubook/create@latest", "."); err != nil {
				return errMsg(fmt.Errorf("failed to setup project: %w", err))
			}
			return stepCompleteMsg{step: 2, success: true}
		},
	)
}

// runCommand executes a shell command with the given context and arguments
func runCommand(ctx context.Context, name string, args ...string) error {
	cmd := exec.CommandContext(ctx, name, args...)

	// Always capture output for better error reporting
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Also show output in debug mode
	if debugMode {
		cmd.Stdout = io.MultiWriter(os.Stdout, &stdout)
		cmd.Stderr = io.MultiWriter(os.Stderr, &stderr)
	}

	err := cmd.Run()
	if err != nil {
		// Include command output in error message for better debugging
		output := strings.TrimSpace(stdout.String() + "\n" + stderr.String())
		if output != "" {
			return fmt.Errorf("%w\nOutput:\n%s", err, output)
		}
	}

	return err
}

func init() {
	createCmd.Flags().BoolVarP(&debugMode, "debug", "d", false, "Enable debug mode (shows detailed output)")
	createCmd.Flags().StringVarP(&targetDir, "dir", "D", ".", "Directory to create the project in")
	rootCmd.AddCommand(createCmd)
}

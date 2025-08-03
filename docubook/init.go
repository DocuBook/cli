// version 0.3.5
package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

// --- Styles ---
var (
	// Define styles for various UI elements
	appStyle = lipgloss.NewStyle().Margin(1, 2)

	// Style for the large, centered welcome message
	welcomeHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FAFAFA")).
				Background(lipgloss.Color("#865DFF")). // Solid purple background
				Padding(1, 4).
				MarginBottom(1)

	// Style for the website link
	websiteStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("39")). // Blue
			MarginBottom(2)

	// Style for the check header
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("42")) // Green

	// Style for the checkmark
	checkStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")) // Green

	// Style for the cross mark
	crossStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("197")) // Red

	// Style for inactive/waiting items
	inactiveStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")) // Gray

	// Style for suggested commands
	commandStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")). // Purple
			Bold(true)

	// Style for installation links
	linkStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("39")). // Blue
			Underline(true)

	// Base style for the help text
	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
)

// --- Cobra Command ---
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Checks for required tools to create a DocuBook project",
	Long:  "Verifies if Node.js (npm, pnpm, yarn) or Bun is installed and provides the next steps.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		p := tea.NewProgram(initialModel())
		if _, err := p.Run(); err != nil {
			return fmt.Errorf("error running program: %w", err)
		}
		return nil
	},
}

// --- Bubble Tea Model & Logic ---

// Define UI stages to control the flow
type uiStage int

const (
	welcomeStage  uiStage = iota // Stage 0: Displaying the welcome message
	checkingStage                // Stage 1: Displaying the check process
)

// packageManager represents a tool we'll check
type packageManager struct {
	name       string
	executable string
	installCmd string
	installed  bool
}

// startCheckingMsg is a message to transition from welcome to checking
type startCheckingMsg struct{}

// checkResultMsg is sent when a check is complete
type checkResultMsg struct {
	index     int
	installed bool
}

// model is the state of our TUI application
type model struct {
	stage    uiStage // Current UI stage
	spinner  spinner.Model
	packages []packageManager
	index    int
	done     bool
}

// initialModel prepares the initial state
func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Line
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return model{
		stage:   welcomeStage, // Start from the welcome stage
		spinner: s,
		packages: []packageManager{
			{name: "npm", executable: "npm", installCmd: "npx @docubook/create@latest"},
			{name: "pnpm", executable: "pnpm", installCmd: "pnpm dlx @docubook/create@latest"},
			{name: "yarn", executable: "yarn", installCmd: "yarn dlx @docubook/create@latest"},
			{name: "bun", executable: "bun", installCmd: "bunx @docubook/create@latest"},
		},
		index: 0,
	}
}

// command to start checking after a delay
func startChecking() tea.Cmd {
	// Pause for 2.5 seconds before starting the check
	return tea.Tick(2500*time.Millisecond, func(t time.Time) tea.Msg {
		return startCheckingMsg{}
	})
}

// command to check for a program's existence
func checkInstalled(pm packageManager, index int) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(500 * time.Millisecond)
		_, err := exec.LookPath(pm.executable)
		return checkResultMsg{index: index, installed: err == nil}
	}
}

// Init is the first command to run
func (m model) Init() tea.Cmd {
	// Start by showing the welcome screen, then send a message to start checking
	return startChecking()
}

// Update handles all messages and updates the state
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	case spinner.TickMsg:
		// Only update the spinner if in the checking stage
		if m.stage == checkingStage {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}

	case startCheckingMsg:
		// Switch to the checking stage and start the spinner & first check
		m.stage = checkingStage
		return m, tea.Batch(
			m.spinner.Tick,
			checkInstalled(m.packages[0], 0),
		)

	case checkResultMsg:
		// Update the status of the package that was just checked
		m.packages[msg.index].installed = msg.installed

		// Move to the next package
		m.index++
		if m.index >= len(m.packages) {
			// All checks are done
			m.done = true
			return m, tea.Quit
		}

		// Start checking the next package
		return m, checkInstalled(m.packages[m.index], m.index)
	}

	return m, nil
}

// View renders the UI based on the current state
func (m model) View() string {
	// If in the welcome stage
	if m.stage == welcomeStage {
		welcomeText := welcomeHeaderStyle.Render("DocuBook CLI")
		websiteText := websiteStyle.Render("Find more information and guides at https://docubook.pro")
		statusText := inactiveStyle.Render("Starting environment check...")

		return appStyle.Render(welcomeText + "\n" + websiteText + "\n" + statusText)
	}

	// View for the checking and done stages
	var b strings.Builder

	if !m.done {
		b.WriteString(headerStyle.Render("Checking for required tools...") + "\n\n")
		// Render progress list while checking
		for i := 0; i < len(m.packages); i++ {
			if i == m.index {
				b.WriteString(fmt.Sprintf("  %s Checking for: %s...\n", m.spinner.View(), m.packages[i].name))
			} else if i < m.index {
				if m.packages[i].installed {
					b.WriteString(fmt.Sprintf("  %s Found: %s\n", checkStyle.Render("✓"), m.packages[i].name))
				} else {
					b.WriteString(fmt.Sprintf("  %s Not found: %s\n", crossStyle.Render("✗"), m.packages[i].name))
				}
			} else {
				b.WriteString(fmt.Sprintf("  %s Waiting for: %s\n", inactiveStyle.Render("?"), m.packages[i].name))
			}
		}
	} else {
		// Render final results when done
		b.WriteString(headerStyle.Render("Check complete. Results:") + "\n\n")

		// Render the final list of results
		var availableCmds []string
		for _, pkg := range m.packages {
			if pkg.installed {
				b.WriteString(fmt.Sprintf("  %s Found %s\n", checkStyle.Render("✓"), pkg.name))
				availableCmds = append(availableCmds, pkg.installCmd)
			} else {
				b.WriteString(fmt.Sprintf("  %s Not Found %s\n", crossStyle.Render("✗"), pkg.name))
			}
		}
		b.WriteString("\n")

		if len(availableCmds) > 0 {
			b.WriteString(checkStyle.Render("✅ Ready to go!"))
			b.WriteString("\n\nYou can create a new project with one of these commands:\n\n")
			for _, cmd := range availableCmds {
				b.WriteString(fmt.Sprintf("  %s\n", commandStyle.Render(cmd)))
			}
		} else {
			b.WriteString(crossStyle.Render("❌ Required tools not found."))
			b.WriteString("\n\nPlease install Node.js or Bun to continue:\n\n")
			b.WriteString(fmt.Sprintf("  - Node.js (includes npm): %s\n", linkStyle.Render("https://nodejs.org/")))
			b.WriteString(fmt.Sprintf("  - Bun: %s\n", linkStyle.Render("https://bun.sh/")))
		}
	}

	b.WriteString(helpStyle.Render("\nPress 'q' to quit."))
	return appStyle.Render(b.String())
}

// init() function is executed when the package is loaded
func init() {
	rootCmd.AddCommand(initCmd)
}

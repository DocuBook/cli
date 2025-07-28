// version 0.3.2
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

// Styles
var (
	// Definisikan gaya untuk berbagai elemen UI
	appStyle = lipgloss.NewStyle().Margin(1, 2)

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("42")) // Hijau

	// Gaya untuk teks yang berhasil
	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")) // Hijau

	// Gaya untuk teks error atau tidak ditemukan
	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("197")) // Merah

	// Gaya untuk teks status atau abu-abu
	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")) // Abu-abu

	// Gaya untuk perintah yang disarankan
	commandStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")). // Ungu
			Bold(true)

	// Gaya untuk link
	linkStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("39")). // Biru
			Underline(true)
)

// initCmd mendefinisikan perintah 'init'
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

// packageManager merepresentasikan satu tool yang akan kita cek
type packageManager struct {
	name       string // Nama yang ditampilkan (e.g., "npm")
	executable string // Nama command yang akan dicek (e.g., "npm")
	installCmd string // Perintah instalasi yang disarankan
	checking   bool   // True jika sedang dicek
	installed  bool   // True jika ditemukan
	checked    bool   // True jika sudah selesai dicek
}

// checkResultMsg adalah pesan yang dikirim setelah pengecekan selesai
type checkResultMsg struct {
	index     int
	installed bool
}

// model adalah state dari aplikasi TUI kita
type model struct {
	spinner  spinner.Model
	packages []packageManager
	results  []string // Menyimpan perintah instalasi yang valid
	index    int      // Index dari package yang sedang dicek
	done     bool     // True jika semua pengecekan selesai
}

// initialModel menyiapkan state awal
func initialModel() model {
	s := spinner.New()
	// Mengubah gaya spinner agar mirip dengan contoh
	s.Spinner = spinner.Line
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return model{
		spinner: s,
		packages: []packageManager{
			{name: "npm", executable: "npm", installCmd: "npx @docubook/create@latest"},
			{name: "pnpm", executable: "pnpm", installCmd: "pnpm dlx @docubook/create@latest"},
			{name: "yarn", executable: "yarn", installCmd: "yarn dlx @docubook/create@latest"},
			{name: "bun", executable: "bun", installCmd: "bunx @docubook/create@latest"},
		},
		// Mulai pengecekan dari index 0
		index: 0,
	}
}

// checkInstalled adalah tea.Cmd untuk mengecek keberadaan sebuah command
func checkInstalled(pm packageManager, index int) tea.Cmd {
	return func() tea.Msg {
		// Simulasi kerja agar spinner terlihat lebih lama
		time.Sleep(400 * time.Millisecond)
		_, err := exec.LookPath(pm.executable)
		return checkResultMsg{index: index, installed: err == nil}
	}
}

// Init adalah command pertama yang dijalankan
func (m model) Init() tea.Cmd {
	// Mulai spinner dan pengecekan pertama
	return tea.Batch(
		m.spinner.Tick,
		checkInstalled(m.packages[0], 0),
	)
}

// Update menangani semua pesan dan memperbarui state
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case checkResultMsg:
		// Update status package yang baru selesai dicek
		m.packages[msg.index].checking = false
		m.packages[msg.index].checked = true
		m.packages[msg.index].installed = msg.installed

		// Jika terinstal, tambahkan perintahnya ke hasil
		if msg.installed {
			m.results = append(m.results, m.packages[msg.index].installCmd)
		}

		// Pindah ke package selanjutnya
		m.index++
		if m.index >= len(m.packages) {
			// Semua sudah dicek, selesai.
			m.done = true
			return m, tea.Quit
		}

		// Mulai pengecekan package selanjutnya
		m.packages[m.index].checking = true
		return m, checkInstalled(m.packages[m.index], m.index)

	default:
		return m, nil
	}
}

// View merender UI berdasarkan state saat ini
func (m model) View() string {
	var b strings.Builder

	// Header
	b.WriteString(headerStyle.Render("✨ DocuBook Environment Check"))
	b.WriteString("\n\n")

	// Tampilan saat pengecekan berlangsung
	if !m.done {
		b.WriteString("Checking for required tools...\n\n")
		for i, pkg := range m.packages {
			if i < m.index {
				// Item yang sudah selesai dicek
				if pkg.installed {
					b.WriteString(fmt.Sprintf("  %s Found %s\n", successStyle.Render("✓"), pkg.name))
				} else {
					b.WriteString(fmt.Sprintf("  %s %s not found\n", errorStyle.Render("✗"), pkg.name))
				}
			} else if i == m.index {
				// Item yang sedang dicek
				b.WriteString(fmt.Sprintf("  %s Checking for %s...\n", m.spinner.View(), pkg.name))
			} else {
				// Item yang belum dicek
				b.WriteString(fmt.Sprintf("  %s %s\n", statusStyle.Render("?"), pkg.name))
			}
		}
	}

	// Tampilan akhir setelah semua selesai dicek
	if m.done {
		// Jika ada tool yang ditemukan
		if len(m.results) > 0 {
			b.WriteString(successStyle.Render("✅ Good to go!"))
			b.WriteString("\n\nYou can now create a new project with one of these commands:\n\n")
			for _, cmd := range m.results {
				b.WriteString(fmt.Sprintf("  %s\n", commandStyle.Render(cmd)))
			}
		} else {
			// Jika tidak ada tool yang ditemukan
			b.WriteString(errorStyle.Render("❌ No required tools found."))
			b.WriteString("\n\nPlease install Node.js or Bun to continue:\n\n")
			b.WriteString(fmt.Sprintf("  - Node.js (includes npm): %s\n", linkStyle.Render("https://nodejs.org/")))
			b.WriteString(fmt.Sprintf("  - Bun: %s\n", linkStyle.Render("https://bun.sh/")))
		}
	}

	b.WriteString(statusStyle.Render("\nPress 'q' to quit."))

	return appStyle.Render(b.String())
}

// init() function dieksekusi saat package di-load
func init() {
	// Tambahkan initCmd ke rootCmd
	rootCmd.AddCommand(initCmd)
}

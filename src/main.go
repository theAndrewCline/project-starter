package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(initialModel())

	m, err := p.StartReturningModel()

	if err != nil {
		log.Fatal(err)
	}

	if m, ok := m.(model); ok && m.completed {
		d := m.codeDir.Value()

		fmt.Printf("\n Starting %s... \n", d)

		err := createNewDir(d)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Successfully initialized %s!", d)
	}
}

func createNewDir(dirname string) error {
	newPath := filepath.Join(repoDir, dirname)
	err := os.MkdirAll(newPath, fs.ModePerm)
	return err
}

const repoDir = "/Users/acline/code"

type model struct {
	completed bool
	codeDir   textinput.Model
	err       error
}

type errMsg error
type tickMsg struct{}

func initialModel() model {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return model{
		completed: false,
		codeDir:   ti,
		err:       nil,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyEnter:
			m.completed = true
			return m, tea.Quit
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	m.codeDir, cmd = m.codeDir.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return fmt.Sprintf("What's your project's name?\n\n%s\n\n%s", m.codeDir.View(), "(Esc to quit)\n")
}

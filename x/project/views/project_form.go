package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/prnk28/gh-pm/internal/ghcli"
	"github.com/prnk28/gh-pm/internal/tui"
)

// ProjectForm represents the data collected from the project creation form
type ProjectForm struct {
	Title        string
	Organization string
	Description  string
	Submitted    bool
}

// NewProjectForm creates a new project form using Huh
func NewProjectForm() (*ProjectForm, error) {
	form := &ProjectForm{}

	// In your form creation code
	orgs, err := ghcli.ListOrgs()
	if err != nil {
		// Handle error
	}

	// Create options for the organization select field
	orgOptions := make([]huh.Option[string], 0, len(orgs)+1)
	orgOptions = append(orgOptions, huh.NewOption[string]("Personal Project", ""))
	for _, org := range orgs {
		orgOptions = append(orgOptions, huh.NewOption[string](org, org))
	}

	// Create the form using Huh
	f := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("Create New Project").
				Description("Fill out the form below to create a new GitHub project.").
				Next(true).
				NextLabel("Start"),
		),

		huh.NewGroup(
			huh.NewInput().
				Title("Project Title").
				Description("Enter a name for your project").
				Placeholder("My Awesome Project").
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("project title cannot be empty")
					}
					return nil
				}).
				Value(&form.Title),
			// Use the options in your form
			huh.NewSelect[string]().
				Title("Organization").
				Description("Select an organization for this project").
				Options(orgOptions...).
				Value(&form.Organization),

			huh.NewText().
				Title("Description").
				Description("Provide a brief description of the project").
				Placeholder("This project is for...").
				CharLimit(400).
				Lines(5).
				Value(&form.Description),

			huh.NewConfirm().
				Title("Create Project?").
				Description("Are you ready to create this project?").
				Affirmative("Yes, create it").
				Negative("No, cancel").
				Value(&form.Submitted),
		),
	).WithTheme(customTheme())

	// Run the form
	err = f.Run()
	if err != nil {
		return nil, err
	}

	return form, nil
}

// customTheme returns a custom theme for the form
func customTheme() *huh.Theme {
	t := huh.ThemeCharm()

	// Customize the theme colors to match your existing UI
	t.Focused.Base = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Background(lipgloss.Color("#2D3142"))
	t.Focused.Title = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true)
	t.Blurred.Title = lipgloss.NewStyle().Foreground(lipgloss.Color("#DDDDDD")).Bold(true)

	return t
}

// FormatSummary returns a formatted summary of the form data
func (f *ProjectForm) FormatSummary() string {
	var sb strings.Builder

	sb.WriteString(tui.Header("Project Created Successfully"))
	sb.WriteString("\n\n")

	// Format the project details
	detailStyle := lipgloss.NewStyle().PaddingLeft(2)
	titleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)

	sb.WriteString(titleStyle.Render("Title:"))
	sb.WriteString(" ")
	sb.WriteString(detailStyle.Render(f.Title))
	sb.WriteString("\n\n")

	sb.WriteString(titleStyle.Render("Organization:"))
	sb.WriteString(" ")
	org := f.Organization
	if org == "" {
		org = "(Personal Project)"
	}
	sb.WriteString(detailStyle.Render(org))
	sb.WriteString("\n\n")

	if f.Description != "" {
		sb.WriteString(titleStyle.Render("Description:"))
		sb.WriteString("\n")
		sb.WriteString(detailStyle.Render(f.Description))
		sb.WriteString("\n\n")
	}

	sb.WriteString(tui.Footer("Press Enter to continue"))

	return sb.String()
}

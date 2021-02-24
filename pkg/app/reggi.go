package app

import (
	"fmt"
	"os"

	"github.com/byxorna/regtest/pkg/ui"
	"github.com/byxorna/regtest/pkg/version"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     version.Name,
	Short:   "Reggi is an interactive regular expression tester",
	Version: version.Version,
	RunE: func(cmd *cobra.Command, args []string) error {
		model, err := ui.New(args)
		if err != nil {
			return fmt.Errorf("unable to initialize model: %v", err)
		}
		p := tea.NewProgram(*model)
		p.EnterAltScreen()
		err = p.Start()
		p.ExitAltScreen()
		return err
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

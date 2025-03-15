package migration

import (
	"path"
	"strings"

	"github.com/mekramy/goconsole"
	"github.com/spf13/cobra"
)

func cmdNew(m Migration, option *cliOption) *cobra.Command {
	return &cobra.Command{
		Use:   "new [name]",
		Short: "Create a new migration file with default stages in the output path",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if option.root == "" {
				goconsole.Message().
					Red("Create").Italic().
					Print("output path must be specified using the WithOutputPath option")
				return
			}

			name := strings.TrimSpace(path.Base(args[0]))
			dir := path.Dir(args[0])
			if dir == "." {
				dir = ""
			}

			if name == "" {
				goconsole.Message().
					Red("Create").Italic().
					Print("file name cannot be empty")
				return
			}

			err := CreateMigrationFile(
				path.Join(option.root, dir),
				name, m.Extension(),
				option.stages.Elements()...,
			)
			if err != nil {
				goconsole.Message().
					Red("Create").
					Italic().Print(err.Error())
				return
			}

			goconsole.Message().
				Green("Create").Italic().
				Printf(`"%s" migration file created`, name)
		},
	}
}

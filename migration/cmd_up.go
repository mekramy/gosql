package migration

import (
	"fmt"
	"strings"

	"github.com/mekramy/goconsole"
	"github.com/spf13/cobra"
)

func cmdUp(m Migration, option *cliOption) *cobra.Command {
	upCmd := &cobra.Command{}
	upCmd.Use = "up [stage1, stage2, ...]"
	upCmd.Short = "applies migrations"
	upCmd.Flags().StringP("name", "n", "", "migration name")
	upCmd.Run = func(cmd *cobra.Command, args []string) {
		stages := append([]string{}, args...)
		if len(stages) == 0 {
			stages = option.stages.Elements()
		}

		if len(stages) == 0 {
			goconsole.Message().
				Red("Up").Italic().
				Print("no stage stage specified")
			return
		}

		options := make([]MigrationOption, 0)
		if name := getFlag(cmd, "name"); name != "" {
			options = append(options, OnlyFiles(name))
		}

		result, err := m.Up(stages, options...)
		if err != nil {
			goconsole.Message().Red("Up").Italic().Print(err.Error())
			return
		}

		goconsole.PrintF("@Bwb{ Migrate Summery: }\n")
		if result.IsEmpty() {
			goconsole.Message().Indent().Italic().Print("nothing to migrate")
		} else {
			for stage, files := range result.GroupByStage() {
				goconsole.PrintF("@BUb{%s} @b{Stage} @Ib{(%d Files)}:\n", strings.ToTitle(stage), len(files))
				for _, file := range files {
					goconsole.PrintF("    @g{UP:} @I{%s}\n", file.Name)
				}

				fmt.Println()
			}
		}
	}

	return upCmd
}

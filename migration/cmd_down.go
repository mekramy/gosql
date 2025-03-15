package migration

import (
	"fmt"
	"strings"

	"github.com/mekramy/goconsole"
	"github.com/spf13/cobra"
)

func cmdDown(m Migration, option *cliOption) *cobra.Command {
	downCmd := &cobra.Command{}
	downCmd.Use = "down [stage1, stage2, ...]"
	downCmd.Short = "rollback migrations"
	downCmd.Flags().StringP("name", "n", "", "migration name")
	downCmd.Run = func(cmd *cobra.Command, args []string) {
		stages := append([]string{}, args...)
		if len(stages) == 0 {
			stages = option.stages.Elements()
		}

		if len(stages) == 0 {
			goconsole.Message().
				Red("Down").Italic().
				Print("no stage stage specified")
			return
		}

		options := make([]MigrationOption, 0)
		if name := getFlag(cmd, "name"); name != "" {
			options = append(options, OnlyFiles(name))
		}

		result, err := m.Down(stages, options...)
		if err != nil {
			goconsole.Message().Red("Down").Italic().Print(err.Error())
			return
		}

		goconsole.PrintF("@Bwb{ Rollback Summery: }\n")
		if result.IsEmpty() {
			goconsole.Message().Indent().Italic().Print("nothing to roll back")
		} else {
			for stage, files := range result.GroupByStage() {
				goconsole.PrintF("@BUb{%s} @b{Stage} @Ib{(%d Files)}:\n", strings.ToTitle(stage), len(files))
				for _, file := range files {
					goconsole.PrintF("    @g{DOWN:} @I{%s}\n", file.Name)
				}

				fmt.Println()
			}
		}
	}

	return downCmd
}

package migration

import (
	"fmt"
	"strings"

	"github.com/mekramy/goconsole"
	"github.com/spf13/cobra"
)

func cmdRefresh(m Migration, option *cliOption) *cobra.Command {
	reCmd := &cobra.Command{}
	reCmd.Use = "refresh [stage1, stage2, ...]"
	reCmd.Short = "refresh migrations"
	reCmd.Flags().StringP("name", "n", "", "migration name")
	reCmd.Run = func(cmd *cobra.Command, args []string) {
		stages := append([]string{}, args...)
		if len(stages) == 0 {
			stages = option.refreshes.Elements()
		}

		if len(stages) == 0 {
			goconsole.Message().
				Red("Refresh").Italic().
				Print("no stage stage specified")
			return
		}

		options := make([]MigrationOption, 0)
		if name := getFlag(cmd, "name"); name != "" {
			options = append(options, OnlyFiles(name))
		}

		result, err := m.Refresh(stages, options...)
		if err != nil {
			goconsole.Message().Red("Refresh").Italic().Print(err.Error())
			return
		}

		goconsole.PrintF("@Bwb{ Refresh Summery: }\n")
		if result.IsEmpty() {
			goconsole.Message().Indent().Italic().Print("nothing to refresh")
		} else {
			for stage, files := range result.GroupByStage() {
				goconsole.PrintF("@BUb{%s} @b{Stage} @Ib{(%d Files)}:\n", strings.ToTitle(stage), len(files))
				for _, file := range files {
					goconsole.PrintF("    @g{REFRESH:} @I{%s}\n", file.Name)
				}

				fmt.Println()
			}
		}
	}

	return reCmd
}

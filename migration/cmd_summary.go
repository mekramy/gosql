package migration

import (
	"fmt"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/mekramy/goconsole"
	"github.com/spf13/cobra"
)

func cmdSummary(m Migration) *cobra.Command {
	return &cobra.Command{
		Use:   "summary",
		Short: "show migration summary",
		Run: func(cmd *cobra.Command, args []string) {
			summary, err := m.Summary()
			if err != nil {
				goconsole.Message().Red("Summary").Italic().Print(err.Error())
				return
			}

			if summary.IsEmpty() {
				goconsole.Message().Blue("Summary").Italic().Print("nothing migrated!")
				return
			}

			goconsole.PrintF("@Bwb{ Migration Summery: }\n")
			for stage, files := range summary.GroupByStage() {
				goconsole.PrintF("@BUb{%s} @b{Stage} @Ib{(%d Files)}:\n", strings.ToTitle(stage), len(files))
				for _, file := range files {
					goconsole.PrintF("    @g{%s}: @I{%s}\n", file.Name, humanize.Time(file.CreatedAt))
				}

				fmt.Println()
			}
		},
	}
}

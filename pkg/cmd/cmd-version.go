package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/sniperkit/cxdig/pkg/core"
)

var (
	// Values are injected at build time from CI
	softwareVersion string
	buildDate       string
	commitHash      string
	commitID        string
	commitUnix      string
	buildVersion    = "2015.6.2-6-gfd7e2d1-dev"
	buildTime       = "2015-06-16-0431 UTC"
	buildCount      string
	buildUnix       string
	branchName      string
)

func printVersion() {
	if softwareVersion != "" {
		// don't use core.Info() to avoid beinf muted by quiet mode
		fmt.Println(softwareVersion)
	} else {
		core.Warn("version is undefined")
	}
	if buildDate != "" {
		core.Infof("Built on %s with %s\n", buildDate, runtime.Version())
	}
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		printVersion()
	},
}

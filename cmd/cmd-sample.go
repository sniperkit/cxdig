package cmd

import (
	"codexray/cxdig/core"
	"codexray/cxdig/core/progress"
	"codexray/cxdig/repos"
	"codexray/cxdig/repos/vcs"
	"errors"

	"github.com/spf13/cobra"
)

var sampleCmd = &cobra.Command{
	Use:   "sample",
	Short: "Repeated source code analysis over time",
	Long:  "Run a sampling tool on the source code at different points in time (sampling frequency)",
	RunE:  cmdSample,
}

type execOptions struct {
	limit int
	freq  string
	cmd   string
}

var execOpts execOptions

func cmdSample(cmd *cobra.Command, args []string) error {

	path, err := getRepositoryPathFromCmdArgs(args)
	if err != nil {
		return err
	}

	repo, err := vcs.OpenRepository(path)
	if err != nil {
		return err
	}

	freq, err := repos.DecodeSamplingFreq(execOpts.freq)
	if err != nil {
		return err
	}

	if execOpts.cmd == "" {
		return errors.New("the command to be executed for each sample is missing")
	}
	tool := repos.NewExternalTool(execOpts.cmd)

	core.Infof("Sampling project '%s' with rate '%s'", repo.Name(), execOpts.freq)
	pb := &progress.ProgressBar{}
	return repo.SampleWithCmd(tool, freq, execOpts.limit, pb)
}

func init() {
	sampleCmd.Flags().IntVarP(&execOpts.limit, "limit", "l", 0, "Set the number of commits used")
	sampleCmd.Flags().StringVarP(&execOpts.freq, "freq", "f", "1w", "Set the frequence separating the commits treated (must be of the form : 10c, 2d, 1m, 3y, etc.")
	sampleCmd.Flags().StringVarP(&execOpts.cmd, "cmd", "c", "", "command to be executed for each sample")
}
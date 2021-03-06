package gitlog

import (
	"bytes"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/sniperkit/cxdig/pkg/core"
	"github.com/sniperkit/cxdig/pkg/repos"
	"github.com/sniperkit/cxdig/pkg/types"
)

type GitRepository struct {
	absPath string
}

func NewGitRepository(path string) *GitRepository {
	absPath, err := filepath.Abs(path)
	if err != nil {
		core.DieOnError(err)
	}

	return &GitRepository{
		absPath: absPath,
	}
}

func (r *GitRepository) SampleWithCmd(tool repos.ExternalTool, rate repos.SamplingRate, commits []types.CommitInfo, samples []types.SampleInfo, p core.Progress) error {
	core.Info("Checking repository status...")
	if !CheckGitStatus(r.absPath) {
		return errors.New("the git repository is not clean, commit your changes or track untracked files and retry")
	}
	return r.walkCommitsWithCommand(tool, commits, samples, p, rate)
}

func (r *GitRepository) Name() repos.ProjectName {
	name := filepath.Base(r.absPath)
	return repos.ProjectName(name)
}
func (r *GitRepository) GetAbsPath() string {
	return r.absPath
}

func (r *GitRepository) walkCommitsWithCommand(tool repos.ExternalTool, commits []types.CommitInfo, samples []types.SampleInfo, p core.Progress, rate repos.SamplingRate) error {
	currentBranch, err := r.getCurrentBranch()
	if err != nil {
		return err
	}

	// TODO: make sure the first commit ID is the current commit ID in the repo
	// restore initial state of the repo
	defer func() {
		p.Done()
		core.Info("Restoring original repository state...")
		_, err := CheckOutOnCommit(r.absPath, currentBranch)
		if err != nil {
			panic(err)
		}
		if err = ClearUntrackedFiles(r.absPath); err != nil {
			panic(err)
		}
	}()
	core.Info("Executing command on each sample...")
	p.Init(len(samples))

	commitIndex := 0
	treatment := 0
	for _, sample := range samples {
		if p != nil {
			if p.IsCancelled() {
				break
			}
			p.Increment()
		}
		for j := commitIndex; j < len(commits); j++ {
			if commits[j].CommitID == sample.CommitID {
				CheckOutOnCommit(r.absPath, commits[j].CommitID.String())
				if err != nil {
					return err
				}
				if err = ClearUntrackedFiles(r.absPath); err != nil {
					return err
				}

				cmd := tool.BuildCmd(r.absPath, r.Name(), commits[j], rate, sample)
				var stderr bytes.Buffer
				cmd.Stderr = &stderr

				// TODO: evaluate CombinedOutput()
				out, err := cmd.Output()
				if err != nil && !p.IsCancelled() {
					// TODO: better error message + use defer on ResetOnCommit
					return errors.Wrap(err, "something wrong happen when running command on commit "+commits[j].CommitID.String())
				}
				logrus.Debug(string(out))
				commitIndex = j
				treatment++
				break
			}
		}
	}
	return nil
}

func (r *GitRepository) ExtractCommits() ([]types.CommitInfo, error) {
	commits, err := ExtractCommitsFromRepository(r.absPath)
	if err != nil {
		return nil, err
	}

	// TODO: check error handling
	commits = GetGitCommitsParents(commits, r.absPath)
	commits = FindMainParentOfCommits(commits, r.absPath)
	return commits, nil
}

func (r *GitRepository) getCurrentBranch() (string, error) {
	rtn, _ := RunGitCommandOnDir(r.absPath, []string{"branch"}, false)
	currentBranch := ""
	for _, branch := range rtn {
		if strings.HasPrefix(branch, "*") {
			currentBranch = strings.TrimSpace(strings.TrimPrefix(branch, "*"))
		}
	}
	if currentBranch == "" {
		return "", errors.New("Current branch could not be found, maybe you are in 'detached HEAD' state?")
	}
	return currentBranch, nil
}

func (r *GitRepository) HasLocalModifications() (bool, error) {
	output, err := RunGitCommandOnDir(r.absPath, []string{"clean", "-ndX"}, false)
	if err != nil {
		return false, err
	}
	if len(output) < 1 || output[0] == "" {
		return false, nil
	}
	return true, nil
}

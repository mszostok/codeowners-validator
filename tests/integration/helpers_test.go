// +build integration

package integration

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func normalizeTimeDurations(in string) string {
	duration := regexp.MustCompile(`\(\d+(\.\d+)?(ns|us|µs|ms|s|m|h)\)`)
	return duration.ReplaceAllString(in, "(<duration>)")
}

func CloneRepo(t *testing.T, url string, branch string) (string, func()) {
	t.Helper()

	repoDir, err := ioutil.TempDir("", strings.ReplaceAll(t.Name(), "/", "-"))
	require.NoError(t, err)

	cleanup := func() {
		err = os.RemoveAll(repoDir)
		require.NoError(t, err)
	}

	_, err = git.PlainClone(repoDir, false, &git.CloneOptions{
		URL:           url,
		ReferenceName: plumbing.NewBranchReferenceName(branch),
		SingleBranch:  true,
		Depth:         1,
	})
	require.NoError(t, err)

	return repoDir, cleanup
}

type Executor struct {
	envs       map[string]string
	timeout    time.Duration
	binaryPath string
}

func Exec() *Executor {
	return &Executor{
		envs: map[string]string{},
	}
}

// WithEnv adds given env. Overrides if previously existed
func (s *Executor) WithEnv(key string, value string) *Executor {
	s.envs[key] = value
	return s
}

type ExecuteOutput struct {
	Stdout   string
	ExitCode int
}

func (s *Executor) AwaitResultAtMost(timeout time.Duration) (*ExecuteOutput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var stdout, stderr bytes.Buffer

	cmd := exec.CommandContext(ctx, s.binaryPath)
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	for k, v := range s.envs {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	err := cmd.Run()
	exitCode := cmd.ProcessState.ExitCode()
	if exitCode != 0 && exitCode != 2 && exitCode != 3 {
		return nil, errors.Wrapf(err, "while executing binary [stdout: %q] [stderr: %q]",
			stdout.String(), stderr.String())
	}
	return &ExecuteOutput{
		ExitCode: exitCode,
		Stdout:   stdout.String(),
	}, nil
}

func (s *Executor) WithTimeout(timeout time.Duration) *Executor {
	s.timeout = timeout
	return s
}

func (s *Executor) Binary(binaryPath string) *Executor {
	s.binaryPath = binaryPath
	return s
}

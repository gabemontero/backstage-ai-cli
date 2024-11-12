package cli

import (
	"github.com/gabemontero/backstage-ai-cli/test/stub"
	"github.com/spf13/cobra"
	"strings"
	"testing"
)

func TestNewCmd(t *testing.T) {
	cmd := NewCmd()

	for _, tc := range []struct {
		args           []string
		generatesError bool
		generatesHelp  bool
		errorStr       string
		outStr         string
	}{
		{
			args:          []string{"new-model"},
			generatesHelp: true,
		},
		{
			args:           []string{"new-model", "kserve"},
			generatesError: true,
			errorStr:       "need to specify an owner and lifecycle setting",
		},
		{
			args:           []string{"new-model", "kubeflow"},
			generatesError: true,
			errorStr:       "need to specify an owner and lifecycle setting",
		},
		{
			args:          []string{"new-model", "help", "kserve"},
			generatesHelp: true,
		},
		{
			args:          []string{"new-model", "help", "kubeflow"},
			generatesHelp: true,
		},
		{
			args:          []string{"fetch-model"},
			generatesHelp: true,
		},
		{
			args:          []string{"fetch-model", "help"},
			generatesHelp: true,
		},
		{
			args:          []string{"fetch-model", "help", "location"},
			generatesHelp: true,
		},
		{
			args:           []string{"fetch-model", "location", "help"},
			generatesError: true,
			errorStr:       "unsupported protocol scheme",
		},
		{
			args:          []string{"fetch-model", "help"},
			generatesHelp: true,
		},
		{
			args:          []string{"fetch-model", "help", "locations"},
			generatesHelp: true,
		},
		{
			args:           []string{"fetch-model", "locations", "foo"},
			generatesError: true,
			errorStr:       "unsupported protocol scheme",
		},
		{
			args:          []string{"fetch-model", "help", "components"},
			generatesHelp: true,
		},
		{
			args:           []string{"fetch-model", "components", "foo"},
			generatesError: true,
			errorStr:       "unsupported protocol scheme",
		},
		{
			args:          []string{"fetch-model", "help", "resources"},
			generatesHelp: true,
		},
		{
			args:           []string{"fetch-model", "resources", "help"},
			generatesError: true,
			errorStr:       "unsupported protocol scheme",
		},
		{
			args:          []string{"fetch-model", "help", "apis"},
			generatesHelp: true,
		},
		{
			args:           []string{"fetch-model", "apis", "foo"},
			generatesError: true,
			errorStr:       "unsupported protocol scheme",
		},
		{
			args:          []string{"fetch-model", "help", "entities"},
			generatesHelp: true,
		},
	} {
		subCmd, stdout, stderr, err := stub.ExecuteCommandC(cmd, tc.args...)
		switch {
		case err == nil && tc.generatesError:
			t.Errorf("error should have been generated for '%s'", strings.Join(tc.args, " "))
		case err != nil && !tc.generatesError:
			t.Errorf("error generated unexpectedly for '%s': %s", strings.Join(tc.args, " "), err.Error())
		case err != nil && tc.generatesError && !strings.Contains(stderr, tc.errorStr):
			t.Errorf("unexpected error output for '%s'- got '%s' but expected '%s'", strings.Join(tc.args, " "), stderr, tc.errorStr)
		case tc.generatesHelp && !testHelpOK(stdout, subCmd):
			t.Errorf("unexpected help output for '%s' - got '%s' but expected '%s'", strings.Join(tc.args, " "), stdout, subCmd.Long)
		}
	}
}

func testHelpOK(stdout string, cmd *cobra.Command) bool {
	if strings.Contains(stdout, cmd.Long) {
		return true
	}
	return false
}

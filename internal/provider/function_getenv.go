package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/kohirens/stdlib/cli"
	"runtime"
	"strings"
)

var _ function.Function = &getenvFunction{}

type getenvFunction struct{}

func NewgetenvFunction() function.Function {
	return &getenvFunction{}
}

func (f *getenvFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "getenv"
}

func (f *getenvFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "Get the value of an environment variable",
		Description: "Given the name of an environment variable name as a string, if set, will return its value.",

		Parameters: []function.Parameter{
			function.StringParameter{
				Name:        "name",
				Description: "Name of an environment variable value to return.",
			},
		},
		Return: function.StringReturn{},
	}
}

func (f *getenvFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var name string

	resp.Error = req.Arguments.Get(ctx, &name)
	if resp.Error != nil {
		return
	}

	value, ok := printEnv(name)
	if !ok {
		// Intentionally not including the Go parse error in the return diagnostic, as the message is based on a Go-specific
		// reference time that may be unfamiliar to practitioners
		tflog.Error(ctx, fmt.Sprintf("could not find environment variable %s", name))

		return
	}

	sv := types.StringValue(value)

	resp.Error = resp.Result.Set(ctx, &sv)
}

// PrintEnv we'll need to use the environment shell to get the environment variable since Terraform forbids getting
// variables other than TF_VAR_* with os.LookupEnv.
func printEnv(name string) (string, bool) {
	ok := true
	var value []byte
	var err error
	var exitCode int
	isWindows := runtime.GOOS == "windows"

	if isWindows {
		value, err, exitCode, _ = cli.RunCommand(".", "Powershell", []string{"echo", "$Env:" + name})
	} else {
		value, err, exitCode, _ = cli.RunCommand(".", "printenv", []string{name})
	}

	if exitCode != 0 || err != nil || isWindows && len(value) == 0 {
		ok = false
	}

	return strings.TrimRight(string(value), "\r\n"), ok
}

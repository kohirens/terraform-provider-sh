// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"runtime"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/kohirens/stdlib/cli"
)

// Ensure the implementation satisfies the expected interfaces.
var ( // Statements like these are asserting that the struct fulfils all the interface requirements.
	_ datasource.DataSource = &varsDataSource{}
)

func NewDataSourceExtractEnvVars() datasource.DataSource {
	return &varsDataSource{}
}

// varsDataSource Defines the input attributes.
type varsDataSource struct{}

// extractEnvVarsDataModel Defines the output attributes.
type varsDataModel struct {
	Names    types.List `tfsdk:"names"`
	Values   types.Map  `tfsdk:"values"`
	Required types.Bool `tfsdk:"required"`
}

// Configure adds the provider configured client to the data source.
func (d *varsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}
}

// Metadata returns the data source type name.
func (d *varsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vars"
}

// Schema defines the schema for the data source.
func (d *varsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The `sh_vars` data source takes a list of environment variable names and returns them as key-value pairs.",
		Attributes: map[string]schema.Attribute{
			"names": schema.ListAttribute{
				ElementType:         types.StringType,
				Required:            true,
				MarkdownDescription: "Names of the environment variables to extract.",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
			},
			"required": schema.BoolAttribute{
				MarkdownDescription: "Requires all variables listed to be present and extracted successfully.",
				Optional:            true,
			},
			"values": schema.MapAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				MarkdownDescription: "Map of variables extracted from the environment.",
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *varsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var envVars varsDataModel

	diags := req.Config.Get(ctx, &envVars)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read environment variables `names` input attribute value.
	var names []types.String
	resp.Diagnostics.Append(envVars.Names.ElementsAs(ctx, &names, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Read environment variables `required` input attribute value.
	required := envVars.Required.ValueBool()
	errMessages := ""
	values := make(map[string]string, len(names))

	// Look up each environment variables response body to model
	for _, name := range names {
		envName := name.ValueString()
		value, ok := PrintEnv(envName)
		tflog.Debug(ctx, fmt.Sprintf("env name: %v; ok: %v\n", envName, ok))
		if !ok {
			if required {
				errMessages += fmt.Sprintf("environment variable %v is not set\n", envName)
			}
			continue
		}

		// update the state with the current value.
		values[envName] = value
	}

	envVars.Values, diags = types.MapValueFrom(ctx, types.StringType, values)
	if resp.Diagnostics.HasError() {
		return
	}

	if required && errMessages != "" {
		resp.Diagnostics.AddError(
			"the following values are required:",
			errMessages,
		)
		return
	}

	// Set state
	diags = resp.State.Set(ctx, &envVars)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// PrintEnv we'll need to use the environment shell to get the environment variable since Terraform forbids getting
// variables other than TF_VAR_* with os.LookupEnv.
func PrintEnv(name string) (string, bool) {
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

// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

const (
	typeName = "sh"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &shProvider{}
)

// New is a helper function to simplify provider instantiation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &shProvider{
			version: version,
		}
	}
}

// shProvider is the provider implementation.
type shProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Metadata returns the provider type name.
func (p *shProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = typeName
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *shProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{}
	resp.Schema.Description = "Extract values from the sh environment where Terraform runs; no input for the configuration is needed."
}

// Configure prepares a sh API client for data sources and resources.
func (p *shProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// No configuration needed.
}

// DataSources defines the data sources implemented in the provider.
func (p *shProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewDataSourceExtractEnvVars,
	}
}

// Resources defines the resources implemented in the provider.
func (p *shProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}

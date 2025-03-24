// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccShDataSource(t *testing.T) {
	_ = os.Setenv("TEST_HOME", "/root")
	defer os.Unsetenv("TEST_HOME")
	fixtureName := "data.sh_vars.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "sh_vars" "test" {names=["TEST_HOME"]}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify the first coffee to ensure all attributes are set
					resource.TestCheckResourceAttr(fixtureName, "values.TEST_HOME", "/root"),
					resource.TestCheckResourceAttr(fixtureName, "values.%", "1"),
				),
			},
		},
	})
}

func TestParseEnv(t *testing.T) {
	_ = os.Setenv("TEST_VAR", "1234")
	defer os.Unsetenv("TEST_VAR")

	cases := []struct {
		name   string
		want   string
		wantOk bool
	}{
		{"TEST_VAR", "1234", true},
		{"NOT_SET", "", false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, ok := printEnv(c.name)

			if ok != c.wantOk {
				t.Errorf("PrintEnv() got ok = %v, wanted = %v", ok, c.wantOk)
				return
			}

			if got != c.want {
				t.Errorf("PrintEnv() got %v; want %v", got, c.want)
				return
			}
		})
	}
}

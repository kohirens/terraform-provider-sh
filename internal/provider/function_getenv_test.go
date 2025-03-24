package provider

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	_ = os.Setenv("TEST_ENV", "123")
	_ = os.Setenv("TEST_EMPTY", "")
	defer os.Unsetenv("TEST_ENV")
	defer os.Unsetenv("TEST_EMPTY")
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
                output "test" {
                    value = provider::sh::getenv("TEST_ENV")
                }
                `,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValue("test", knownvalue.StringExact("123")),
					},
				},
			},
			{
				Config: `
			   output "test" {
			       value = provider::sh::getenv("NOT_SET")
			   }
			   `,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				Config: `
			   output "test" {
			       value = provider::sh::getenv("TEST_EMPTY")
			   }
			   `,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

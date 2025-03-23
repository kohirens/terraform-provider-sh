package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccShDataSource(t *testing.T) {
	fixtureName := "data.sh_vars.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "sh_vars" "test" {names=["HOME"]}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify the first coffee to ensure all attributes are set
					resource.TestCheckResourceAttr(fixtureName, `values.HOME`, `/root`),
					resource.TestCheckResourceAttr(fixtureName, `values.%`, "1"),
					//resource.TestCheckResourceAttrWith(fixtureName, "values.HOME", func(value string) error {
					//	if value != "HOME" {
					//		return fmt.Errorf("expected HOME, got %s", value)
					//	}
					//	return nil
					//}),
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
			got, ok := PrintEnv(c.name)

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

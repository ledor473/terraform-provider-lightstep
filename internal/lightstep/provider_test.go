package lightstep

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = New()
	testAccProviders = map[string]*schema.Provider{
		"lightstep": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := New().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = New()
}

func testAccPreCheck(t *testing.T) {
	for _, e := range []string{"LIGHTSTEP_API_KEY", "LIGHTSTEP_ORGANIZATION", "LIGHTSTEP_PROJECT"} {
		if v := os.Getenv(e); v == "" {
			t.Fatalf("%s must be set for acceptance tests", e)
		}
	}
}

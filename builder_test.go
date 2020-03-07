package docs

import (
	"github.com/threeq/docs/swagger"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDescription(t *testing.T) {
	api := New(
		Description("blah"),
	)
	assert.Equal(t, "blah", api.Info.Description)
}

func TestVersion(t *testing.T) {
	api := New(
		Version("blah"),
	)
	assert.Equal(t, "blah", api.Info.Version)
}

func TestTermsOfService(t *testing.T) {
	api := New(
		TermsOfService("blah"),
	)
	assert.Equal(t, "blah", api.Info.TermsOfService)
}

func TestTitle(t *testing.T) {
	api := New(
		Title("blah"),
	)
	assert.Equal(t, "blah", api.Info.Title)
}

func TestContactEmail(t *testing.T) {
	api := New(
		ContactEmail("blah"),
	)
	assert.Equal(t, "blah", api.Info.Contact.Email)
}

func TestLicense(t *testing.T) {
	api := New(
		License("name", "url"),
	)
	assert.Equal(t, "name", api.Info.License.Name)
	assert.Equal(t, "url", api.Info.License.URL)
}

func TestBasePath(t *testing.T) {
	api := New(
		BasePath("/"),
	)
	assert.Equal(t, "/", api.BasePath)
}

func TestSchemes(t *testing.T) {
	api := New(
		Schemes("blah"),
	)
	assert.Equal(t, []string{"blah"}, api.Schemes)
}

func TestTag(t *testing.T) {
	api := New(
		Tag("name", "desc",
			TagDescription("ext-desc"),
			TagURL("ext-url"),
		),
	)

	expected := swagger.Tag{
		Name:        "name",
		Description: "desc",
		Docs: &swagger.Docs{
			Description: "ext-desc",
			URL:         "ext-url",
		},
	}
	assert.Equal(t, expected, api.Tags[0])
}

func TestHost(t *testing.T) {
	api := New(
		Host("blah"),
	)
	assert.Equal(t, "blah", api.Host)
}

func TestSecurityScheme(t *testing.T) {
	api := New(
		SecurityScheme("basic", swagger.BasicSecurity()),
		SecurityScheme("apikey", swagger.APIKeySecurity("Authorization", "header")),
	)
	assert.Len(t, api.SecurityDefinitions, 2)
	assert.Contains(t, api.SecurityDefinitions, "basic")
	assert.Contains(t, api.SecurityDefinitions, "apikey")
	assert.Equal(t, "header", api.SecurityDefinitions["apikey"].In)
}

func TestSecurity(t *testing.T) {
	api := New(
		Security("basic"),
	)
	assert.Len(t, api.Security.Requirements, 1)
	assert.Contains(t, api.Security.Requirements[0], "basic")
}

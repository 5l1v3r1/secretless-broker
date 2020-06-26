package generic

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"
	"text/template"

	"github.com/cyberark/secretless-broker/pkg/secretless/plugin/connector"
	validation "github.com/go-ozzo/ozzo-validation"
)

type config struct {
	CredentialPatterns map[string]*regexp.Regexp
	Headers            map[string]*template.Template
	QueryParams 	   map[string]*template.Template
	ForceSSL           bool
}

// validate validates that the given creds satisfy the CredentialValidations of
// the config.
func (c *config) validate(credsByID connector.CredentialValuesByID) error {
	for requiredCred, pattern := range c.CredentialPatterns {
		credVal, ok := credsByID[requiredCred]
		if !ok {
			return fmt.Errorf("missing required credential: %q", requiredCred)
		}
		if !pattern.Match(credVal) {
			return fmt.Errorf(
				"credential %q doesn't match pattern %q", requiredCred, pattern,
			)
		}
	}
	return nil
}

// Render returns the config's `config` section templates filled in with the
// given credentialValues.
func render(
	template map[string]*template.Template,
	credsByID connector.CredentialValuesByID,
) (map[string]string, error) {
	errs := validation.Errors{}
	args := make(map[string]string)

	// Creds must be strings to work with templates
	credStringsByID := make(map[string]string)
	for credName, credBytes := range credsByID {
		credStringsByID[credName] = string(credBytes)
	}

	for arg, tmpl := range template {
		builder := &strings.Builder{}
		if err := tmpl.Execute(builder, credStringsByID); err != nil {
			errs[arg] = fmt.Errorf("couldn't render template: %q", err)
			continue
		}
		args[arg] = builder.String()
	}

	if err := errs.Filter(); err != nil {
		return nil, err
	}

	return args, nil
}

// newConfig takes a ConfigYAML, validates it, and converts it into a
// generic.config struct -- which is what our application wants to work with.
func newConfig(cfgYAML *ConfigYAML) (*config, error) {
	errs := validation.Errors{}

	cfg := &config{
		CredentialPatterns: make(map[string]*regexp.Regexp),
		ForceSSL:           cfgYAML.ForceSSL,
	}

	// Validate and save regexps
	for cred, reStr := range cfgYAML.CredentialValidations {
		re, err := regexp.Compile(reStr)
		if err != nil {
			errs[cred] = fmt.Errorf("invalid regex: %q", err)
			continue
		}
		cfg.CredentialPatterns[cred] = re
	}

	cfg.Headers, errs = validateAndSaveTemplateString(cfgYAML.Headers, errs)
	cfg.QueryParams, errs = validateAndSaveTemplateString(cfgYAML.QueryParams, errs)

	if err := errs.Filter(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func validateAndSaveTemplateString(templates map[string]string,
	errs validation.Errors) (map[string]*template.Template, validation.Errors) {
	parsedTemplates := make(map[string]*template.Template)
	// Validate and save template strings
	for tmplName, tmplStr := range templates {
		tmpl := newHTTPTemplate(tmplName)
		// Ignore pointer to receiver returned by Parse(): it's just "tmpl".
		_, err := tmpl.Parse(tmplStr)
		if err != nil {
			errs[tmplName] = fmt.Errorf("invalid query param template: %q", err)
			continue
		}
		parsedTemplates[tmplName] = tmpl
	}
	return parsedTemplates, errs
}

// templateFuncs is a map holding the custom functions available for use within
// a header template string.  We can easily add new functions as needed or
// requested.
var templateFuncs = template.FuncMap{
	"base64": func(str string) string {
		return base64.StdEncoding.EncodeToString([]byte(str))
	},
}

func newHTTPTemplate(name string) *template.Template {
	return template.New(name).Funcs(templateFuncs)
}

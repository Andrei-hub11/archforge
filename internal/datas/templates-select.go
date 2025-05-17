package datas

// TemplatesSelect maps template IDs to template names
var TemplatesSelect = map[string]string{
	"1": "clean-arch-keycloak-pg-dapper",
	"2": "clean-arch-keycloak-pg-ef",
	"3": "webapi",
}

// TemplateDescriptions maps template IDs to their descriptions
var TemplateDescriptions = map[string]string{
	"1": "Clean Architecture with Keycloak and PostgreSQL using Dapper",
	"2": "Clean Architecture with Keycloak and PostgreSQL using Entity Framework",
	"3": "Basic ASP.NET Core Web API",
}

// GetTemplateOptions returns formatted options for selection
func GetTemplateOptions() []string {
	options := make([]string, len(TemplatesSelect))
	i := 0
	for id := range TemplatesSelect {
		desc := TemplateDescriptions[id]
		options[i] = id + " - " + desc
		i++
	}
	return options
}

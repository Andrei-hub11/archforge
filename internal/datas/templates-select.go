package datas

// TemplatesSelect mapeia IDs para nomes de templates
var TemplatesSelect = map[string]string{
	"1": "clean-arch-keycloak-pg-dapper",
	"2": "clean-arch-keycloak-pg-ef",
	"3": "webapi",
}

// TemplateDescriptions mapeia IDs para descrições de templates
var TemplateDescriptions = map[string]string{
	"1": "Clean Architecture com Keycloak e PostgreSQL usando Dapper",
	"2": "Clean Architecture com Keycloak e PostgreSQL usando Entity Framework",
	"3": "Web API básica com ASP.NET Core",
}

// GetTemplateOptions retorna as opções formatadas para seleção
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

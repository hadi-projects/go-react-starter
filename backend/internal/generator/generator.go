package generator

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
)

type Field struct {
	Name       string `yaml:"name" json:"name"`
	Type       string `yaml:"type" json:"type"`
	Binding    string `yaml:"binding" json:"binding"`
	Searchable bool   `yaml:"searchable" json:"searchable"`
	Unique     bool   `yaml:"unique" json:"unique"`
}

func (f Field) NameGo() string {
	return ToCamelCase(f.Name)
}

func (f Field) NameLowerGo() string {
	return ToLowerCamelCase(f.Name)
}

func (f Field) NameJson() string {
	return f.Name
}

func (f Field) NameSql() string {
	return f.Name
}

func (f Field) TypeGo() string {
	switch f.Type {
	case "string", "wysiwyg", "file", "image", "video", "audio", "enum":
		return "string"
	case "int":
		return "int"
	case "float":
		return "float64"
	case "date", "time", "datetime":
		return "time.Time"
	case "boolean":
		return "bool"
	case "json":
		return "interface{}"
	default:
		return "string"
	}
}

func (f Field) GormType() string {
	switch f.Type {
	case "string":
		return "type:varchar(255);not null"
	case "wysiwyg":
		return "type:longtext"
	case "int":
		return "type:int"
	case "float":
		return "type:decimal(10,2)"
	case "boolean":
		return "type:boolean"
	case "enum":
		return "type:varchar(50)"
	default:
		return "type:varchar(255)"
	}
}

type ModuleConfig struct {
	ModuleName string  `yaml:"module_name" json:"module_name"`
	TableName  string  `yaml:"table_name" json:"table_name"`
	AuditLog   bool    `yaml:"audit_log" json:"audit_log"`
	Icon       string  `yaml:"icon" json:"icon"`
	Fields     []Field `yaml:"fields" json:"fields"`
}

type Generator struct {
	Config  ModuleConfig
	BaseDir string
}

func NewGenerator(configPath string, baseDir string) (*Generator, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config ModuleConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &Generator{
		Config:  config,
		BaseDir: baseDir,
	}, nil
}

func NewGeneratorFromConfig(config ModuleConfig, baseDir string) *Generator {
	return &Generator{
		Config:  config,
		BaseDir: baseDir,
	}
}

func (g *Generator) Generate() error {
	templates := map[string]string{
		"entity.go.tmpl":       filepath.Join(g.BaseDir, "internal/entity", strings.ToLower(g.Config.ModuleName)+".go"),
		"dto.go.tmpl":          filepath.Join(g.BaseDir, "internal/dto", strings.ToLower(g.Config.ModuleName)+"_dto.go"),
		"repository.go.tmpl":   filepath.Join(g.BaseDir, "internal/repository", strings.ToLower(g.Config.ModuleName)+"_repository.go"),
		"service.go.tmpl":      filepath.Join(g.BaseDir, "internal/service", strings.ToLower(g.Config.ModuleName)+"_service.go"),
		"handler.go.tmpl":      filepath.Join(g.BaseDir, "internal/handler", strings.ToLower(g.Config.ModuleName)+"_handler.go"),
		"service_test.go.tmpl": filepath.Join(g.BaseDir, "internal/service", strings.ToLower(g.Config.ModuleName)+"_service_test.go"),
		// Frontend UI templates
		"frontend_api.js.tmpl":   filepath.Join(g.BaseDir, "../frontend/src/api", strings.ToLower(g.Config.ModuleName)+".js"),
		"frontend_page.jsx.tmpl": filepath.Join(g.BaseDir, "../frontend/src/pages/admin", ToCamelCase(g.Config.ModuleName)+"Page.jsx"),
	}

	data := map[string]interface{}{
		"ModuleName":           ToCamelCase(g.Config.ModuleName),
		"ModuleNameLower":      strings.ToLower(g.Config.ModuleName),
		"ModuleNameLowerCamel": ToLowerCamelCase(g.Config.ModuleName),
		"TableName":            g.Config.TableName,
		"Fields":               g.Config.Fields,
		"AuditLog":             g.Config.AuditLog,
		"HasSearchableFields":  g.hasSearchableFields(),
	}

	for tmplName, outputPath := range templates {
		if err := g.renderTemplate(tmplName, outputPath, data); err != nil {
			return err
		}
	}

	if err := g.registerRouter(); err != nil {
		fmt.Printf("Warning: Failed to register router: %v\n", err)
	}

	if err := g.registerMigration(); err != nil {
		fmt.Printf("Warning: Failed to register migration: %v\n", err)
	}

	if err := g.registerFrontend(); err != nil {
		fmt.Printf("Warning: Failed to register frontend: %v\n", err)
	}

	return nil
}

func (g *Generator) registerRouter() error {
	routerPath := filepath.Join(g.BaseDir, "internal/router/router.go")
	privateRouterPath := filepath.Join(g.BaseDir, "internal/router/private_router.go")

	repoInit := fmt.Sprintf("\t%sRepo := customRepository.New%sRepository(db)\n\t// [GENERATOR_INSERT_REPOSITORY]", strings.ToLower(g.Config.ModuleName), ToCamelCase(g.Config.ModuleName))
	serviceInit := fmt.Sprintf("\t%sService := customService.New%sService(%sRepo, r.cache)\n\t// [GENERATOR_INSERT_SERVICE]", strings.ToLower(g.Config.ModuleName), ToCamelCase(g.Config.ModuleName), strings.ToLower(g.Config.ModuleName))
	handlerInit := fmt.Sprintf("\t%sHandler := customHandler.New%sHandler(%sService)\n\t// [GENERATOR_INSERT_HANDLER]", strings.ToLower(g.Config.ModuleName), ToCamelCase(g.Config.ModuleName), strings.ToLower(g.Config.ModuleName))
	handlerParam := fmt.Sprintf("\t\t%sHandler,\n\t\t// [GENERATOR_INSERT_HANDLER_PARAM]", strings.ToLower(g.Config.ModuleName))

	if err := g.insertAtMarker(routerPath, "// [GENERATOR_INSERT_REPOSITORY]", repoInit); err != nil {
		return err
	}
	if err := g.insertAtMarker(routerPath, "// [GENERATOR_INSERT_SERVICE]", serviceInit); err != nil {
		return err
	}
	if err := g.insertAtMarker(routerPath, "// [GENERATOR_INSERT_HANDLER]", handlerInit); err != nil {
		return err
	}
	if err := g.insertAtMarker(routerPath, "// [GENERATOR_INSERT_HANDLER_PARAM]", handlerParam); err != nil {
		return err
	}

	handlerParamPrivate := fmt.Sprintf("\t%sHandler customHandler.%sHandler,\n\t// [GENERATOR_INSERT_HANDLER_PARAM]", strings.ToLower(g.Config.ModuleName), ToCamelCase(g.Config.ModuleName))
	groupInit := fmt.Sprintf(`	%s := v1.Group("/%s")
	%s.Use(middleware.AuthMiddleware(r.config.JWT.Secret))
	{
		%s.POST("", %sHandler.Create)
		%s.GET("", %sHandler.GetAll)
		%s.GET("/export", %sHandler.Export)
		%s.GET("/:id", %sHandler.GetByID)
		%s.PUT("/:id", %sHandler.Update)
		%s.DELETE("/:id", %sHandler.Delete)
	}
	// [GENERATOR_INSERT_GROUP]`, strings.ToLower(g.Config.ModuleName), g.Config.TableName, strings.ToLower(g.Config.ModuleName), strings.ToLower(g.Config.ModuleName), strings.ToLower(g.Config.ModuleName), strings.ToLower(g.Config.ModuleName), strings.ToLower(g.Config.ModuleName), strings.ToLower(g.Config.ModuleName), strings.ToLower(g.Config.ModuleName), strings.ToLower(g.Config.ModuleName), strings.ToLower(g.Config.ModuleName), strings.ToLower(g.Config.ModuleName), strings.ToLower(g.Config.ModuleName), strings.ToLower(g.Config.ModuleName), strings.ToLower(g.Config.ModuleName))

	if err := g.insertAtMarker(privateRouterPath, "// [GENERATOR_INSERT_HANDLER_PARAM]", handlerParamPrivate); err != nil {
		return err
	}
	if err := g.insertAtMarker(privateRouterPath, "// [GENERATOR_INSERT_GROUP]", groupInit); err != nil {
		return err
	}

	return nil
}

func (g *Generator) registerMigration() error {
	migratePath := filepath.Join(g.BaseDir, "cmd/migrate/migrate.go")
	migrationInit := fmt.Sprintf("\t\t&customEntity.%s{},\n\t\t// [GENERATOR_INSERT_MIGRATION]", ToCamelCase(g.Config.ModuleName))
	return g.insertAtMarker(migratePath, "// [GENERATOR_INSERT_MIGRATION]", migrationInit)
}

func (g *Generator) registerFrontend() error {
	appPath := filepath.Join(g.BaseDir, "../frontend/src/App.jsx")
	sidebarPath := filepath.Join(g.BaseDir, "../frontend/src/layouts/AdminLayout.jsx")

	// Inject Route
	routeImport := fmt.Sprintf("import %sPage from './pages/admin/%sPage';\n// [GENERATOR_INSERT_IMPORT]", ToCamelCase(g.Config.ModuleName), ToCamelCase(g.Config.ModuleName))
	routeDefinition := fmt.Sprintf("\t\t\t\t\t<Route path=\"admin/%s\" element={<%sPage />} />\n\t\t\t\t\t// [GENERATOR_INSERT_ROUTE]", strings.ToLower(g.Config.ModuleName), ToCamelCase(g.Config.ModuleName))

	if err := g.insertAtMarker(appPath, "// [GENERATOR_INSERT_IMPORT]", routeImport); err != nil {
		return err
	}
	if err := g.insertAtMarker(appPath, "// [GENERATOR_INSERT_ROUTE]", routeDefinition); err != nil {
		return err
	}

	// Inject Sidebar Item
	iconSvg := g.Config.Icon
	if iconSvg == "" {
		iconSvg = `<path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />`
	} else if !strings.Contains(iconSvg, "<path") {
		// If it's just the 'd' attribute, wrap it
		iconSvg = fmt.Sprintf(`<path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="%s" />`, iconSvg)
	}

	sidebarItem := fmt.Sprintf(`                                { label: '%s', path: '/admin/%s', permission: 'get-%s', icon: (
                                    <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        %s
                                    </svg>
                                ) },
                                // [GENERATOR_INSERT_ADMIN_ITEM]`, ToCamelCase(g.Config.ModuleName), strings.ToLower(g.Config.ModuleName), strings.ToLower(g.Config.ModuleName), iconSvg)

	if err := g.insertAtMarker(sidebarPath, "// [GENERATOR_INSERT_ADMIN_ITEM]", sidebarItem); err != nil {
		return err
	}

	return nil
}

func (g *Generator) insertAtMarker(filePath string, marker string, content string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	body := string(data)
	if !strings.Contains(body, marker) {
		return fmt.Errorf("marker %s not found in %s", marker, filePath)
	}

	// Avoid duplicate insertion
	firstLine := strings.Split(strings.TrimSpace(content), "\n")[0]
	if strings.Contains(body, firstLine) {
		return nil
	}

	newBody := strings.Replace(body, marker, content, 1)
	return os.WriteFile(filePath, []byte(newBody), 0644)
}

func (g *Generator) renderTemplate(tmplName string, outputPath string, data interface{}) error {
	tmplPath := filepath.Join(g.BaseDir, "internal/generator/templates", tmplName)
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return err
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return err
	}

	return os.WriteFile(outputPath, buf.Bytes(), 0644)
}

func (g *Generator) hasSearchableFields() bool {
	for _, f := range g.Config.Fields {
		if f.Searchable {
			return true
		}
	}
	return false
}

func ToCamelCase(s string) string {
	parts := strings.Split(s, "_")
	for i := range parts {
		if parts[i] == "" {
			continue
		}
		parts[i] = strings.Title(parts[i])
	}
	return strings.Join(parts, "")
}

func ToLowerCamelCase(s string) string {
	if s == "" {
		return ""
	}
	camel := ToCamelCase(s)
	return strings.ToLower(camel[:1]) + camel[1:]
}

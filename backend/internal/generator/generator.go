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
	Name       string `yaml:"name"`
	Type       string `yaml:"type"`
	Binding    string `yaml:"binding"`
	Searchable bool   `yaml:"searchable"`
	Unique     bool   `yaml:"unique"`
}

func (f Field) NameGo() string {
	return ToCamelCase(f.Name)
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
		return "type:text"
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
	ModuleName string  `yaml:"module_name"`
	TableName  string  `yaml:"table_name"`
	AuditLog   bool    `yaml:"audit_log"`
	Fields     []Field `yaml:"fields"`
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
	}

	data := map[string]interface{}{
		"ModuleName":          g.Config.ModuleName,
		"ModuleNameLower":     strings.ToLower(g.Config.ModuleName),
		"TableName":           g.Config.TableName,
		"Fields":              g.Config.Fields,
		"AuditLog":            g.Config.AuditLog,
		"HasSearchableFields": g.hasSearchableFields(),
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

	return nil
}

func (g *Generator) registerRouter() error {
	routerPath := filepath.Join(g.BaseDir, "internal/router/router.go")
	privateRouterPath := filepath.Join(g.BaseDir, "internal/router/private_router.go")

	repoInit := fmt.Sprintf("\t%sRepo := repository.New%sRepository(db)\n\t// [GENERATOR_INSERT_REPOSITORY]", strings.ToLower(g.Config.ModuleName), g.Config.ModuleName)
	serviceInit := fmt.Sprintf("\t%sService := service.New%sService(%sRepo, r.cache)\n\t// [GENERATOR_INSERT_SERVICE]", strings.ToLower(g.Config.ModuleName), g.Config.ModuleName, strings.ToLower(g.Config.ModuleName))
	handlerInit := fmt.Sprintf("\t%sHandler := handler.New%sHandler(%sService)\n\t// [GENERATOR_INSERT_HANDLER]", strings.ToLower(g.Config.ModuleName), g.Config.ModuleName, strings.ToLower(g.Config.ModuleName))
	handlerParam := fmt.Sprintf("\t\t\t%sHandler,\n\t\t\t// [GENERATOR_INSERT_HANDLER_PARAM]", strings.ToLower(g.Config.ModuleName))

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

	handlerParamPrivate := fmt.Sprintf("\t%sHandler handler.%sHandler,\n\t// [GENERATOR_INSERT_HANDLER_PARAM]", strings.ToLower(g.Config.ModuleName), g.Config.ModuleName)
	groupInit := fmt.Sprintf(`	%s := v1.Group("/%s")
	%s.Use(middleware.AuthMiddleware(r.config.JWT.Secret))
	{
		%s.POST("", %sHandler.Create)
		%s.GET("", %sHandler.GetAll)
		%s.GET("/:id", %sHandler.GetByID)
		%s.PUT("/:id", %sHandler.Update)
		%s.DELETE("/:id", %sHandler.Delete)
	}
	// [GENERATOR_INSERT_GROUP]`, strings.ToLower(g.Config.ModuleName), g.Config.TableName, strings.ToLower(g.Config.ModuleName), strings.ToLower(g.Config.ModuleName), strings.ToLower(g.Config.ModuleName), strings.ToLower(g.Config.ModuleName), strings.ToLower(g.Config.ModuleName), strings.ToLower(g.Config.ModuleName), strings.ToLower(g.Config.ModuleName), strings.ToLower(g.Config.ModuleName), strings.ToLower(g.Config.ModuleName), strings.ToLower(g.Config.ModuleName), strings.ToLower(g.Config.ModuleName))

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
	migrationInit := fmt.Sprintf("\t\t&entity.%s{},\n\t\t// [GENERATOR_INSERT_MIGRATION]", g.Config.ModuleName)
	return g.insertAtMarker(migratePath, "// [GENERATOR_INSERT_MIGRATION]", migrationInit)
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
	if strings.Contains(body, strings.Split(content, "\n")[0]) {
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

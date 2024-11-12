package codegen

import (
	"strconv"
	"strings"

	genstrings "github.com/aryanugroho/gost/pkg/strings"
)

type Codegen struct {
	AppName        string         `yaml:"app_name"`
	Port           string         `yaml:"port"`
	Infrastructure Infrastructure `yaml:"infrastructure"`
	Entities       []Entities     `yaml:"entities"`
}
type Db struct {
	Type     string `yaml:"type"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	Schema   string `yaml:"schema"`
	Ssl      bool   `yaml:"ssl"`
}
type Mq struct {
	Type string `yaml:"type"`
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}
type Sentry struct {
	Dsn         string `yaml:"dsn"`
	Environment string `yaml:"environment"`
	Release     string `yaml:"release"`
}
type Infrastructure struct {
	Db     Db     `yaml:"db"`
	Mq     Mq     `yaml:"mq"`
	Sentry Sentry `yaml:"sentry"`
}
type Entities struct {
	Name   string  `yaml:"name"`
	Fields []Field `yaml:"fields"`
}

type Field struct {
	Name     string `yaml:"name"`
	Type     string `yaml:"type"`
	Length   int    `yaml:"length"`
	Validate string `yaml:"validate"`
	Index    bool   `yaml:"index"`
}

func (c *Codegen) TemplateReplaces(entityName, param string) string {
	// Process template replacements
	rep := strings.NewReplacer(
		"[entity]", strings.ToLower(entityName),
		"[Entity]", strings.Title(strings.ToLower(genstrings.SnakeCaseToCamelCase(entityName))),
		// "[fields]", strings.Join(entityFields, "\n"),
		// "[indexScripts]", strings.Join(indexScripts, "\n"),
		// "[finderMethods]", strings.Join(finderMethods, "\n"),
		// "[crudMethods]", strings.Join(crudMethods, "\n"),
		"[port]", c.Port,
		"[appName]", c.AppName,
		"[driver]", c.Infrastructure.Db.Type,
		"[db_port]", strconv.Itoa(c.Infrastructure.Db.Port),
		"[db_user]", c.Infrastructure.Db.User,
		"[db_pass]", c.Infrastructure.Db.Password,
		"[db_name]", c.Infrastructure.Db.Database,
	)
	nc := rep.Replace(string(param))
	return nc
}

func (c *Codegen) FileReplaces(entityName, param string) string {
	rep := strings.NewReplacer(
		"sample", strings.ToLower(entityName),
		".go.tmpl", ".go",
		".yml.tmpl", ".yml",
		".mod.tmpl", ".mod",
		".tmpl", "",
		".gitignore.tmpl", ".gitignore",
		".sql.tmpl", ".sql",
	)
	nn := rep.Replace(param)
	return nn
}

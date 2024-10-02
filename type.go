package main

type Codegen struct {
	Name           string         `yaml:"name"`
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
	Name string `yaml:"name"`
}

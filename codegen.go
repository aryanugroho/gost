package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-yaml/yaml"
	"github.com/otiai10/copy"
)

func main() {
	Generate()
}

func Generate() {
	output := flag.String("output", "", "project output path directory")
	goPath := os.Getenv("GOPATH")
	fmt.Println("Current $GOPATH: ", goPath)
	if goPath == "" {
		fmt.Println("$GOPATH is not set, setting a new $GOPATH to $HOME/go")
		homePath := os.Getenv("HOME")
		os.Setenv("GOPATH", fmt.Sprintf("%s/go", homePath))
		goPath = os.Getenv("GOPATH")
		fmt.Println("New $GOPATH: ", goPath)
	}
	flag.Parse()

	// read codegen.yml
	var codegen Codegen
	codegenFile, err := ioutil.ReadFile("./codegen.yml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(codegenFile, &codegen)
	if err != nil {
		panic(err)
	}
	fmt.Println("reading codegen.yml")
	fmt.Printf("%+v \n", codegen)

	// component map
	sampleFile := "sample.go.tmpl"
	components := make(map[string]string)
	components["usecase"] = fmt.Sprint("/internal/usecase/", sampleFile)
	components["model"] = fmt.Sprint("/internal/model/", sampleFile)
	components["sqlstore"] = fmt.Sprint("/internal/infrastructure/sqlstore/", sampleFile)
	components["rest"] = fmt.Sprint("/internal/presenter/rest/api1/", sampleFile)

	// copy base
	basePath := fmt.Sprint(*output, "/", codegen.Name)
	targetDir := fmt.Sprint(goPath, "/src/", basePath)
	fmt.Println("generating baseline on", targetDir, "...")

	templateDir := "./_template"
	copy.Copy(templateDir, targetDir)

	// detect sqlstore type and keep it in codegen
	if codegen.Infrastructure.Db.Type == "postgres" {
		err = os.Remove(fmt.Sprint(targetDir, "/internal/infrastructure/sqlstore/sqlstore.mysql.go.tmpl"))
		if err != nil {
			panic(err)
		}
		err = os.Rename(fmt.Sprint(targetDir, "/internal/infrastructure/sqlstore/sqlstore.postgres.go.tmpl"), fmt.Sprint(targetDir, "/internal/infrastructure/sqlstore/sqlstore.go.tmpl"))
		if err != nil {
			panic(err)
		}
		err = os.Remove(fmt.Sprint(targetDir, "/internal/infrastructure/sqlstore/migration.mysql.go.tmpl"))
		if err != nil {
			panic(err)
		}
		err = os.Rename(fmt.Sprint(targetDir, "/internal/infrastructure/sqlstore/migration.postgres.go.tmpl"), fmt.Sprint(targetDir, "/internal/infrastructure/sqlstore/migration.go.tmpl"))
		if err != nil {
			panic(err)
		}
		err = os.Remove(fmt.Sprint(targetDir, "/docker-compose-mysql.yml.tmpl"))
		if err != nil {
			panic(err)
		}
		err = os.Rename(fmt.Sprint(targetDir, "/docker-compose-postgres.yml.tmpl"), fmt.Sprint(targetDir, "/docker-compose.yml.tmpl"))
		if err != nil {
			panic(err)
		}
	} else if codegen.Infrastructure.Db.Type == "mysql" {
		err = os.Remove(fmt.Sprint(targetDir, "/internal/infrastructure/sqlstore/sqlstore.postgres.go.tmpl"))
		if err != nil {
			panic(err)
		}
		err = os.Rename(fmt.Sprint(targetDir, "/internal/infrastructure/sqlstore/sqlstore.mysql.go.tmpl"), fmt.Sprint(targetDir, "/internal/infrastructure/sqlstore/sqlstore.go.tmpl"))
		if err != nil {
			panic(err)
		}
		err = os.Remove(fmt.Sprint(targetDir, "/internal/infrastructure/sqlstore/migration.postgres.go.tmpl"))
		if err != nil {
			panic(err)
		}
		err = os.Rename(fmt.Sprint(targetDir, "/internal/infrastructure/sqlstore/migration.mysql.go.tmpl"), fmt.Sprint(targetDir, "/internal/infrastructure/sqlstore/migration.go.tmpl"))
		if err != nil {
			panic(err)
		}
		err = os.Remove(fmt.Sprint(targetDir, "/docker-compose-postgres.yml.tmpl"))
		if err != nil {
			panic(err)
		}
		err = os.Rename(fmt.Sprint(targetDir, "/docker-compose-mysql.yml.tmpl"), fmt.Sprint(targetDir, "/docker-compose.yml.tmpl"))
		if err != nil {
			panic(err)
		}
	}

	// copy components
	fmt.Println("generating additional entities...")
	for _, entity := range codegen.Entities {
		entityFileName := entity.Name

		// generate components
		for component, componentPath := range components {
			fmt.Println("generating", entityFileName, component, "component...")

			componentSrc := fmt.Sprint(templateDir, componentPath)
			componentOut := fmt.Sprint(targetDir, componentPath)
			err = copy.Copy(componentSrc, componentOut)
			if err != nil {
				panic(err)
			}

			// generate fields is not supported yet
			// for _, field := range entity.Fields {
			// 	if field.Type == "string" {
			// 		field.Type = "string"
			// 	} else if field.Type == "int" {
			// 		field.Type = "int"
			// 	} else if field.Type == "bool" {
			// 		field.Type = "bool"
			// 	} else if field.Type == "float" {
			// 		field.Type = "float"
			// 	} else if field.Type == "time" {
			// 		field.Type = "time.Time"
			// 	} else if field.Type == "[]string" {
			// 		field.Type = "[]string"
			// 	} else if field.Type == "[]int" {
			// 		field.Type = "[]int"
			// 	} else if field.Type == "[]bool" {
			// 		field.Type = "[]bool"
			// 	} else if field.Type == "[]float" {
			// 		field.Type = "[]float"
			// 	} else if field.Type == "[]time" {
			// 		field.Type = "[]time.Time"
			// 	}
			// }
		}

		err := walkBuildFiles(targetDir, basePath, entity.Name, codegen)
		if err != nil {
			panic(err)
		}

	}
	copy.Copy(fmt.Sprintf("%s/cmd/server/main.go", targetDir), fmt.Sprintf("%s/cmd/%s/main.go", targetDir, codegen.Name))
	os.RemoveAll(fmt.Sprintf("%s/cmd/server/", targetDir))
	exec.Command("gofmt", "-s", "-w", targetDir)
	fmt.Println("updating dependencies...")
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = targetDir
	cmd.Output()
	fmt.Println(codegen.Name, "is generated successfully")

}

func walkBuildFiles(dir, proj, entityName string, codegen Codegen) error {
	return filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() && (strings.Contains(f.Name(), ".tmpl") || strings.Contains(f.Name(), ".go")) {
			// read file
			read, err := ioutil.ReadFile(path)
			if err != nil {
				fmt.Println("error reading file: ", err)
				return err
			}

			// prepare replacer
			rep := strings.NewReplacer(
				"[base_project]", proj,
				"[entity]", strings.ToLower(entityName),
				"[Entity]", strings.Title(strings.ToLower(snakeCaseToCamelCase(entityName))),
				"[port]", codegen.Port,
				"[appName]", codegen.Name,
				"[driver]", codegen.Infrastructure.Db.Type,
				"[db_port]", strconv.Itoa(codegen.Infrastructure.Db.Port),
				"[db_user]", codegen.Infrastructure.Db.User,
				"[db_pass]", codegen.Infrastructure.Db.Password,
				"[db_name]", codegen.Infrastructure.Db.Database,
			)

			nc := rep.Replace(string(read))
			err = ioutil.WriteFile(path, []byte(nc), 0)
			if err != nil {
				fmt.Println("error writing file: ", err)
				return err
			}

			//rename file
			rep = strings.NewReplacer(
				"sample", strings.ToLower(entityName),
				".go.tmpl", ".go",
				".yml.tmpl", ".yml",
				".mod.tmpl", ".mod",
				".tmpl", "",
				".gitignore.tmpl", ".gitignore",
				".sql.tmpl", ".sql",
			)

			nn := rep.Replace(path)
			err = os.Rename(path, nn)
			if err != nil {
				fmt.Println("error renaming file: ", err)
				return err
			}
		}

		return nil
	})
}

func snakeCaseToCamelCase(inputUnderScoreStr string) (camelCase string) {
	isToUpper := false

	for k, v := range inputUnderScoreStr {
		if k == 0 {
			camelCase = strings.ToUpper(string(inputUnderScoreStr[0]))
		} else {
			if isToUpper {
				camelCase += strings.ToUpper(string(v))
				isToUpper = false
			} else {
				if v == '_' {
					isToUpper = true
				} else {
					camelCase += string(v)
				}
			}
		}
	}
	return

}

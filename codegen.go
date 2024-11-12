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
	"time"

	"github.com/otiai10/copy"
	"gopkg.in/yaml.v2"
)

var (
	indexScripts  []string
	crudMethods   []string
	finderMethods []string
	entityFields  []string
)

func main() {
	output := flag.String("output", "", "project output path directory")
	flag.Parse()

	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		homePath := os.Getenv("HOME")
		os.Setenv("GOPATH", fmt.Sprintf("%s/go", homePath))
		goPath = os.Getenv("GOPATH")
	}

	codegen := loadCodegenConfig("./codegen.yaml")

	basePath := fmt.Sprint(*output, "/", codegen.Name)
	targetDir := fmt.Sprint(goPath, "/src/", basePath)
	templateDir := "./_template"

	err := copy.Copy(templateDir, targetDir)
	if err != nil {
		panic(err)
	}

	err = walkBuildFiles(targetDir, basePath, "", codegen)
	if err != nil {
		panic(err)
	}

	finalizeProject(targetDir, codegen.Name)
}

func loadCodegenConfig(path string) Codegen {
	var codegen Codegen
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	if err := yaml.Unmarshal(data, &codegen); err != nil {
		panic(err)
	}
	return codegen
}

func finalizeProject(targetDir, projectName string) {
	copy.Copy(fmt.Sprintf("%s/cmd/server/main.go", targetDir), fmt.Sprintf("%s/cmd/%s/main.go", targetDir, projectName))
	os.RemoveAll(fmt.Sprintf("%s/cmd/server/", targetDir))
	exec.Command("gofmt", "-s", "-w", targetDir).Run()

	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = targetDir
	if _, err := cmd.Output(); err != nil {
		fmt.Println("go mod tidy failed:", err)
	}
	fmt.Printf("%s is generated successfully\n", projectName)
}

// createMigrationFiles generates the migration files for the given entity
func createMigrationFiles(entity Entities) error {
	timestamp := time.Now().Unix()
	migrationDir := "migrations" // Adjust as needed

	fileNameUp := fmt.Sprintf("%s/%d_create_%s_table.up.sql", migrationDir, timestamp, entity.Name)
	fileNameDown := fmt.Sprintf("%s/%d_create_%s_table.down.sql", migrationDir, timestamp, entity.Name)

	// Create .up.sql file
	upFile, err := os.Create(fileNameUp)
	if err != nil {
		return err
	}
	defer upFile.Close()

	// Create table skeleton
	upFileContent := fmt.Sprintf("CREATE TABLE %s (\n", entity.Name)
	for _, field := range entity.Fields {
		fieldDef := fmt.Sprintf("    %s %s", field.Name, field.Type)
		if field.Length > 0 {
			fieldDef = fmt.Sprintf("    %s %s(%d)", field.Name, field.Type, field.Length)
		}
		if field.Index {
			fieldDef += " INDEX"
		}
		fieldDef += ",\n"
		upFileContent += fieldDef
	}
	upFileContent = strings.TrimSuffix(upFileContent, ",\n") // Remove trailing comma
	upFileContent += "\n);"

	_, err = upFile.WriteString(upFileContent)
	if err != nil {
		return err
	}

	// Create .down.sql file
	downFile, err := os.Create(fileNameDown)
	if err != nil {
		return err
	}
	defer downFile.Close()

	downFileContent := fmt.Sprintf("DROP TABLE IF EXISTS %s;", entity.Name)
	_, err = downFile.WriteString(downFileContent)
	if err != nil {
		return err
	}

	fmt.Printf("Migration files created for entity %s: %s and %s\n", entity.Name, fileNameUp, fileNameDown)
	return nil
}

func getType(field Field) string {
	goType := mapCustomTypeToGoType(field.Type)
	return fmt.Sprintf("%s `json:\"%s\" validate:\"%s\"`", goType, field.Name, field.Validate)
}

func walkBuildFiles(dir, proj, entityName string, codegen Codegen) error {
	return filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("error encountering file:", err)
			return err
		}
		if !f.IsDir() && (strings.Contains(f.Name(), ".tmpl") || strings.Contains(f.Name(), ".go")) {
			read, err := ioutil.ReadFile(path)
			if err != nil {
				fmt.Println("error reading file:", err)
				return err
			}

			// Clear any previous values (important if walkBuildFiles is called multiple times)
			indexScripts = nil
			crudMethods = nil
			finderMethods = nil
			entityFields = nil

			entityFields, indexScripts, finderMethods, crudMethods = generateEntityInfo(codegen.Entities)

			// Process template replacements
			rep := strings.NewReplacer(
				"[base_project]", proj,
				"[entity]", strings.ToLower(entityName),
				"[Entity]", strings.Title(strings.ToLower(snakeCaseToCamelCase(entityName))),
				"[fields]", strings.Join(entityFields, "\n"),
				"[indexScripts]", strings.Join(indexScripts, "\n"),
				"[finderMethods]", strings.Join(finderMethods, "\n"),
				"[crudMethods]", strings.Join(crudMethods, "\n"),
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
				fmt.Println("error writing file:", err)
				return err
			}

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
				fmt.Println("error renaming file:", err)
				return err
			}

			// Create migration files for each entity
			for _, entity := range codegen.Entities {
				err := createMigrationFiles(entity)
				if err != nil {
					fmt.Printf("Error creating migration files for entity %s: %v\n", entity.Name, err)
				}
			}
		}
		return nil
	})
}

func generateEntityInfo(entities []Entities) (entityFields, indexScripts, finderMethods, crudMethods []string) {
	for _, entity := range entities {
		for _, field := range entity.Fields {
			field.Type = getType(field)
			// Build entity field string
			entityField := fmt.Sprintf("%s %s", field.Name, field.Type)
			entityFields = append(entityFields, entityField)
			// Generate finder methods and indexes
			if field.Index {
				// Generate index logic
				indexLogic := fmt.Sprintf("CREATE INDEX idx_%s_%s ON %s (%s);", entity.Name, field.Name, strings.ToLower(entity.Name), field.Name)
				indexScripts = append(indexScripts, indexLogic)
			}
			// Generate finder methods
			finderMethod := fmt.Sprintf("func Find%sBy%s(%s %s) (*%s, error) { \n // finder logic here... \n}", entity.Name, strings.Title(field.Name), field.Name, field.Type, entity.Name)
			finderMethods = append(finderMethods, finderMethod)
		}

		// Generate CRUD logic
		crudCreate := fmt.Sprintf("func Create%s(entity *%s) error { \n // create logic here... \n}", entity.Name, entity.Name)
		crudRead := fmt.Sprintf("func Get%sByID(id int) (*%s, error) { \n // read logic here... \n}", entity.Name, entity.Name)
		crudUpdate := fmt.Sprintf("func Update%s(entity *%s) error { \n // update logic here... \n}", entity.Name, entity.Name)
		crudDelete := fmt.Sprintf("func Delete%sByID(id int) error { \n // delete logic here... \n}", entity.Name)

		crudMethods = append(crudMethods, crudCreate, crudRead, crudUpdate, crudDelete)
	}
	return
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

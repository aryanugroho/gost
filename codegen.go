package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/aryanugroho/gost/codegen"
	genstrings "github.com/aryanugroho/gost/pkg/strings"
	"github.com/otiai10/copy"
	"gopkg.in/yaml.v2"
)

func main() {
	output := flag.String("output", "", "project output path directory")
	lang := flag.String("lang", "", "project language")
	flag.Parse()

	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		homePath := os.Getenv("HOME")
		os.Setenv("GOPATH", fmt.Sprintf("%s/go", homePath))
		goPath = os.Getenv("GOPATH")
	}

	codegen := loadCodegenConfig("./codegen.yaml")

	basePath := fmt.Sprintf("%s/%s", *output, codegen.AppName)
	targetDir := fmt.Sprint(goPath, "/src/", basePath)
	templateDir := fmt.Sprintf("./_%s_template", *lang)

	fmt.Printf("Configuration: %v\n", codegen)
	fmt.Printf("Generating project using template: %s, into: %s\n", templateDir, targetDir)

	// copy all from template as a baseline
	err := copy.Copy(templateDir, targetDir)
	if err != nil {
		log.Fatalf("error copying template files: %v", err)
	}

	// generate entity files based on the configuration
	for _, entity := range codegen.Entities {
		_, _, _, _ = generateEntityInfo(codegen.Entities)
		entityName := genstrings.SnakeCaseToCamelCase(entity.Name)

		// generate entity files
		copy.Copy(fmt.Sprintf("%s/internal/usecase/[entity].go", templateDir), fmt.Sprintf("%s/internal/usecase/%s.go", targetDir, entityName))
		copy.Copy(fmt.Sprintf("%s/internal/repository/[entity].go", templateDir), fmt.Sprintf("%s/internal/repository/%s.go", targetDir, entityName))
		copy.Copy(fmt.Sprintf("%s/internal/delivery/http/[entity].go", templateDir), fmt.Sprintf("%s/internal/delivery/http/%s.go", targetDir, entityName))

		// write entity fields to the entity file
		//entityFieldsStr := strings.Join(entityFields, "\n")
	}

	err = walkBuildFiles(targetDir, codegen.Entities[0].Name, &codegen)
	if err != nil {
		log.Fatalf("error walking build files: %v", err)
	}

	finalizeProject(targetDir, codegen.AppName)
}

func loadCodegenConfig(path string) codegen.Codegen {
	var codegen codegen.Codegen
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("error reading file: %v", err)
	}
	if err := yaml.Unmarshal(data, &codegen); err != nil {
		log.Fatalf("error unmarshalling yaml: %v", err)
	}
	return codegen
}

func finalizeProject(targetDir, projectName string) {
	exec.Command("gofmt", "-s", "-w", targetDir).Run()

	exec.Command("go", "mod", "init")
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = targetDir
	if _, err := cmd.Output(); err != nil {
		fmt.Println("go mod tidy failed:", err)
	}
	fmt.Printf("%s is generated successfully\n", projectName)
}

// createMigrationFiles generates the migration files for the given entity
func createMigrationFiles(entity codegen.Entities) error {
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

func getType(field codegen.Field) string {
	goType := codegen.MapCustomTypeToGoType(field.Type)
	return fmt.Sprintf("%s `json:\"%s\" validate:\"%s\"`", goType, field.Name, field.Validate)
}

func fileRename(string) string {
	return ""
}

func walkBuildFiles(dir, entityName string, codegen *codegen.Codegen) error {
	// generate entity

	return filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("error encountering file:", err)
			return err
		}
		if strings.Contains(f.Name(), ".tmpl") {
			read, err := ioutil.ReadFile(path)
			if err != nil {
				fmt.Println("error reading file:", err)
				return err
			}

			nc := codegen.TemplateReplaces(entityName, string(read))
			err = ioutil.WriteFile(path, []byte(nc), 0)
			if err != nil {
				fmt.Println("error writing file:", err)
				return err
			}

			nn := codegen.FileReplaces(entityName, path)
			err = os.Rename(path, nn)
			if err != nil {
				fmt.Println("error renaming file:", err)
				return err
			}
		}
		return nil
	})
}

func generateEntityInfo(entities []codegen.Entities) ([]string, []string, []string, []string) {
	var sqlIndexScripts []string
	var finderMethods []string
	var crudMethods []string
	var entityFields []string
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
				sqlIndexScripts = append(sqlIndexScripts, indexLogic)
				finderMethods = append(finderMethods, "FindBy"+field.Name)
			}

		}
	}
	return entityFields, sqlIndexScripts, finderMethods, crudMethods
}

package main

import (
    "strings"
)

func replaceTemplatePlaceholders(content string, replacements map[string]string) string {
    keys := make([]string, 0, len(replacements))
    values := make([]string, 0, len(replacements))
    
    for k, v := range replacements {
        keys = append(keys, k)
        values = append(values, v)
    }

    rep := strings.NewReplacer(keys...)
    return rep.Replace(content)
}

func replaceTemplateExtensions(path, entityName string) string {
    rep := strings.NewReplacer(
        "sample", strings.ToLower(entityName),
        ".go.tmpl", ".go",
        ".yml.tmpl", ".yml",
        ".mod.tmpl", ".mod",
        ".tmpl", "",
        ".gitignore.tmpl", ".gitignore",
        ".sql.tmpl", ".sql",
    )
    return rep.Replace(path)
}

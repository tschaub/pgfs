package models

import "fmt"

func column(table, name string) string {
	return fmt.Sprintf("%s.%s", table, name)
}

func alias(name, alias string) string {
	return fmt.Sprintf("%s as %s", name, alias)
}

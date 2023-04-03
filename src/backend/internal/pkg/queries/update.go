package queries

import (
	"strconv"
	"strings"
)

func CreateUpdateQuery(entityName string, fields map[string]any) (string, []any) {
	query := `update ` + entityName + ` set (`

	keys := make([]string, 0, len(fields))
	values := make([]any, 0, len(fields))
	ids := make([]string, 0, len(fields))
	id := 1
	for key, value := range fields {
		keys = append(keys, key)
		values = append(values, value)
		ids = append(ids, "$"+strconv.Itoa(id))
		id++
	}
	query += strings.Join(keys, ", ") + ") = (" + strings.Join(ids, ", ") + ")"

	return query, values
}

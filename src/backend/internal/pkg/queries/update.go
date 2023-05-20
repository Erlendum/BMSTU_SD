package queries

import (
	"go.mongodb.org/mongo-driver/bson"
	"strconv"
	"strings"
)

func CreatePostgresUpdateQuery(entityName string, fields map[string]any) (string, []any) {
	query := `update ` + entityName + ` set `

	if len(fields) > 1 {
		query += `(`
	}

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
	query += strings.Join(keys, ", ")

	if len(fields) > 1 {
		query += ") = ("
	} else {
		query += " = "
	}

	query += strings.Join(ids, ", ")
	if len(fields) > 1 {
		query += ")"
	}

	return query, values
}

func CreateMongoUpdateQuery(fields map[string]any) bson.D {
	values := bson.D{}
	for key, value := range fields {
		values = append(values, bson.E{Key: key, Value: value})
	}
	return bson.D{{"$set", values}}
}

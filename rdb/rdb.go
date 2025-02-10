// This packages creates an abstraction layer for interacting with
// relational databases
package rdb

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"time"

	_ "github.com/tursodatabase/go-libsql"
)

var (
	ErrTypeOfModelIsNil  = errors.New("type of model is <nil/>")
	ErrModelMustBeStruct = errors.New("model must be a struct")
)

type ORM struct {
	db *sql.DB
}

type RawFn func(db *sql.DB)

func Open(driver, database string) (ORM, error) {
	db, err := sql.Open(driver, database)
	if err != nil {
		return ORM{}, err
	}
	return ORM{db}, nil
}

// Saves model to tablename in the database
func (r ORM) Save(model interface{}, tablename string) error {
	v := reflect.ValueOf(model).Elem()
	t := v.Type()

	if t == nil {
		return ErrTypeOfModelIsNil
	}

	if t.Kind() != reflect.Struct {
		return ErrModelMustBeStruct
	}

	var inserts []string
	var values []string
	for i := 0; i < t.NumField(); i += 1 {
		field := t.Field(i)
		tag_notation := field.Tag.Get("sql")
		tag_notation = strings.ToLower(tag_notation)

		if tag_notation == "" {
			continue
		}

		var columnname string
		var columntype string
		tags := strings.Split(tag_notation, ",")

		//WARNING: could crash here validate tags input
		columnname = tags[0]
		if len(tags) > 1 {
			columntype = strings.ToUpper(tags[1])
		} else {
			columntype = "TEXT"
		}

		if columnname == "omit" {
			continue
		}

		//TODO: handle flatten inner structure
		if columnname == "flatten" {
			continue
		}

		switch columntype {
		case "DATETIME":
			dtval := handle_time(v.Field(i), columntype, tags[2:])
			if dtval != "" {
				inserts = append(inserts, columnname)
				values = append(values, dtval)
			}
		case "INT", "INTEGER":
			field := v.Field(i)

			switch field.Kind() {
			case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
				if slices.Index(tags, "primary key") != -1 || slices.Index(tags, "autoincrement") != -1 {
					continue
				}
				inserts = append(inserts, columnname)
				values = append(values, fmt.Sprintf(`%d`, v.Field(i).Int()))
			case reflect.Struct:
				timeval := handle_time(field, columntype, tags[2:])
				if timeval != "" {
					inserts = append(inserts, columnname)
					values = append(values, timeval)
				}
			}

		default:
			inserts = append(inserts, columnname)
			values = append(values, fmt.Sprintf(`'%s'`, v.Field(i).String()))
		}

	}

	query := fmt.Sprintf("INSERT INTO %s (%s)\n VALUES (%s);\n", tablename, strings.Join(inserts, ","), strings.Join(values, ","))

	if _, err := r.db.Exec(query); err != nil {
		return err
	}

	return nil
}

// Creates a table in the datase mapped to model. If tablename name is ""  then the
// struct name is used
func (r ORM) CreateTable(model interface{}, tablename string) error {
	var t reflect.Type
	tt, ok := model.(reflect.Type)
	if ok {
		t = tt
	} else {
		t = reflect.TypeOf(model)
	}

	if tablename == "" {
		tablename = strings.ToLower(t.Name())
		if !strings.HasSuffix(tablename, "s") {
			tablename = tablename + "s"
		}
	}

	var columns []string
	for i := 0; i < t.NumField(); i += 1 {
		field := t.Field(i)
		sql_tag := field.Tag.Get("sql")

		var column_line string

		if sql_tag == "omit" {
			continue
		}

		if sql_tag == "" {
			if field.Type.Kind() == reflect.Struct {

				if field.Type == reflect.TypeOf(time.Time{}) {
					column_name := sanitize_keyword(column_line)
					columns = append(columns, fmt.Sprintf("%s INTEGER", column_name))
					continue
				}

				_, err := process_inner_struct(field.Type, r.db, "")
				if err != nil {
					return err
				}
				continue
			}

			column_name := sanitize_keyword(field.Name)

			column_line = fmt.Sprintf("%s TEXT", column_name)
			columns = append(columns, column_line)
			continue
		}

		tags := strings.Split(sql_tag, ",")

		if len(tags) == 1 {

			if field.Type.Kind() == reflect.Struct {
				if tags[0] == "flatten" {
					flatten_columns, err := process_inner_struct(field.Type, r.db, "flatten")
					if err != nil {
						return err
					}
					columns = append(columns, flatten_columns...)
				}
				continue
			}

			column_line = fmt.Sprintf("%s TEXT", tags[0])
			columns = append(columns, column_line)
			continue
		}

		column_line = strings.Join(tags, " ")
		columns = append(columns, column_line)

	}

	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s);", tablename, strings.Join(columns, ","))

	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func process_inner_struct(t reflect.Type, db *sql.DB, tag string) ([]string, error) {
	var columns []string
	for i := 0; i < t.NumField(); i += 1 {
		field := t.Field(i)
		sql_tag := field.Tag.Get("sql")

		if sql_tag == "omit" {
			continue
		}

		var column_line string

		if sql_tag == "" {
			if field.Type.Kind() == reflect.Struct {
				column_lines, err := process_inner_struct(field.Type, db, sql_tag)
				if err != nil {
					return []string{}, err
				}
				if len(column_lines) > 0 {
					table_name := strings.ToLower(t.Name())
					if !strings.HasSuffix(table_name, "s") {
						table_name = table_name + "s"
					}
					query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s);", table_name, strings.Join(columns, ","))
					_, err := db.Exec(query)
					if err != nil {
						return []string{}, err
					}
				}
				continue
			}

			column_name := sanitize_keyword(field.Name)

			column_line = fmt.Sprintf("%s TEXT", column_name)
			columns = append(columns, column_line)
			continue
		}

		tags := strings.Split(sql_tag, ",")

		if len(tags) == 1 {
			column_line = fmt.Sprintf("%s TEXT", tags[0])
			columns = append(columns, column_line)
			continue
		}

		column_line = strings.Join(tags, " ")
		columns = append(columns, column_line)

	}

	if tag != "flatten" {
		table_name := strings.ToLower(t.Name())
		if !strings.HasSuffix(table_name, "s") {
			table_name = table_name + "s"
		}
		query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s);", table_name, strings.Join(columns, ","))

		_, err := db.Exec(query)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}

	return columns, nil
}

// Handles time data type
func handle_time(field reflect.Value, columntype string, tags []string) string {
	var omitzero bool

	for _, v := range tags {
		if strings.Index(v, "default") != -1 {
			omitzero = true
			break
		}
	}

	if field.Type() == reflect.TypeOf(time.Time{}) {
		t := field.Interface().(time.Time)
		if t.IsZero() && omitzero {
			return ""
		}

		switch columntype {
		case "DATETIME":
			return fmt.Sprintf(`'%s'`, t.String())
		case "INT", "INTEGER":
			return fmt.Sprintf(`%d`, t.Unix())
		//NOTE: not sure what other cases to catch yet
		default:
			return ""
		}

	}

	return ""
}

// Chekcs if the column_name is a SQL keyword and wraps it
func sanitize_keyword(column_name string) string {
	column_name = strings.ToLower(column_name)
	_, ok := keywords[column_name]
	if !ok {
		return column_name
	}
	return fmt.Sprintf("[%s]", column_name)
}

// Raw query
func (orm ORM) Raw(f RawFn) {
	f(orm.db)
}

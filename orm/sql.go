package orm

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
)

// GetColNames 获取列名列表， modelPtr为数据对象的指针，tagName为成员描述符（即，"db", "gorm"等）
func GetColNames(modelPtr interface{}, tagName string) ([]string, error) {
	if reflect.ValueOf(modelPtr).Kind() != reflect.Ptr {
		return nil, errors.New("need a pointer")
	}

	var cols []string
	t := reflect.TypeOf(modelPtr).Elem()
	for i := 0; i < t.NumField(); i++ {
		cols = append(cols, t.Field(i).Tag.Get(tagName))
	}

	return cols, nil
}

// GetColumns 获取列字段指针
func GetColumns(modelPtr interface{}) ([]interface{}, error) {
	if reflect.ValueOf(modelPtr).Kind() != reflect.Ptr {
		return nil, errors.New("need a pointer")
	}

	var cols []interface{}
	value := reflect.ValueOf(modelPtr).Elem()
	for i := 0; i < value.NumField(); i++ {
		cols = append(cols, value.Field(i).Addr().Interface())
	}

	return cols, nil
}

// Query 执行sql语句，modelPtr为数据对象的指针
func Query(ctx context.Context, db *sql.DB, sql string, modelPtr interface{}) ([]interface{}, error) {
	if reflect.ValueOf(modelPtr).Kind() != reflect.Ptr {
		return nil, errors.New("need a pointer")
	}

	rows, err := db.QueryContext(ctx, sql)
	if err != nil {
		return nil, err
	}

	defer func() {
		err = rows.Close()
	}()

	var records []interface{}
	for rows.Next() {
		cp := reflect.New(reflect.Indirect(reflect.ValueOf(modelPtr)).Type())
		record := cp.Interface()

		fields, _ := GetColumns(record)
		if err := rows.Scan(fields...); err != nil {
			continue
		}

		records = append(records, record)
	}

	return records, nil
}

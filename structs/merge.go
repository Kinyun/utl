package structs

import (
	"database/sql"
	"github.com/mitchellh/mapstructure"
	"github.com/matthewhartstonge/argon2"
	"reflect"
	"strconv"
	"strings"
	"fmt"
)

// Merge receives two structs, and merges them excluding fields with tag name: `structs`, value "-"
func Merge(dst, src interface{}) {
	s := reflect.ValueOf(src)
	d := reflect.ValueOf(dst)
	if s.Kind() != reflect.Ptr || d.Kind() != reflect.Ptr {
		return
	}
	for i := 0; i < s.Elem().NumField(); i++ {
		v := s.Elem().Field(i)
		fieldName := s.Elem().Type().Field(i).Name
		skip := s.Elem().Type().Field(i).Tag.Get("json")
		if skip == "-" {
			continue
		}
		if v.Kind() > reflect.Float64 &&
			v.Kind() != reflect.String &&
			v.Kind() != reflect.Struct &&
			v.Kind() != reflect.Ptr &&
			v.Kind() != reflect.Slice {
			continue
		}
		if v.Kind() == reflect.Ptr {
			// Field is pointer check if it's nil or set
			if !v.IsNil() {
				// Field is set assign it to dest

				if d.Elem().FieldByName(fieldName).Kind() == reflect.Ptr {
					d.Elem().FieldByName(fieldName).Set(v)
					continue
				}
				f := d.Elem().FieldByName(fieldName)
				if f.IsValid() {
					f.Set(v.Elem())
				}
			}
			continue
		}
		d.Elem().FieldByName(fieldName).Set(v)
	}
}

func MergeRow(rows *sql.Rows, dst interface{}) {
	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		fmt.Println(err.Error())
	}
	// Get column type
	columns_type, err := rows.ColumnTypes()
	if err != nil {
		fmt.Println(err.Error())
	}

	// Make a slice for the values
	values := make([]sql.RawBytes, len(columns))

	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	// get RawBytes from data
	err = rows.Scan(scanArgs...)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	// Now do something with the data.
	// Here we just print each column as a string.
	m := make(map[string]interface{})
	//var value string
	for i, col := range values {
		// Here we can check if the value is nil (NULL value)
		if col != nil {
			if strings.ToLower(columns_type[i].DatabaseTypeName()) == "varchar" {
				m[columns[i]] = string(col)
			} else if strings.ToLower(columns_type[i].DatabaseTypeName()) == "decimal" {
				s, _ := strconv.ParseFloat(string(col), 64)
				m[columns[i]] = s
			} else if strings.ToLower(columns_type[i].DatabaseTypeName()) == "int" {
				m[columns[i]], err = strconv.Atoi(string(col))
			} else {
				m[columns[i]], err = strconv.Atoi(string(col))
				if err != nil {
					m[columns[i]] = string(col)
				}
			}
		}
	}

	config := &mapstructure.DecoderConfig{
		TagName: "json",
		Result:  &dst,
	}
	decoder, _ := mapstructure.NewDecoder(config)
	err = decoder.Decode(m)
	if err != nil {
		fmt.Println(err.Error())
	}
}

//Check Different struct between data from database and req update from front-end
//then create query Set (MySQL) using prepare statement based on comparision result and set the value on interface
func DifSqlSet(src, req interface{}, setUpdate *strings.Builder, binds *[]interface{}) {
	bind := []interface{}{}
	s := reflect.ValueOf(src)
	r := reflect.ValueOf(req)
	if s.Kind() != reflect.Ptr || r.Kind() != reflect.Ptr {
		return
	}
	for i := 0; i < s.Elem().NumField(); i++ {
		v := s.Elem().Field(i)
		fieldName := s.Elem().Type().Field(i).Name
		vr := r.Elem().FieldByName(fieldName)
		tagName := s.Elem().Type().Field(i).Tag.Get("json")
		if tagName == "-" {
			continue
		}
		if v.Kind() > reflect.Float64 &&
			v.Kind() != reflect.String &&
			v.Kind() != reflect.Struct &&
			v.Kind() != reflect.Ptr &&
			v.Kind() != reflect.Slice {
			continue
		}

		if v.Interface() != vr.Interface() {
			if tagName == "password" {
				raws, _ := argon2.Decode([]byte(v.Interface().(string)))
				ok, _ := raws.Verify([]byte(vr.Interface().(string)))
				if ok {
					continue
				}
			}
			if len(setUpdate.String()) == 0 {
				setUpdate.WriteString(" SET ")
			} else {
				setUpdate.WriteString(", ")
			}
			setUpdate.WriteString(tagName)
			setUpdate.WriteString(" = ? ")
			if tagName == "password" {
				secArgon2 := argon2.DefaultConfig()
				raw, _ := secArgon2.Hash([]byte(vr.Interface().(string)), nil)
				bind = append(bind, string(raw.Encode()))
			} else {
				bind = append(bind, ConvertValue(vr))
			}
		}
	}
	*binds = append(*binds, bind...)
}

func ConvertValue(val reflect.Value) interface{} {
	v := val.Interface()
	switch v.(type) {
	case int:
		return v.(int)
	case int32:
		return v.(int32)
	case int64:
		return v.(int64)
	case float32:
		return v.(float32)
	case float64:
		return v.(float64)
	default:
		return v.(string)
	}

}
package main

import (
	"errors"
	"fmt"
	"strings"
)

const (
	DDL_TYPE_UNKNOWN = 0 + iota
	DDL_TYPE_UINT8
	DDL_TYPE_UINT16
	DDL_TYPE_UINT32
	DDL_TYPE_UINT64
	DDL_TYPE_INT8
	DDL_TYPE_INT16
	DDL_TYPE_INT32
	DDL_TYPE_INT64
	DDL_TYPE_BOOL
	DDL_TYPE_STRING
	DDL_TYPE_INT
	DDL_TYPE_FLOAT32
	DDL_TYPE_FLOAT64
)

type TableFiled struct {
	FieldName    string
	FieldType    int
	FieldComment string
}

type TableIndex struct {
	IndexName string
	FiledName []string
}

type TableDDL struct {
	TableName string
	PkName    string
	PkType    int
	Engine    string
	Charset   string
	Indexs    []*TableIndex
	Fileds    []*TableFiled
}

func Type2Str(tp int) string {
	switch tp {
	case DDL_TYPE_INT:
		return "int"
	case DDL_TYPE_INT8:
		return "int8"
	case DDL_TYPE_INT16:
		return "int16"
	case DDL_TYPE_INT32:
		return "int32"
	case DDL_TYPE_INT64:
		return "int64"
	case DDL_TYPE_UINT8:
		return "uint8"
	case DDL_TYPE_UINT16:
		return "uint16"
	case DDL_TYPE_UINT32:
		return "uint32"
	case DDL_TYPE_UINT64:
		return "uint64"
	case DDL_TYPE_STRING:
		return "string"
	case DDL_TYPE_FLOAT32:
		return "float32"
	case DDL_TYPE_FLOAT64:
		return "float64"
	case DDL_TYPE_UNKNOWN:
		return "string"
	}
	return "string"
}

func ParseDDL(ddl_string string) (*TableDDL, error) {
	lines := strings.Split(strings.Trim(ddl_string, " "), "\n")
	ddl := new(TableDDL)
	for i, line := range lines {
		//获取表名
		if i == 0 {
			f := strings.Split(strings.Trim(line, " "), " ")
			if len(f) < 3 || strings.ToLower(f[0]) != "create" || strings.ToLower(f[1]) != "table" {
				return nil, errors.New("ddl error fist line")
			}
			ddl.TableName = strings.Trim(f[2], "`")
			continue
		}
		if len(line) == 0 {
			continue
		}

		parts := strings.Split(strings.Trim(line, " ,"), " ")
		if len(parts) < 3 {
			return nil, fmt.Errorf("[%s] error", line)
		}

		//主键
		if strings.ToUpper(parts[0]) == "PRIMARY" && strings.ToUpper(parts[1]) == "KEY" {
			ddl.PkName = strings.Trim(parts[2], "(`)")
			continue
		}

		//索引
		if strings.ToUpper(parts[0]) == "KEY" {
			index := new(TableIndex)
			index.IndexName = strings.Trim(parts[1], "`")
			tmp := strings.Split(strings.Trim(parts[2], "()"), ",")
			for _, t := range tmp {
				index.FiledName = append(index.FiledName, strings.Trim(t, "`"))
			}
			ddl.Indexs = append(ddl.Indexs, index)
			continue
		}

		//唯一主键
		if strings.ToUpper(parts[0]) == "UNIQUE" {
			continue
		}

		// 表信息
		if parts[0] == ")" {
			//ddl.Engine = strings.Split(parts[1], "=")[1]
			//ddl.Charset = strings.Split(parts[3], "=")[1]
			continue
		}

		//字段
		field := new(TableFiled)
		for j, _ := range parts {
			if j == len(parts)-1 {
				continue
			}

			if j == 0 {
				field.FieldName = strings.Trim(parts[j], "`")
				continue
			}
			if j == 1 {
				tp := strings.ToUpper(strings.Split(parts[j], "(")[0])
				sign := strings.ToUpper(parts[j+1])

				switch tp {
				case "INT":
					field.FieldType = DDL_TYPE_INT32
				case "TINYINT":
					field.FieldType = DDL_TYPE_INT8
				case "BIGINT":
					field.FieldType = DDL_TYPE_INT64
				case "VARCHAR":
					field.FieldType = DDL_TYPE_STRING
				case "TEXT":
					field.FieldType = DDL_TYPE_STRING
				case "TIMESTAMP":
					field.FieldType = DDL_TYPE_INT64
				case "FLOAT":
					field.FieldType = DDL_TYPE_FLOAT32
				case "DOUBLE":
					field.FieldType = DDL_TYPE_FLOAT64
				}
				if sign == "UNSIGNED" {
					switch tp {
					case "INT":
						field.FieldType = DDL_TYPE_UINT32
					case "TINYINT":
						field.FieldType = DDL_TYPE_UINT8
					case "BIGINT":
						field.FieldType = DDL_TYPE_UINT64
					}
				}
			}

			if strings.ToUpper(parts[j]) == "COMMENT" {
				field.FieldComment = strings.Trim(parts[j+1], "'")
			}
		}
		ddl.Fileds = append(ddl.Fileds, field)
	}

	//主键类型
	for _, field := range ddl.Fileds {
		if field.FieldName == ddl.PkName {
			ddl.PkType = field.FieldType
			break
		}
	}
	return ddl, nil
}

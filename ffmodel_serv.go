package main

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
)

var MODLE_TPL string = `
package main

import (
	"github.com/astaxie/beego/orm"
	"fmt"
)

// ==========================
type {{.ModelName}} struct {
	{{range .FieldContent}}
	{{.FieldName}}	{{.FieldType}} {{.Tag|unescaped}} {{.Comment}} {{end}}
}

func init() {
	orm.RegisterModel(new({{.ModelName}}))
}

func (obj *{{.ModelName}}) TableName() string {
	return "{{.TableName}}"
}

func Add{{.ModelName}}(obj *{{.ModelName}}) error {
	o := orm.NewOrm()
	_, err := o.Insert(obj)
	return err
}

func Del{{.ModelName}}(id {{.PkType}}) error {
	o := orm.NewOrm()
	_, err := o.Delete(&{{.ModelName}}{Id: id})
	return err
}

func Get{{.ModelName}}(id {{.PkType}}) (*{{.ModelName}}, error) {
	o := orm.NewOrm()
	obj := &{{.ModelName}}{Id: id}
	err := o.Read(obj, "Id")
	return obj, err
}

func Update{{.ModelName}}(obj *{{.ModelName}}) error {
	o := orm.NewOrm()
	_, err := o.Update(obj)
	return err
}

func Query{{.ModelName}}(page, page int) ([]*{{.ModelName}}, int64, error){
	o := orm.NewOrm()
	qs := o.QueryTable("{{.TableName}}")
	if page <= 0 {
		page = 1
	}

	if size <= 0 {
		size = 10
	}

	offset := (page - 1) * size
	var objs []*{{.ModelName}}
	_, err := qs.Limit(size, offset).All(&objs)

	if err != nil {
		return nil, 0, err
	}

	cnt, _ := qs.Count()
	return objs, cnt,nil
}

// ==========================

`

// 定义函数unescaped
func unescaped(x string) interface{} { return template.HTML(x) }

func DDL2Model(ddl_string string, json_tag, orm_tag bool) (string, error) {
	ddl, err := ParseDDL(ddl_string)
	if err != nil {
		return "", err
	}
	buf := bytes.NewBuffer(nil)
	t := template.New("ddl")
	t = t.Funcs(template.FuncMap{"unescaped": unescaped})
	t, err = t.Parse(MODLE_TPL)
	if err != nil {
		return "", err
	}

	mp := make(map[string]interface{})
	mp["TableName"] = ddl.TableName
	mp["ModelName"] = camelString(ddl.TableName)
	mp["PkType"] = Type2Str(ddl.PkType)
	type field_t struct {
		FieldName string
		FieldType string
		Tag       string
		Comment   string
	}

	var fields []*field_t
	for _, field := range ddl.Fileds {
		tmp := new(field_t)
		tmp.FieldName = camelString(strings.ToLower(field.FieldName))
		tmp.FieldType = Type2Str(field.FieldType)
		if field.FieldComment != "" {
			tmp.Comment = "// " + field.FieldComment
		}
		var tags []string
		if json_tag {
			tags = append(tags, fmt.Sprintf("json:\"%s\"", field.FieldName))
		}
		if orm_tag {
			if field.FieldName == ddl.PkName {
				tags = append(tags, fmt.Sprintf("orm:\"pk;column(%s)\"", field.FieldName))
			} else {
				tags = append(tags, fmt.Sprintf("orm:\"column(%s)\"", field.FieldName))
			}
		}
		if len(tags) != 0 {
			tmp.Tag = fmt.Sprintf("`%s`", strings.Join(tags, " "))
		}
		fields = append(fields, tmp)

	}
	mp["FieldContent"] = fields

	err = t.Execute(buf, mp)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

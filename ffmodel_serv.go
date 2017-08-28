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
type {{.ModelName}}Model struct {
	{{range .FieldContent}}
	{{.FieldName}}	{{.FieldType}} {{.Tag|unescaped}} {{.Comment}} {{end}}
}

func init() {
	orm.RegisterModel(new({{.ModelName}}Model))
}

func (*{{.ModelName}}Model) TableName() string {
	return "{{.TableName}}"
}

func Add{{.ModelName}}(obj *{{.ModelName}}Model) error {
	o := orm.NewOrm()
	_, err := o.Insert(obj)
	return err
}

func Del{{.ModelName}}(id {{.PkType}}) error {
	o := orm.NewOrm()
	_, err := o.Delete(&{{.ModelName}}Model{ID: id})
	return err
}

func Get{{.ModelName}}(id {{.PkType}}) (*{{.ModelName}}Model, error) {
	o := orm.NewOrm()
	obj := &{{.ModelName}}Model{ID: id}
	err := o.Read(obj, "ID")
	return obj, err
}

func Update{{.ModelName}}(obj *{{.ModelName}}Model) error {
	o := orm.NewOrm()
	_, err := o.Update(obj)
	return err
}

func Query{{.ModelName}}(keys []string, values []interface{}, page_no, page_count int) ([]*{{.ModelName}}Model, int64, error){
	o := orm.NewOrm()
	qs := o.QueryTable("{{.TableName}}")
	if len(keys) != len(values) {
		return nil, 0, fmt.Errorf("key[%d] value[%d] not equal", len(keys), len(values))
	}
	for i,_ := range keys {
		qs = qs.Filter(keys[i], values[i])
	}
	if page_no == 0 {
		page_no = 0
	}

	if page_count == 0 {
		page_count = 100
	}

	var objs []*{{.ModelName}}Model
	_, err := qs.Limit(page_count, page_count * page_no).All(&objs)

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

func DDL2Model(ddl_string string) (string, error) {
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
		tmp.Tag = fmt.Sprintf("`orm:\"colunm(%s)\" json:\"%s\"`", field.FieldName, strings.ToLower(field.FieldName))

		if field.FieldName == ddl.PkName {
			tmp.Tag = fmt.Sprintf("`orm:\"pk;colunm(%s)\" json:\"%s\"`", field.FieldName, strings.ToLower(field.FieldName))
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

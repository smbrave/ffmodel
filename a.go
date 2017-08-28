
package main

import (
	"github.com/astaxie/beego/orm"
	"fmt"
)

// ==========================
type TbURLModel struct {
	
	ID	int32 `orm:"pk;colunm(id)" json:"id"`  
	IntField	int32 `orm:"colunm(int_field)" json:"int_field"` // 整型 
	DateField	string `orm:"colunm(date_field)" json:"date_field"`  
	DatetimeFiled	string `orm:"colunm(datetime_filed)" json:"datetime_filed"`  
	TimestampFiled	int64 `orm:"colunm(timestamp_filed)" json:"timestamp_filed"`  
	TextField	string `orm:"colunm(text_field)" json:"text_field"`  
	DoubleFiled	float64 `orm:"colunm(double_filed)" json:"double_filed"`  
	TinyInt	int8 `orm:"colunm(tiny_int)" json:"tiny_int"`  
}

func init() {
	orm.RegisterModel(new(TbURLModel))
}

func (*TbURLModel) TableName() string {
	return "tb_url"
}

func AddTbURL(obj *TbURLModel) error {
	o := orm.NewOrm()
	_, err := o.Insert(obj)
	return err
}

func DelTbURL(id int32) error {
	o := orm.NewOrm()
	_, err := o.Delete(&TbURLModel{ID: id})
	return err
}

func GetTbURL(id int32) (*TbURLModel, error) {
	o := orm.NewOrm()
	obj := &TbURLModel{ID: id}
	err := o.Read(obj, "ID")
	return obj, err
}

func UpdateTbURL(obj *TbURLModel) error {
	o := orm.NewOrm()
	_, err := o.Update(obj)
	return err
}

func QueryTbURL(keys []string, values []interface{}, page_no, page_count int) ([]*TbURLModel, int64, error){
	o := orm.NewOrm()
	qs := o.QueryTable("tb_url")
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

	var objs []*TbURLModel
	_, err := qs.Limit(page_count, page_count * page_no).All(&objs)

	if err != nil {
		return nil, 0, err
	}

	cnt, _ := qs.Count()
	return objs, cnt,nil
}

// ==========================



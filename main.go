package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

var (
	flag_orm_tag  = flag.Bool("orm", false, "output orm tag")
	flag_json_tag = flag.Bool("json", false, "output json tag")
	flag_sql_file = flag.String("sql", "t.sql", "input sql file")
)

func main() {
	if len(os.Args) <= 1 {
		return
	}
	flag.Parse()

	content, err := ioutil.ReadFile(*flag_sql_file)
	if err != nil {
		fmt.Printf("%s %s\n", *flag_sql_file, err.Error())
		return
	}

	result, err := DDL2Model(string(content), *flag_json_tag, *flag_orm_tag)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return
	}
	fmt.Printf("%s\n", result)
}

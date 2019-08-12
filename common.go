package main

import (
	"database/sql"
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

func NewJwtParser(key interface{}) func (t *jwt.Token) (interface{}, error){
	return func(t *jwt.Token)(interface{}, error){
		if t.Method.Alg() != "HS256" {
			return nil, fmt.Errorf("unexpected jwt signing method=%v", t.Header["alg"])
		}
		return key, nil
	}
}

func GetTableFields(db *sql.DB,database string,table string) (interface{},error){
	rows,err:=db.Query("select `column_name`,`data_type`,`column_default`,`is_nullable`,`extra` from information_schema.columns where table_schema=? and table_name=? order by ordinal_position",database,table)
	if err != nil{
		return nil,err
	}
	defer rows.Close()
	var field struct{
		Name string `json:"name"`
		Type string `json:"type"`
		Default interface{} `json:"default"`
		NullAble string `json:"nullable"`
		Extra string `json:"extra"`
	}
	var fields []struct{
		Name string `json:"name"`
		Type string `json:"type"`
		Default interface{} `json:"default"`
		NullAble string `json:"nullable"`
		Extra string `json:"extra"`
	}
	for rows.Next(){
		rows.Scan(&field.Name,&field.Type,&field.Default,&field.NullAble,&field.Extra)
		if df,ok:=field.Default.([]byte);ok{
			field.Default=string(df)
		}
		fields= append(fields, field)
	}
	return fields,nil
}

func isTableExist(db *sql.DB,dbName string,table string) bool{
	var count int
	db.QueryRow("select count(*) from information_schema.TABLES t where t.TABLE_SCHEMA =? and t.TABLE_NAME =?",dbName,table).Scan(&count)
	if count>0{
		return true
	}else{
		return false
	}
}

func getPage(){

}
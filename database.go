package main

import (
	"database/sql"
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/qianlidongfeng/toolbox"
)

type Database struct{
	db *sql.DB
}

func NewDatabase(d *sql.DB) *Database{
	return &Database{
		db:d,
	}
}

func (this *Database) GetUserPasswordByName(name string) string{
	var password string
	this.db.QueryRow("select password from users where `name`=?",name).Scan(&password)
	return password
}


func (this *Database) UpdateUserPasswordByName(name string,password string) error{
	_,err:=this.db.Exec("update users set password=? where `name`=?",password,name)
	if err != nil{
		log.Warn(err)
	}
	return err
}

func (this *Database) AddUser(name string,password string,group int) error{
	_,err:=this.db.Exec("insert ignore into users (`name`,password,`group`)value(?,?,?)",name,password,group)
	if err != nil{
		return err
	}
	return nil
}

func (this *Database) AddMenuItem(name string,parent int,parentName string,allow string,creater string) error{
	_,err:=this.db.Exec("insert into menu_items (`name`,parent,parentName,allow,creater,ctime)value(?,?,?,?,?,?)",name,parent,parentName,allow,creater,toolbox.GetTimeSecond())
	if err != nil{
		return err
	}
	return nil
}

func (this *Database) UpdateMenuItem(menuID int,name string,parent int,parentName string,allow string) error{
	_,err:=this.db.Exec("update menu_items set `name`=?,parent=?,parentName=?,allow=? where id=?",name,parent,parentName,allow,menuID)
	if err != nil{
		return err
	}
	return nil
}

func (this *Database) GetMenuItemParentID(menuID int) int{
	var id int
	this.db.QueryRow("select parent from menu_items where `id`=?",menuID).Scan(&id)
	return id
}


func (this *Database) AddViewerItem(name string,parent int,parentName string,allow string,creater string,dbAddress string,dbPassword string,dbName string,dbTable string,dbUser string,fields string) error{
	_,err:=this.db.Exec("insert into viewer_items (`name`,parent,parentName,allow,creater,dbAddress,dbPassword,dbName,dbTable,dbUser,fields,ctime)value(?,?,?,?,?,?,?,?,?,?,?,?)",
		name,parent,parentName,allow,creater,dbAddress,dbPassword,dbName,dbTable,dbUser,fields,toolbox.GetTimeSecond())
	if err != nil{
		return err
	}
	return nil
}

func (this *Database) UpdateViewerItem(viewerID int,name string,parent int,parentName string,allow string,creater string,dbAddress string,dbPassword string,dbName string,dbTable string,dbUser string,fields string) error{
	_,err:=this.db.Exec("update viewer_items set `name`=?,parent=?,parentName=?,allow=?,creater=?,dbAddress=?,dbPassword=?,dbName=?,dbTable=?,dbUser=?,fields=?,ctime=? where id=?",
		name,parent,parentName,allow,creater,dbAddress,dbPassword,dbName,dbTable,dbUser,fields,toolbox.GetTimeSecond(),viewerID)
	if err != nil{
		return err
	}
	return nil
}

func (this *Database) GetViewerItemParentID(viewerID int) int{
	var id int
	this.db.QueryRow("select parent from viewer_items where `id`=?",viewerID).Scan(&id)
	return id
}

func (this *Database) GetViewerItem(parent int,user string) (interface{},error){
	var rows *sql.Rows
	var err error
	if user=="admin"{
		rows,err=this.db.Query("select `id`,`name` from viewer_items where parent=? and id > 3 order by id",parent)
		if err != nil{
			return nil,err
		}
	}else{
		rows,err=this.db.Query("select `id`,`name` from viewer_items where parent=? and json_contains(allow,json_array(?)) and id>3 order by id",parent,user)
		if err != nil{
			return nil,err
		}
	}
	defer rows.Close()
	var items []struct{
		ID int `json:"id"`
		Name string `json:"name"`
	}
	for rows.Next(){
		var item struct{
			ID int `json:"id"`
			Name string `json:"name"`
		}
		rows.Scan(&item.ID,&item.Name)
		items=append(items,item)
	}
	return items,nil
}

func (this *Database) CheckMenuItemExist(name string,parent int,creater string) bool{
	var count int
	this.db.QueryRow("select count(*) from menu_items where `name`=? and parent=? and creater=?",name,parent,creater).Scan(&count)
	if count == 0{
		return false
	}
	return true
}

func (this *Database) CheckViewerItemExist(name string,parent int,creater string) bool{
	var count int
	this.db.QueryRow("select count(*) from viewer_items where `name`=? and parent=? and creater=?",name,parent,creater).Scan(&count)
	if count == 0{
		return false
	}
	return true
}

func (this *Database) GetMenuItem(parent int,user string) (interface{},error){
	var rows *sql.Rows
	var err error
	if user=="admin"{
		rows,err=this.db.Query("select `id`,`name` from menu_items where parent=? order by id",parent)
		if err != nil{
			return nil,err
		}
	}else{
		rows,err=this.db.Query("select `id`,`name` from menu_items where parent=? and json_contains(allow,json_array(?)) order by id",parent,user)
		if err != nil{
			return nil,err
		}
	}
	defer rows.Close()
	var items []struct{
		ID int `json:"id"`
		Name string `json:"name"`
	}
	for rows.Next(){
		var item struct{
			ID int `json:"id"`
			Name string `json:"name"`
		}
		rows.Scan(&item.ID,&item.Name)
		items=append(items,item)
	}
	return items,nil
}

func (this *Database) CheckUserPermissionOfMenuItem(id int,user string) bool{
	var count int
	this.db.QueryRow("select count(*) from menu_items where id=? and json_contains(allow,json_array(?))",id,user).Scan(&count)
	if count>0{
		return true
	}
	return false
}


func (this *Database) CheckUserPermissionOfViewerItem(id int,user string) bool{
	var count int
	this.db.QueryRow("select count(*) from viewer_items where id=? and json_contains(allow,json_array(?))",id,user).Scan(&count)
	if count>0{
		return true
	}
	return false
}

func (this *Database) GetUserGroupByName(name string) int{
	var group int
	this.db.QueryRow("select `group` from users where `name`=?",name).Scan(&group)
	return group
}

func (this *Database) GetViewerFields(viewerID int,name string) (toolbox.MySqlConfig,string){
	var fields string
	var dbConfig toolbox.MySqlConfig
	if name=="admin"{
		this.db.QueryRow("select dbAddress,dbPassword,dbName,dbTable,dbUser,fields from viewer_items where id=?",viewerID).Scan(
			&dbConfig.Address,&dbConfig.PassWord,&dbConfig.DataBase,&dbConfig.Table,&dbConfig.User,&fields)
	}else{
		this.db.QueryRow("select dbAddress,dbPassword,dbName,dbTable,dbUser,fields from viewer_items where id=? and json_contains(allow,json_array(?))",viewerID,name).Scan(
			&dbConfig.Address,&dbConfig.PassWord,&dbConfig.DataBase,&dbConfig.Table,&dbConfig.User,&fields)
	}
	return dbConfig,fields
}

func (this *Database) GetSysViewerFields(viewerID int,name string) (toolbox.MySqlConfig,string){
	var fields string
	var dbConfig toolbox.MySqlConfig
	this.db.QueryRow("select dbAddress,dbPassword,dbName,dbTable,dbUser,fields from viewer_items where id=?",viewerID).Scan(
		&dbConfig.Address,&dbConfig.PassWord,&dbConfig.DataBase,&dbConfig.Table,&dbConfig.User,&fields)
	return dbConfig,fields
}

func (this *Database) GetMenuNameByID(id int) string{
	var name string
	this.db.QueryRow("select `name` from menu_items where `id`=?",id).Scan(&name)
	return name
}

func (this *Database) GetMenuByUser(user string) []map[string]interface{}{
	var count int
	this.db.QueryRow("select count(*) from menu_items where json_contains(allow,json_array(?))",user).Scan(&count)
	var sqlSentence=fmt.Sprintf("select * from menu_items where json_contains(allow,json_array('%s'))",user)
	re,_:=toolbox.SelectMapFromMysql(this.db,sqlSentence)
	return re
}

func (this *Database) UpdateUserLoginTime(user string) error{
	_,err:=this.db.Exec("update users set ltime=? where `name`=?",toolbox.GetTimeSecond(),user)
	return err
}
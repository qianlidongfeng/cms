package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/qianlidongfeng/toolbox"
	"os"
)

func init(){
	dir,err:=toolbox.AppDir()
	if err != nil{
		panic(err)
	}
	os.Chdir(dir)
	G=Global{}
	G.Init()
	InitTables()
}

func InitTables(){
	InitUserTable()
	InitMenuItemTable()
	InitViewerItemTable()
}

func InitUserTable(){
	fileds:=map[string]string{
		"id":`int(11) UNSIGNED NOT NULL AUTO_INCREMENT`,
		"name":`varchar(32) NOT NULL`,
		"password":`varchar(64) NOT NULL`,
		"group":`int(11) NOT NULL`,
		"ltime":`datetime(0)`,
	}
	index:=map[string]toolbox.MysqlIndex{
			"name":toolbox.MysqlIndex{Typ:"unique",Name:"column",Method:"BTREE"},
	}
	if err:=toolbox.CheckAndFixTable(G.db,"users",fileds,index);err!=nil{
		G.loger.Fatal(err)
	}
	if err:=G.DB.AddUser("admin","123456",1);err!=nil{
		G.loger.Fatal(err)
	}
}

func InitMenuItemTable(){
	fileds:=map[string]string{
		"id":`int(11) UNSIGNED NOT NULL AUTO_INCREMENT`,
		"name":`varchar(32) NOT NULL`,
		"parent":`int(11) UNSIGNED NOT NULL DEFAULT 0`,
		"parentName":`varchar(32) NOT NULL DEFAULT ""`,
		"allow":`json NOT NULL`,
		"creater":`varchar(32) NOT NULL`,
		"ctime":`datetime(0)`,
	}
	if err:=toolbox.CheckAndFixTable(G.db,"menu_items",fileds,nil);err!=nil{
		G.loger.Fatal(err)
	}
}

func InitViewerItemTable(){
	fileds:=map[string]string{
		"id":`int(11) UNSIGNED NOT NULL AUTO_INCREMENT`,
		"name":`varchar(32) NOT NULL`,
		"parent":`int(11) UNSIGNED NOT NULL DEFAULT 0`,
		"parentName":`varchar(32) NOT NULL DEFAULT ""`,
		"allow":`json NOT NULL`,
		"creater":`varchar(32) NOT NULL`,
		"dbAddress":`varchar(32) NOT NULL`,
		"dbPassword":`varchar(32) NOT NULL DEFAULT ""`,
		"dbName":`varchar(32) NOT NULL`,
		"dbTable":`varchar(32) NOT NULL`,
		"dbUser":`varchar(32) NOT NULL`,
		"fields":`json NOT NULL`,
		"ctime":`datetime(0)`,
	}
	if err:=toolbox.CheckAndFixTable(G.db,"viewer_items",fileds,nil);err!=nil{
		G.loger.Fatal(err)
	}
	_,err:=G.db.Exec("insert ignore into viewer_items set `id`=?,`name`=?,`parent`=?,`parentName`=?,`allow`=?,`creater`=?,`dbAddress`=?,`dbPassword`=?," +
		"`dbName`=?,`dbTable`=?,`dbUser`=?,`fields`=?,`ctime`=?",1,"菜单管理",0,"","[]","admin",G.cfg.DB.Address,G.cfg.DB.PassWord,G.cfg.DB.DataBase,"menu_items",
		G.cfg.DB.User,`{"map": {"菜单ID": "id", "创建者": "creater", "菜单名": "name", "访问者": "allow", "父节点ID": "parent", "创建时间": "ctime", "父节点名": "parentName"}, "extra": {"菜单ID": "auto_increment", "创建者": "", "菜单名": "", "访问者": "", "父节点ID": "", "创建时间": "", "父节点名": ""}, "fields": [{"name": "菜单ID", "field": "id", "order": true}, {"name": "菜单名", "field": "name", "order": false}, {"name": "父节点ID", "field": "parent", "order": true}, {"name": "父节点名", "field": "parentName", "order": false}, {"name": "创建者", "field": "creater", "order": false}, {"name": "访问者", "field": "allow", "order": false}, {"name": "创建时间", "field": "ctime", "order": false}], "script": {"菜单ID": "", "创建者": "", "菜单名": "", "访问者": "ewppbjpmdW5jdGlvbihkYXRhKXt2YXIgYWxsb3c9SlNPTi5wYXJzZShkYXRhKTt2YXIgcmU9Jyc7Zm9yKHZhciBpPTA7aTxhbGxvdy5sZW5ndGg7aSsrKXtyZT1yZSthbGxvd1tpXSsnLCd9O3JlPXJlLnN1YnN0cigwLHJlLmxlbmd0aC0xKTtyZXR1cm4gcmV9LApvdXQ6ZnVuY3Rpb24oZGF0YSl7dmFyIGFsbG93PWRhdGEuc3BsaXQoJywnKTtyZXR1cm4gSlNPTi5zdHJpbmdpZnkoYWxsb3cpO30KfQ==", "父节点ID": "", "创建时间": "", "父节点名": ""}, "default": {"菜单ID": null, "创建者": null, "菜单名": null, "访问者": null, "父节点ID": "0", "创建时间": null, "父节点名": ""}, "nullable": {"菜单ID": "NO", "创建者": "NO", "菜单名": "NO", "访问者": "NO", "父节点ID": "NO", "创建时间": "NO", "父节点名": "NO"}, "primaryKey": "id", "primaryName": "菜单ID"}`,toolbox.GetTimeSecond())
	if err != nil{
		G.loger.Fatal(err)
	}
	_,err=G.db.Exec("insert ignore into viewer_items set `id`=?,`name`=?,`parent`=?,`parentName`=?,`allow`=?,`creater`=?,`dbAddress`=?,`dbPassword`=?," +
		"`dbName`=?,`dbTable`=?,`dbUser`=?,`fields`=?,`ctime`=?",2,"视图管理",0,"","[]","admin",G.cfg.DB.Address,G.cfg.DB.PassWord,G.cfg.DB.DataBase,"viewer_items",
		G.cfg.DB.User,`{"map": {"表名": "dbTable", "视图ID": "id", "创建者": "creater", "视图名": "name", "访问者": "allow", "父节点ID": "parent", "创建时间": "ctime", "字段信息": "fields", "数据库名": "dbName", "父节点名": "parentName", "数据库地址": "dbAddress", "数据库密码": "dbPassword", "数据库用户": "dbUser"}, "extra": {"表名": "", "视图ID": "auto_increment", "创建者": "", "视图名": "", "访问者": "", "父节点ID": "", "创建时间": "", "字段信息": "", "数据库名": "", "父节点名": "", "数据库地址": "", "数据库密码": "", "数据库用户": ""}, "fields": [{"name": "视图ID", "field": "id", "order": true}, {"name": "视图名", "field": "name", "order": false}, {"name": "父节点ID", "field": "parent", "order": true}, {"name": "父节点名", "field": "parentName", "order": false}, {"name": "创建者", "field": "creater", "order": false}, {"name": "访问者", "field": "allow", "order": false}, {"name": "创建时间", "field": "ctime", "order": true}, {"name": "数据库地址", "field": "dbAddress", "order": false}, {"name": "数据库密码", "field": "dbPassword", "order": false}, {"name": "数据库名", "field": "dbName", "order": false}, {"name": "表名", "field": "dbTable", "order": false}, {"name": "数据库用户", "field": "dbUser", "order": false}, {"name": "字段信息", "field": "fields", "order": false}], "script": {"表名": "", "视图ID": "", "创建者": "", "视图名": "", "访问者": "ewppbjpmdW5jdGlvbihkYXRhKXt2YXIgYWxsb3c9SlNPTi5wYXJzZShkYXRhKTt2YXIgcmU9Jyc7Zm9yKHZhciBpPTA7aTxhbGxvdy5sZW5ndGg7aSsrKXtyZT1yZSthbGxvd1tpXSsnLCd9O3JlPXJlLnN1YnN0cigwLHJlLmxlbmd0aC0xKTtyZXR1cm4gcmV9LApvdXQ6ZnVuY3Rpb24oZGF0YSl7dmFyIGFsbG93PWRhdGEuc3BsaXQoJywnKTtyZXR1cm4gSlNPTi5zdHJpbmdpZnkoYWxsb3cpfQp9", "父节点ID": "", "创建时间": "", "字段信息": "", "数据库名": "", "父节点名": "", "数据库地址": "", "数据库密码": "", "数据库用户": ""}, "default": {"表名": null, "视图ID": null, "创建者": null, "视图名": null, "访问者": null, "父节点ID": "0", "创建时间": null, "字段信息": null, "数据库名": null, "父节点名": "", "数据库地址": null, "数据库密码": "", "数据库用户": null}, "nullable": {"表名": "NO", "视图ID": "NO", "创建者": "NO", "视图名": "NO", "访问者": "NO", "父节点ID": "NO", "创建时间": "NO", "字段信息": "NO", "数据库名": "NO", "父节点名": "NO", "数据库地址": "NO", "数据库密码": "NO", "数据库用户": "NO"}, "primaryKey": "id", "primaryName": "视图ID"}`,toolbox.GetTimeSecond())
	if err != nil{
		G.loger.Fatal(err)
	}
	_,err=G.db.Exec("insert ignore into viewer_items set `id`=?,`name`=?,`parent`=?,`parentName`=?,`allow`=?,`creater`=?,`dbAddress`=?,`dbPassword`=?," +
		"`dbName`=?,`dbTable`=?,`dbUser`=?,`fields`=?,`ctime`=?",3,"用户管理",0,"","[]","admin",G.cfg.DB.Address,G.cfg.DB.PassWord,G.cfg.DB.DataBase,"users",
		G.cfg.DB.User,`{"map": {"ID": "id", "用户名": "name", "用户组": "group", "用户密码": "password", "登录时间": "ltime"}, "extra": {"ID": "auto_increment", "用户名": "", "用户组": "", "用户密码": "", "登录时间": ""}, "fields": [{"name": "ID", "field": "id", "order": true}, {"name": "用户名", "field": "name", "order": false}, {"name": "用户密码", "field": "password", "order": false}, {"name": "用户组", "field": "group", "order": false}, {"name": "登录时间", "field": "ltime", "order": true}], "script": {"ID": "", "用户名": "", "用户组": "", "用户密码": "", "登录时间": ""}, "default": {"ID": null, "用户名": null, "用户组": null, "用户密码": null, "登录时间": null}, "nullable": {"ID": "NO", "用户名": "NO", "用户组": "NO", "用户密码": "NO", "登录时间": "YES"}, "primaryKey": "id", "primaryName": "ID"}`,toolbox.GetTimeSecond())
	if err != nil{
		G.loger.Fatal(err)
	}
}
package main

import (
	"database/sql"
	"encoding/json"
	"github.com/casbin/casbin"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/qianlidongfeng/toolbox"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func apitest(c echo.Context) error{
	return c.JSON(200,"")
}

func casbinAuth(c echo.Context) error{
	ce:=casbin.NewEnforcer("casbin_auth_model.conf", "casbin_auth_policy.csv")
	re:=ce.Enforce("alice","data1","read")
	_=re
	return c.String(http.StatusOK, "Welcome "+"!")
}

func createSession(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	realPasswd := G.DB.GetUserPasswordByName(username)
	// Throws unauthorized error
	if password != realPasswd ||realPasswd==""{
		return c.JSON(401,map[string]string{"msg":"用户名不存在或密码错误"})
	}

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)
	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = username
	claims["exp"] = time.Now().Add(time.Hour*12).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(G.cfg.Secret.LoginJwt))
	if err != nil {
		return err
	}
	cookie := new(http.Cookie)
	cookie.Name = "session"
	cookie.Value = t
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.HttpOnly=true
	c.SetCookie(cookie)
	G.DB.UpdateUserLoginTime(username)
	return c.JSON(200,map[string]string{"location":"/system"})
}

func deleteSession(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "session"
	cookie.Value = ""
	cookie.Expires = time.Now().Add(-1)
	cookie.HttpOnly=true
	c.SetCookie(cookie)
	return c.JSON(200,map[string]string{"location":"/login"})
}

func updatePassword(c echo.Context) error{
	old:=c.FormValue("old")
	new:=c.FormValue("new")
	renew:=c.FormValue("renew")
	if new != renew{
		return c.JSON(401,map[string]string{"msg":"两次输入的密码不一样"})
	}
	if len(new)<6{
		return c.JSON(401,map[string]string{"msg":"密码不能小于6位"})
	}
	if len(new)> 20{
		return c.JSON(401,map[string]string{"msg":"密码不能超过20位"})
	}
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	username := claims["name"].(string)
	password:=G.DB.GetUserPasswordByName(username)
	if password==""{
		return c.JSON(401,map[string]string{"msg":"服务器异常，请稍后再试"})
	}
	if old != password{
		return c.JSON(401,map[string]string{"msg":"原密码输入不正确"})
	}
	if old == new{
		return c.JSON(401,map[string]string{"msg":"新密码和旧密码不能相同"})
	}
	err:=G.DB.UpdateUserPasswordByName(username,new)
	if err != nil{
		return c.JSON(401,map[string]string{"msg":"服务器异常，请稍后再试"})
	}
	return c.JSON(200,map[string]string{"msg":"密码修改成功"})
}


func addUser(c echo.Context) error{
	return nil
}

func updateUser(c echo.Context) error{
	return nil
}

func deleteUser(c echo.Context) error{
	return nil
}

func addProject(c echo.Context) error{
	return nil
}

func deleteProject(c echo.Context) error{
	return nil
}

func getProject(c echo.Context) error{
	return nil
}

func addMenuItem(c echo.Context) error{
	itemName:=c.FormValue("item-name")
	if itemName == ""{
		return c.JSON(401,map[string]string{"msg":"菜单项名称不能为空"})
	}
	parentID,err:=strconv.Atoi(c.FormValue("parent"))
	if err != nil{
		return c.JSON(401,map[string]string{"msg":"parent错误"})
	}
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	username := claims["name"].(string)
	group := G.DB.GetUserGroupByName(username)
	if group<1||group>20{
		return c.JSON(401,map[string]interface{}{"msg":"用户组错误"})
	}
	allowlist := strings.Split(strings.TrimRight(username+","+c.FormValue("allow"),","),",")
	if G.DB.CheckMenuItemExist(itemName,parentID,username){
		return c.JSON(401,map[string]string{"msg":"parent下已存在该名字的菜单项"})
	}
	if username!="admin" && parentID!=0 && G.DB.CheckUserPermissionOfMenuItem(parentID,username) == false{
		return c.JSON(401,map[string]string{"msg":"非法操作"})
	}
	parentName:=G.DB.GetMenuNameByID(parentID)
	mAllow:=make(map[string]struct{})
	for _,v := range(allowlist){
		mAllow[v]=struct{}{}
	}
	allowlist=[]string{}
	for k,_:=range mAllow{
		allowlist=append(allowlist,k)
	}
	allow,err:=json.Marshal(allowlist)
	if err != nil{
		return c.JSON(401,map[string]string{"msg":"其他允许访问的用户有误"})
	}
	err=G.DB.AddMenuItem(itemName,parentID,parentName,string(allow),username)
	if err != nil{
		return c.JSON(401,map[string]string{"msg":"服务器异常，请稍后再试"})
	}
	return c.JSON(200,map[string]string{"msg":"创建菜单项成功"})
}

func updateMenuItem(c echo.Context) error{
	itemName:=c.FormValue("item-name")
	if itemName == ""{
		return c.JSON(401,map[string]string{"msg":"菜单项名称不能为空"})
	}
	parentID,err:=strconv.Atoi(c.FormValue("parent"))
	if err != nil{
		return c.JSON(401,map[string]string{"msg":"parent错误"})
	}
	menuID,err:=strconv.Atoi(c.FormValue("menu-id"))
	if err != nil{
		return c.JSON(401,map[string]string{"msg":"viewer-id错误"})
	}
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	username := claims["name"].(string)
	group := G.DB.GetUserGroupByName(username)
	if group<1||group>20{
		return c.JSON(401,map[string]interface{}{"msg":"用户组错误"})
	}
	allowlist := strings.Split(strings.TrimRight(c.FormValue("allow"),","),",")
	parent:=G.DB.GetMenuItemParentID(menuID)
	if parent!= parentID&&G.DB.CheckMenuItemExist(itemName,parentID,username){
		return c.JSON(401,map[string]string{"msg":"parent下已存在该名字的菜单项"})
	}
	if username!="admin" && G.DB.CheckUserPermissionOfMenuItem(parentID,username) == false{
		return c.JSON(401,map[string]string{"msg":"非法操作"})
	}
	if username!="admin" && G.DB.CheckUserPermissionOfMenuItem(menuID,username) == false{
		return c.JSON(401,map[string]string{"msg":"非法操作"})
	}
	parentName:=G.DB.GetMenuNameByID(parentID)
	mAllow:=make(map[string]struct{})
	for _,v := range(allowlist){
		mAllow[v]=struct{}{}
	}
	allowlist=[]string{}
	for k,_:=range mAllow{
		allowlist=append(allowlist,k)
	}
	allow,err:=json.Marshal(allowlist)
	if err != nil{
		return c.JSON(401,map[string]string{"msg":"其他允许访问的用户有误"})
	}
	err=G.DB.UpdateMenuItem(menuID,itemName,parentID,parentName,string(allow))
	if err != nil{
		return c.JSON(401,map[string]string{"msg":"服务器异常，请稍后再试"})
	}
	return c.JSON(200,map[string]string{"msg":"编辑菜单项成功"})
}

func getMenuItem(c echo.Context) error{
	parent,err:=strconv.Atoi(c.FormValue("parent"));
	if err != nil{
		return c.JSON(401,map[string]string{"msg":"parentID错误"})
	}
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	username := claims["name"].(string)
	menuItems,err:=G.DB.GetMenuItem(parent,username)
	if err !=nil{
		return c.JSON(401,map[string]string{"msg":err.Error()})
	}
	return c.JSON(200,map[string]interface{}{"msg":"success","items":menuItems})
}


func addViewerItem(c echo.Context) error{
	itemName:=c.FormValue("item-name")
	if itemName == ""{
		return c.JSON(401,map[string]string{"msg":"视图项名称不能为空"})
	}
	parentID,err:=strconv.Atoi(c.FormValue("parent"))
	if err != nil{
		return c.JSON(401,map[string]string{"msg":"parent错误"})
	}
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	username := claims["name"].(string)
	group := G.DB.GetUserGroupByName(username)
	if group<1||group>20{
		return c.JSON(401,map[string]interface{}{"msg":"用户组错误"})
	}
	if G.DB.CheckViewerItemExist(itemName,parentID,username){
		return c.JSON(401,map[string]string{"msg":"parent下已存在该名字的菜单项"})
	}
	if username!="admin"&&parentID!=0 && G.DB.CheckUserPermissionOfMenuItem(parentID,username) == false{
		return c.JSON(401,map[string]string{"msg":"非法操作"})
	}
	parentName:=c.FormValue("parentName")
	dbAddress:=c.FormValue("address")
	dbPassword:=c.FormValue("password")
	dbName:=c.FormValue("database")
	dbTable:=c.FormValue("table")
	dbUser:=c.FormValue("user")
	if dbAddress==""||dbName==""||dbTable==""||dbUser==""{
		return c.JSON(401,map[string]string{"msg":"数据库填写错误"})
	}
	fields:=c.FormValue("fields")
	mFields :=map[string]interface{}{}
	err=json.Unmarshal([]byte(fields),&mFields)
	if err != nil{
		return c.JSON(401,map[string]string{"msg":err.Error()})
	}
	if mFields["primaryKey"]==nil||mFields["primaryName"]==nil{
		return c.JSON(401,map[string]string{"msg":"缺少主键"})
	}
	if m,ok:=(mFields["fields"]).([]interface{}); ok{
		if len(m)<=0{
			return c.JSON(401,map[string]string{"msg":"字段对应错误"})
		}
	}else{
		return c.JSON(401,map[string]string{"msg":"字段对应错误"})
	}

	allowlist := strings.Split(strings.TrimRight(username+","+c.FormValue("allow"),","),",")
	mAllow:=make(map[string]struct{})
	for _,v := range(allowlist){
		mAllow[v]=struct{}{}
	}
	allowlist=[]string{}
	for k,_:=range mAllow{
		allowlist=append(allowlist,k)
	}
	allow,err:=json.Marshal(allowlist)
	if err != nil{
		return c.JSON(401,map[string]string{"msg":"其他允许访问的用户有误"})
	}
	err=G.DB.AddViewerItem(itemName,parentID,parentName,string(allow),username,dbAddress,dbPassword,dbName,dbTable,dbUser,fields)
	if err != nil{
		return c.JSON(401,map[string]string{"msg":"服务器异常，请稍后再试"})
	}
	return c.JSON(200,map[string]string{"msg":"创建视图项成功"})
}

func updateViewerItem(c echo.Context) error{
	itemName:=c.FormValue("item-name")
	if itemName == ""{
		return c.JSON(401,map[string]string{"msg":"视图项名称不能为空"})
	}
	parentID,err:=strconv.Atoi(c.FormValue("parent"))
	if err != nil{
		return c.JSON(401,map[string]string{"msg":"parent错误"})
	}
	viewerID,err:=strconv.Atoi(c.FormValue("viewer-id"))
	if err != nil{
		return c.JSON(401,map[string]string{"msg":"viewer-id错误"})
	}
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	username := claims["name"].(string)
	group := G.DB.GetUserGroupByName(username)
	if group<1||group>20{
		return c.JSON(401,map[string]interface{}{"msg":"用户组错误"})
	}
	parent:=G.DB.GetViewerItemParentID(viewerID)
	if parentID != parent && G.DB.CheckViewerItemExist(itemName,parentID,username){
		return c.JSON(401,map[string]string{"msg":"parent下已存在该名字的视图项"})
	}
	if username!="admin"&&parentID!=0 && G.DB.CheckUserPermissionOfMenuItem(parentID,username) == false{
		return c.JSON(401,map[string]string{"msg":"非法操作"})
	}
	if username!="admin" && G.DB.CheckUserPermissionOfViewerItem(viewerID,username) == false{
		return c.JSON(401,map[string]string{"msg":"非法操作"})
	}
	parentName:=c.FormValue("parentName")
	dbAddress:=c.FormValue("address")
	dbPassword:=c.FormValue("password")
	dbName:=c.FormValue("database")
	dbTable:=c.FormValue("table")
	dbUser:=c.FormValue("user")
	if dbAddress==""||dbName==""||dbTable==""||dbUser==""{
		return c.JSON(401,map[string]string{"msg":"数据库填写错误"})
	}
	fields:=c.FormValue("fields")
	mFields :=map[string]interface{}{}
	err=json.Unmarshal([]byte(fields),&mFields)
	if err != nil{
		return c.JSON(401,map[string]string{"msg":err.Error()})
	}
	if mFields["primaryKey"]==nil||mFields["primaryName"]==nil{
		return c.JSON(401,map[string]string{"msg":"缺少主键"})
	}
	if m,ok:=(mFields["fields"]).([]interface{}); ok{
		if len(m)<=0{
			return c.JSON(401,map[string]string{"msg":"字段对应错误"})
		}
	}else{
		return c.JSON(401,map[string]string{"msg":"字段对应错误"})
	}

	allowlist := strings.Split(strings.TrimRight(c.FormValue("allow"),","),",")
	mAllow:=make(map[string]struct{})
	for _,v := range(allowlist){
		mAllow[v]=struct{}{}
	}
	allowlist=[]string{}
	for k,_:=range mAllow{
		allowlist=append(allowlist,k)
	}
	allow,err:=json.Marshal(allowlist)
	if err != nil{
		return c.JSON(401,map[string]string{"msg":"其他允许访问的用户有误"})
	}
	err=G.DB.UpdateViewerItem(viewerID,itemName,parentID,parentName,string(allow),username,dbAddress,dbPassword,dbName,dbTable,dbUser,fields)
	if err != nil{
		return c.JSON(401,map[string]string{"msg":"服务器异常，请稍后再试"})
	}
	return c.JSON(200,map[string]string{"msg":"编辑视图项成功"})
}

func getViewerItem(c echo.Context) error{
	parent,err:=strconv.Atoi(c.FormValue("parent"));
	if err != nil{
		return c.JSON(401,map[string]string{"msg":"parentID错误"})
	}
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	username := claims["name"].(string)
	viewerItems,err:=G.DB.GetViewerItem(parent,username)
	if err !=nil{
		return c.JSON(401,map[string]string{"msg":err.Error()})
	}
	return c.JSON(200,map[string]interface{}{"msg":"success","items":viewerItems})
}

func getDatabaseField(c echo.Context) error{
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	username := claims["name"].(string)
	group:=G.DB.GetUserGroupByName(username)
	if group==0 || group>20{
		return c.JSON(401,map[string]string{"msg":"没有权限"})
	}
	addr:=c.FormValue("address")
	if addr==""{
		return c.JSON(401,map[string]string{"msg":"数据库地址不能为空"})
	}
	password:=c.FormValue("password")
	if password==""{
		return c.JSON(401,map[string]string{"msg":"数据库密码不能为空"})
	}
	database:=c.FormValue("database")
	if database==""{
		return c.JSON(401,map[string]string{"msg":"数据库名不能为空"})
	}
	table:=c.FormValue("table")
	if table==""{
		return c.JSON(401,map[string]string{"msg":"表名不能为空"})
	}
	dbuser:=c.FormValue("user")
	if dbuser==""{
		return c.JSON(401,map[string]string{"msg":"用户名名不能为空"})
	}
	var DBConfig toolbox.MySqlConfig
	DBConfig.Address=addr
	DBConfig.PassWord=password
	DBConfig.DataBase=database
	DBConfig.User=dbuser
	db,err:=toolbox.InitMysql(DBConfig)
	if err != nil{
		return c.JSON(401,map[string]string{"msg":err.Error()})
	}
	defer db.Close()
	if !isTableExist(db,database,table){
		return c.JSON(401,map[string]string{"msg":"该表不存在"})
	}
	fields,err:=GetTableFields(db,database,table)
	if err != nil{
		return c.JSON(401,map[string]string{"msg":err.Error()})
	}
	return c.JSON(200,fields)
}

func getSysViewerData(c echo.Context) error{
	vid:=c.FormValue("viewer-id")
	viewerID,err:=strconv.Atoi(vid)
	if err != nil{
		return c.JSON(401,map[string]string{"msg":"viewer id错误"})
	}
	limit:=c.FormValue("limit")
	limit=strings.TrimRight(limit," ")
	if limit==""{
		limit="20"
	}
	nLimit,err:=strconv.Atoi(limit)
	if err != nil||nLimit<=0{
		return c.JSON(401,map[string]string{"msg":"limit 错误"})
	}
	filter:=c.FormValue("filter")
	page:=c.FormValue("page")
	if page==""{
		page="1"
	}
	nPage,err:=strconv.Atoi(page)
	if err != nil||nPage<=0{
		return c.JSON(401,map[string]string{"msg":"page 错误"})
	}
	order:=c.FormValue("order")
	orderType:=c.FormValue("order-type")
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	username := claims["name"].(string)
	group:=G.DB.GetUserGroupByName(username)
	if group>10&&group<20&&viewerID==3||group>20{
		return c.JSON(401,map[string]string{"msg":"没有权限"})
	}
	dbConfig,fieldsInfo:=G.DB.GetSysViewerFields(viewerID,username)
	if fieldsInfo==""{
		return c.JSON(401,map[string]string{"msg":"没有相关数据"})
	}
	info :=make(map[string]interface{})
	err=json.Unmarshal([]byte(fieldsInfo),&info)
	if err != nil{
		return c.JSON(401,map[string]string{"msg":"解析字段错误"})
	}
	primaryKey,ok:=info["primaryKey"].(string)
	if !ok{
		return c.JSON(401,map[string]string{"msg":"解析字段错误"})
	}
	fields,ok:= info["fields"].([]interface{})
	if !ok{
		return c.JSON(401,map[string]string{"msg":"解析字段错误"})
	}
	sqlSentence:="select "
	for _,v:=range fields{
		fieldInfo,ok:=v.(map[string]interface{})
		if !ok{
			return c.JSON(401,map[string]string{"msg":"解析字段错误"})
		}
		name,ok:= fieldInfo["name"].(string)
		if !ok{
			return c.JSON(401,map[string]string{"msg":"解析字段错误"})
		}
		field,ok:=fieldInfo["field"].(string)
		if !ok{
			return c.JSON(401,map[string]string{"msg":"解析字段错误"})
		}
		sqlSentence=sqlSentence+"`"+field+"`"+" as "+"`"+name+"`"+","
	}
	if order==""{
		order=primaryKey
	}
	oderSentence:=" order by "+order +" "+orderType
	limitSentence:=" limit "+strconv.Itoa((nPage-1)*nLimit)+","+limit
	where:=""
	mfilter:=[]map[string]string{}
	err=json.Unmarshal([]byte(filter),&mfilter)
	if err != nil{
		return c.JSON(401,map[string]string{"msg":"filter错误"})
	}
	where =" where json_contains(allow,json_array('"+username+"')) "
	if len(mfilter)>0{
		where= where + " and "
		mmap,ok:=info["map"].(map[string]interface{})
		if !ok{
			return c.JSON(401,map[string]string{"msg":"filter错误"})
		}
		for _,v:=range mfilter{
			name:=v["name"]
			value:=v["value"]
			relation:=v["relation"]
			field,ok:=mmap[name].(string)
			if !ok{
				return c.JSON(401,map[string]string{"msg":"filter错误"})
			}
			where = where+ "`"+field+"`"+relation+"'"+value+"'"+" and "
		}
		where=strings.TrimRight(where," and ")
	}
	sqlSentence=strings.TrimRight(sqlSentence,",")+" from "+dbConfig.Table+where+ oderSentence + limitSentence
	dbCacheKey:=toolbox.MD5([]byte(dbConfig.Address+";"+dbConfig.DataBase+";"+dbConfig.User))
	var db *sql.DB
	if value,ok:=G.dbCache.Load(dbCacheKey);!ok{
		db,err=toolbox.InitMysql(dbConfig)
		if err !=nil{
			return c.JSON(401,map[string]string{"msg":"连接数据库失败"})
		}
		G.dbCache.Store(dbCacheKey,db)
	}else{
		db=value.(*sql.DB)
	}
	result,err:=toolbox.SelectArrayFromMysql(db,sqlSentence)
	if err !=nil{
		return c.JSON(401,map[string]string{"msg":"解析数据错误"})
	}
	if err != nil{
		return c.JSON(401,map[string]string{"msg":err.Error()})
	}
	var rowCount int
	db.QueryRow("select count(*) from "+dbConfig.Table + where).Scan(&rowCount)
	pageTotal:=rowCount/nLimit
	if rowCount%nLimit!=0{
		pageTotal=pageTotal+1
	}
	pages,err:=toolbox.UI.GetPagesBarInfo(nPage,pageTotal,10)
	return c.JSON(200,map[string]interface{}{"info":info,"data":result,"pages":pages})
}

func getViewerData(c echo.Context) error{
	vid:=c.FormValue("viewer-id")
	viewerID,err:=strconv.Atoi(vid)
	if err != nil{
		return c.JSON(401,map[string]string{"msg":"viewer id错误"})
	}
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	username := claims["name"].(string)
	if viewerID>0&&viewerID<4&&username!="admin"{
		return getSysViewerData(c)
	}
	limit:=c.FormValue("limit")
	limit=strings.TrimRight(limit," ")
	if limit==""{
		limit="20"
	}
	nLimit,err:=strconv.Atoi(limit)
	if err != nil||nLimit<=0{
		return c.JSON(401,map[string]string{"msg":"limit 错误"})
	}
	filter:=c.FormValue("filter")
	page:=c.FormValue("page")
	if page==""{
		page="1"
	}
	nPage,err:=strconv.Atoi(page)
	if err != nil||nPage<=0{
		return c.JSON(401,map[string]string{"msg":"page 错误"})
	}
	order:=c.FormValue("order")
	orderType:=c.FormValue("order-type")
	dbConfig,fieldsInfo:=G.DB.GetViewerFields(viewerID,username)
	if fieldsInfo==""{
		return c.JSON(401,map[string]string{"msg":"没有相关数据"})
	}
	info :=make(map[string]interface{})
	err=json.Unmarshal([]byte(fieldsInfo),&info)
	if err != nil{
		return c.JSON(401,map[string]string{"msg":"解析字段错误"})
	}
	primaryKey,ok:=info["primaryKey"].(string)
	if !ok{
		return c.JSON(401,map[string]string{"msg":"解析字段错误"})
	}
	fields,ok:= info["fields"].([]interface{})
	if !ok{
		return c.JSON(401,map[string]string{"msg":"解析字段错误"})
	}
	sqlSentence:="select "
	for _,v:=range fields{
		fieldInfo,ok:=v.(map[string]interface{})
		if !ok{
			return c.JSON(401,map[string]string{"msg":"解析字段错误"})
		}
		name,ok:= fieldInfo["name"].(string)
		if !ok{
			return c.JSON(401,map[string]string{"msg":"解析字段错误"})
		}
		field,ok:=fieldInfo["field"].(string)
		if !ok{
			return c.JSON(401,map[string]string{"msg":"解析字段错误"})
		}
		sqlSentence=sqlSentence+"`"+field+"`"+" as "+"`"+name+"`"+","
	}
	if order==""{
		order=primaryKey
	}
	oderSentence:=" order by "+order +" "+orderType
	limitSentence:=" limit "+strconv.Itoa((nPage-1)*nLimit)+","+limit
	where:=""
	mfilter:=[]map[string]string{}
	err=json.Unmarshal([]byte(filter),&mfilter)
	if err != nil{
		return c.JSON(401,map[string]string{"msg":"filter错误"})
	}
	if len(mfilter)>0{
		where= " where "
		mmap,ok:=info["map"].(map[string]interface{})
		if !ok{
			return c.JSON(401,map[string]string{"msg":"filter错误"})
		}
		for _,v:=range mfilter{
			name:=v["name"]
			value:=v["value"]
			relation:=v["relation"]
			field,ok:=mmap[name].(string)
			if !ok{
				return c.JSON(401,map[string]string{"msg":"filter错误"})
			}
			where = where+ "`"+field+"`"+relation+"'"+value+"'"+" and "
		}
		where=strings.TrimRight(where," and ")
	}

	sqlSentence=strings.TrimRight(sqlSentence,",")+" from "+dbConfig.Table+where+ oderSentence + limitSentence
	dbCacheKey:=toolbox.MD5([]byte(dbConfig.Address+";"+dbConfig.DataBase+";"+dbConfig.User))
	var db *sql.DB
	if value,ok:=G.dbCache.Load(dbCacheKey);!ok{
		db,err=toolbox.InitMysql(dbConfig)
		if err !=nil{
			return c.JSON(401,map[string]string{"msg":"连接数据库失败"})
		}
		G.dbCache.Store(dbCacheKey,db)
	}else{
		db=value.(*sql.DB)
	}
	result,err:=toolbox.SelectArrayFromMysql(db,sqlSentence)
	if err !=nil{
		return c.JSON(401,map[string]string{"msg":"解析数据错误"})
	}
	if err != nil{
		return c.JSON(401,map[string]string{"msg":err.Error()})
	}
	var rowCount int
	db.QueryRow("select count(*) from "+dbConfig.Table + where).Scan(&rowCount)
	pageTotal:=rowCount/nLimit
	if rowCount%nLimit!=0{
		pageTotal=pageTotal+1
	}
	pages,err:=toolbox.UI.GetPagesBarInfo(nPage,pageTotal,10)
	return c.JSON(200,map[string]interface{}{"info":info,"data":result,"pages":pages})
}

func deleteSysViewerData(c echo.Context) error{
	params:=c.FormValue("param")
	vid:=c.FormValue("viewer-id")
	viewerID,err:=strconv.Atoi(vid)
	if err != nil{
		return c.JSON(401,map[string]string{"msg":"viewer id错误"})
	}
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	username := claims["name"].(string)
	group:=G.DB.GetUserGroupByName(username)
	if group>10 && group<20&&viewerID==3||group>20{
		return c.JSON(401,map[string]string{"msg":"没有权限"})
	}
	dbConfig,fieldsInfo:=G.DB.GetSysViewerFields(viewerID,username)
	if fieldsInfo==""{
		return c.JSON(401,map[string]string{"msg":"获取字段错误"})
	}
	info:=map[string]interface{}{}
	err=json.Unmarshal([]byte(fieldsInfo),&info)
	if err != nil{
		return c.JSON(401,map[string]string{"msg":"解析字段错误"})
	}
	key,ok:=info["primaryKey"].(string)
	if !ok{
		return c.JSON(401,map[string]string{"msg":"解析字段错误"})
	}
	values:=[]string{}
	err=json.Unmarshal([]byte(params),&values)
	if err!=nil{
		return c.JSON(401,map[string]string{"msg":"解析参数错误"})
	}
	in:="("
	for _,v :=range values{
		in=in+"'"+v+"'"+","
	}
	in = strings.TrimRight(in,",")
	in=in+")"
	dbCacheKey:=toolbox.MD5([]byte(dbConfig.Address+";"+dbConfig.DataBase+";"+dbConfig.User))
	var db *sql.DB
	if value,ok:=G.dbCache.Load(dbCacheKey);!ok{
		db,err=toolbox.InitMysql(dbConfig)
		if err !=nil{
			return c.JSON(401,map[string]string{"msg":"连接数据库失败"})
		}
		G.dbCache.Store(dbCacheKey,db)
	}else{
		db=value.(*sql.DB)
	}
	where:=" where json_contains(allow,json_array('"+username+"'))"
	sqlSentence:="delete from "+dbConfig.Table +where+" and "+key+" in "+in
	_,err=db.Exec(sqlSentence)
	if err!=nil{
		return c.JSON(401,map[string]string{"msg":"写入数据库失败"})
	}
	return c.JSON(200,map[string]string{"msg":"success"})
}

func deleteViewerData(c echo.Context) error{
	params:=c.FormValue("param")
	vid:=c.FormValue("viewer-id")
	viewerID,err:=strconv.Atoi(vid)
	if err != nil{
		return c.JSON(401,map[string]string{"msg":"viewer id错误"})
	}
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	username := claims["name"].(string)
	if viewerID>0 && viewerID<4 && username!= "admin"{
		return deleteSysViewerData(c)
	}
	dbConfig,fieldsInfo:=G.DB.GetViewerFields(viewerID,username)
	if fieldsInfo==""{
		return c.JSON(401,map[string]string{"msg":"获取字段错误"})
	}
	info:=map[string]interface{}{}
	err=json.Unmarshal([]byte(fieldsInfo),&info)
	if err != nil{
		return c.JSON(401,map[string]string{"msg":"解析字段错误"})
	}
	key,ok:=info["primaryKey"].(string)
	if !ok{
		return c.JSON(401,map[string]string{"msg":"解析字段错误"})
	}
	values:=[]string{}
	err=json.Unmarshal([]byte(params),&values)
	if err!=nil{
		return c.JSON(401,map[string]string{"msg":"解析参数错误"})
	}
	in:="("
	for _,v :=range values{
		in=in+"'"+v+"'"+","
	}
	in = strings.TrimRight(in,",")
	in=in+")"
	dbCacheKey:=toolbox.MD5([]byte(dbConfig.Address+";"+dbConfig.DataBase+";"+dbConfig.User))
	var db *sql.DB
	if value,ok:=G.dbCache.Load(dbCacheKey);!ok{
		db,err=toolbox.InitMysql(dbConfig)
		if err !=nil{
			return c.JSON(401,map[string]string{"msg":"连接数据库失败"})
		}
		G.dbCache.Store(dbCacheKey,db)
	}else{
		db=value.(*sql.DB)
	}

	sqlSentence:="delete from "+dbConfig.Table + " where "+key+" in "+in
	_,err=db.Exec(sqlSentence)
	if err!=nil{
		return c.JSON(401,map[string]string{"msg":"写入数据库失败"})
	}
	return c.JSON(200,map[string]string{"msg":"success"})
}

func updateSysViewerData(c echo.Context) error{
	vid:=c.FormValue("viewer-id")
	viewerID,err:=strconv.Atoi(vid)
	if err != nil{
		return c.JSON(401,map[string]string{"msg":"viewer id错误"})
	}
	primary:=c.FormValue("primary")
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	username := claims["name"].(string)
	group:=G.DB.GetUserGroupByName(username)
	if group>10&&group<20&&viewerID==3||group>20{
		return c.JSON(401,map[string]string{"msg":"没有权限"})
	}
	data:=c.FormValue("data")
	dbConfig,fieldsInfo:=G.DB.GetSysViewerFields(viewerID,username)
	if fieldsInfo==""{
		return c.JSON(401,map[string]string{"msg":"获取字段错误"})
	}
	info:=map[string]interface{}{}
	err=json.Unmarshal([]byte(fieldsInfo),&info)
	if err != nil{
		return c.JSON(401,map[string]string{"msg":"解析字段错误"})
	}
	primaryKey,ok:=info["primaryKey"].(string)
	if !ok{
		return c.JSON(401,map[string]string{"msg":"解析字段错误"})
	}
	m,ok:=info["map"].(map[string]interface{})
	if !ok{
		return c.JSON(401,map[string]string{"msg":"解析字段错误"})
	}
	mFields:=make(map[string]string)
	err=json.Unmarshal([]byte(data),&mFields)
	if err != nil{
		return c.JSON(401,map[string]string{"msg":"解析字段错误"})
	}
	set:=""
	for k,v:=range mFields{
		key,ok:=m[k].(string)
		if !ok{
			return c.JSON(401,map[string]string{"msg":"解析字段错误"})
		}
		set=set+"`"+key+"`"+"="+"'"+v+"'"+","
	}
	set = strings.TrimRight(set,",")
	set = set+" "
	where:=" where "+primaryKey+"='"+primary+"'"
	where = where + " and json_contains(allow,json_array('"+username+"'))"
	sqlSentence:="update "+dbConfig.Table +" set "+set+where
	dbCacheKey:=toolbox.MD5([]byte(dbConfig.Address+";"+dbConfig.DataBase+";"+dbConfig.User))
	var db *sql.DB
	if value,ok:=G.dbCache.Load(dbCacheKey);!ok{
		db,err=toolbox.InitMysql(dbConfig)
		if err !=nil{
			return c.JSON(401,map[string]string{"msg":"连接数据库失败"})
		}
		G.dbCache.Store(dbCacheKey,db)
	}else{
		db=value.(*sql.DB)
	}
	_,err=db.Exec(sqlSentence)
	if err!=nil{
		return c.JSON(401,map[string]string{"msg":"写入数据库失败"})
	}
	return c.JSON(200,map[string]string{"msg":"操作成功"})
}

func updateViewerData(c echo.Context) error{
	vid:=c.FormValue("viewer-id")
	viewerID,err:=strconv.Atoi(vid)
	if err != nil{
		return c.JSON(401,map[string]string{"msg":"viewer id错误"})
	}
	primary:=c.FormValue("primary")
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	username := claims["name"].(string)
	if viewerID>0&&viewerID<4&&username != "admin"{
		return updateSysViewerData(c)
	}
	data:=c.FormValue("data")
	dbConfig,fieldsInfo:=G.DB.GetViewerFields(viewerID,username)
	if fieldsInfo==""{
		return c.JSON(401,map[string]string{"msg":"获取字段错误"})
	}
	info:=map[string]interface{}{}
	err=json.Unmarshal([]byte(fieldsInfo),&info)
	if err != nil{
		return c.JSON(401,map[string]string{"msg":"解析字段错误"})
	}
	primaryKey,ok:=info["primaryKey"].(string)
	if !ok{
		return c.JSON(401,map[string]string{"msg":"解析字段错误"})
	}
	m,ok:=info["map"].(map[string]interface{})
	if !ok{
		return c.JSON(401,map[string]string{"msg":"解析字段错误"})
	}
	mFields:=make(map[string]string)
	err=json.Unmarshal([]byte(data),&mFields)
	if err != nil{
		return c.JSON(401,map[string]string{"msg":"解析字段错误"})
	}
	set:=""
	for k,v:=range mFields{
		key,ok:=m[k].(string)
		if !ok{
			return c.JSON(401,map[string]string{"msg":"解析字段错误"})
		}
		set=set+"`"+key+"`"+"="+"'"+v+"'"+","
	}
	set = strings.TrimRight(set,",")
	set = set+" "
	where:=" where "+primaryKey+"='"+primary+"'"

	sqlSentence:="update "+dbConfig.Table +" set "+set+where
	dbCacheKey:=toolbox.MD5([]byte(dbConfig.Address+";"+dbConfig.DataBase+";"+dbConfig.User))
	var db *sql.DB
	if value,ok:=G.dbCache.Load(dbCacheKey);!ok{
		db,err=toolbox.InitMysql(dbConfig)
		if err !=nil{
			return c.JSON(401,map[string]string{"msg":"连接数据库失败"})
		}
		G.dbCache.Store(dbCacheKey,db)
	}else{
		db=value.(*sql.DB)
	}
	_,err=db.Exec(sqlSentence)
	if err!=nil{
		return c.JSON(401,map[string]string{"msg":"写入数据库失败"})
	}
	return c.JSON(200,map[string]string{"msg":"操作成功"})
}

func addViewerData(c echo.Context) error{
	fields:=map[string]string{}
	err:=json.Unmarshal([]byte(c.FormValue("fields")),&fields)
	if err != nil{
		return c.JSON(401,map[string]string{"msg":err.Error()})
	}
	vid:=c.FormValue("viewer-id")
	viewerID,err:=strconv.Atoi(vid)
	if err != nil{
		return c.JSON(401,map[string]string{"msg":"viewer id错误"})
	}
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	username := claims["name"].(string)
	dbConfig,fieldsInfo:=G.DB.GetViewerFields(viewerID,username)
	if fieldsInfo==""{
		return c.JSON(401,map[string]string{"msg":"获取字段错误"})
	}
	info:=map[string]interface{}{}
	err=json.Unmarshal([]byte(fieldsInfo),&info)
	if err != nil{
		return c.JSON(401,map[string]string{"msg":"解析字段错误"})
	}
	mmap,ok:=info["map"].(map[string]interface{})
	if !ok{
		return c.JSON(401,map[string]string{"msg":"解析字段错误"})
	}
	set:=" set "
	for k,v:=range fields{
		key,ok:=mmap[k].(string)
		if !ok{
			return c.JSON(401,map[string]string{"msg":"解析字段错误"})
		}
		set=set+"`"+key+"`="+"'"+v+"'"+","
	}
	set=strings.TrimRight(set,",")
	sqlSentence:="insert into "+dbConfig.Table + set
	dbCacheKey:=toolbox.MD5([]byte(dbConfig.Address+";"+dbConfig.DataBase+";"+dbConfig.User))
	var db *sql.DB
	if value,ok:=G.dbCache.Load(dbCacheKey);!ok{
		db,err=toolbox.InitMysql(dbConfig)
		if err !=nil{
			return c.JSON(401,map[string]string{"msg":"连接数据库失败"})
		}
		G.dbCache.Store(dbCacheKey,db)
	}else{
		db=value.(*sql.DB)
	}
	_,err=db.Exec(sqlSentence)
	if err!=nil{
		return c.JSON(401,map[string]string{"msg":err.Error()})
	}
	return c.JSON(200,map[string]string{"msg":"操作成功"})
}
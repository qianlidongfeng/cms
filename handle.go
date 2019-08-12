package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)


func test(c echo.Context) error{
	return c.Render(200,"test.html",nil)
}

func login(c echo.Context) error{
	session,err:=c.Cookie("session")
	if err==nil{
		auth:=session.Value
		token,err:=jwt.Parse(auth, NewJwtParser([]byte(G.cfg.Secret.LoginJwt)))
		if err == nil && token.Valid{
			return c.Redirect(RedirectCode,"/system")
		}else{
			return c.Render(200,"login.html",nil)
		}
	}
	return c.Render(200,"login.html",nil)
}

func index(c echo.Context) error {
	return c.Redirect(RedirectCode,"/system")
}

func system(c echo.Context) error {
	data:=make(map[string]interface{})
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	username := claims["name"].(string)
	data["username"]=username
	group := G.DB.GetUserGroupByName(username)
	if group==0{
		return c.JSON(401,map[string]interface{}{"msg":"用户组错误"})
	}
	data["group"]=group
	menuItems,err:=G.DB.GetMenuItem(0,username)
	if err != nil{
		return c.JSON(401,map[string]interface{}{"msg":err.Error()})
	}
	data["menuItems"]=menuItems
	viewerItems,err:=G.DB.GetViewerItem(0,username)
	if err != nil{
		return c.JSON(401,map[string]interface{}{"msg":err.Error()})
	}
	data["viewerItems"]=viewerItems
	return c.Render(200,"system.html",data)
}


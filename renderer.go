package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/qianlidongfeng/execho"
	"io"
)

type Renderer struct{
	execho.IRenderer
}

func NewRenderer(pattern string,debug bool) execho.IRenderer{
	var iRenderer execho.IRenderer
	if debug{
		iRenderer = execho.NewDebugRenderer(pattern)
	}else{
		iRenderer = execho.NewRenderer()
		iRenderer.ParseGlob(pattern)
	}
	return iRenderer
}


type RenderMiddleWare struct{

}

func NewRenderMiddleWare() RenderMiddleWare{
	return RenderMiddleWare{}
}


func (this *RenderMiddleWare) AddUserNamefunc(next execho.RenderFunc) execho.RenderFunc{
	return func(w io.Writer, name string, data interface{}, c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		username := claims["name"].(string)
		if data==nil{
			data = map[string]interface{}{"username":username}
		}else{
			data.(map[string]interface{})["username"]=username
		}
		return next(w,name,data,c)
	}
}
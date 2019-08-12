package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/qianlidongfeng/execho"
	"github.com/qianlidongfeng/loger"
)

type loginJwtConfig struct{
	Jwtkey []byte
	TokenLookup string
}

type MiddleWare struct{
	JwtKey []byte
	loger loger.Loger
}

func NewMiddleWare() MiddleWare{
	return MiddleWare{
	}
}

func (this *MiddleWare) Recover()echo.MiddlewareFunc{
	exMiddleWare:=execho.NewMiddleWare()
	return exMiddleWare.Recover(this.loger)
}

func (this *MiddleWare) LoginJwtAuth() echo.MiddlewareFunc{
	c := middleware.DefaultJWTConfig
	c.SigningKey = this.JwtKey
	c.TokenLookup = "cookie:session"
	c.ContextKey="user"
	c.ErrorHandlerWithContext = func(e error, context echo.Context) error {
		return context.Redirect(RedirectCode,"/login")
	}
	return middleware.JWTWithConfig(c)
}


func (this *MiddleWare) NoCache() echo.MiddlewareFunc{
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("CacheControl","no-cache")
			return next(c)
		}
	}
}
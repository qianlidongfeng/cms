package main

import (
	"crypto/tls"
	"github.com/labstack/echo"
	"net/http"
	"time"
)

func InitEcho(e *echo.Echo){
	e.DisableHTTP2=true
	if G.cfg.Renderer.Enable{
		e.Renderer=G.renderer
	}
	//recover
	e.Use(G.middleware.Recover())
	//nocache
	//e.Use(G.middleware.NoCache())
	e.GET("/", index)
	e.GET("/test",test)
	e.GET("/apitest",apitest)
	e.GET("/login", login)
	e.POST("/session",createSession)
	e.DELETE("/session",deleteSession)
	e.Static("/static", "static")
	api:=e.Group("api")
	api.Use(G.middleware.LoginJwtAuth())
	api.PUT("/password",updatePassword)
	api.POST("/user",addUser)
	api.PUT("/user",updateUser)
	api.DELETE("/user",deleteUser)
	api.POST("/project",addProject)
	api.DELETE("/project",deleteProject)
	api.GET("/project",getProject)
	api.GET("/menuitem",getMenuItem)
	api.POST("/menuitem",addMenuItem)
	api.PUT("/menuitem",updateMenuItem)
	api.POST("/vieweritem",addViewerItem)
	api.PUT("/vieweritem",updateViewerItem)
	api.GET("/vieweritem",getViewerItem)
	api.GET("/databaseField",getDatabaseField)
	api.GET("/viewerdata",getViewerData)
	api.DELETE("/viewerdata",deleteViewerData)
	api.PUT("/viewerdata",updateViewerData)
	api.POST("/viewerdata",addViewerData)
	r := e.Group("/system")
	//loginAuth
	r.Use(G.middleware.LoginJwtAuth())
	r.GET("", system)
}

func main() {
	s := &http.Server{
		Addr:  G.cfg.Httpserver.Address,
		ReadTimeout:  G.cfg.Httpserver.ReadTimeOut*time.Millisecond,
		WriteTimeout: G.cfg.Httpserver.WriteTimeOut*time.Millisecond,
	}
	if G.cfg.Httpserver.Https{
		crt, err := tls.LoadX509KeyPair(G.cfg.Httpserver.CertFile, G.cfg.Httpserver.KeyFile)
		if err != nil {
			panic(err)
			G.loger.Fatal(err)
			return
		}
		s.TLSConfig= &tls.Config{
			MinVersion:               tls.VersionTLS12,
			CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
			Certificates: []tls.Certificate{crt},
		}
		s.TLSNextProto=make(map[string]func(*http.Server, *tls.Conn, http.Handler))
	}
	e := echo.New()
	InitEcho(e)
	G.loger.Fatal(e.StartServer(s))
}
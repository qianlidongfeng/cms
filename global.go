package main

import (
	"database/sql"
	"github.com/qianlidongfeng/execho"
	"github.com/qianlidongfeng/loger"
	"github.com/qianlidongfeng/toolbox"
	"os"
	"runtime"
	"strings"
	"sync"
)

var (
	G Global
	RedirectCode = 302
)

type Global struct{
	cfg Config
	loger loger.Loger
	db *sql.DB
	DB *Database
	slash string
	middleware MiddleWare
	renderMiddleWare RenderMiddleWare
	renderer execho.IRenderer
	dbCache sync.Map
}

func (this *Global) Init(){
	configFile,err:=toolbox.GetConfigFile()
	if err != nil{
		panic(err)
	}
	err=toolbox.LoadConfig(configFile,&this.cfg)
	if err != nil{
		panic(err)
	}
	if this.cfg.Debug == false{
		toolbox.RediRectOutPutToLog()
	}
	this.loger,err=loger.NewLoger(this.cfg.Log)
	if err != nil{
		panic(err)
	}
	this.db,err=toolbox.InitMysql(this.cfg.DB)
	if err != nil{
		this.loger.Fatal(err)
	}
	this.DB=NewDatabase(this.db)
	if appdir,err:=toolbox.AppDir();err != nil{
		this.loger.Fatal(err)
	}else if err:=os.Chdir(appdir);err!=nil{
		this.loger.Fatal(err)
	}
	if runtime.GOOS == "windows"{
		this.slash="\\"
	}else{
		this.slash="/"
	}
	this.initMiddleWare()
	this.initRenderer()
}


func (this *Global) initMiddleWare(){
	this.middleware=NewMiddleWare()
	this.middleware.JwtKey=[]byte(this.cfg.Secret.LoginJwt)
	this.middleware.loger=this.loger
	this.renderMiddleWare=NewRenderMiddleWare()
}


func (this *Global) initRenderer(){
	if this.cfg.Renderer.Enable{
		if this.cfg.Renderer.Type=="files"{
			pattern:=this.cfg.Renderer.Pattern
			pattern=strings.ReplaceAll(pattern,"\\",this.slash)
			pattern=strings.ReplaceAll(pattern,"/",this.slash)
			//this.renderer.Funcs(template.FuncMap{"include":renderer.include})
			this.renderer=NewRenderer(pattern,this.cfg.Debug)
			//this.renderer.Use(this.renderMiddleWare.AddUserNamefunc)
		}
	}
}
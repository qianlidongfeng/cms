package main

import (
	"github.com/qianlidongfeng/httpserver"
	"github.com/qianlidongfeng/loger"
	"github.com/qianlidongfeng/toolbox"
)

type SecretConfig struct{
	LoginJwt string
}

type RendererConfig struct{
	Enable bool
	Type string
	Pattern string
}

type Config struct{
	Debug bool
	Httpserver httpserver.Config
	Log loger.Config
	DB toolbox.MySqlConfig
	Secret SecretConfig
	Renderer RendererConfig
}

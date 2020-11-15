package main

import (
	"github.com/gorilla/mux"
	"mariadb"
	"os"
	"server"
	"strings"
	"tools"
)

var db bool

func main() {
	host := "127.0.0.1"
	port := uint(3010)
	dbUser := "root"
	dbPswd := ""
	dbHost := "127.0.0.1"
	dbPort := uint(3306)
	dbBase := "cep"

	sysConf := make(map[string]string, 2)
	for _, arg := range os.Args {
		key := strings.Split(arg, "=")
		if len(key) != 2 {
			continue
		}
		sysConf[tools.Lower(key[0])] = key[1]
	}

	if len(sysConf["host"]) > 0 {
		host = sysConf["host"]
	}
	if len(sysConf["port"]) > 0 {
		port = tools.Stui(sysConf["port"])
	}
	if len(sysConf["dbuser"]) > 0 {
		dbUser = sysConf["dbUser"]
	}
	if len(sysConf["dbpswd"]) > 0 {
		dbPswd = sysConf["dbpswd"]
	}
	if len(sysConf["dbhost"]) > 0 {
		dbHost = sysConf["dbhost"]
	}
	if len(sysConf["dbport"]) > 0 {
		dbPort = tools.Stui(sysConf["dbport"])
	}
	if len(sysConf["dbbase"]) > 0 {
		dbBase = sysConf["dbbase"]
	}

	db = mariadb.New(mariadb.Connect{
		User: dbUser,
		Pswd: dbPswd,
		Host: dbHost,
		Port: dbPort,
		Base: dbBase,
	})

	r := mux.NewRouter()
	routes(r)

	server.Start(server.Server{
		Name:   "CEP - api",
		Host:   host,
		Port:   port,
		Router: r,
	})
}

func routes(r *mux.Router) {
	server.AddRoute(server.Route{Url: "/", Method: get, Handler: home, Router: r})
	server.AddRoute(server.Route{Url: "/status", Method: get, Handler: home, Router: r})
	server.AddRoute(server.Route{Url: "/buscar/{cep:[0-9]+}", Method: get, Handler: pesquisar, Router: r})
	server.AddRoute(server.Route{Url: "/buscar/{cep:[0-9]+}", Method: post, Handler: pesquisar, Router: r})
}

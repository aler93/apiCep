package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"server"
	"tools"
)

func home(w http.ResponseWriter, r *http.Request) {
	response := server.Json{
		Status: 200,
		Msg:    "Ok",
		Data:   map[string]string{"Title": "API CEP", "Author": "Alisson Naimayer"},
	}

	if !db {
		response.Status = 500
		response.Msg = "Erro: Conexão com banco de dados recusada"
		response.Data = map[string]string{"Projeto": "API CEP", "Autor": "Alisson Naimayer"}
	}

	server.JSON(w, response)
}

func pesquisar(w http.ResponseWriter, r *http.Request) {
	response := server.Json{
		Status: 200,
		Msg:    "Ok",
	}

	cep := mux.Vars(r)["cep"]
	if len(cep) != 8 {
		response.Status = 400
		response.Msg = "CEP inválido, o CEP precisa ter 8 dígitos. "
		response.Data = cep
	}

	if len(cep) == 8 {
		masked := cep[0:5] + "-" + cep[5:8]
		json := fetchCepAsync(masked)
		if tools.EmptyString(json.CEP) {
			response.Status = 404
			response.Msg = "CEP não encontrado, verifique e tente novamente"
		}
		response.Data = json
	}

	server.JSON(w, response)
}

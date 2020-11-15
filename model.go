package main

import (
	"mariadb"
	"tools"
)

type Cep struct {
	ID          uint64 `json:"id"`
	CEP         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	UF          string `json:"uf"`
	IBGE        string `json:"ibge"`
	GIA         string `json:"gia"`
	DDD         string `json:"ddd"`
	SIAFI       string `json:"siafi"`
}

func create(row Cep) uint64 {
	var ibge interface{}
	ibge = tools.Sti(row.IBGE)
	if ibge == 0 {
		ibge = nil
	}
	var gia interface{}
	gia = tools.Sti(row.GIA)
	if gia == 0 {
		gia = nil
	}
	var ddd interface{}
	ddd = tools.Sti(row.DDD)
	if ddd == 0 {
		ddd = nil
	}
	var siafi interface{}
	siafi = tools.Sti(row.SIAFI)
	if siafi == 0 {
		siafi = nil
	}
	var complemento interface{}
	complemento = row.Complemento
	if tools.EmptyString(row.Complemento) {
		complemento = nil
	}
	query, vars := mariadb.Insert(mariadb.Create{
		Table: "viacep",
		Values: map[string]interface{}{
			"cep":         row.CEP,
			"logradouro":  row.Logradouro,
			"complemento": complemento,
			"bairro":      row.Bairro,
			"localidade":  row.Localidade,
			"uf":          row.UF,
			"ibge":        ibge,
			"gia":         gia,
			"ddd":         ddd,
			"siafi":       siafi,
		},
	})

	res := mariadb.Execute(query, vars...)
	id, err := res.LastInsertId()
	tools.StayCalm(err)

	return uint64(id)
}

func update(row Cep) uint64 {
	var ibge interface{}
	ibge = tools.Sti(row.IBGE)
	if ibge == 0 {
		ibge = nil
	}
	var gia interface{}
	gia = tools.Sti(row.GIA)
	if gia == 0 {
		gia = nil
	}
	var ddd interface{}
	ddd = tools.Sti(row.DDD)
	if ddd == 0 {
		ddd = nil
	}
	var siafi interface{}
	siafi = tools.Sti(row.SIAFI)
	if siafi == 0 {
		siafi = nil
	}
	var complemento interface{}
	complemento = row.Complemento
	if tools.EmptyString(row.Complemento) {
		complemento = nil
	}
	query, vars := mariadb.Update(mariadb.Edit{
		Table: "viacep",
		Values: map[string]interface{}{
			"logradouro":  row.Logradouro,
			"complemento": complemento,
			"bairro":      row.Bairro,
			"localidade":  row.Localidade,
			"uf":          row.UF,
			"ibge":        ibge,
			"gia":         gia,
			"ddd":         ddd,
			"siafi":       siafi,
		},
		Where: map[string]interface{}{"cep": row.CEP},
	})

	res := mariadb.Execute(query, vars...)
	id, err := res.LastInsertId()
	tools.StayCalm(err)

	return uint64(id)
}

package main

import (
	"database/sql"
	"encoding/json"
	"mariadb"
	"net/http"
	"sync"
	"time"
	"tools"
)

func fetchCepAsync(cep string) Cep {
	var wg sync.WaitGroup
	var f Cep
	var viacep Cep

	go func() {
		callViacep(cep, &viacep)

		compare(f, viacep)
	}()

	wg.Add(1)
	go func(group *sync.WaitGroup, f *Cep) {
		query, vars := mariadb.Select(mariadb.Read{
			Table:    "viacep",
			Cols:     []string{"ID", "cep", "logradouro", "complemento", "bairro", "localidade", "uf", "ibge", "gia", "ddd", "siafi"},
			WhereOld: map[string]interface{}{"cep": cep},
		})
		res := mariadb.Row(query, vars...)

		var complemento sql.NullString
		var gia sql.NullString
		tools.StayCalm(
			res.Scan(&f.ID, &f.CEP, &f.Logradouro, &complemento, &f.Bairro, &f.Localidade, &f.UF, &f.IBGE, &gia, &f.DDD, &f.SIAFI),
		)
		if complemento.Valid {
			f.Complemento = complemento.String
		}
		if gia.Valid {
			f.GIA = gia.String
		}

		group.Done()
	}(&wg, &f)

	wg.Wait()

	if f.ID <= 0 {
		for i := 0; i < 100; i++ {
			if !tools.EmptyString(viacep.CEP) {
				create(viacep)
				f = viacep
				break
			}

			tools.Wait(time.Second / 10)
		}
	}

	return f
}

/*func fetchCep(cep string) Cep {
	query, vars := mariadb.Select(mariadb.Read{
		Table:    "viacep",
		WhereOld: map[string]interface{}{"cep": cep},
	})
	res := mariadb.Row(query, vars...)

	var f Cep
	var complemento sql.NullString
	var gia sql.NullString
	tools.StayCalm(
		res.Scan(&f.ID, &f.CEP, &f.Logradouro, &complemento, &f.Bairro, &f.Localidade, &f.UF, &f.IBGE, &gia, &f.DDD, &f.SIAFI),
	)
	if complemento.Valid {
		f.Complemento = complemento.String
	}
	if gia.Valid {
		f.GIA = gia.String
	}

	if f.ID <= 0 {
		callViacep(cep, &f)

		go func() {
			if !tools.EmptyString(f.CEP) {
				id := create(f)

				f.ID = id
			}
		}()
	}

	return f
}*/

func callViacep(cep string, target *Cep) {
	url := "https://viacep.com.br/ws/" + cep + "/json"

	resp, err := http.Get(url)
	tools.StayCalm(err)
	// defer tools.StayCalm(resp.Body.Close())

	r := json.NewDecoder(resp.Body)
	tools.StayCalm(r.Decode(target))
}

func compare(local Cep, viacep Cep) {
	diff := false
	if local.Logradouro != viacep.Logradouro {
		diff = true
	}
	if local.Complemento != viacep.Complemento {
		diff = true
	}
	if local.Bairro != viacep.Bairro {
		diff = true
	}
	if local.Localidade != viacep.Localidade {
		diff = true
	}
	if local.UF != viacep.UF {
		diff = true
	}
	if local.IBGE != viacep.IBGE {
		diff = true
	}
	if local.GIA != viacep.GIA {
		diff = true
	}
	if local.DDD != viacep.DDD {
		diff = true
	}
	if local.SIAFI != viacep.SIAFI {
		diff = true
	}

	if diff {
		update(viacep)
	}
}

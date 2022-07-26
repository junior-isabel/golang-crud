package servidor

import (
	"crud/banco"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type usuario struct {
	ID    uint32 `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func CriarUsuario(w http.ResponseWriter, r *http.Request) {

	body, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		w.Write([]byte("error no servidor"))
		return
	}

	var usuario usuario

	if erro := json.Unmarshal(body, &usuario); erro != nil {
		w.Write([]byte("erro ao converte para struct usuario"))
		return
	}

	db, erro := banco.Conectar()

	if erro != nil {
		w.Write([]byte("Erro na conexão com servidor"))
		fmt.Println(erro)
		return
	}

	defer db.Close()

	statement, erro := db.Prepare("insert into usuarios (nome, email) values(?,?)")
	if erro != nil {
		w.Write([]byte("Error a criar o template string de inserção"))
		fmt.Println(erro)
		return
	}
	defer statement.Close()

	insercao, erro := statement.Exec(usuario.Name, usuario.Email)
	if erro != nil {
		w.Write([]byte("Corpo da requisão invalido"))
		fmt.Println(erro)
		return
	}

	usuarioId, erro := insercao.LastInsertId()

	if erro != nil {
		w.Write([]byte("Não foi possivel pega o id do ultimo usuário"))
		fmt.Println(erro)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("usuario criado com sucesso id: %d", usuarioId)))
}

func ListUsuarios(w http.ResponseWriter, r *http.Request) {

	db, erro := banco.Conectar()

	if erro != nil {
		w.Write([]byte("Erro na conexão com servidor"))
		fmt.Println(erro)
		return
	}

	defer db.Close()
	var usuarios []usuario
	linhas, erro := db.Query("select * from usuarios")

	if erro != nil {
		w.Write([]byte("sql invalida"))
		fmt.Println(erro)
		return
	}

	defer linhas.Close()

	for linhas.Next() {
		var usuario usuario
		erro := linhas.Scan(&usuario.ID, &usuario.Name, &usuario.Email)
		if erro != nil {
			w.Write([]byte("erro ao preencher a struct usuario com dados da base de dados"))
			fmt.Println(erro)
			return
		}

		usuarios = append(usuarios, usuario)
	}

	w.WriteHeader(http.StatusOK)

	if erro := json.NewEncoder(w).Encode(usuarios); erro != nil {

		w.Write([]byte("Erro ao converte usuarios para jsons"))
		fmt.Println(erro)
		return
	}
}

func ListUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	ID, erro := strconv.ParseUint(parametros["id"], 10, 32)

	if erro != nil {
		w.Write([]byte("parametro id não existe"))
		fmt.Println(erro)
		return
	}

	db, erro := banco.Conectar()

	if erro != nil {
		w.Write([]byte("Erro na conexão comm banco de dados"))
		fmt.Println(erro)
		return
	}

	defer db.Close()

	linha, erro := db.Query("select * from usuarios where id=?", ID)
	if erro != nil {
		w.Write([]byte("Erro não criar consulta esquele"))
		fmt.Println(erro)
		return
	}
	var usuario usuario
	defer linha.Close()
	if linha.Next() {
		if erro := linha.Scan(&usuario.ID, &usuario.Name, &usuario.Email); erro != nil {
			w.Write([]byte("erro ao recoperar usuario da base de dados"))
			fmt.Println(erro)
			return
		}
	}
	if usuario.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("usuario solicitado não existe"))
		fmt.Println("usuario solicitado não existe")
		return
	}

	if erro := json.NewEncoder(w).Encode(usuario); erro != nil {
		w.Write([]byte("erro ao converte para json"))
		fmt.Println(erro)
	}

}

func AtualizarUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	ID, erro := strconv.ParseUint(parametros["id"], 10, 32)

	if erro != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Parametro invalido"))
		fmt.Println(erro)
		return
	}

	body, erro := ioutil.ReadAll(r.Body)

	if erro != nil {
		w.Write([]byte("Erro ao ler o corpo da requisição"))
		fmt.Println(erro)
		return
	}
	var usuario usuario

	if erro := json.Unmarshal(body, &usuario); erro != nil {
		w.Write([]byte("Erro ao fazer parse nos dados"))
		fmt.Println(erro)
		return
	}

	db, erro := banco.Conectar()

	if erro != nil {
		w.Write([]byte("Erro com a conexão da base de dados"))
		fmt.Println(erro)
		return
	}
	defer db.Close()

	statement, erro := db.Prepare("update usuarios set nome = ?, email = ? where id = ?")

	if erro != nil {
		w.Write([]byte("Erro ao criar sql"))
		fmt.Println(erro)
		return
	}

	defer statement.Close()

	if _, erro := statement.Exec(usuario.Name, usuario.Email, ID); erro != nil {
		w.Write([]byte(""))
		fmt.Println(erro)

	}

}

func EliminarUsuario(w http.ResponseWriter, r *http.Request) {

	parametros := mux.Vars(r)

	ID, erro := strconv.ParseUint(parametros["id"], 10, 32)

	if erro != nil {
		w.Write([]byte("Parametro invalid"))
		return
	}

	db, erro := banco.Conectar()

	if erro != nil {
		w.Write([]byte("Erro ao criar conexão com base de dados"))
		return
	}

	defer db.Close()

	statement, erro := db.Prepare("delete from usuarios where id = ?")

	if erro != nil {
		w.Write([]byte("Erro ao criar sql consulta"))
		return
	}

	defer statement.Close()

	if linha, erro := statement.Exec(ID); erro != nil {
		w.Write([]byte("Não foi possivel executar a operação"))
		return
	} else if rows, erro := linha.RowsAffected(); erro != nil {
		w.Write([]byte("Erro ao consultar numero de linha elimada"))
		return
	} else if rows == 0 {

		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("usuario invalido, não foi possivel remover este usuario"))
	}

}

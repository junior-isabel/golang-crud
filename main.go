package main

import (
	"crud/servidor"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/usuarios", servidor.CriarUsuario).Methods(http.MethodPost)
	router.HandleFunc("/usuarios", servidor.ListUsuarios).Methods(http.MethodGet)
	router.HandleFunc("/usuarios/{id}", servidor.ListUsuario).Methods(http.MethodGet)
	router.HandleFunc("/usuarios/{id}", servidor.AtualizarUsuario).Methods(http.MethodPut)
	router.HandleFunc("/usuarios/{id}", servidor.EliminarUsuario).Methods(http.MethodDelete)

	fmt.Println("excutando na porta :3000")
	log.Fatal(http.ListenAndServe(":3000", router))

}

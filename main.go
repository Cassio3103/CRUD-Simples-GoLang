package main

import (
	"golang/config"
	handles "golang/handlers"
	"golang/models"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	/*Iniciar conexões com banco de dados e definir rotas*/
	dbConection := config.SetUpDatabase()

	// Fechar conexão com o banco
	defer dbConection.Close()

	_, err := dbConection.Exec(models.CreateTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	// SETANDO ROTAS
	router := mux.NewRouter()

	taskHandler := handles.NewTaskHandler(dbConection)

	router.HandleFunc("/tasks", taskHandler.ReadTasks).Methods("GET")
	router.HandleFunc("/tasks", taskHandler.CreateTask).Methods("POST")
	router.HandleFunc("/tasks/{id}", taskHandler.UpdateTask).Methods("PUT")
	router.HandleFunc("/tasks/{id}", taskHandler.DeleteTask).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))

}

package main

import (
	"golang/config"
	"golang/models"
	"log"
	"net/http"
)

func main() {

	/*Iniciar conexões com banco de dados e definir rotas*/
	dbConection := config.SetUpDatabase()

	_, err := dbConection.Exec(models.CreateTableSQL)

	if err != nil {
		log.Fatal(err)
	}

	// Fechar conexão com o banco
	defer dbConection.Close()

	log.Fatal(http.ListenAndServe(":8080", nil))

}

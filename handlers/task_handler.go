package handles

import (
	"database/sql"
	"encoding/json"
	"golang/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type TaskHandler struct {
	DB *sql.DB
}

// CONSTRUTOR PARA TASKHANDLER
func NewTaskHandler(db *sql.DB) *TaskHandler {
	return &TaskHandler{DB: db}
}

// READ
func (taskHandler *TaskHandler) ReadTasks(
	writer http.ResponseWriter,
	request *http.Request,
) {

	rows, err := taskHandler.DB.Query("SELECT * FROM tasks")
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	var tasks []models.Task

	for rows.Next() {
		var task models.Task

		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
		)

		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		tasks = append(tasks, task)
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(tasks)
}

// CREATE
func (taskHandler *TaskHandler) CreateTask(
	writer http.ResponseWriter,
	request *http.Request,
) {

	var task models.Task

	err := json.NewDecoder(request.Body).Decode(&task)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = taskHandler.DB.Exec(
		"INSERT INTO tasks (title, description, status) VALUES ($1, $2, $3)",
		task.Title,
		task.Description,
		task.Status,
	)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)

	json.NewEncoder(writer).Encode(task)
}

// UPDATE
func (taskHandler *TaskHandler) UpdateTask(
	writer http.ResponseWriter,
	request *http.Request,
) {

	params := mux.Vars(request)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(writer, "Invalid ID", http.StatusBadRequest)
		return
	}

	var task models.Task

	err = json.NewDecoder(request.Body).Decode(&task)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = taskHandler.DB.Exec(
		"UPDATE tasks SET title=$1, description=$2, status=$3 WHERE id=$4",
		task.Title,
		task.Description,
		task.Status,
		id,
	)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")

	json.NewEncoder(writer).Encode(map[string]string{
		"message": "Task updated successfully",
	})
}

// DELETE
func (taskHandler *TaskHandler) DeleteTask(
	writer http.ResponseWriter,
	request *http.Request,
) {

	params := mux.Vars(request)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(writer, "Invalid ID", http.StatusBadRequest)
		return
	}

	_, err = taskHandler.DB.Exec(
		"DELETE FROM tasks WHERE id=$1",
		id,
	)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")

	json.NewEncoder(writer).Encode(map[string]string{
		"message": "Task deleted successfully",
	})
}

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type task struct {
	ID      int    `json:"id"`
	Name    string `json:"nombre"`
	Content string `json:"contenido"`
}

type allTask []task

var tasks = allTask{
	{
		ID:      1,
		Name:    "tarea 1",
		Content: "Some content",
	},
}

func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", indexRout)
	router.HandleFunc("/tasks", getTasks).Methods("GET")
	router.HandleFunc("/tasks", createTaskt).Methods("POST")
	router.HandleFunc("/tasks/{id}", getTask).Methods("GET")
	router.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE")
	router.HandleFunc("/tasks/{id}", updateTask).Methods("PUT")

	http.ListenAndServe(":8080", router)
	// log.Fatal(http.ListenAndServe(":8080", router))
}

func indexRout(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "Hola Mundo")
}

func createTaskt(res http.ResponseWriter, req *http.Request) {
	var newTask task
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Fprintf(res, "Inserte una tarea valida")
	}

	// Agregamos el json del cliente en forma de bytes a nuestra variable
	json.Unmarshal(reqBody, &newTask)
	// Creamos el ID
	newTask.ID = len(tasks) + 1
	// Agregamos a nuestra lista de tareas
	tasks = append(tasks, newTask)
	// Respondemos al cliente con la tarea recien creada
	res.Header().Set("Content-type", "application/json")
	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(newTask)
}

func getTasks(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-type", "application/json")
	json.NewEncoder(res).Encode(tasks)
}

func getTask(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	taskId, err := strconv.Atoi(vars["id"])

	if err != nil {
		fmt.Fprintf(res, "ID Invalido")
		return
	}

	for _, task := range tasks {
		if task.ID == taskId {
			res.Header().Set("Content-type", "application/json")
			json.NewEncoder(res).Encode(task)
		}
	}
}

func deleteTask(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	taskId, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Fprintf(res, "ID Invalido")
		return
	}

	for index, task := range tasks {
		if task.ID == taskId {
			tasks = append(tasks[:index], tasks[:index+1]...)
			res.Header().Set("Content-type", "application/json")
			json.NewEncoder(res).Encode(tasks)
		}
	}
}

func updateTask(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	taskId, err := strconv.Atoi(vars["id"])
	var updatedTask task

	if err != nil {
		fmt.Fprintf(res, "ID Invalido")
		return
	}

	reqBody, err := ioutil.ReadAll(req.Body)

	if err != nil {
		fmt.Fprintf(res, "Por favor ingrese datos validos")
	}

	json.Unmarshal(reqBody, &updatedTask)

	for index, tarea := range tasks {
		if tarea.ID == taskId {
			tasks = append(tasks[:index], tasks[index+1:]...)
			updatedTask.ID = taskId
			tasks = append(tasks, updatedTask)

			fmt.Fprintf(res, "la tarea con %v sea ctualizó correctamente", taskId)
		}
	}
}

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// TYPES
type Task struct {
	Title      string
	Category   string
	IsComplete bool
}

type ToDoApp []Task

// FILE MANIPULATION

func writeDataFromJson(tasks ToDoApp) error {
	file, err := os.OpenFile("tasks.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		return err
	}

	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(tasks)
	if err != nil {
		return err
	}

	return nil
}

func readDataFromJson() ToDoApp {
	file, err := os.OpenFile("tasks.json", os.O_RDWR|os.O_CREATE, 0664)
	if err != nil {
		fmt.Println("Ha habido un error para abrir el archivo", err)
		return ToDoApp{}
	}
	defer file.Close()

	var taskApp ToDoApp
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&taskApp)

	if err != nil {
		fmt.Println("Ha habido un error de codificacion")
		return ToDoApp{}
	}

	return taskApp
}

// TASKS FUNCTIONS

func showHelp() {
	fmt.Println(`APLICACION TODO CLI

-h	help		Imprime a pantalla este menu de ayuda
-a	add task	Agrega una tarea, todo -a <Title> <Category>
-l	list tasks	Imprimir todas las tareas en pantalla
-u	update tasks	Actualiza el title y categoria de tarea, todo -u <newTitle> <newCategory>
-d	delete task	Borra una tarea, todo -d <Id>
-s	change status	Pasa la tarea de False a True y viceversa, todo -s <Id>`)
}

func createTask(title, category string) {
	tasks := readDataFromJson()
	newTask := Task{
		Title:      title,
		Category:   category,
		IsComplete: false,
	}

	tasks = append(tasks, newTask)
	err := writeDataFromJson(tasks)

	if err != nil {
		fmt.Printf("Ha ocurrido un error: %v", err)
	}

	tasks = ToDoApp{newTask}
	formatedText(tasks)
}

func deleteTask(id string) error {
	tasks := readDataFromJson()
	newTasks := ToDoApp{}

	number, err := strconv.Atoi(id)
	if err != nil {
		return errors.New("Argumento no es un numero valido")

	}

	if len(tasks) == 0 {
		return errors.New("No hay tareas para poder borrar")
	}

	for i, task := range tasks {
		if i == number {
			continue
		}
		newTasks = append(newTasks, task)
	}
	err = writeDataFromJson(newTasks)
	if err != nil {
		return err
	}

	return nil
}

func updateTask(id, newTitle, category string) error {
	tasks := readDataFromJson()

	if len(tasks) == 0 {
		return errors.New("No hay tareas creadas, utiliza todo -c <Title> <Category> para crear una")
	}

	idUpdatedTask, err := strconv.Atoi(id)
	if err != nil || len(tasks) < idUpdatedTask {
		return errors.New("El argumento <id> no es un numero o fuera de rango")
	}

	newTask := Task{Title: newTitle, Category: category, IsComplete: tasks[idUpdatedTask].IsComplete}

	tasks[idUpdatedTask] = newTask

	err = writeDataFromJson(tasks)

	if err != nil {
		return errors.New("No se ha podido escribir en el archivo tasks.json")
	}

	tasks = ToDoApp{newTask}
	fmt.Println("Se ha actualizado la tarea")
	formatedText(tasks)
	return nil
}

func completeTask(id string) error {
	tasks := readDataFromJson()

	if len(tasks) == 0 {
		return errors.New("No hay tareas agregadas")
	}

	idCompleteTask, err := strconv.Atoi(id)
	if err != nil {
		return errors.New("Ha habido un error en convertir el id a un int")
	}

	if len(tasks) < idCompleteTask {
		return errors.New("id fuera de rango.")
	}

	tasks[idCompleteTask].IsComplete = !tasks[idCompleteTask].IsComplete

	err = writeDataFromJson(tasks)

	if err != nil {
		return err
	}

	fmt.Println("\nHaz cambiado el status de la tarea")
	formatedText(tasks)
	return nil
}

func maxLeght(tasks ToDoApp) int {
	maxNumber := 0
	for _, task := range tasks {
		v := reflect.ValueOf(task)
		for i := 0; i < v.NumField(); i++ {
			valueField := v.Field(i)
			fieldValueStr := fmt.Sprintf("%v", valueField)
			if len(fieldValueStr) > maxNumber {
				maxNumber = len(fieldValueStr)
			}
		}
	}
	return maxNumber

}

func formatedText(tasks ToDoApp) {
	w := maxLeght(tasks)
	if w < 8 {
		w = 8
	}
	bar := w*4 + 3*4 + 1
	baseString := strings.Repeat("-", bar)
	fmt.Printf("%v\n", baseString)
	fmt.Printf("| %-*s | %-*s | %-*s | %-*s |\n", w, "ID", w, "TITLE", w, "CATEGORY", w, "STATUS")
	fmt.Printf("%v\n", baseString)
	fmt.Printf("%v\n", baseString)
	for i, task := range tasks {
		fmt.Printf("| %-*d | %-*s | %-*s | %-*v |\n", w, i, w, task.Title, w, task.Category, w, task.IsComplete)
		fmt.Println(baseString)
	}
}

func listTask() {
	tasks := readDataFromJson()
	//maxLeght(tasks)
	formatedText(tasks)
}

func main() {
	if len(os.Args) <= 1 {
		listTask()
		return
	}

	switch action := os.Args[1]; action {
	case "-h":
		showHelp()
	case "-l":
		listTask()
	case "-a":
		if len(os.Args) < 4 {
			fmt.Println("Invalid Args: todo -a <title> <category>")
			return
		}
		createTask(os.Args[2], os.Args[3])
	case "-d":
		if len(os.Args) < 3 {
			fmt.Println("Invalid Args: todo -d <id>")
			return
		}
		err := deleteTask(os.Args[2])
		if err != nil {
			fmt.Printf("Ha habido un error: %v", err)
		}
	case "-u":

		if len(os.Args) < 5 {
			fmt.Println("Faltan argumentos, todo -u <Id> <newTitle> <newCategory>")
			return
		}
		err := updateTask(os.Args[2], os.Args[3], os.Args[4])
		if err != nil {
			fmt.Printf("Ha habido un error: %v", err)
		}
	case "-s":
		if len(os.Args) < 3 {
			fmt.Println("Faltan Argumentos: todo -s <Id>")
			return
		}
		err := completeTask(os.Args[2])
		if err != nil {
			fmt.Printf("Ha habido un error: %v", err)
		}
	}

}

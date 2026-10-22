package main

//importing libraries
import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//model definition
type Todo struct {
	//ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey;"`
	ID         uint `autoIncrement;primaryKey;"`
	Desc       string
	Status     string    `gorm:"check:,name like '%Completed%'"`
	Created_At time.Time `gorm:"autoCreateTime"`
	Updated_At time.Time `gorm:"autoUpdateTime"`
}

//Database connection settings are read from the environment (see .env.example).
//A .env file in the working directory is loaded if present; real deployments
//should set the variables directly and never commit a .env file.
func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func dbConfig() (host, user, password, name, port string) {
	_ = godotenv.Load()
	host = envOr("DB_HOST", "localhost")
	user = envOr("DB_USER", "postgres")
	name = envOr("DB_NAME", "mydb")
	port = envOr("DB_PORT", "5432")
	password = os.Getenv("DB_PASSWORD")
	if password == "" {
		log.Fatal("DB_PASSWORD is not set; copy .env.example to .env or export it")
	}
	return
}

var todo = Todo{
	Desc:   "Task1",
	Status: "Completed",
}

//type ToDo struct {
//	Desc   string
//	status string
//}

//Global variable declaration
var db *gorm.DB
var err error

func main() {

	//Connecting with database server
	host, user, password, name, port := dbConfig()
	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s", host, user, name, password, port)
	db, err = gorm.Open(postgres.Open(dbURI), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Succes!")
	}
	//Migrating the schema automatically to keep it updated
	db.AutoMigrate(&Todo{})
	//fmt.Println("Migrated")
	//rows := db.Create(&todo)
	/*if rows.RowsAffected == 0 {
		fmt.Println("Not created")
	} else {
		fmt.Println("Created")
	}*/
	router := mux.NewRouter()
	//Endpoints
	//Handler Functions
	fmt.Println("1.List all tasks /v1/alltask\n2.List with Completed(0) ,Not Completed filters(1)  /v1/alltodo/{filter condition-0 or 1} \n3.Get By Id /v1/todo/id \n4.Create a Task /v1/todo \n5.Bulk Create /v1/todo/bulk \n6.Update Task /v1/todo/id \n7.Delete a Task /v1/todo/id")
	router.HandleFunc("/v1/alltodo/{id}", getTasks).Methods("GET")
	router.HandleFunc("/v1/alltask", getAllTask).Methods("GET")
	router.HandleFunc("/v1/todo/{id}", getById).Methods("GET")
	router.HandleFunc("/v1/todo", createTask).Methods("POST")
	router.HandleFunc("/v1/todo/bulk", createBulk).Methods("POST")
	router.HandleFunc("/v1/todo/{id}", updateTask).Methods("PATCH")
	router.HandleFunc("/v1/todo/{id}", deleteTask).Methods("DELETE")
	http.ListenAndServe(":8080", router)
}
func getAllTask(w http.ResponseWriter, r *http.Request) {
	var tasks []Todo
	db.Find(&tasks)
	json.NewEncoder(w).Encode(&tasks)
}
func getTasks(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var tasks []Todo
	i, errs := strconv.Atoi(params["id"])
	if errs != nil {
		fmt.Println("Unable to convert to integer")
		json.NewEncoder(w).Encode(&tasks)
		return

	}
	//json.NewDecoder(r.Body).Decode(a)

	if i == 0 {
		db.Where("status=?", "Completed").Find(&tasks)

	} else {
		db.Where("status=?", "Not Completed").Find(&tasks)
	}

	json.NewEncoder(w).Encode(&tasks)
}
func getById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var task Todo
	db.First(&task, params["id"])
	json.NewEncoder(w).Encode(task)

}
func createTask(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	json.NewDecoder(r.Body).Decode(&todo)
	createdTask := db.Create(&todo)
	err = createdTask.Error
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}
	json.NewEncoder(w).Encode(todo)

}
func createBulk(w http.ResponseWriter, r *http.Request) {
	var tasks []Todo
	json.NewDecoder(r.Body).Decode(&tasks)
	createdResult := db.Create(&tasks)
	err = createdResult.Error
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}
	json.NewEncoder(w).Encode(&tasks)
}
func updateTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var task Todo
	var todo Todo
	json.NewDecoder(r.Body).Decode(&task)
	updatedResult := db.Model(&todo).Where("id=?", params["id"]).Updates(&task)
	err = updatedResult.Error
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}
	json.NewEncoder(w).Encode(task)
}
func deleteTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var task Todo
	db.First(&task, params["id"])
	db.Delete(&task)
	json.NewEncoder(w).Encode(&task)
}

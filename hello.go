package hello

import (
        "fmt"
        "log"
        "net/http"
        "encoding/json"
        "cloud.google.com/go/datastore"
        "google.golang.org/appengine"
        "github.com/gorilla/mux"
        "time"
)
type Task struct {
        Description    string    `datastore:"description"`
        Created time.Time `datastore:"created"`
        id      int64
}
func saveData(w http.ResponseWriter, r *http.Request){
        fmt.Println("Endpoint Hit: saveData")
        input := r.URL.Query().Get("input")
        fmt.Println(input)

        ctx := appengine.NewContext(r)
        projectID := "api-project-426361742627"
        client, err := datastore.NewClient(ctx, projectID)
        if err != nil {
                log.Fatalf("Failed to create client: %v", err)
        }
        kind := "Task"
        task := Task{
                Description: input,
                Created: time.Now(),
        }
        taskKey := datastore.IncompleteKey(kind, nil)

        if _, err := client.Put(ctx, taskKey, &task); err != nil {
                log.Fatalf("Failed to save task: %v", err)
        }

        fmt.Printf("Saved %v: %v\n", taskKey, task.Description)
        json.NewEncoder(w).Encode(task)
}
func returnAll(w http.ResponseWriter, r *http.Request){
        fmt.Println("Endpoint Hit: returnAll")

        ctx := appengine.NewContext(r)
        projectID := "api-project-426361742627"
        client, err := datastore.NewClient(ctx, projectID)
        if err != nil {
                log.Fatalf("Failed to create client: %v", err)
        }

        query := datastore.NewQuery("Task").Order("description")

        var tasks []*Task

        keys, err := client.GetAll(ctx, query, &tasks)
        for i, key := range keys {
                tasks[i].id = key.ID
        }

        json.NewEncoder(w).Encode(tasks)
}

func homePage(w http.ResponseWriter, r *http.Request){
        fmt.Fprintf(w, "Welcome to the HomePage!")
        fmt.Println("Endpoint Hit: homePage")
}
func handleRequests() {
        myRouter := mux.NewRouter()
        myRouter.HandleFunc("/", homePage)
        myRouter.HandleFunc("/save", saveData)
        myRouter.HandleFunc("/retrieve", returnAll)
        http.Handle("/", myRouter)
}
func init() {
        handleRequests()
}

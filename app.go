package main

import (
    "fmt"
    "log"
    "net/http"
    "os"

    "github.com/urfave/negroni"
    "github.com/gorilla/mux"
    "github.com/tobyjsullivan/ues-projection-svc/projection"
    "github.com/tobyjsullivan/log-sdk/reader"
    "github.com/satori/go.uuid"
    "encoding/json"
)

var (
    logger *log.Logger
    state *projection.Projection
)

func init() {
    var err error

    logger = log.New(os.Stdout, "[service] ", 0)

    logger.Println("Parsing SERVICE_LOG_ID...")
    logId := reader.LogID{}
    err = logId.Parse(os.Getenv("SERVICE_LOG_ID"))
    if err != nil {
        panic("Error parsing SERVICE_LOG_ID. "+err.Error())
    }

    logger.Println("Creating log reader client...")
    client, err := reader.New(&reader.ClientConfig{
        ServiceAddress: os.Getenv("LOG_READER_API"),
        Logger: logger,
    })
    if err != nil {
        panic("Error creating log reader client. "+err.Error())
    }

    logger.Println("Subscribing projection to log...")
    state = projection.New()
    client.Subscribe(logId, reader.EventID{}, state.Apply, true)
    logger.Println("Hydration complete.")
}

func main() {
    r := buildRoutes()

    n := negroni.New()
    n.UseHandler(r)

    port := os.Getenv("PORT")
    if port == "" {
        port = "3000"
    }

    n.Run(":" + port)
}

func buildRoutes() http.Handler {
    r := mux.NewRouter()
    r.HandleFunc("/", statusHandler).Methods("GET")
    r.HandleFunc("/accounts/{accountId}", getAccountHandler).Methods("GET")

    return r
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "The service is online!\n")
}

func getAccountHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    paramAcctId := vars["accountId"]

    if paramAcctId == "" {
        http.Error(w, "Must specify accountId in path.", http.StatusBadRequest)
        return
    }
    
    acctId := uuid.UUID{}
    if err := acctId.UnmarshalText([]byte(paramAcctId)); err != nil {
        http.Error(w, "Error parsing accountId. "+err.Error(), http.StatusBadRequest)
        return
    }

    acct := state.FindAccount(acctId)
    if acct == nil {
        http.Error(w, "No such account.", http.StatusNotFound)
        return
    }

    response := struct{
        AccountID string `json:"accountId"`
    } {
        AccountID: acct.ID.String(),
    }
    encoder := json.NewEncoder(w)
    if err := encoder.Encode(&response); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

package main

import (
	"log-service/data"
	"net/http"
)

type jsonPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	//read the json
	var req jsonPayload
	_ = app.readJSON(w, r, &req)

	event := data.LogEntry{
		Name: req.Name,
		Data: req.Data,
	}

	err := app.Models.LogEntry.Insert(event)
	if err != nil {
		app.errorJSON(w, err)
		return 
	}

	resp := jsonResponse{
		Error: false,
		Message: "logged",
	}

	app.writeJSON(w, 200, resp)
}
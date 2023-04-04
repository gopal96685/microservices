package main

import (
	"log"
	"net/http"
)

func (app *Config) SendMail(w http.ResponseWriter, r *http.Request) {
	log.Println("sendmail api working ...")
	type mailMsg struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Subject string `json:"subject"`
		Msg     string `json:"msg"`
	}

	var reqPayload mailMsg
	err := app.readJSON(w, r, &reqPayload)
	if err != nil {
		app.errorJSON(w, err)
		log.Println(err)
		return
	}
	log.Println(reqPayload)

	msg := Message{
		From:    reqPayload.From,
		To:      reqPayload.To,
		Subject: reqPayload.Subject,
		Data:    reqPayload.Msg,
	}

	err = app.Mailer.SendSMPTPMsg(msg)
	if err != nil {
		app.errorJSON(w, err)
		log.Println(err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "sent to " + reqPayload.To,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

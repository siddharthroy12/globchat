package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func (app *application) createReportHandler(w http.ResponseWriter, r *http.Request) {
	user := app.getUserFromRequst(r)

	var inputs struct {
		MessageId int    `json:"message_id"`
		Reason    string `json:"reason"`
	}

	err := app.readJSONFromRequest(w, r, &inputs)
	if err != nil {
		app.badRequestResponse(w, r, fmt.Errorf("invalid inputs"))
		return
	}

	exists, err := app.messageModel.Exists(inputs.MessageId)

	if !exists {
		app.badRequestResponse(w, r, fmt.Errorf("message doesn't exists"))
		return
	}

	if err != nil {
		app.serverErrorResponse(w, r, err, "checking if message exists")
		return
	}

	exists, err = app.reportModel.Exists(user.ID, inputs.MessageId)

	if err != nil {
		app.serverErrorResponse(w, r, err, "checking if report exists")
		return
	}

	if exists {
		app.badRequestResponse(w, r, fmt.Errorf("already exists"))
		return
	}

	err = app.reportModel.Create(user.ID, inputs.MessageId, inputs.Reason)

	if err != nil {
		app.serverErrorResponse(w, r, err, "creating report")
		return
	}

	app.writeJSON(w, 200, envelope{"message": "report created"}, nil)

}

func (app *application) deleteReportHandler(w http.ResponseWriter, r *http.Request) {
	reportId, err := strconv.Atoi(r.URL.Query().Get("reportId"))

	if err != nil {
		app.badRequestResponse(w, r, fmt.Errorf("reportId must be a valid number"))
		return
	}

	err = app.reportModel.RemoveById(reportId)

	if err != nil {
		app.badRequestResponse(w, r, fmt.Errorf("report not found"))
		return
	}

	app.writeJSON(w, 200, envelope{"message": "report deleted"}, nil)
}

func (app *application) resolveReportHandler(w http.ResponseWriter, r *http.Request) {
	reportId, err := strconv.Atoi(r.URL.Query().Get("reportId"))

	if err != nil {
		app.badRequestResponse(w, r, fmt.Errorf("reportId must be a valid number"))
		return
	}

	report, err := app.reportModel.GetByID(reportId)

	if err != nil {
		app.badRequestResponse(w, r, fmt.Errorf("report not found"))
		return
	}

	message, err := app.messageModel.GetByID(report.MessageId)

	if err != nil {
		app.serverErrorResponse(w, r, err, "get message by id")
		return
	}

	err = app.deleteMessage(message)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	app.writeJSON(w, 200, envelope{"message": "report deleted"}, nil)
}

func (app *application) queryReportsHandler(w http.ResponseWriter, r *http.Request) {

}

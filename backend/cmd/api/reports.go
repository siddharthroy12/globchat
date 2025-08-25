package main

import (
	"fmt"
	"net/http"
	"strconv"

	"globechat.live/internal/models"
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
	// Get search parameter (optional)
	search := r.URL.Query().Get("search")

	// Get pageSize parameter with default value
	pageSizeStr := r.URL.Query().Get("page_size")
	pageSize := 20 // default page size
	if pageSizeStr != "" {
		var err error
		pageSize, err = strconv.Atoi(pageSizeStr)
		if err != nil || pageSize <= 0 {
			app.badRequestResponse(w, r, fmt.Errorf("page_size must be a valid positive number"))
			return
		}
		// Set a reasonable maximum page size to prevent abuse
		if pageSize > 100 {
			pageSize = 100
		}
	}

	// Get pageIndex parameter with default value
	pageIndexStr := r.URL.Query().Get("page")
	pageIndex := 0 // default page index (first page)
	if pageIndexStr != "" {
		var err error
		pageIndex, err = strconv.Atoi(pageIndexStr)
		if err != nil || pageIndex < 0 {
			app.badRequestResponse(w, r, fmt.Errorf("page must be a valid non-negative number"))
			return
		}
	}

	// Create the query struct
	query := models.ReportQuery{
		Search:    search,
		PageSize:  pageSize,
		PageIndex: pageIndex,
	}

	// Execute the query
	result, err := app.reportModel.Query(query)
	if err != nil {
		app.serverErrorResponse(w, r, err, "query reports")
		return
	}

	// Return the results
	app.writeJSON(w, 200, envelope{"result": result}, nil)
}

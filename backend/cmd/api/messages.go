package main

import (
	"fmt"
	"net/http"
	"strconv"

	"globechat.live/internal/models"
)

func (app *application) createMessageHandler(w http.ResponseWriter, r *http.Request) {

	user := app.getUserFromRequst(r)

	var input struct {
		ThreadId int    `json:"thread_id"`
		Text     string `json:"text"`
		Image    string `json:"image"`
	}

	err := app.readJSONFromRequest(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	message, err := app.messageModel.Create(input.Text, input.Image, input.ThreadId, user.ID, false)

	if err != nil {
		app.serverErrorResponse(w, r, err, "create message")
		return
	}

	app.roomManager.notifyRoom(input.ThreadId, WebsocketConnectionMessage{
		Type:   "new-message",
		RoomID: input.ThreadId,
		Data:   message,
	})

	app.writeJSON(w, 200, envelope{"message": message}, nil)
}

func (app *application) deleteMessageHandler(w http.ResponseWriter, r *http.Request) {
	user := app.getUserFromRequst(r)

	messageId, err := strconv.Atoi(r.URL.Query().Get("messageId"))

	if err != nil {
		app.badRequestResponse(w, r, fmt.Errorf("messageId must be a valid number"))
		return
	}

	message, err := app.messageModel.GetByID(messageId)

	if err != nil {
		app.notFoundResponse(w, r, fmt.Errorf("message not found"))
		return
	}

	if message.UserId != user.ID {
		app.badRequestResponse(w, r, fmt.Errorf("you do not own this message naughty boy"))
		return
	}

	if message.IsFirst {
		app.badRequestResponse(w, r, fmt.Errorf("you are trying to do something funny i know it"))
		return
	}

	err = app.messageModel.Delete(messageId)

	if err != nil {
		app.serverErrorResponse(w, r, err, "delete message")
		return
	}

	app.roomManager.notifyRoom(message.ThreadId, WebsocketConnectionMessage{
		Type:   "delete-message",
		RoomID: message.ThreadId,
		Data:   message,
	})

	app.writeJSON(w, 200, envelope{"message": "message deleted"}, nil)
}

func (app *application) getMessagesHandler(w http.ResponseWriter, r *http.Request) {
	threadId, err := strconv.Atoi(r.URL.Query().Get("threadId"))
	if err != nil {
		app.badRequestResponse(w, r, fmt.Errorf("theadId must be a valid number"))
		return
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		app.badRequestResponse(w, r, fmt.Errorf("limit must be a valid number"))
		return
	}
	messageId, err := strconv.Atoi(r.URL.Query().Get("messageId"))

	if err != nil {
		messages, err := app.messageModel.GetByThreadID(threadId, limit)
		if err != nil {
			app.serverErrorResponse(w, r, err, "get messages for thread id")
			return
		}

		app.writeJSON(w, 200, envelope{"messages": messages}, nil)
		return
	}

	direction := r.URL.Query().Get("direction")

	var messages []models.Message
	if direction == "after" {
		messages, err = app.messageModel.GetAfterID(threadId, messageId, limit)
		if err != nil {
			app.serverErrorResponse(w, r, err, "get messages before thread id")
			return
		}
	} else {
		messages, err = app.messageModel.GetBeforeID(threadId, messageId, limit)
		if err != nil {
			app.serverErrorResponse(w, r, err, "get messages before thread id")
			return
		}
	}

	app.writeJSON(w, 200, envelope{"messages": messages}, nil)
}

func (app *application) reportMessageHandler(w http.ResponseWriter, r *http.Request) {

	messageId, err := strconv.Atoi(r.URL.Query().Get("messageId"))

	if err != nil {
		app.badRequestResponse(w, r, fmt.Errorf("messageId must be a valid number"))
		return
	}

	err = app.messageModel.IncreaseReported(messageId)

	if err != nil {
		app.notFoundResponse(w, r, fmt.Errorf("message not found"))
		return
	}

	app.writeJSON(w, 200, envelope{"message": "message deleted"}, nil)
}

func (app *application) queryMessagesHandler(w http.ResponseWriter, r *http.Request) {
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
	query := models.MessageQuery{
		Search:    search,
		PageSize:  pageSize,
		PageIndex: pageIndex,
	}

	// Execute the query
	result, err := app.messageModel.Query(query)
	if err != nil {
		app.serverErrorResponse(w, r, err, "query messages")
		return
	}

	// Return the results
	app.writeJSON(w, 200, envelope{"result": result}, nil)
}

package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"globechat.live/internal/models"
)

func (app *application) createThreadHandler(w http.ResponseWriter, r *http.Request) {
	user := app.getUserFromRequst(r)

	var input struct {
		Lat     float64 `json:"lat"`
		Long    float64 `json:"long"`
		Message string  `json:"message"`
	}

	err := app.readJSONFromRequest(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if len(input.Message) == 0 {
		app.badRequestResponse(w, r, fmt.Errorf("message is empty"))
		return
	}

	thread, err := app.threadModel.Create(input.Message, input.Lat, input.Long, user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err, "create thread")
		return
	}
	_, err = app.messageModel.Create(input.Message, "", thread.ID, user.ID, true)
	if err != nil {
		app.serverErrorResponse(w, r, err, "create message")
		app.threadModel.Delete(thread.ID)
		return
	}

	app.writeJSON(w, 200, envelope{"thread": (thread)}, nil)
}

func (app *application) getThreadsHandler(w http.ResponseWriter, r *http.Request) {

	lat, err := strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
	if err != nil {
		app.badRequestResponse(w, r, fmt.Errorf("send valid lat"))
		return
	}

	long, err := strconv.ParseFloat(r.URL.Query().Get("long"), 64)
	if err != nil {
		app.badRequestResponse(w, r, fmt.Errorf("send valid long"))
		return
	}
	km, err := strconv.ParseFloat(r.URL.Query().Get("km"), 64)
	if err != nil {
		app.badRequestResponse(w, r, fmt.Errorf("send valid km"))
		return
	}

	if km > 50 {
		app.badRequestResponse(w, r, fmt.Errorf("km too big"))
		return
	}

	threads, err := app.threadModel.GetByLocationRadius(lat, long, km)

	if err != nil {
		app.serverErrorResponse(w, r, err, "fetching threads")
		return
	}

	app.writeJSON(w, 200, envelope{"threads": threads}, nil)
}

func (app *application) getRandomThread(w http.ResponseWriter, r *http.Request) {

	thread, err := app.threadModel.GetRandomThread()

	if err != nil {
		app.serverErrorResponse(w, r, err, "fetching threads")
		return
	}

	app.writeJSON(w, 200, envelope{"thread": (thread)}, nil)
}

func (app *application) deleteThreadHandler(w http.ResponseWriter, r *http.Request) {
	user := app.getUserFromRequst(r)

	threadId, err := strconv.Atoi(r.URL.Query().Get("threadId"))

	if err != nil {
		app.badRequestResponse(w, r, fmt.Errorf("theadId must be a valid number"))
		return
	}

	thread, err := app.threadModel.GetById(threadId)

	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFoundResponse(w, r, fmt.Errorf("thread not found"))
			return
		}
		app.serverErrorResponse(w, r, err, "delete thread")
		return
	}

	if thread.UserId != user.ID {
		app.badRequestResponse(w, r, fmt.Errorf("you do not own this thread naughty boy"))
		return
	}
	err = app.messageModel.DeleteByThreadID(threadId)
	if err != nil {
		app.serverErrorResponse(w, r, err, "delete messages using thread id")
		return
	}
	err = app.threadModel.Delete(threadId)
	if err != nil {
		app.serverErrorResponse(w, r, err, "delete thread")
		return
	}

	app.roomManager.notifyRoom(threadId, WebsocketConnectionMessage{
		Type:   "delete-thread",
		RoomID: threadId,
		Data:   thread,
	})
	app.writeJSON(w, 200, envelope{"message": "thread deleted"}, nil)
}

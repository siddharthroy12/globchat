package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"globechat.live/internal/models"
)

const MinDistanceBetweenThreads = 0.05 // in km

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

	threads, err := app.threadModel.GetByLocationRadius(input.Lat, input.Long, MinDistanceBetweenThreads)

	if err != nil {
		app.serverErrorResponse(w, r, err, "fetching nearby threads which creating one")
		return
	}

	if len(threads) > 0 {
		app.badRequestResponse(w, r, fmt.Errorf("too close to other threads"))
		return
	}

	thread, err := app.threadModel.Create(input.Message, input.Lat, input.Long, user.ID)
	if err != nil {
		if errors.Is(err, models.ErrTextTooLong) {
			app.badRequestResponse(w, r, err)
			return
		}
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

	minLat, err := strconv.ParseFloat(r.URL.Query().Get("minLat"), 64)
	if err != nil {
		app.badRequestResponse(w, r, fmt.Errorf("send valid minLat"))
		return
	}

	maxLat, err := strconv.ParseFloat(r.URL.Query().Get("maxLat"), 64)
	if err != nil {
		app.badRequestResponse(w, r, fmt.Errorf("send valid maxLat"))
		return
	}

	minLong, err := strconv.ParseFloat(r.URL.Query().Get("minLong"), 64)
	if err != nil {
		app.badRequestResponse(w, r, fmt.Errorf("send valid minLong"))
		return
	}

	maxLong, err := strconv.ParseFloat(r.URL.Query().Get("maxLong"), 64)
	if err != nil {
		app.badRequestResponse(w, r, fmt.Errorf("send valid maxLong"))
		return
	}

	threads, err := app.threadModel.GetByBounds(minLat, minLong, maxLat, maxLong, 500)

	if err != nil {
		if errors.Is(err, models.ErrTooManyItems) {
			app.badRequestResponse(w, r, err)
			return
		}
		app.serverErrorResponse(w, r, err, "fetching threads")
		return
	}

	app.writeJSON(w, 200, envelope{"threads": threads}, nil)
}

func (app *application) getThreadByIDHandler(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("id")

	threadId, err := strconv.Atoi(id)
	if err != nil {
		app.badRequestResponse(w, r, fmt.Errorf("id must be a valid number"))
		return
	}

	thread, err := app.threadModel.GetById(threadId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFoundResponse(w, r, fmt.Errorf("thread not found"))
			return
		}
		app.serverErrorResponse(w, r, err, "fetching thread")
		return
	}

	app.writeJSON(w, 200, envelope{"thread": thread}, nil)
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

	if thread.UserId != user.ID && !user.IsAdmin {
		app.badRequestResponse(w, r, fmt.Errorf("you do not own this thread naughty boy"))
		return
	}
	err = app.deleteThread(threadId)
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

func (app *application) deleteThread(threadId int) error {
	err := app.threadModel.Delete(threadId)
	if err != nil {
		return err
	}

	app.roomManager.notifyRoom(threadId, WebsocketConnectionMessage{
		Type:   "delete-thread",
		RoomID: threadId,
		Data:   "",
	})

	return nil
}

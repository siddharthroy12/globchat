package main

import (
	"fmt"
	"net/http"
	"strconv"

	"globechat.live/internal/models"
)

func createThreadObject(thread models.Thread) map[string]any {
	return map[string]any{
		"id":         thread.ID,
		"lat":        thread.Lat,
		"long":       thread.Long,
		"user_id":    thread.UserId,
		"user_name":  thread.Username,
		"user_image": thread.UserImage,
		"created_at": thread.CreatedAt,
	}
}

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

	thread, err := app.threadModel.Create(input.Lat, input.Long, user.ID)

	if err != nil {
		app.serverErrorResponse(w, r, err, "create thread")
		return
	}

	app.writeJSON(w, 200, envelope{"thread": createThreadObject(thread)}, nil)
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

	var res = make([]map[string]any, 0)

	for _, thread := range threads {
		res = append(res, createThreadObject(*thread))
	}

	app.writeJSON(w, 200, envelope{"threads": res}, nil)
}

func (app *application) getRandomThread(w http.ResponseWriter, r *http.Request) {

	thread, err := app.threadModel.GetRandomThread()

	if err != nil {
		app.serverErrorResponse(w, r, err, "fetching threads")
		return
	}

	app.writeJSON(w, 200, envelope{"thread": createThreadObject(thread)}, nil)
}

func (app *application) deleteThreadHandler(w http.ResponseWriter, r *http.Request) {
	user := app.getUserFromRequst(r)

	threadId := r.URL.Query().Get("threadId")

	i, err := strconv.Atoi(threadId)

	if err != nil {
		app.badRequestResponse(w, r, fmt.Errorf("theadId must be a valid number"))
	}

	if len(threadId) == 0 {
		app.badRequestResponse(w, r, fmt.Errorf("theadId is missing"))
		return
	}

	thread, err := app.threadModel.GetById(i)

	if err != nil {
		app.notFoundResponse(w, r, fmt.Errorf("thread not found"))
		return
	}

	if thread.UserId != user.ID {
		app.badRequestResponse(w, r, fmt.Errorf("you do not own this thread naughty boy"))
		return
	}

	app.threadModel.Delete(i)

	app.writeJSON(w, 200, envelope{"message": "thread deleted"}, nil)
}

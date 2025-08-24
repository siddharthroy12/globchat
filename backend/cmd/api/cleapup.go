package main

import "time"

func (app *application) startCleanupRoutine() {
	ticker := time.NewTicker(24 * time.Hour) // Check every hour

	go func() {
		for range ticker.C {
			err := app.cleanupExpiredThreads()
			if err != nil {
				app.logger.Error(err.Error(), "action", "cleaning expired threads")
			}
		}
	}()
}

func (app *application) cleanupExpiredThreads() error {
	ids, err := app.threadModel.GetExpiredIds()
	if err != nil {
		return err
	}

	for _, id := range ids {
		err = app.deleteThread(id)
		if err != nil {
			return err
		}
		app.roomManager.notifyRoom(id, WebsocketConnectionMessage{
			Type:   "delete-thread",
			RoomID: id,
			Data:   "",
		})
	}

	return err
}

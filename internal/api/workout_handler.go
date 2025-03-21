package api

import (
	"database/sql"
	"encoding/json"
	"github.com/nickemma/internal/store"
	"github.com/nickemma/internal/utils"
	"log"
	"net/http"
)

// decoupling our database
type WorkoutHandler struct {
	workoutStore store.WorkoutStore
	logger       *log.Logger
}

func NewWorkoutHandler(workoutStore store.WorkoutStore, logger *log.Logger) *WorkoutHandler {
	return &WorkoutHandler{
		workoutStore: workoutStore,
		logger:       logger,
	}
}

// HandlerGetWorkoutByID Get the workout
func (wh *WorkoutHandler) HandlerGetWorkoutByID(w http.ResponseWriter, r *http.Request) {
	workoutId, err := utils.ReadJSON(r)
	if err != nil {
		wh.logger.Printf("ERROR: readIdParams: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid workout id"})
		return
	}
	workout, err := wh.workoutStore.GetWorkoutByID1(workoutId)
	if err != nil {
		wh.logger.Printf("ERROR: GetworkoutById: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"workout": workout})
}

// HandlerCreateWorkout Create a workout
func (wh *WorkoutHandler) HandlerCreateWorkout(w http.ResponseWriter, r *http.Request) {
	var workout store.Workout
	err := json.NewDecoder(r.Body).Decode(&workout)

	if err != nil {
		wh.logger.Printf("ERROR: decodingworkout: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request sent"})
		return
	}

	createdWorkout, err := wh.workoutStore.CreateWorkout(&workout)
	if err != nil {
		wh.logger.Printf("ERROR: CreateWorkout: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"workout": createdWorkout})
}

// HandleUpdateWorkoutById update a workout
func (wh *WorkoutHandler) HandleUpdateWorkoutById(w http.ResponseWriter, r *http.Request) {
	workoutId, err := utils.ReadJSON(r)
	if err != nil {
		wh.logger.Printf("ERROR: readIdParams: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request sent"})
		return
	}

	existingWorkout, err := wh.workoutStore.GetWorkoutByID(workoutId)
	if err != nil {
		wh.logger.Printf("ERROR: GetWorkoutByID: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	if existingWorkout == nil {
		http.NotFound(w, r)
		return
	}
	// at this point we have our workout
	var updateWorkoutRequest struct {
		Title           *string              `json:"title"`
		Description     *string              `json:"description"`
		DurationMinutes *int                 `json:"duration_minutes"`
		CaloriesBurned  *int                 `json:"calories_burned"`
		Entries         []store.WorkoutEntry `json:"entries"`
	}
	err = json.NewDecoder(r.Body).Decode(&updateWorkoutRequest)
	if err != nil {
		wh.logger.Printf("ERROR: decodingerror: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid bad request payload"})
		return
	}
	if updateWorkoutRequest.Title != nil {
		existingWorkout.Title = *updateWorkoutRequest.Title
	}
	if updateWorkoutRequest.Description != nil {
		existingWorkout.Description = *updateWorkoutRequest.Description
	}
	if updateWorkoutRequest.DurationMinutes != nil {
		existingWorkout.DurationMinutes = *updateWorkoutRequest.DurationMinutes
	}
	if updateWorkoutRequest.CaloriesBurned != nil {
		existingWorkout.CaloriesBurned = *updateWorkoutRequest.CaloriesBurned
	}
	if updateWorkoutRequest.Entries != nil {
		existingWorkout.Entries = updateWorkoutRequest.Entries
	}
	err = wh.workoutStore.UpdateWorkout(existingWorkout)
	if err != nil {
		wh.logger.Printf("ERROR: updateworkout: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "error updating workout"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"workout": existingWorkout})
}

// // HandleDeleteWorkoutById update a workout
func (wh *WorkoutHandler) HandleDeleteWorkoutById(w http.ResponseWriter, r *http.Request) {
	workoutId, err := utils.ReadJSON(r)
	if err != nil {
		wh.logger.Printf("ERROR: readIdParams: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid workout id"})
		return
	}

	err = wh.workoutStore.DeleteWorkout(workoutId)
	if err == sql.ErrNoRows {
		wh.logger.Printf("ERROR: deleteworkout: %v", err)
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "workout not found"})
		return
	}
	if err != nil {
		wh.logger.Printf("ERROR: deleteworkout: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "error deleting workout"})
		return
	}
	utils.WriteJSON(w, http.StatusNoContent, utils.Envelope{})
}

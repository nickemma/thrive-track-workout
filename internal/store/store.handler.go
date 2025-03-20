package store

import "database/sql"

type Workout struct {
	ID              int            `json:"id"`
	Title           string         `json:"title"`
	Description     string         `json:"description"`
	DurationMinutes int            `json:"duration_minutes"`
	CaloriesBurned  int            `json:"calories_burned"`
	Entries         []WorkoutEntry `json:"entries"`
}

type WorkoutEntry struct {
	ID              int      `json:"id"`
	ExcerciseName   string   `json:"excercise_name"`
	Sets            int      `json:"sets"`
	Reps            *int     `json:"reps"`
	DurationSeconds *int     `json:"duration_seconds"`
	Weight          *float64 `json:"weight"`
	Notes           string   `json:"notes"`
	OrderIndex      int      `json:"order_index"`
}

type PostgresWorkoutStore struct {
	db *sql.DB
}

func NewPostgresWorkoutStore(db *sql.DB) *PostgresWorkoutStore {
	return &PostgresWorkoutStore{db: db}
}

type WorkoutStore interface {
	CreateWorkout(*Workout) (*Workout, error)
	GetWorkoutByID(id int64) (*Workout, error)
}

// CreateWorkout Creating a workout transaction
func (pg *PostgresWorkoutStore) CreateWorkout(workout *Workout) (*Workout, error) {
	// Beginning the sql transaction and commiting it to the database
	tx, err := pg.db.Begin()
	if err != nil {
		return nil, err
	}
	// rolling back our transaction in case of failed transactions or error
	defer tx.Rollback()

	// inserting the data into our database
	query := `
INSERT INTO workouts (title, description, duration_minutes, calories_burned) 
VALUES ($1, $2, $3, $4) 
RETURNING id
`
	err = tx.QueryRow(query, workout.Title, workout.Description, workout.DurationMinutes, workout.CaloriesBurned).Scan(&workout.ID)

	if err != nil {
		return nil, err
	}

	for _, entry := range workout.Entries {
		query := `
INSERT INTO workout_entries (workout_id, excercise_name, sets, reps, duration_seconds, weight, notes, order_index) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
RETURNING id
    `
		err = tx.QueryRow(query, workout.ID, entry.ExcerciseName, entry.Sets, entry.Reps, entry.DurationSeconds, entry.Weight, entry.Notes, entry.OrderIndex).Scan(&entry.ID)
		if err != nil {
			return nil, err
		}
	}

	// commiting the transaction
	err = tx.Commit()

	if err != nil {
		return nil, err
	}

	return workout, nil
}

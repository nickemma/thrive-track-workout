package store

import (
	"database/sql"
	"errors"
)

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
	ExerciseName    string   `json:"exercise_name"`
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
	GetWorkoutByID1(id int64) (*Workout, error)
	UpdateWorkout(*Workout) error
	DeleteWorkout(id int64) error
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

	for i := range workout.Entries {
		entry := &workout.Entries[i]
		query := `
INSERT INTO workout_entries (workout_id, exercise_name, sets, reps, duration_seconds, weight, notes, order_index) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
RETURNING id
    `
		err = tx.QueryRow(query, workout.ID, entry.ExerciseName, entry.Sets, entry.Reps, entry.DurationSeconds, entry.Weight, entry.Notes, entry.OrderIndex).Scan(&entry.ID)
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

// GetWorkoutById getting the workout by id
func (pg *PostgresWorkoutStore) GetWorkoutByID(id int64) (*Workout, error) {
	workout := &Workout{}
	query := `
	SELECT id, title, description, duration_minutes, calories_burned
    FROM workouts
    WHERE id = $1;
`
	err := pg.db.QueryRow(query, id).Scan(&workout.ID, &workout.Title, &workout.Description, &workout.DurationMinutes, &workout.CaloriesBurned)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	entryQuery := `
   SELECT id, exercise_name, sets, reps, duration_seconds, weight, notes, order_index
FROM workout_entries
WHERE workout_id = $1
ORDER BY order_index
`
	rows, err := pg.db.Query(entryQuery, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var entry WorkoutEntry
		err = rows.Scan(
			&entry.ID,
			&entry.ExerciseName,
			&entry.Sets,
			&entry.Reps,
			&entry.DurationSeconds,
			&entry.Weight,
			&entry.Notes,
			&entry.OrderIndex,
		)
		if err != nil {
			return nil, err
		}
		workout.Entries = append(workout.Entries, entry)
	}
	return workout, nil
}

// GetWorkoutByID1 getting the workout by id another method
func (pg *PostgresWorkoutStore) GetWorkoutByID1(id int64) (*Workout, error) {
	query := `
        SELECT 
            w.id, w.title, w.description, w.duration_minutes, w.calories_burned,
            e.id, e.exercise_name, e.sets, e.reps, e.duration_seconds, e.weight, e.notes, e.order_index
        FROM workouts w
        LEFT JOIN workout_entries e ON w.id = e.workout_id
        WHERE w.id = $1
        ORDER BY e.order_index
    `

	rows, err := pg.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workout Workout
	var hasEntries bool

	for rows.Next() {
		var entry WorkoutEntry
		err := rows.Scan(
			&workout.ID, &workout.Title, &workout.Description,
			&workout.DurationMinutes, &workout.CaloriesBurned,
			&entry.ID, &entry.ExerciseName, &entry.Sets, &entry.Reps,
			&entry.DurationSeconds, &entry.Weight, &entry.Notes, &entry.OrderIndex,
		)
		if err != nil {
			return nil, err
		}

		// Only append valid entries (LEFT JOIN may return NULLs)
		if entry.ID != 0 {
			workout.Entries = append(workout.Entries, entry)
			hasEntries = true
		}
	}

	if workout.ID == 0 { // No workout found
		return nil, nil
	}

	// Handle case where no entries exist
	if !hasEntries {
		workout.Entries = []WorkoutEntry{}
	}

	return &workout, nil
}

// UpdateWorkout Update a workout
func (pg *PostgresWorkoutStore) UpdateWorkout(workout *Workout) error {
	tx, err := pg.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	query := `
UPDATE workouts 
SET title = $1, description = $2, duration_minutes = $3, calories_burned = $4
WHERE id = $5
`

	result, err := tx.Exec(query, workout.Title, workout.Description, workout.DurationMinutes, workout.CaloriesBurned, workout.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	_, err = tx.Exec(`DELETE FROM workout_entries WHERE workout_id = $1`, workout.ID)
	if err != nil {
		return err
	}

	for _, entry := range workout.Entries {
		query := `
INSERT INTO workout_entries (workout_id, exercise_name, sets, reps, duration_seconds, weight, notes, order_index) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
`
		_, err := tx.Exec(
			query,
			workout.ID,
			entry.ExerciseName,
			entry.Sets,
			entry.Reps,
			entry.DurationSeconds,
			entry.Weight,
			entry.Notes,
			entry.OrderIndex,
		)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

// DeleteWorkout deletes a workout
func (pg *PostgresWorkoutStore) DeleteWorkout(id int64) error {
	query := `
DELETE FROM workouts
WHERE id = $1
`
	res, err := pg.db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

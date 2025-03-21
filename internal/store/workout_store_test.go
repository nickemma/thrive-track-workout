package store

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func SetupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("pgx", "host=localhost user=root password=postgres dbname=postgres port=5433 sslmode=disable")

	if err != nil {
		t.Fatal("opening test db: %w", err)
	}

	// run migrations for the test database
	err = Migrate(db, "../../migration/")
	if err != nil {
		t.Fatalf("migrating test db error: %v", err)
	}

	_, err = db.Exec(`TRUNCATE workouts, workout_entries CASCADE`)
	if err != nil {
		t.Fatalf("truncating table error: %v", err)
	}
	return db
}
func TestCreate(t *testing.T) {
	db := SetupTestDB(t)
	defer db.Close()

	store := NewPostgresWorkoutStore(db)
	// table driven test
	test := []struct {
		name    string
		workout *Workout
		wantErr bool
	}{{name: "valid workout", workout: &Workout{
		Title:           "push ups",
		Description:     "upper body day",
		DurationMinutes: 60,
		CaloriesBurned:  200,
		Entries: []WorkoutEntry{
			{
				ExerciseName: "Bench press",
				Sets:         3,
				Reps:         IntPointer(10),
				Weight:       FloatPointer(135.5),
				Notes:        "Awesome today",
				OrderIndex:   1,
			},
		},
	},
		wantErr: false,
	},
		{
			name: "workout with invalid code",
			workout: &Workout{
				Title:           "full body",
				Description:     "complete workout",
				DurationMinutes: 90,
				CaloriesBurned:  500,
				Entries: []WorkoutEntry{
					{
						ExerciseName: "Plank",
						Sets:         3,
						Reps:         IntPointer(60),
						Notes:        "keep form",
						OrderIndex:   1,
					},
					{
						ExerciseName:    "Plank",
						Sets:            3,
						Reps:            IntPointer(12),
						DurationSeconds: IntPointer(60),
						Weight:          FloatPointer(185.6),
						Notes:           "full depth",
						OrderIndex:      2,
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			createdWorkout, err := store.CreateWorkout(tt.workout)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.workout.Title, createdWorkout.Title)
			assert.Equal(t, tt.workout.Description, createdWorkout.Description)
			assert.Equal(t, tt.workout.DurationMinutes, createdWorkout.DurationMinutes)
			assert.Equal(t, tt.workout.CaloriesBurned, createdWorkout.CaloriesBurned)

			retrieved, err := store.GetWorkoutByID(int64(createdWorkout.ID))
			require.NoError(t, err)

			assert.Equal(t, createdWorkout.ID, retrieved.ID)
			assert.Equal(t, len(tt.workout.Entries), len(retrieved.Entries))

			for i := range retrieved.Entries {
				assert.Equal(t, tt.workout.Entries[i].ExerciseName, retrieved.Entries[i].ExerciseName)
				assert.Equal(t, tt.workout.Entries[i].Sets, retrieved.Entries[i].Sets)
				assert.Equal(t, tt.workout.Entries[i].OrderIndex, retrieved.Entries[i].OrderIndex)

			}
		})
	}
}

func IntPointer(i int) *int {
	return &i
}
func FloatPointer(i float64) *float64 {
	return &i
}

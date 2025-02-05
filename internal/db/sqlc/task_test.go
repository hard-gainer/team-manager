package db

import (
	"context"
	"testing"
	"time"

	"github.com/hard-gainer/task-tracker/internal/util"
	"github.com/stretchr/testify/require"
)

func createRandomTask(t *testing.T) Task {
    arg := CreateTaskParams{
        Title:       util.RandomString(10),
        Description: util.RandomString(20),
        DueTo:       time.Now().Add(24 * time.Hour),
        Status:      "ASSIGNED",
        Priority:    "LOW",
        ProjectID:   util.ToNullInt4(1),
        AssignedTo:  util.ToNullInt4(1),
    }

    task, err := testStore.CreateTask(context.Background(), arg)
    require.NoError(t, err)
    require.NotEmpty(t, task)

    require.Equal(t, arg.Title, task.Title)
    require.Equal(t, arg.Description, task.Description)
    require.WithinDuration(t, arg.DueTo, task.DueTo, time.Millisecond)
    require.Equal(t, arg.Status, task.Status)
    require.Equal(t, arg.Priority, task.Priority)
    require.Equal(t, arg.ProjectID, task.ProjectID)
    require.Equal(t, arg.AssignedTo, task.AssignedTo)

    require.NotZero(t, task.ID)
    require.NotZero(t, task.CreatedAt)

    return task
}

func TestCreateTask(t *testing.T) {
    createRandomTask(t)
}

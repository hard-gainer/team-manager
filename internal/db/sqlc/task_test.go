package db

import (
	"context"
	"testing"
	"time"

	"github.com/hard-gainer/team-manager/internal/util"
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

func TestGetTask(t *testing.T) {
    task1 := createRandomTask(t)
    task2, err := testStore.GetTask(context.Background(), task1.ID)
    require.NoError(t, err)
    require.NotEmpty(t, task2)

    require.Equal(t, task1.ID, task2.ID)
    require.Equal(t, task1.Title, task2.Title)
    require.Equal(t, task1.Description, task2.Description)
    require.WithinDuration(t, task1.CreatedAt, task2.CreatedAt, time.Millisecond)
    require.WithinDuration(t, task1.DueTo, task2.DueTo, time.Millisecond)
    require.Equal(t, task1.Status, task2.Status)
    require.Equal(t, task1.Priority, task2.Priority)
    require.Equal(t, task1.ProjectID, task2.ProjectID)
    require.Equal(t, task1.AssignedTo, task2.AssignedTo)
}

func TestListTasks(t *testing.T) {
    for i := 0; i < 10; i++ {
        createRandomTask(t)
    }

    tasks, err := testStore.ListTasks(context.Background())
    require.NoError(t, err)
    require.NotEmpty(t, tasks)
    require.GreaterOrEqual(t, len(tasks), 10)

    for _, task := range tasks {
        require.NotEmpty(t, task)
    }
}

func TestListProjectTasks(t *testing.T) {
    projectID := util.ToNullInt4(1)
    
    for i := 0; i < 5; i++ {
        arg := CreateTaskParams{
            Title:       util.RandomString(10),
            Description: util.RandomString(20),
            DueTo:       time.Now().Add(24 * time.Hour),
            Status:      "ASSIGNED",
            Priority:    "LOW",
            ProjectID:   projectID,
            AssignedTo:  util.ToNullInt4(1),
        }
        _, err := testStore.CreateTask(context.Background(), arg)
        require.NoError(t, err)
    }

    tasks, err := testStore.ListProjectTasks(context.Background(), projectID)
    require.NoError(t, err)
    require.NotEmpty(t, tasks)
    
    for _, task := range tasks {
        require.NotEmpty(t, task)
        require.Equal(t, projectID, task.ProjectID)
    }
}

func TestUpdateTaskStatus(t *testing.T) {
    task1 := createRandomTask(t)

    arg := UpdateTaskStatusParams{
        ID:     task1.ID,
        Status: "STARTED",
    }

    task2, err := testStore.UpdateTaskStatus(context.Background(), arg)
    require.NoError(t, err)
    require.NotEmpty(t, task2)

    require.Equal(t, task1.ID, task2.ID)
    require.Equal(t, arg.Status, task2.Status)
    require.Equal(t, task1.Title, task2.Title)
    require.Equal(t, task1.Description, task2.Description)
    require.WithinDuration(t, task1.DueTo, task2.DueTo, time.Millisecond)
}

func TestUpdateTaskPriority(t *testing.T) {
    task1 := createRandomTask(t)

    arg := UpdateTaskPriorityParams{
        ID:       task1.ID,
        Priority: "HIGH",
    }

    task2, err := testStore.UpdateTaskPriority(context.Background(), arg)
    require.NoError(t, err)
    require.NotEmpty(t, task2)

    require.Equal(t, task1.ID, task2.ID)
    require.Equal(t, arg.Priority, task2.Priority)
    require.Equal(t, task1.Title, task2.Title)
    require.Equal(t, task1.Description, task2.Description)
    require.WithinDuration(t, task1.DueTo, task2.DueTo, time.Millisecond)
}

func TestUpdateTaskTitle(t *testing.T) {
    task1 := createRandomTask(t)

    newTitle := util.RandomString(10)
    arg := UpdateTaskTitleParams{
        ID:    task1.ID,
        Title: newTitle,
    }

    task2, err := testStore.UpdateTaskTitle(context.Background(), arg)
    require.NoError(t, err)
    require.NotEmpty(t, task2)

    require.Equal(t, task1.ID, task2.ID)
    require.Equal(t, newTitle, task2.Title)
    require.Equal(t, task1.Description, task2.Description)
    require.WithinDuration(t, task1.DueTo, task2.DueTo, time.Millisecond)
}

func TestDeleteTask(t *testing.T) {
    task1 := createRandomTask(t)
    err := testStore.DeleteTask(context.Background(), task1.ID)
    require.NoError(t, err)

    task2, err := testStore.GetTask(context.Background(), task1.ID)
    require.Error(t, err)
    require.Empty(t, task2)
}
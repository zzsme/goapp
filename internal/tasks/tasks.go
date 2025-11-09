package tasks

import (
	"fmt"
	"time"

	"goapp/internal/app"
)

// TaskFunc represents a function that can be executed as a task
type TaskFunc func() error

// tasks is a map of task names to task functions
var tasks = map[string]TaskFunc{
	"cleanup":     cleanupTask,
	"data-sync":   dataSyncTask,
	"send-emails": sendEmailsTask,
}

// RunTask runs a task by name
func RunTask(taskName string) {
	// Check if task exists
	taskFunc, exists := tasks[taskName]
	if !exists {
		availableTasks := ""
		for name := range tasks {
			availableTasks += name + ", "
		}
		if len(availableTasks) > 0 {
			availableTasks = availableTasks[:len(availableTasks)-2] // Remove trailing comma and space
		}

		app.Error("Task not found", "task", taskName, "available", availableTasks)
		fmt.Printf("Task '%s' not found. Available tasks: %s\n", taskName, availableTasks)
		return
	}

	// Execute task
	app.Info("Running task", "task", taskName)
	fmt.Printf("Running task: %s\n", taskName)

	startTime := time.Now()
	err := taskFunc()
	duration := time.Since(startTime)

	if err != nil {
		app.Error("Task failed", "task", taskName, "error", err, "duration", duration)
		fmt.Printf("Task '%s' failed: %v (took %v)\n", taskName, err, duration)
		return
	}

	app.Info("Task completed", "task", taskName, "duration", duration)
	fmt.Printf("Task '%s' completed successfully (took %v)\n", taskName, duration)
}

// cleanupTask is an example task
func cleanupTask() error {
	// Simulate work
	time.Sleep(1 * time.Second)
	fmt.Println("Cleaning up old data...")

	// Example: You could clean up old records in the database
	// err := app.DB.Where("created_at < ?", time.Now().AddDate(0, -3, 0)).Delete(&models.OldData{}).Error

	return nil
}

// dataSyncTask is an example task
func dataSyncTask() error {
	// Simulate work
	time.Sleep(2 * time.Second)
	fmt.Println("Syncing data from external API...")

	// Example: You could fetch data from an external API and update your database

	return nil
}

// sendEmailsTask is an example task
func sendEmailsTask() error {
	// Simulate work
	time.Sleep(1500 * time.Millisecond)
	fmt.Println("Sending scheduled emails...")

	// Example: You could send newsletter emails to users

	return nil
}

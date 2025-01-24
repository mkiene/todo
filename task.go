package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
)

// Task represents a Taskwarrior task
type Task struct {
	ID          int          `json:"id"`                    // Task ID (not persistent)
	Description string       `json:"description"`           // Main task description
	Status      string       `json:"status"`                // Task status (e.g., pending, completed)
	Due         string       `json:"due,omitempty"`         // Due date (ISO 8601 format)
	Entry       string       `json:"entry"`                 // Entry timestamp
	Modified    string       `json:"modified"`              // Last modified timestamp
	Start       string       `json:"start,omitempty"`       // Start timestamp (if task is started)
	End         string       `json:"end,omitempty"`         // Completion timestamp
	Project     string       `json:"project,omitempty"`     // Associated project name
	Priority    string       `json:"priority,omitempty"`    // Task priority (L, M, H)
	Tags        []string     `json:"tags,omitempty"`        // List of tags
	Depends     []string     `json:"depends,omitempty"`     // List of dependent task UUIDs
	Urgency     float64      `json:"urgency,omitempty"`     // Calculated urgency value
	Annotations []Annotation `json:"annotations,omitempty"` // Annotations
	Recur       string       `json:"recur,omitempty"`       // Recurrence interval
	Wait        string       `json:"wait,omitempty"`        // Wait timestamp
	Scheduled   string       `json:"scheduled,omitempty"`   // Scheduled start time
	Until       string       `json:"until,omitempty"`       // Expiration date for recurring tasks
	UUID        string       `json:"uuid"`                  // Unique identifier (persistent)
	Parent      string       `json:"parent,omitempty"`      // Parent task UUID (recurring tasks)
	Imask       int          `json:"imask,omitempty"`       // Recurrence mask (for internal use)
	Mask        int          `json:"mask,omitempty"`        // Masked tasks in recurring sets
}

// Annotation represents annotations for a task
type Annotation struct {
	Entry       string `json:"entry"`       // Timestamp of the annotation
	Description string `json:"description"` // Content of the annotation
}

func get_tasks() ([]*Task, error) {
	// Execute Taskwarrior command
	cmd := exec.Command("task", "export")
	var out bytes.Buffer
	cmd.Stdout = &out

	// Run the command and check for errors
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("task command failed: %v\n", err)
	}

	// Parse the JSON output into a slice of Task structs
	var tasks []*Task
	if err := json.Unmarshal(out.Bytes(), &tasks); err != nil {
		return nil, fmt.Errorf("failed to parse JSON output: %v\n", err)
	}

	return tasks, nil
}

func find_task(description string) (*Task, error) {

	if description == "none" {
		return &Task{}, nil
	}

	tasks, err := get_tasks()
	if err != nil {
		return nil, err
	}

	for _, task := range tasks {
		if task.Description == description {
			return task, nil
		}
	}

	return nil, fmt.Errorf("couldn't find task '%v'.", description)
}

func get_tags() ([]string, error) {
	tasks, err := get_tasks()
	if err != nil {
		return nil, err
	}

	var tags []string

	for _, task := range tasks {
		for _, task_tag := range task.Tags {
			is_duplicate := false
			for _, tag := range tags {
				if tag == task_tag {
					is_duplicate = true
					break
				}
			}
			if !is_duplicate {
				tags = append(tags, task_tag)
			}
		}
	}

	return tags, nil
}

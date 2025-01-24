package main

import "fmt"

type Project struct {
	Name  string // Project name
	Tasks []*Task // List of tasks in the project
}

func get_projects() ([]*Project, error) {
	task_map := make(map[string][]*Task)
	tasks, err := get_tasks()
	if err != nil {
		return nil, err
	}

	for _, task := range tasks {
		project := task.Project
		task_map[project] = append(task_map[project], task)
	}

	var projects []*Project

	for name, tasks := range task_map {
		projects = append(projects, &Project{
			Name:  name,
			Tasks: tasks,
		})
	}

	return projects, nil
}

func find_project(name string) (*Project, error) {

	if name == "none" {
		return &Project{}, nil
	}

	projects, err := get_projects()
	if err != nil {
		return nil, err
	}

	for _, project := range projects {
		if project.Name == name {
			return project, nil
		}
	}

	return nil, fmt.Errorf("couldn't find project '%v'.", name)
}

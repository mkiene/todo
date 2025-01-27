package main

import (
	"fmt"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mkiene/huh"
)

func create_task() error {
	projects, err := get_projects()
	if err != nil {
		return err
	}

	tasks, err := get_tasks()
	if err != nil {
		return err
	}

	fetched_tags, err := get_tags()
	if err != nil {
		return err
	}

	fetched_tags = append(fetched_tags, "new")

	var project_names []string
	var task_names []string

	for _, project := range projects {
		if project.Name != "" {
			project_names = append(project_names, project.Name)
		}
	}

	for _, task := range tasks {
		if task.Description != "" {
			task_names = append(task_names, task.Description)
		}
	}

	project_names = append(project_names, "none", "new")
	task_names = append(task_names, "none")

	project_options := huh.NewOptions(project_names...)
	tag_options := huh.NewOptions(fetched_tags...)

	project_name := ""
	new_proj_name := ""
	dependency_name := ""
	var dependency *Task
	description := ""
	annotation := ""
	var tags []string
	new_tag_name := ""
	due_date := ""
	due_time := "PT23H59M"

	project_selection := huh.NewSelect[string]().
		Title("Choose a Project").
		Options(project_options...).
		Value(&project_name)

	new_project_input := huh.NewInput().
		Title("#####").
		Value(&new_proj_name)

	dependency_selection := huh.NewSelect[string]().
		Title("Choose a dependency").
		OptionsFunc(func() []huh.Option[string] {

			var child_tasks []string

			if project_name != "" &&
				project_name != "new" &&
				project_name != "none" {

				project, err := find_project(project_name)
				if err != nil {
					return nil
				}

				for _, task := range project.Tasks {
					if task.Status == "completed" {
						continue
					}
					child_tasks = append(child_tasks, task.Description)
				}
			}

			child_tasks = append(child_tasks, "none")

			opts := huh.NewOptions(child_tasks...)

			return opts

		}, &project_name).
		Value(&dependency_name)

	description_input := huh.NewInput().
		Title("Enter a Description").
		Placeholder("Lorem ipsum dolor sit amet.").
		Value(&description)

	annotation_input := huh.NewInput().
		Title("Enter an annotation").
		Placeholder("Lorem ipsum dolor sit amet.").
		Value(&annotation)

	tag_selection := huh.NewMultiSelect[string]().
		Title("Choose a tag").
		Options(tag_options...).
		Value(&tags)

	new_tag_input := huh.NewInput().
		Title("Enter a new tag(s)").
		Value(&new_tag_name)

	due_date_input := huh.NewInput().
		Title("Enter a due date in 'yyyy-mm-dd' format").
		Placeholder("yyyy-mm-dd").
		Value(&due_date)

	due_time_input := huh.NewInput().
		Title("Enter a time in 'hh-mm-ss' format").
		Value(&due_time)

	form := huh.NewForm(
		huh.NewGroup(
			project_selection,
		),
		huh.NewGroup(
			new_project_input,
		).WithHideFunc(func() bool {
			if project_name == "new" {
				new_project_input.Title("Enter a project name")
				return false
			}
			new_project_input.Title("#####")
			return true
		}),

		huh.NewGroup(
			dependency_selection,
		).WithHideFunc(func() bool {
			if project_name == "new" {
				return true
			}

			var child_tasks []string

			if project_name != "" &&
				project_name != "new" &&
				project_name != "none" {

				project, err := find_project(project_name)
				if err != nil {
					return true
				}

				for _, task := range project.Tasks {
					child_tasks = append(child_tasks, task.Description)
				}
			}

			if len(child_tasks) < 1 {
				return true
			}

			return false
		}),

		huh.NewGroup(
			description_input,
			annotation_input,
			tag_selection,
		),
		huh.NewGroup(
			new_tag_input,
		).WithHideFunc(func() bool {
			for _, tag := range tags {
				if tag == "new" {
					return false
				}
			}
			return true
		}),
		huh.NewGroup(
			due_date_input,
			due_time_input,
		),
	).WithTheme(huh.ThemeBase()).
		WithProgramOptions(tea.WithAltScreen())

	if err := form.Run(); err != nil {
		return err
	}

	if project_name == "none" {
		project_name = ""
	}

	if new_proj_name != "" {
		project_name = new_proj_name
	}

	tags = append(tags, new_tag_name)

	if dependency_name == "none" || dependency_name == "" {
		dependency = &Task{}
	} else {
		dependency, err = find_task(dependency_name)
		if err != nil {
			dependency = &Task{}
		}
	}

	cmd := exec.Command(
		"task",
		"add",
		`"`+description+`"`,
		"project:"+project_name,
		"due:"+due_date+"+"+due_time,
		"depends:"+dependency.UUID,
	)

	// Run the command and check for errors
	err = cmd.Run()
	if err != nil {
		fmt.Println("ERROR INITIAL COMMAND")
		fmt.Println(err)
		return fmt.Errorf("task command failed: %v\n", err)
	}

	task_object, err := find_task(description)
	if err != nil {
		return err
	}

	cmd = exec.Command("task",
		task_object.UUID,
		"annotate",
		`"`+annotation+`"`)

	// Run the command and check for errors
	err = cmd.Run()
	if err != nil {
		fmt.Println("ERROR ANNOTATION")
		fmt.Println(err)
		return fmt.Errorf("annotation command failed: %v\n", err)
	}

	for _, tag := range tags {

		if tag == "new" || tag == "" {
			continue
		}

		cmd = exec.Command("task",
			task_object.UUID,
			"modify",
			"+"+tag)

		// Run the command and check for errors
		err = cmd.Run()
		if err != nil {
			fmt.Println("ERROR TAG ADDITION(S)")
			fmt.Println(err)
			return fmt.Errorf("tag command failed: %v\n", err)
		}
	}

	fmt.Printf("created new task: '%v'!\ndue: %v\nannotation: %v\ntag(s): %v\n", description, due_date, annotation, tags)

	return nil
}

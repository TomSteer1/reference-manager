package main

import (
	"fmt"
	"encoding/csv"
	"os"
	"bufio"
	"strings"
	"strconv"
)

type Project struct {
	id string
	title string
	references map[string]Reference
}

type Reference struct {
	id string
	reftype string
	title string
	year string
	month string
	url string
	authors []string
	publisher string
}

var projects map[string]Project
var references map[string]Reference

func main() {	
	// Clear screen
	fmt.Print("\033[2J\033[H")
	setReferences()
	fmt.Println("Reference manager")
	fmt.Println("\t- Loaded", len(projects), "project(s)")
	fmt.Println("\t- Loaded", len(references), "reference(s)")
	menu()
}

func setReferences()  {
	projects = make(map[string]Project)
	references = make(map[string]Reference)
	referencesData := loadCSV("references.csv")
	projectsData := loadCSV("projects.csv")
	for _, referenceData := range referencesData {
		references[referenceData[0]] = Reference{
			id: referenceData[0],
			reftype: referenceData[1],
			title: referenceData[2],
			year: referenceData[3],
			month: referenceData[4],
			url: referenceData[5],
			publisher: referenceData[6],
		}
		authors := make([]string, len(referenceData)-7)
		for i := 7; i < len(referenceData); i++ {
			authors[i-7] = strings.TrimSpace(referenceData[i])
		}
		ref := references[referenceData[0]]
		ref.authors = authors
		references[referenceData[0]] = ref
	}

	for _, projectData := range projectsData {
		projects[projectData[0]] = Project{
			id: projectData[0],
			title: projectData[1],
			references: make(map[string]Reference, len(projectData)-2),
		}
		for i := 2; i < len(projectData); i++ {
			projects[projectData[0]].references[projectData[i]] = references[projectData[i]]
		}
	}
}

func loadCSV(fileName string) ([][]string) {
	// Open the file
	f, err := os.Open(fileName)
	handleErr(err)

	// Parse the file
	r := csv.NewReader(bufio.NewReader(f))
	r.FieldsPerRecord = -1
	r.Read() // Skip the header
	records, err := r.ReadAll()
	handleErr(err)

	return records
}

func menu(again ...bool) {
	fmt.Println("Main menu")
	fmt.Println("\t1. Select a project")
	fmt.Println("\t2. Select a reference")
	fmt.Println("\tPress enter to exit")
	count := 5 + len(again)
	switch input := getInput(); input {
	case "1":
		clearLines(count)
		selectProject()
	case "2":
		clearLines(count)
		selectReference()
	case "":
		clearLines(count + 2)
		fmt.Println("Exiting...")
		os.Exit(0)
	default:
		clearLines(count)
		fmt.Println("Invalid input")
		menu(true)
	}
}

func selectProject(again ...bool) {
	fmt.Println("Select a project")
	for id, project := range projects {
		fmt.Println("\t"+id, ":", project.title)
	}
	count := len(projects) + 4 + len(again)
	fmt.Println("\t0 : Create new project")
	fmt.Println("\tPress enter to return")
	switch input := getInput(); input {
	case "":
		clearLines(count)
		menu()
	case "0":
		clearLines(count)
		createProject()
	default:
		if _, ok := projects[input]; ok {
			clearLines(count)
			projectMenu(input)
		} else {
			clearLines(count)
			fmt.Println("Invalid input")
			selectProject(true)
		}
	}
}

func projectMenu(projectId string, again ...bool) {
	fmt.Println("Selected project:", projects[projectId].title)
	fmt.Println("\t1. View references")
	fmt.Println("\t2. Add reference")
	fmt.Println("\t3. Remove reference")
	fmt.Println("\t4. Export references")
	fmt.Println("\t5. Edit project")
	fmt.Println("\t6. Delete project")
	fmt.Println("\tPress enter to return")
	count := 9 + len(again)
	switch input := getInput(); input {
	case "1":
		clearLines(count)
		viewReferences(projectId)
	case "2":
		clearLines(count)
		addReference(projectId)
	case "3":
		clearLines(count)
		removeReference(projectId)
	case "4":
		clearLines(count)
		exportReferences(projectId)
	case "5":
		clearLines(count)
		editProject(projectId)
	case "6":
		clearLines(count)
		deleteProject(projectId)
	case "":
		clearLines(count)
		selectProject()
	default:
		clearLines(count)
		fmt.Println("Invalid input")
		projectMenu(projectId, true)
	}
}

func viewReferences(projectId string) {
	fmt.Println("References for project:", projects[projectId].title)
	for id, reference := range projects[projectId].references {
		fmt.Println("\t"+id, ":", reference.title)
	}
	fmt.Println("Press enter to go back")
	getInput()
	clearLines(len(projects[projectId].references) + 3)
	projectMenu(projectId)
}

func addReference(projectId string, again ...bool) {
	fmt.Println("Select a reference to add")
	count := 4 + len(again)
	for id, reference := range references {
		if _, ok := projects[projectId].references[id]; !ok {
			fmt.Println("\t"+id, ":", reference.title)
			count++
		}
	}
	fmt.Println("\t0 : New reference")
	fmt.Println("\tPress enter to return")
		switch input := getInput(); input {
	case "0":
		clearLines(count)
		ref := newReference()
		projects[projectId].references[ref.id] = ref
		saveProjects()
		projectMenu(projectId)
	case "":
		clearLines(count)
		projectMenu(projectId)
	default:
		if _, ok := references[input]; ok {
			// Check if the reference is already in the project 
			if _, ok := projects[projectId].references[input]; !ok {
				clearLines(count)
				projects[projectId].references[input] = references[input]
				saveProjects()
				fmt.Println("Reference added")
				fmt.Println("Press enter to return")
				getInput()
				clearLines(3)
				projectMenu(projectId)
			} else {
				clearLines(count)
				fmt.Println("Reference already in project")
				addReference(projectId, true)
			}
		} else {
			clearLines(count)
			fmt.Println("Invalid input")
			addReference(projectId, true)
		}
	}
}

func removeReference(projectId string, again ...bool) {
	fmt.Println("Select a reference to remove")
	for id, reference := range projects[projectId].references {
		fmt.Println("\t" + id, ":", reference.title)
	}
	fmt.Println("\tPress enter to return")
	count := len(projects[projectId].references) + 3 + len(again)
	switch input := getInput(); input {
	case "":
		clearLines(count)
		projectMenu(projectId)
	default:
		if _, ok := projects[projectId].references[input]; ok {
			clearLines(count)
			delete(projects[projectId].references, input)
			saveProjects()
			fmt.Println("Reference removed")
			fmt.Println("Press enter to return")
			getInput()
			clearLines(3)
			projectMenu(projectId)
		} else {
			clearLines(count)
			fmt.Println("Invalid input")
			removeReference(projectId, true)
		}
	}
}

func newReference() Reference {
	fmt.Println("Enter title:")
	title := getInput()
	fmt.Println("Enter year:")
	year := getInput()
	fmt.Println("Enter month:")
	month := getInput()
	fmt.Println("Enter publisher:")
	publisher := getInput()
	fmt.Println("Enter url:")
	url := getInput()
	var authors []string
	count := 4 
	for true {
		fmt.Println("Enter author: (Press enter to finish)")
		author := getInput()
		if len(authors) == 0 && author == "" {
			fmt.Println("Please enter at least one author")
			count++
			continue
		} else if author == "" {
			break
		}
		authors = append(authors, author)
	}
	// Remove spaces from id
	id := authors[0] + year
	id = strings.Replace(id, " ", "", -1)
	if _, ok := references[id]; ok {
		id += "a"
	}
	references[id] = Reference{
		id: id,
		reftype: "article",
		title: title,
		year: year,
		month: month,
		url: url,
		authors: authors,
		publisher: publisher,
	}
	saveReferences()
	clearLines(count + len(authors))
	return references[id]
}

func saveProjects() {
	file, err := os.OpenFile("projects.csv", os.O_WRONLY|os.O_TRUNC, 0644)
	handleErr(err)
	defer file.Close()
	writer := csv.NewWriter(file)
	err = writer.Write([]string{"id", "title", "references"})
	for _, project := range projects {
		line := []string{project.id, project.title}
		for _, reference := range project.references {
			line = append(line, reference.id)
		}
		err := writer.Write(line)
		handleErr(err)
	}
	writer.Flush()
}

func saveReferences() {
	// Save references to existing csv 
	file, err := os.OpenFile("references.csv", os.O_WRONLY|os.O_TRUNC, 0644)
	handleErr(err)
	defer file.Close()
	writer := csv.NewWriter(file)
	err = writer.Write([]string{"id", "reftype", "title", "year", "month", "url","publisher", "authors"})
	for _, reference := range references {
		line := []string{reference.id, reference.reftype, reference.title, reference.year, reference.month, reference.url, reference.publisher}
		line = append(line, reference.authors...)
		err := writer.Write(line)
		handleErr(err)
	}
	writer.Flush()
}

func selectReference(again ...bool) {
	fmt.Println("Select a reference")
	count := len(references) + 3 + len(again)
	for id, reference := range references {
		fmt.Println("\t"+id, ":", reference.title)
	}
	fmt.Println("\tPress enter to return")
	switch input := getInput(); input {
	case "":
		clearLines(count)
		menu()
	default:
		if _, ok := references[input]; ok {
			clearLines(count)
			referenceMenu(input)
		} else {
			clearLines(count)
			fmt.Println("Invalid input")
			selectReference(true)
		}
	}
}

func referenceMenu(referenceId string, again ...bool) {
	fmt.Println("Selected reference:", references[referenceId].title)
	fmt.Println("\t1. View reference")
	fmt.Println("\t2. Edit reference")
	fmt.Println("\t3. Delete reference")
	fmt.Println("\tPress enter to return")
	count := 6 + len(again)
	switch input := getInput(); input {
	case "1":
		clearLines(count)
		viewReference(referenceId)
	case "2":
		clearLines(count)
		editReference(referenceId)
	case "3":
		clearLines(count)
		deleteReference(referenceId)
	case "":
		clearLines(count)
		selectReference()	
	default:
		clearLines(count)
		fmt.Println("Invalid input")
		referenceMenu(referenceId, true)
	}
}

func editReference(referenceId string, again ...bool) {
	fmt.Println("Editing reference:", references[referenceId].title)
	fmt.Println("\t1. Edit title")
	fmt.Println("\t2. Edit year")
	fmt.Println("\t3. Edit month")
	fmt.Println("\t4. Edit url")
	fmt.Println("\t5. Edit authors")
	fmt.Println("\t6. Edit publisher")
	fmt.Println("\tPress enter to return")
	count := 8 + len(again)
	ref := references[referenceId]
	switch input := getInput(); input {
	case "1":
		clearLines(count)
		fmt.Println("Enter new title")
		ref.title = getInput()
		references[referenceId] = ref
		saveReferences()
		clearLines(2)
		fmt.Println("Title changed")
		fmt.Println("Press enter to continue")
		getInput()
		clearLines(4)
		editReference(referenceId)
	case "2":
		clearLines(count)
		fmt.Println("Enter new year")
		ref.year = getInput()
		references[referenceId] = ref
		saveReferences()
		clearLines(2)
		fmt.Println("Year changed")
		fmt.Println("Press enter to continue")
		getInput()
		clearLines(4)
		editReference(referenceId)
	case "3":
		clearLines(count)
		fmt.Println("Enter new month")
		ref.month = getInput()
		references[referenceId] = ref
		saveReferences()
		clearLines(2)
		fmt.Println("Month changed")
		fmt.Println("Press enter to continue")
		getInput()
		clearLines(4)
		editReference(referenceId)
	case "4":
		clearLines(count)
		fmt.Println("Enter new url")
		ref.url = getInput()
		references[referenceId] = ref
		saveReferences()
		clearLines(2)
		fmt.Println("Url changed")
		fmt.Println("Press enter to continue")
		getInput()
		clearLines(4)
		editReference(referenceId)
	case "5":
		clearLines(count)
		fmt.Println("\t1. Add author")
		fmt.Println("\t2. Remove author")
		fmt.Println("\tPress enter to return")
		switch input := getInput(); input {
		case "1":
			clearLines(4)
			fmt.Println("Enter new author")
			ref.authors = append(ref.authors, getInput())
			references[referenceId] = ref
			saveReferences()
			clearLines(2)
			fmt.Println("Author added")
			fmt.Println("Press enter to continue")
			getInput()
			clearLines(4)
			editReference(referenceId)
		case "2":
			clearLines(4)
			fmt.Println("Enter author to remove")
			for i, author := range references[referenceId].authors {
				fmt.Println(i, ":", author)
			}
			fmt.Println("Press enter to return")
			switch input := getInput(); input {
			case "":
				clearLines(3)
				editReference(referenceId)
			default:
				if i, err := strconv.Atoi(input); err == nil {
					if i < len(references[referenceId].authors) {
						ref.authors = append(references[referenceId].authors[:i], references[referenceId].authors[i+1:]...)
						references[referenceId] = ref
						saveReferences()
						clearLines(len(references[referenceId].authors) + 4)
						fmt.Println("Author removed")
						fmt.Println("Press enter to continue")
						getInput()
						clearLines(4)
						editReference(referenceId)
					} else {
						clearLines(3)
						fmt.Println("Invalid input")
						editReference(referenceId)
					}
				} else {
					clearLines(3)
					fmt.Println("Invalid input")
					editReference(referenceId)
				}
			}
		}
	case "6":
		clearLines(count)
		fmt.Println("Enter new publisher")
		ref.publisher = getInput()
		references[referenceId] = ref
		saveReferences()
		clearLines(2)
		fmt.Println("Publisher changed")
		fmt.Println("Press enter to continue")
		getInput()
		clearLines(4)
		editReference(referenceId)
	case "":
		clearLines(9)
		referenceMenu(referenceId)
	default:
		clearLines(3)
		fmt.Println("Invalid input")
		editReference(referenceId)
	}
}

func deleteReference(referenceId string, again ...bool) {
	fmt.Println("Deleting reference:", references[referenceId].title)
	fmt.Println("Are you sure you want to delete this reference? (y/n)")
	count := 3 + len(again)
	switch input := getInput(); input {
	case "y":
		clearLines(count)
		for _, project := range projects {
			delete(project.references, referenceId)
		}
		delete(references, referenceId)
		saveProjects()
		saveReferences()
		fmt.Println("Reference deleted")
		fmt.Println("Press enter to continue")
		getInput()
		clearLines(2)
		selectReference()
	case "n":
		clearLines(count)
		referenceMenu(referenceId)
	default:
		clearLines(count)
		fmt.Println("Invalid input")
		deleteReference(referenceId, true)
	}
}


func exportReferences(projectId string) {
	output := ""
	fmt.Println("Exporting references for project:", projects[projectId].title)
	count := 5 + len(projects[projectId].references) * 8
	fmt.Print("\n")
	for _, reference := range projects[projectId].references {
		// Export with biblatex
		referenceId := reference.id
		output += fmt.Sprintf("@%s{%s,\n", references[referenceId].reftype, referenceId)
		output += fmt.Sprintf("\ttitle = {%s},\n", references[referenceId].title)
		output += fmt.Sprintf("\tauthor = {%s},\n", strings.Join(references[referenceId].authors, " and "))
		output += fmt.Sprintf("\tyear = {%s},\n", references[referenceId].year)
		output += fmt.Sprintf("\tmonth = {%s},\n", references[referenceId].month)
		output += fmt.Sprintf("\turl = {%s},\n", references[referenceId].url)
		output += fmt.Sprintf("\tpublisher = {%s},\n", references[referenceId].publisher)
		output += "}\n"
	}
	fmt.Println(output)
	fmt.Println("Press enter to continue")
	getInput()
	clearLines(count)
	projectMenu(projectId)
}

func viewReference(referenceId string) {
	fmt.Println("Title:", references[referenceId].title)
	fmt.Println("Authors:", strings.Join(references[referenceId].authors, ", "))
	fmt.Println("Publisher:", references[referenceId].publisher)
	fmt.Println("Year:", references[referenceId].year)
	fmt.Println("Month:", references[referenceId].month)
	fmt.Println("URL:", references[referenceId].url)
	fmt.Println("Press enter to continue")
	getInput()
	clearLines(8)
	referenceMenu(referenceId)
}

func createProject() {
	fmt.Println("Enter project title:")
	title := getInput()
	fmt.Println("Enter id:")
	id := getInput()
	if _, ok := projects[id]; ok {
		clearLines(4)
		fmt.Println("Project with this id already exists")
	}else
	{
		projects[id] = Project{
			id: id,
			title: title,
			references: make(map[string]Reference),
		}
		saveProjects()
		clearLines(4)
		fmt.Println("Project created")
	}
	fmt.Println("Press enter to continue")
	getInput()
	clearLines(3)
	selectProject()
}

func editProject(projectId string) {
	fmt.Println("Editing project:", projects[projectId].title)
	fmt.Println("1. Edit title")
	fmt.Println("2. Edit id")
	fmt.Println("Press enter to return")
	project := projects[projectId]
	switch input := getInput(); input {
	case "1":
		clearLines(4)
		fmt.Println("Enter new title")
		project.title = getInput()
		projects[projectId] = project
		saveProjects()
		clearLines(2)
		fmt.Println("Title changed")
		fmt.Println("Press enter to continue")
		getInput()
		clearLines(4)
		editProject(projectId)
	case "2":
		clearLines(4)
		fmt.Println("Enter new id")
		newId := getInput()
		if _, ok := projects[newId]; ok {
			clearLines(2)
			fmt.Println("Project with this id already exists")
			fmt.Println("Press enter to continue")
			getInput()
			clearLines(4)
			editProject(projectId)
		}else
		{
			project.id = newId
			projects[newId] = project
			delete(projects, projectId)
		}
		saveProjects()
		clearLines(2)
		fmt.Println("Id changed")
		fmt.Println("Press enter to continue")
		getInput()
		clearLines(4)
		editProject(newId)
	default:
		clearLines(4)
		projectMenu(projectId)
	}
}

func deleteProject(projectId string, again ...bool) {
	fmt.Println("Are you sure you want to delete project", projects[projectId].title, "?")
	fmt.Println("y/n")
	count := 3 + len(again)
	switch input := getInput(); input {
	case "y":
		clearLines(count)
		delete(projects, projectId)
		saveProjects()
		fmt.Println("Project deleted")
		fmt.Println("Press enter to continue")
		getInput()
		clearLines(2)
		selectProject()
	case "n":
		clearLines(count)
		projectMenu(projectId)
	default:
		clearLines(count)
		fmt.Println("Invalid input")
		deleteProject(projectId,true)
	}
}


// HELPER FUNCTIONS


func clearLines(count int) bool {
	for i := 0; i < count; i++ {
		fmt.Print("\033[2K\033[A")
	}
	fmt.Print("\033[2K")
	return true
}

func handleErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func getInput() string {
	in := bufio.NewScanner(os.Stdin)
	in.Scan()
	input := in.Text()
	return input
}

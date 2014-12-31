package main

import(
  "bytes"
  "encoding/json"
  "fmt"
  "net/http"
)

const redmineTimeFormat = "2006-01-02"

func uploadReportItems(reportItems []ReportItem) {
  host := getHost("REDMINE_HOST", "Enter your redmine host")
  apiKey := getVarOrInput("REDMINE_TOKEN", "Enter your api token")

  projectId := getProject(host, apiKey)
  postUrl := fmt.Sprintf("%vtime_entries.json", host)
  for _, reportItem := range reportItems {
    body := map[string]map[string]interface{}{
      "time_entry": map[string]interface{} {
        "project_id": projectId,
        "spent_on": reportItem.Date.Format(redmineTimeFormat),
        "hours": reportItem.Hours,
        "comments": reportItem.Description,
      },
    }
    jsonBody, err := json.Marshal(body)
    if err != nil { panic(err) }

    request, err := http.NewRequest("POST", postUrl, bytes.NewBuffer(jsonBody))
    if err != nil { panic(err) }
    request.Header.Set("X-Redmine-API-Key", apiKey)
    request.Header.Set("Content-Type", "application/json")
    request.Header.Set("Accept", "application/json")

    client := &http.Client{}
    _, err = client.Do(request)
    if err != nil { panic(err) }
    fmt.Print(".")
  }
  fmt.Println("")
}

func getProject(host, apiKey string) (result int) {
  fullUrl := fmt.Sprintf("%vprojects.json", host)
  request, err := http.NewRequest("GET", fullUrl, nil)
  if err != nil { panic(err) }

  request.Header.Set("X-Redmine-API-Key", apiKey)

  var jsonOutput map[string]interface{}

  json.Unmarshal(getResponseBody(request), &jsonOutput)

  projects := map[int]string{}
  jsonProjects := jsonOutput["projects"].([]interface{})
  for _, jsonProject := range jsonProjects {
    project := jsonProject.(map[string]interface{})
    projectId := int(project["id"].(float64))
    projects[projectId] = project["name"].(string)
  }

  fmt.Println("Please select project:")
  for projectId, projectName := range projects {
    fmt.Printf("%v - %v\n", projectName, projectId)
  }

  fmt.Scanf("%d", &result)
  return result
}

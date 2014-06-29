package main

import(
  "bytes"
  "encoding/json"
  "fmt"
  "net/http"
)

const redmineTimeFormat = "2006/01/02"

func uploadReportItems(reportItems []ReportItem) {
  host := getHost("Enter your redmine host")
  var apiKey string
  fmt.Println("Enter your api token")
  fmt.Scanf("%s", &apiKey)

  fmt.Println("Please select project:")
  for projectId, projectName := range getProjects(host, apiKey) {
    fmt.Printf("%v - %v\n", projectName, projectId)
  }
  var projectId int
  fmt.Scanf("%d", &projectId)

  client := &http.Client{}
  postUrl := fmt.Sprintf("%vtime_entries", host)
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

    _, err = client.Do(request)
    if err != nil { panic(err) }
    fmt.Print('.')
  }
}

func getProjects(host, apiKey string) map[int]string {
  client := &http.Client{}
  fullUrl := fmt.Sprintf("%vprojects.json", host)
  request, err := http.NewRequest("GET", fullUrl, nil)
  if err != nil { panic(err) }

  request.Header.Set("X-Redmine-API-Key", apiKey)
  response, err := client.Do(request)
  if err != nil { panic(err) }

  var jsonOutput map[string]interface{}

  buf := new(bytes.Buffer)
  buf.ReadFrom(response.Body)
  json.Unmarshal(buf.Bytes(), &jsonOutput)

  projects := map[int]string{}
  jsonProjects := jsonOutput["projects"].([]interface{})
  for _, jsonProject := range jsonProjects {
    project := jsonProject.(map[string]interface{})
    projectId := int(project["id"].(float64))
    projects[projectId] = project["name"].(string)
  }
  return projects
}

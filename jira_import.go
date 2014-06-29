package main

import (
  "github.com/moovweb/gokogiri"
  "fmt"
  "net/http"
  "strconv"
  "time"
)

const (
  jiraTimeLayout = "02/Jan/06"
  htmlTimeLayout = "02/Jan/06 15:04"
  jiraUrlLayout = "%vsecure/TempoUserBoard!reportPrint.jspa?v=1&periodType=FLEX&periodView=DATES&from=%v&to=%v"
)

type ReportItem struct {
  Description string
  Hours float64
  Date time.Time
}

func getReportItems() []ReportItem {
  request, err := http.NewRequest("GET", getJiraUrl(), nil)
  if err != nil { panic(err) }

  login := getString("Enter your jira login")
  password := getString("Enter your jira password")
  request.SetBasicAuth(login, password)

  fmt.Println("Loading your reports from jira...")
  doc, err := gokogiri.ParseHtml(getResponseBody(request))
  if err != nil { panic("Failed to proccess page") }

  nodes, err := doc.Root().Search(toXpath(".tempo-table .level-3"))
  if err != nil { panic("Invalid html document") }

  reportItems := []ReportItem{}
  for _, node := range nodes {
    parentId := fmt.Sprintf("section-%v", node.Attr("rel"))
    parentNode := doc.NodeById(parentId)

    issueName := stripString(getContent(parentNode, "td.summary a"))
    issueDetails := stripString(getContent(node,"td.summary span.tempo-inline-edit"))
    comments := fmt.Sprintf("%v (%v)", issueName, issueDetails)

    issueTime, err := time.Parse(htmlTimeLayout, getContent(node, "td span.tempo-calendar-display"))
    if err != nil { panic(err) }

    hours, _ := strconv.ParseFloat(getContent(node, "td.sum span.tempo-inline-edit"), 64)

    reportItems = append(reportItems, ReportItem{comments, hours, issueTime})
  }

  return reportItems
}

func getJiraUrl() string {
  host := getHost("Hi! For starting trasfer input hostname of your jira")
  var startDate time.Time
  var endDate time.Time

  for {
    startDate = getTime("Enter a first day", getBegginingOfMonth())
    endDate = getTime("Enter a last day", time.Now())

    if endDate.After(startDate) { break }

    fmt.Println("Start time is greater than finish time")
  }
  return fmt.Sprintf(jiraUrlLayout, host, startDate.Format(jiraTimeLayout), endDate.Format(jiraTimeLayout))
}

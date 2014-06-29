package main

import (
  "bytes"
  "github.com/moovweb/gokogiri"
  gcss "github.com/moovweb/gokogiri/css"
  "fmt"
  "net/http"
  "regexp"
  "strconv"
  "strings"
  "time"
)

const (
  terminalTimeLayout = "2/1/06"
  jiraTimeLayout = "02/Jan/06"
  htmlTimeLayout = "02/Jan/06 15:04"
  jiraUrlLayout = "%vsecure/TempoUserBoard!reportPrint.jspa?v=1&periodType=FLEX&periodView=DATES&from=%v&to=%v"
)
var blankRegex = regexp.MustCompile("[\t\n ]+")

type ReportItem struct {
  Description string
  Hours float64
  Date time.Time
}

func getReportItems() []ReportItem {
  client := &http.Client{}
  request, err := http.NewRequest("GET", getJiraUrl(), nil)
  if err != nil { panic(err) }

  var login string
  fmt.Println("Enter your jira login:")
  fmt.Scanf("%s", &login)
  var password string
  fmt.Println("Enter your jira password:")
  fmt.Scanf("%s", &password)
  request.SetBasicAuth(login, password)

  fmt.Println("Loading your reports from jira...")
  response, err := client.Do(request)
  if err != nil { panic(err) }

  buf := new(bytes.Buffer)
  buf.ReadFrom(response.Body)
  body := buf.Bytes()

  doc, err := gokogiri.ParseHtml(body)
  if err != nil {
    panic("Failed to proccess page")
  }

  xpath := gcss.Convert(".tempo-table .level-3", gcss.GLOBAL)
  nodes, err := doc.Root().Search(xpath)
  if err != nil {
    panic("Failed to search an issue records")
  }

  reportItems := []ReportItem{}
  for _, node := range nodes {
    parentId := fmt.Sprintf("section-%v", node.Attr("rel"))
    parentNode := doc.NodeById(parentId)
    variants, _ := parentNode.Search(gcss.Convert("td.summary a", gcss.LOCAL))
    issueName := stripString(variants[0].Content())
    variants, _ = node.Search(gcss.Convert("td.summary span.tempo-inline-edit", gcss.LOCAL))
    issueDetails := stripString(variants[0].Content())
    description := fmt.Sprintf("%v (%v)", issueName, issueDetails)

    variants, _ = node.Search(gcss.Convert("td span.tempo-calendar-display", gcss.LOCAL))
    issueTime, _ := time.Parse(htmlTimeLayout, variants[0].Content())

    variants, _ = node.Search(gcss.Convert("td.sum span.tempo-inline-edit", gcss.LOCAL))
    hours, _ := strconv.ParseFloat(variants[0].Content(), 64)
    reportItems = append(reportItems, ReportItem{description, hours, issueTime})
  }

  return reportItems
}

func stripString(input string) string {
  return strings.TrimSpace(blankRegex.ReplaceAllLiteralString(input, " "))
}

func getHost(prompt string) string {
  fmt.Println(prompt)
  var host string
  fmt.Scanf("%s", &host)
  if !strings.HasSuffix(host, "/") { host = host + "/" }
  return host
}

func getJiraUrl() string {
  host := getHost("Hi! For starting trasfer input hostname of your jira")
  var startDate time.Time
  var endDate time.Time

  for {
    startDate = time.Now()
    startDate = startDate.AddDate(0, 0, (1 - startDate.Day()))
    startDate = getTime("Enter a first day", startDate)

    endDate = time.Now()
    endDate = getTime("Enter a last day", endDate)

    if endDate.After(startDate) { break }

    fmt.Println("Start time is greater than finish time")
  }
  return fmt.Sprintf(jiraUrlLayout, host, startDate.Format(jiraTimeLayout), endDate.Format(jiraTimeLayout))
}

func getTime(prompt string, defTime time.Time) time.Time {
  var rawInput string

  fmt.Println(fmt.Sprintf("%v [default is %v]", prompt, defTime.Format(terminalTimeLayout)))
  for {
    fmt.Scanf("%s", &rawInput)
    if len(rawInput) == 0 {
      return defTime
    } else {
      result, err := time.Parse(terminalTimeLayout, rawInput)
      if err == nil { return result }
    }
    fmt.Println("Invalid input. Please try again")
  }
}

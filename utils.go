package main

import(
  "bytes"
  "fmt"
  "net/http"
  gcss "github.com/moovweb/gokogiri/css"
  gxtml "github.com/moovweb/gokogiri/xml"
  "os"
  "regexp"
  "strings"
  "time"
)

const terminalTimeLayout = "2/1/06"
var blankRegex = regexp.MustCompile("[\t\n ]+")

func getBegginingOfMonth() time.Time {
  now := time.Now()
  return time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
}

func getContent(node gxtml.Node, cssQuery string) string {
  result, err := node.Search(toXpath(cssQuery))
  if err != nil { panic(fmt.Errorf("Failed to find %v node", cssQuery)) }
  return result[0].Content()
}

func getHost(varName string, prompt string) (host string) {
  host = getVarOrInput(varName, prompt)
  if !strings.HasSuffix(host, "/") { host = host + "/" }
  if !strings.HasPrefix(host, "http") { host = "http://" + host }
  return host
}

func getResponseBody(request *http.Request) []byte {
  client := &http.Client{}
  response, err := client.Do(request)
  if err != nil { panic(err) }

  buf := new(bytes.Buffer)
  buf.ReadFrom(response.Body)
  return buf.Bytes()
}

func getTime(prompt string, defTime time.Time) time.Time {
  prompt = fmt.Sprintf("%v [default is %v]", prompt, defTime.Format(terminalTimeLayout))
  for {
    rawInput := getString(prompt)
    if len(rawInput) == 0 {
      return defTime
    } else {
      result, err := time.Parse(terminalTimeLayout, rawInput)
      if err == nil { return result }
    }
    fmt.Println("Invalid input. Please try again")
  }
}

func getVarOrInput(varName string, prompt string) (result string) {
  result = os.Getenv(varName)
  if result == "" {
    result = getString(prompt)
  }
  return result
}

func getString(prompt string) (result string) {
  fmt.Printf("%v:\n", prompt)
  fmt.Scanf("%s", &result)
  return result
}

func logger(key interface{}, message interface{}) {
  fmt.Printf("[%v] %v\n", key, message)
}

func stripString(input string) string {
  return strings.TrimSpace(blankRegex.ReplaceAllLiteralString(input, " "))
}

func toXpath(cssQuery string) string {
  return gcss.Convert(cssQuery, gcss.LOCAL)
}

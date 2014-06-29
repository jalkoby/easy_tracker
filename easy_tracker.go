package main

import "fmt"

func main() {
  uploadReportItems(getReportItems())
}

func logger(key string, message interface{}) {
  fmt.Printf("[%v] %v\n", key, message)
}

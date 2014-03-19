/*
Example rules files
[
	{"Referer":"github.com/blog", "Process":"echo 'hi'", "Interval":600}
]

*/


package main

import (
       "bytes"
       "encoding/json"
       "flag"
       "fmt"
       "log"
       "net/http"
       "os"
       "os/exec"
       "io/ioutil"
)

// Rule provided by configuration file
type Rule struct {
     Endpoint string
     Process interface{}	// Process to execute (points to shell script)
     Interval int
}

var (
    ruleFile = flag.String("rules", "", "rule definition file")
    rules []*Rule
)

func handler(w http.ResponseWriter, r *http.Request) {
     // Check for X-Github-Event
     eventType := r.Header.Get("X-Github-Event")
     if eventType != "" {
     	// Process Event
	var results interface{}
     	hah, _ := ioutil.ReadAll(r.Body)
  
	err := json.Unmarshal(hah, &results);
     	if err != nil {
	}
	repository := results.(map[string]interface{})["repository"]
	url := repository.(map[string]interface{})["url"]
	 for _, r := range rules {
	     if r.Endpoint == url {
     	     	log.Println("Found rule for url:", url)
		array := r.Process.([]interface{})
		b := make([]string, len(array))
		for i := range array {
		    b[i] = array[i].(string)
		}
		slice := b[1:len(array)]
		cmd := exec.Command(b[0], slice...)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Run()
		fmt.Printf(out.String())
	     }
     	 }
     }
}

func loadRules(file string) error {
     log.Printf("Loading rules")
     f, err := os.Open(file)
     if err != nil {
     	log.Printf("Error loading rules file")
     	return err
     }
     defer f.Close()
     err = json.NewDecoder(f).Decode(&rules)
     if err != nil {
     	fmt.Println(err)
     	return nil
     }
     for _, r := range rules {
     	 log.Println("Loaded rule for:", r.Endpoint)
     }
     return nil 
}

func main() {
     flag.Parse()
     fmt.Println("Starting Webhook Server on port 8118")
     loadRules(*ruleFile)
     http.HandleFunc("/webhook", handler)
     http.ListenAndServe(":8118", nil)
}
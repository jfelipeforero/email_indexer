package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"	
	"strings"
	"time"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
  "github.com/go-chi/cors"
)

func main(){    
  r := chi.NewRouter()

  r.Use(middleware.Logger)
  r.Use(middleware.RealIP) 
  r.Use(middleware.Recoverer)
  // Set up CORS middleware
  corsMiddleware := cors.New(cors.Options{
    AllowedOrigins: []string{"*"},
    AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
  })
  r.Use(corsMiddleware.Handler)

  r.Use(middleware.Timeout(60 * time.Second))

  r.Get("/search",func(w http.ResponseWriter, r *http.Request){
    query := r.URL.Query().Get("q") 
    search := search(query)
    jsonBytes, err := json.Marshal(search)
    if err!= nil {
      log.Panic(err)
    }
    w.Write(jsonBytes)
  })

  http.ListenAndServe(":3000", r)
}

// This function sends a request to the ZincSearch API with the query 
// provided as argument and returns an array of interfaces
func search(term string) []interface{} {   
  // We define the query, formatting the string of the query, so that it 
  // receives a query and 
  query := fmt.Sprintf(`{
    "search_type": "match",
    "query":
    {
      "term": "%s",
      "start_time": "2022-06-02T14:28:31.894Z",
      "end_time": "2023-12-02T15:28:31.894Z"
    },
    "from": 0,
    "max_results": 50,
    "_source": ["Subject","From","To","Message"]
  }`, term) 

  req, err := http.NewRequest("POST", "http://localhost:4080/api/emails/_search", strings.NewReader(query)) 
  if err != nil {
    log.Fatal(err)
  }
  req.SetBasicAuth("admin", "password")
  req.Header.Set("Content-Type", "application/json") 

  resp, err := http.DefaultClient.Do(req)
  if err != nil {
    log.Fatal(err)
  } 
  defer resp.Body.Close()
  log.Println(resp.StatusCode)
  body, err := io.ReadAll(resp.Body)
  if err != nil {
    log.Fatal(err)
  } 

  return filter(body)  
}

// Filter the emails received by the ZincSearch response and
// return an array of emails
func filter(respBody []byte) []interface{} {      
  results := SplitEmails(respBody)

  var searches []interface{} 
  for i := 0; i < len(results); i++ {
    var search []interface{}
    var data map[string]interface{} 
    err := json.Unmarshal([]byte(results[i]), &data) 
    if err!=nil{
      log.Fatal(err)
    }     
    fields := data["_source"]  
    subject := fields.(map[string]interface{})["Subject"]
    from := fields.(map[string]interface{})["From"]
    to := fields.(map[string]interface{})["To"] 
    message := fields.(map[string]interface{})["Message"]

    search = append(search,subject)
    search = append(search,from)
    search = append(search,to)
    search = append(search,message)
    searches  = append(searches,search)  

  } 
  return searches 
}

// SplitEmails separetes every email from each other storing
// them in a new strings array
func SplitEmails (respBody []byte) []string  { 
  // Removes all characters up to the first email
  removeFields := strings.SplitAfter(string(respBody),`"hits":[`)  
  removeFields[1] = removeFields[1][:len(removeFields[1])-3]  
  emailsSlice := removeFields[1]
  var emails []string

  finished := false 
  for finished != true {   
    pos := strings.Index(emailsSlice[2:],`,{"_index":"emails"`)       
    if pos != -1 { 
      emails = append(emails, emailsSlice[:pos+2]) 
      emailsSlice = emailsSlice[pos+3:]   
    }else { 
      emails = append(emails, emailsSlice)   
      finished = true
      }
    } 
  return emails
}

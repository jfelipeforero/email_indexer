package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	//"reflect"
	"strings"
	"time"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
  "github.com/go-chi/cors"
)

var routes = flag.Bool("routes", false, "Generate router documentation")

func splitRec (slice string, emails []string) []string  { 
  finished := false 
  for finished != true {   
    pos := strings.Index(slice[2:],`,{"_index":"emails"`)       
    if pos != -1 { 
      emails = append(emails, slice[:pos+2]) 
      slice = slice[pos+3:]   
    }else { 
      emails = append(emails, slice)   
      finished = true
      }
    } 
  return emails
}

func filter(respBody []byte) []interface{} {     
  first := strings.SplitAfter(string(respBody),`"hits":[`)  
  first[1] = first[1][:len(first[1])-3] 

  var emails []string

  results := splitRec(first[1],emails)
 
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
  fmt.Println(searches)

  return searches 
}

func search(term string) []interface{} {  
  fmt.Println(term)
  query := fmt.Sprintf(`{
        "search_type": "match",
        "query":
        {
            "term": "%s",
            "start_time": "2022-06-02T14:28:31.894Z",
            "end_time": "2023-12-02T15:28:31.894Z"
        },
        "from": 0,
        "max_results": 40,
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

func main(){
    //flag.Parse()  
    url := "http://localhost:4080/api/index"

    r := chi.NewRouter()
    
    r.Use(middleware.Logger)
    r.Use(middleware.RealIP)
    r.Use(middleware.Logger)
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
      fmt.Println("Are we even here!?")
      w.Write(jsonBytes)
    })
    
    r.Get("/hola", func(w http.ResponseWriter, r *http.Request){
        w.Write([]byte("hola"))
        req, err := http.NewRequest("GET",url,nil)
          if err != nil {
            log.Fatal(err) 
          }    

          req.SetBasicAuth("admin", "password")
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

          fmt.Println(string(body))
    })

    

    http.ListenAndServe(":3000", r)
}

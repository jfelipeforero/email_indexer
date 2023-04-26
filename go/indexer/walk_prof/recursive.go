package main

import (
	"bufio"
    "runtime/pprof"
    "encoding/json"	
	"fmt"	
	"log"
	"os"
	"path/filepath"
	"strings"
)

func check(e error){
    if e != nil {
        panic(e)
    }
}

func main () {
    args := os.Args
    f, err := os.Create("profile.pb.gz")
    if err != nil {
        log.Fatal(err)
    }
    err = pprof.StartCPUProfile(f)
        if err != nil {
             log.Fatal(err)
    }
    defer pprof.StopCPUProfile()
    walk(args[1]) 
}

func walk (path string) {
    // It creates the file where the emails will be stored as json
    emails, err := os.Create("emails.ndjson")
    check(err)
    writer := bufio.NewWriter(emails) 
    filepath.Walk(path, func(path string, info os.FileInfo, err error) error { 
    check(err)

    fmt.Println("Indexing data...")
    if info.IsDir() != true && filepath.Ext(path) != ".txt"{
        file, err := os.Open(path)
        if err != nil {
            log.Fatal(err)
        } 
        scanner := bufio.NewScanner(file)
  
        // Following the standard that I detected in the olympics.ndjson in the zinc example, I write the filepath
        // so every email is preceded by a default index
        jsonIndex := `{ "index" : { "_index" : "emails" } }`+ "\n"
        _, err = writer.WriteString(string(jsonIndex))  
        check(err)
        writer.Flush()

        // This variable is declared to store the message inside the emails, the content itself
        var content string  
        var info string
        counter := 0

        for scanner.Scan(){  
            line := scanner.Text()
            index := strings.Index(line, ":") 
            if index != -1 && counter != 15 {
                result := line[:index]
                message := line[index+1:] 
                messageFmt, err := json.Marshal(message)
                check(err)
             
                rawMessage := string(fmt.Sprintf(`"%s": %s`, result, string(messageFmt)))
                info += rawMessage + ", "
                
                counter++

            }else { 
                content += "\n" + line
            }
        } 
        jsonMessage, err := json.Marshal(content)
        hola := string(jsonMessage) 

        rawMessage := string(fmt.Sprintf(`"%s": %s`, "Message",hola)) 
        final := info + rawMessage 
        sisi  := fmt.Sprintf(`{%s}`, final)
        check(err)
        data := string(sisi)  
        _, err = writer.WriteString(data + "\n")
        writer.Flush()
    }    
    return nil
 })
 fmt.Println("Emails indexed succesfully")
}

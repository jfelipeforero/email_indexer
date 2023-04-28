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
    // Receive the path of the directory to be traversed.
    args := os.Args
    // Profilling
    f, err := os.Create("profile.pb.gz")
    if err != nil {
        log.Fatal(err)
    }
    err = pprof.StartCPUProfile(f)
        if err != nil {
             log.Fatal(err)
    }
    defer pprof.StopCPUProfile()

    TraverseDirectory(args[1]) 
}

func TraverseDirectory (path string) {
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
  
        // Following the standard format that I detected in the olympics.ndjson in the ZincSearch example, I write the filepath
        // so every email is preceded by a default index
        jsonIndex := `{ "index" : { "_index" : "emails" } }`+ "\n"
        _, err = writer.WriteString(string(jsonIndex))  
        check(err)
        writer.Flush()

        // This variable is declared to store the message inside the emails, the content itself
        var content string  
        // info is used to store the information of the first 15 lines containing information 
        // about the email
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
        // This Marshal is needed to convert some of the characters in the email
        // into a valid json character
        jsonMessage, err := json.Marshal(content)
        formattedJson := string(jsonMessage) 

        // Define formatted string with a placeholder for the Message and another for the
        // content of the email 
        rawMessage := string(fmt.Sprintf(`"%s": %s`, "Message",formattedJson)) 
        
        infoMessage := info + rawMessage 
        // Wrap the result into a {} to make the entire email
        jsonEmail  := fmt.Sprintf(`{%s}`, infoMessage)
        check(err)
        data := string(jsonEmail)  
        _, err = writer.WriteString(data + "\n")
        writer.Flush()
    }    
    return nil
 })
 fmt.Println("Emails indexed succesfully")
}

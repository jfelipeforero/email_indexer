package main

import (	
    "fmt"
	"os"
)

func check(err error){
    if err != nil{
        panic("error")
    }
}
func main() {
    args := os.Args[1]
    walkRec(args)
}

func walkRec(path string){
    entries, err := os.ReadDir(path) 
    check(err)
    for _, entry := range entries {
        fmt.Println(entry.Name()) 
        if entry.IsDir(){
            walkRec(fmt.Sprintf(`%s+/%s`,path,entry.Name()))
        }
    }
}

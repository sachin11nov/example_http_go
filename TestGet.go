package main

import (
        "io"
        "log"
        "net/http"
        "os"
)

func main() {
        /*create a client*/
        client := &http.Client{}
        req, _ := http.NewRequest("GET", "http://tekbuds.com", nil)

        /*Set request headers if any*/
        req.Header.Set("Authorization", "Basic Y3BzaWFtOmY1M2IwY2RyYWI5ZTU0NTAwOGVhYWE=")

        /*make http(s) Call*/
        res, err := client.Do(req)

        /* Or make a simple http Get
        res, err := http.Get("http://tekbuds.com")
        */

        /*Print output to console*/
        if err != nil {
                log.Fatal(err)
        } else {
                defer res.Body.Close()
                _, err := io.Copy(os.Stdout, res.Body)
                if err != nil {
                        log.Fatal(err)
                }
        }
}

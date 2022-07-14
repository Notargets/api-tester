package main

import (
    "fmt"
    "github.com/notargets/api-tester/loadrequestor"
    "github.com/notargets/api-tester/webserver"
    "time"
)

func main() {
    fmt.Printf("starting web server at port 5050, CTRL-C or other signal to stop\n")
    go webserver.Start(5050)
    time.Sleep(time.Second)
    loadRequestor := loadrequestor.NewLoadRequestor([]string{
        "http://localhost:5050/countAccess",
        "http://localhost:5050/countAccess",
        "http://localhost:5050/countAccess",
    })
    loadRequestor.SubmitWorkLoop()
}

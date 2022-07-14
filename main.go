package main

import (
    "github.com/notargets/api-tester/loadrequestor"
)

func main() {
    // fmt.Printf("starting web server at port 5050\n")
    // go webserver.Start(5050)
    // time.Sleep(time.Second)
    loadRequestor := loadrequestor.NewLoadRequestor([]string{
        "http://localhost:3000/company/30/offices",
        "http://localhost:3000/company",
        "http://localhost:3000/v3/company/stats",
    }, map[string]string{"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOjMzLCJpYXQiOjE2NTc4MzQwNTAsImV4cCI6MTY1ODQzODg1MH0.surFlcIzr5U_m4dAlb3cNqB3FhViwrE2Y4DoXUPNgEI"})
    loadRequestor.SubmitWorkLoop(true)
}

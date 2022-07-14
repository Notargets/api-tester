package loadrequestor

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
)

type LoadRequestor struct {
    Requests []string // one request per Requestor, each request is a URL with parameters
}

func NewLoadRequestor(requests []string) (lr *LoadRequestor) {
    return &LoadRequestor{
        Requests: requests,
    }
}

func (lr *LoadRequestor) SubmitWorkLoop() {
    var (
        workerChan    = make(chan Response, len(lr.Requests))
        sigChan       = make(chan os.Signal)
        finishedCount int
        startTime     time.Time
        // Limit the rate to 10/second
        rate_limit    = 10
        rate_duration = time.Duration(float64(1000/rate_limit)) * time.Millisecond
    )
    signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
    // submit work initially for all requests
    for id, url := range lr.Requests {
        go WebRequest(id, url, workerChan)
    }
    startTime = time.Now()
    startLoop := time.Now()
    go func() {
        for {
            select {
            case result := <-workerChan:
                if time.Since(startTime) < rate_duration {
                    time.Sleep(rate_duration - time.Since(startTime))
                }
                finishedCount++
                id := result.ID
                // fmt.Printf("id: %d, result: [%s]\n", result.ID, result.Result)
                startTime = time.Now()
                go WebRequest(id, lr.Requests[id], workerChan)
            }
        }
    }()
    select {
    case <-sigChan:
        fmt.Printf("Got interrupt, exiting...\n")
        break
    }
    duration := time.Since(startLoop)
    rate := float64(finishedCount) / duration.Seconds()
    fmt.Printf("Rate is limited to %d requests per second\n", rate_limit)
    fmt.Printf("%d requests completed in %5.3f seconds, %5.3f per second, %5.3f per minute\n",
        finishedCount, duration.Seconds(), rate, rate*60.)
}

type Response struct {
    ID     int
    Result string
}

func WebRequest(requestID int, URL string, workerChan chan Response) {
    var (
        err  error
        resp *http.Response
        body []byte
    )

    // time.Sleep(time.Second)

    if resp, err = http.Get(URL); err != nil {
        workerChan <- Response{requestID, err.Error()}
    } else {
        defer func() { _ = resp.Body.Close() }()
        if body, err = ioutil.ReadAll(resp.Body); err != nil {
            panic(err)
        }
        workerChan <- Response{requestID, string(body)}
    }
}

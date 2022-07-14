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
    Headers  map[string]string
}

func NewLoadRequestor(requests []string, headersA ...map[string]string) (lr *LoadRequestor) {
    lr = &LoadRequestor{
        Requests: requests,
    }
    if len(headersA) != 0 {
        lr.Headers = headersA[0]
    }
    return
}

func (lr *LoadRequestor) SubmitWorkLoop(printResultsA ...bool) {
    var (
        workerChan    = make(chan Response, len(lr.Requests))
        sigChan       = make(chan os.Signal)
        finishedCount int
        startTime     time.Time
        // Limit the rate to 10/second
        rate_limit    = 10
        rate_duration = time.Duration(float64(1000/rate_limit)) * time.Millisecond
        printResults  bool
        results       = make(map[string]string)
    )
    if len(printResultsA) != 0 {
        printResults = printResultsA[0]
    }
    signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
    fmt.Printf("starting worker loop, CTRL-C or other signal to stop\n")
    // submit work initially for all requests
    for id, url := range lr.Requests {
        go WebRequest(id, url, lr.Headers, workerChan)
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
                if printResults {
                    URL := lr.Requests[id]
                    if val, ok := results[URL]; ok {
                        if val != result.Result {
                            fmt.Printf("Differing result[%d]:\n[%s]\n", result.ID, result.Result)
                            results[URL] = result.Result
                        }
                    } else {
                        fmt.Printf("New result[%d]:\n[%s]\n", result.ID, result.Result)
                        results[URL] = result.Result
                    }
                }
                startTime = time.Now()
                go WebRequest(id, lr.Requests[id], lr.Headers, workerChan)
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

func WebRequest(requestID int, URL string, Headers map[string]string, workerChan chan Response) {
    var (
        err  error
        resp *http.Response
        body []byte
    )
    client := &http.Client{}
    req, _ := http.NewRequest("GET", URL, nil)
    if Headers != nil {
        for key, val := range Headers {
            req.Header.Set(key, val)
        }
    }
    if resp, err = client.Do(req); err != nil {
        workerChan <- Response{requestID, err.Error()}
    } else {
        defer func() { _ = resp.Body.Close() }()
        if body, err = ioutil.ReadAll(resp.Body); err != nil {
            panic(err)
        }
        workerChan <- Response{requestID, string(body)}
    }
}

package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "strconv"
    "sync"
    "time"
        "runtime"
)

// Fetch photos handler
func fetchPhotosHandler(w http.ResponseWriter, r *http.Request) {
    nbRequestsStr := r.URL.Query().Get("nbRequests")
    nbRequests, err := strconv.Atoi(nbRequestsStr)
    if err != nil || nbRequests <= 0 {
        http.Error(w, "Invalid number of requests. Please provide a positive integer.", http.StatusBadRequest)
        return
    }

    // Set maximum number of threads
    runtime.GOMAXPROCS(runtime.NumCPU())

    startTime := time.Now()

    urls := make([]string, nbRequests)
    for i := 0; i < nbRequests; i++ {
        urls[i] = fmt.Sprintf("https://jsonplaceholder.typicode.com/photos/%d", i+1)
    }

    var wg sync.WaitGroup
    processingTimes := make([]map[string]interface{}, nbRequests)
    results := make([]interface{}, nbRequests)
    var mu sync.Mutex

    for i, url := range urls {
        wg.Add(1)
        go func(i int, url string) {
            defer wg.Done()
            requestStartTime := time.Now()
            resp, err := http.Get(url)
            if err != nil {
                log.Printf("Error fetching data from %s: %v", url, err)
                return
            }
            defer resp.Body.Close()

            body, err := ioutil.ReadAll(resp.Body)
            if err != nil {
                log.Printf("Error reading response body from %s: %v", url, err)
                return
            }

            requestEndTime := time.Now()
            processingTime := requestEndTime.Sub(requestStartTime).Milliseconds()

            var result interface{}
            err = json.Unmarshal(body, &result)
            if err != nil {
                log.Printf("Error parsing response body from %s: %v", url, err)
                return
            }

            mu.Lock()
            processingTimes[i] = map[string]interface{}{
                "url":            url,
                "processingTime": processingTime,
            }
            results[i] = result
            mu.Unlock()
        }(i, url)
    }

    wg.Wait()

    maxProcessingTime := int64(0)
    for _, pt := range processingTimes {
        if pt != nil && pt["processingTime"].(int64) > maxProcessingTime {
            maxProcessingTime = pt["processingTime"].(int64)
        }
    }

    response := map[string]interface{}{
        "totalProcessingTime": time.Since(startTime).Milliseconds(),
        "maxProcessingTime":   maxProcessingTime,
//         "processingTimes":     processingTimes,
//         "results":             results,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

// Heavy computation handler
func heavyComputationHandler(w http.ResponseWriter, r *http.Request) {
    startTime := time.Now()
    iterations := 1000000000
    numGoroutines := runtime.NumCPU()
    var wg sync.WaitGroup
    results := make([]int64, numGoroutines)

    for g := 0; g < numGoroutines; g++ {
        wg.Add(1)
        go func(g int) {
            defer wg.Done()
            partialResult := int64(0)
            for i := g; i < iterations; i += numGoroutines {
                partialResult += int64(i)
            }
            results[g] = partialResult
        }(g)
    }

    wg.Wait()

    totalResult := int64(0)
    for _, partialResult := range results {
        totalResult += partialResult
    }

    totalProcessingTime := time.Since(startTime).Milliseconds()

    response := map[string]interface{}{
        "totalProcessingTime": totalProcessingTime,
        "result":              totalResult,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}


func helloHandler(w http.ResponseWriter, r *http.Request) {
    // Get the current timestamp
    currentTime := time.Now()

    // Add 5 minutes to the current time
    futureTime := currentTime.Add(5 * time.Minute)

    // Format the timestamps for display
    currentTimeFormatted := currentTime.Format(time.RFC1123)
    futureTimeFormatted := futureTime.Format(time.RFC1123)

    // Set response header and write the response
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "Hello World!\nCurrent Time: %s\nFuture Time (after 5 mins): %s", currentTimeFormatted, futureTimeFormatted)
}



func main() {
    http.HandleFunc("/", fetchPhotosHandler)
        http.HandleFunc("/hello", helloHandler)
    http.HandleFunc("/heavy", heavyComputationHandler)

    port := "3001"
    fmt.Printf("Server running at http://localhost:%s\n", port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}

package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type TestConfig struct {
	Name        string            `json:"name"`
	URI         string            `json:"uri"`
	Requests    int               `json:"requests"`
	Duration    int               `json:"duration"`
	Method      string            `json:"method"`
	Body        interface{}       `json:"body"`
	Headers     map[string]string `json:"headers"`
	Concurrency int               `json:"concurrency"`
}

type TestResult struct {
	Name       string        `json:"name"`
	Method     string        `json:"method"`
	Requests   int           `json:"requests"`
	Duration   float64       `json:"duration"`
	RPS        float64       `json:"rps"`
	AvgLatency time.Duration `json:"avg_latency"`
	MaxLatency time.Duration `json:"max_latency"`
	MinLatency time.Duration `json:"min_latency"`
}

func sendRequest(test TestConfig, latencies *[]time.Duration, wg *sync.WaitGroup) {
	defer wg.Done()

	var req *http.Request
	var err error

	body := ""

	if test.Method == http.MethodPost || test.Method == http.MethodPut {

		jsonBody, err := json.Marshal(test.Body)

		body = string(jsonBody)

		if err != nil {
			fmt.Println("Error serializer body:", err)
			return
		}
	}

	req, err = http.NewRequest(test.Method, test.URI, strings.NewReader(body))

	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	for key, value := range test.Headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	start := time.Now()

	resp, err := client.Do(req)

	duration := time.Since(start)

	if err != nil {
		fmt.Println("Request error:", err)
	} else {
		*latencies = append(*latencies, duration)
	}

	if resp.StatusCode >= 400 {
		fmt.Println("Unsuccessful HTTP response: %w", resp.Status)
	}
}

func runTest(test TestConfig, results chan<- TestResult, wg *sync.WaitGroup) {
	defer wg.Done()

	var wgRequests sync.WaitGroup

	latencies := make([]time.Duration, 0)

	startTime := time.Now()

	if test.Duration > 0 {

		ticker := time.NewTicker(time.Duration(test.Duration) * time.Second)

		done := make(chan bool)

		go func() {
			<-ticker.C
			done <- true
		}()

		for {
			select {
			case <-done:
				break
			default:
				wgRequests.Add(1)
				go sendRequest(test, &latencies, &wgRequests)
			}
		}
	} else {

		concurrencyChan := make(chan struct{}, test.Concurrency)

		for i := 0; i < test.Requests; i++ {

			concurrencyChan <- struct{}{}

			wgRequests.Add(1)

			go func() {
				sendRequest(test, &latencies, &wgRequests)
				<-concurrencyChan
			}()
		}
	}

	wgRequests.Wait()

	duration := time.Since(startTime).Seconds()

	totalLatency := time.Duration(0)

	maxLatency, minLatency := time.Duration(0), time.Duration(0)

	for _, latency := range latencies {
		totalLatency += latency
		if latency > maxLatency {
			maxLatency = latency
		}
		if minLatency == 0 || latency < minLatency {
			minLatency = latency
		}
	}

	avgLatency := totalLatency / time.Duration(len(latencies))
	rps := float64(len(latencies)) / duration

	results <- TestResult{
		Name:       test.Name,
		Method:     test.Method,
		Requests:   len(latencies),
		Duration:   duration,
		RPS:        rps,
		AvgLatency: avgLatency,
		MaxLatency: maxLatency,
		MinLatency: minLatency,
	}
}

func saveResults(results []TestResult) {

	file, err := os.Create("results.json")

	if err != nil {
		fmt.Println("Error creating results file:", err)
		return
	}

	defer file.Close()

	data, err := json.MarshalIndent(results, "", "  ")

	if err != nil {
		fmt.Println("Error converting results to JSON:", err)
		return
	}

	file.Write(data)

	fmt.Println("Results saved in results.json")
}

func loadConfig(filename string) ([]TestConfig, error) {

	data, err := os.ReadFile(filename)

	if err != nil {
		return nil, fmt.Errorf("error reading configuration file: %w", err)
	}

	var tests []TestConfig

	err = json.Unmarshal(data, &tests)

	if err != nil {
		return nil, fmt.Errorf("error parsing configuration file: %w", err)
	}

	return tests, nil
}

func main() {
	tests, err := loadConfig("config.json")

	if err != nil {
		fmt.Println(err)
		return
	}

	results := make(chan TestResult, len(tests))

	var wg sync.WaitGroup

	for _, test := range tests {
		wg.Add(1)
		go runTest(test, results, &wg)
	}

	wg.Wait()

	close(results)

	var testResults []TestResult

	for result := range results {

		testResults = append(testResults, result)

		fmt.Printf("Prueba '%s': RPS = %.2f, Latencia Promedio = %s\n", result.Name, result.RPS, result.AvgLatency)
	}

	saveResults(testResults)
}

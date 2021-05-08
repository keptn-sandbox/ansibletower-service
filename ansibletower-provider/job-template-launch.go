package provider

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	job_template "github.com/keptn-sandbox/ansibletower-service/ansibletower-provider/api-response-types/job_template"
	jobs "github.com/keptn-sandbox/ansibletower-service/ansibletower-provider/api-response-types/jobs"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
)

func LaunchJobTemplate(data *keptnv2.ActionTriggeredEventData) (string, error) {
	var err error

	serviceHost := os.Getenv("ANSIBLETOWER_HOST")
	serviceToken := os.Getenv("ANSIBLETOWER_TOKEN")

	// set url parts for launching job templates
	serviceAPIContext := "api"
	serviceAPIVersion := "v2"
	serviceAPIEndpoint := "job_templates"
	serviceAPIFunction := "launch"

	// get values from actionsOnOpen
	valuesMap, ok := data.Action.Value.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("Expected type map[string]interface{};  got %T", data.Action.Value)
	}

	var jobTemplateID string

	// check if value is supported
	for k, v := range valuesMap {
		switch k {
		case "JobTemplate":
			jobTemplateID = v.(string)
		default:
			log.Printf("Received an unsupported key '%v' with value '%v' from actionsOnOpen", k, v)
		}
	}

	// assemble request by concatenating all strings with slash
	// the empty string leaves a trailing slash (tower's nginx by default won't redirect without the slash)
	// example: https://{{host}}/api/v2/job_templates/9/launch/
	requestProtocol := "https://"
	requestUrl := requestProtocol + strings.Join([]string{
		serviceHost,
		serviceAPIContext,
		serviceAPIVersion,
		serviceAPIEndpoint,
		jobTemplateID,
		serviceAPIFunction,
		""}, "/")

	requestMethod := "POST"
	requestAuthorization := "Bearer " + serviceToken
	requestContentType := "application/json"
	requestMaxTimeout := 10 * time.Second

	// what about self-signed: skip or import?
	// serviceCertificate := os.Getenv("ANSIBLETOWER_CERTIFICATE")

	// log.Printf("requestUrl: %s", requestUrl)
	// log.Printf("requestMethod: %s", requestMethod)
	// log.Printf("requestAuthorization: %s", requestAuthorization)
	// log.Printf("requestContentType: %s", requestContentType)

	// check that everything is ok with the url
	_, err = url.Parse(requestUrl)
	if err != nil {
		return "", fmt.Errorf("URL could not be assembled: %s", err.Error())
	}

	// create request
	req, err := http.NewRequest(requestMethod, requestUrl, nil)
	req.Header.Add("Authorization", requestAuthorization)
	req.Header.Add("Content-Type", requestContentType)

	if err != nil {
		return "", fmt.Errorf("Error creating the request: %s", err.Error())
	}

	// // get ca certificate
	// caCert := []byte(``)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // create a certificate pool and append the certificate
	// caCertPool := x509.NewCertPool()
	// caCertPool.AppendCertsFromPEM(caCert)

	// TODO create the TLS object with the ca pool
	tlsConfig := &tls.Config{
		// RootCAs: caCertPool,
		InsecureSkipVerify: true,
	}

	// prepare the transport object
	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}

	// receive response
	client := &http.Client{
		Transport: transport,
		Timeout:   requestMaxTimeout,
	}

	statusCode, resp := ExecuteLaunchRequest(client, req)

	log.Println(resp.URL)

	errorMessage := fmt.Sprintf("failed to execute service action. status code %d", statusCode)

	// known status codes
	if statusCode == 201 {
		log.Println("action execute successfully. status code:", statusCode)

		return resp.URL, nil

	} else if statusCode == 400 {
		log.Println("required passwords were not provided.")
		return "", fmt.Errorf(errorMessage)

	} else if statusCode == 401 {
		log.Println("server is requesting log, check authorization header.")
		return "", fmt.Errorf(errorMessage)

	} else if statusCode == 403 {
		log.Println("provided credential or inventory are not allowed to be used by the user.")
		return "", fmt.Errorf(errorMessage)

	} else if statusCode == 405 {
		log.Println("job cannot be launched.")
		return "", fmt.Errorf(errorMessage)

	} else {
		// unknown response code
		return "", fmt.Errorf("Got unknown response. status code: %d", statusCode)
	}

}

func WaitJobEnd(jobURL string) {
	var interval = 10
	var maxTimeout = 30.0

	log.Println("Tick interval", time.Duration(interval)*time.Second)
	done := make(chan bool)

	go worker(done, interval, maxTimeout, jobURL)

	log.Println("Now waiting for the job to end...")

	select {
	case <-done:
		return
	}
}

func worker(done chan bool, interval int, maxTimeout float64, jobURL string) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)

	for t := range ticker.C {
		log.Println("Checking job finished tick", t)
		if IsFinished(maxTimeout, jobURL) {
			ticker.Stop()
			done <- true
		}
	}
}

func IsFinished(maxTimeout float64, jobURL string) bool {

	// log.Println("Querying Ansible Tower to check if job has finished")

	// get environment variables
	serviceHost := os.Getenv("ANSIBLETOWER_HOST")
	serviceToken := os.Getenv("ANSIBLETOWER_TOKEN")

	// assemble request by concatenating all strings with slash
	// the empty string leaves a trailing slash (tower's nginx by default won't redirect without the slash)
	// example: https://{{host}}/api/v2/jobs/129/
	requestProtocol := "https://"
	requestUrl := requestProtocol + serviceHost + jobURL

	requestMethod := "GET"
	requestAuthorization := "Bearer " + serviceToken
	requestContentType := "application/json"
	requestMaxTimeout := 10 * time.Second

	// create request
	req, err := http.NewRequest(requestMethod, requestUrl, nil)
	req.Header.Add("Authorization", requestAuthorization)
	req.Header.Add("Content-Type", requestContentType)

	// log.Printf("requestUrl: %s", requestUrl)
	// log.Printf("requestMethod: %s", requestMethod)
	// log.Printf("requestAuthorization: %s", requestAuthorization)
	// log.Printf("requestContentType: %s", requestContentType)

	if err != nil {
		log.Printf("Error creating the request: %s", err.Error())
	}

	// // get ca certificate
	// caCert := []byte(``)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // create a certificate pool and append the certificate
	// caCertPool := x509.NewCertPool()
	// caCertPool.AppendCertsFromPEM(caCert)

	// TODO create the TLS object with the ca pool
	tlsConfig := &tls.Config{
		// RootCAs: caCertPool,
		InsecureSkipVerify: true,
	}

	// prepare the transport object
	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}

	// receive response
	client := &http.Client{
		Transport: transport,
		Timeout:   requestMaxTimeout,
	}

	// launch job template
	_, resp := ExecuteJobsRequest(client, req)

	// log.Println("job", resp.ID)
	log.Println("job status", resp.Status)
	// log.Println("name", resp.Name)
	// log.Println("elapsed", resp.Elapsed)

	if resp.Elapsed >= maxTimeout {
		log.Println("maxTimeout reached, interrupting...")
		return true
	}

	switch resp.Status {
	case "new":
	case "pending":
	case "waiting":
	case "running":
		return false // keep checking
	case "successful":
	case "failed":
	case "error":
	case "canceled":
		return true // stop checking
	default:
		return true
	}
	return true
}

func ExecuteJobsRequest(client *http.Client, req *http.Request) (int, jobs.JobsResponse) {
	statusCode, body := ExecuteRequest(client, req)

	var jsonResponse jobs.JobsResponse

	err := json.Unmarshal(body, &jsonResponse)
	if err != nil {
		log.Println("Error when executing jobs request:", err)
	}

	return statusCode, jsonResponse
}

func ExecuteLaunchRequest(client *http.Client, req *http.Request) (int, job_template.LaunchResponse) {
	statusCode, body := ExecuteRequest(client, req)

	var jsonResponse job_template.LaunchResponse

	err := json.Unmarshal(body, &jsonResponse)
	if err != nil {
		log.Println("Error when executing launch request:", err)
	}

	return statusCode, jsonResponse
}

func ExecuteRequest(client *http.Client, req *http.Request) (int, []byte) {
	resp, err := client.Do(req)

	if err != nil {
		log.Printf("Response error: %s", err.Error())
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("ERror when reading the response body:", err)
	}

	return resp.StatusCode, body
}

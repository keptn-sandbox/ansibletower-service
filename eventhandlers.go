package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	keptn "github.com/keptn/go-utils/pkg/lib"
)

/**
* Here are all the handler functions for the individual event
  See https://github.com/keptn/spec/blob/0.1.3/cloudevents.md for details on the payload

  -> "sh.keptn.event.configuration.change"
  -> "sh.keptn.events.deployment-finished"
  -> "sh.keptn.events.tests-finished"
  -> "sh.keptn.event.start-evaluation"
  -> "sh.keptn.events.evaluation-done"
  -> "sh.keptn.event.problem.open"
	-> "sh.keptn.events.problem"
	-> "sh.keptn.event.action.triggered"
*/

//
// Handles ConfigurationChangeEventType = "sh.keptn.event.configuration.change"
// TODO: add in your handler code
//
func HandleConfigurationChangeEvent(myKeptn *keptn.Keptn, incomingEvent cloudevents.Event, data *keptn.ConfigurationChangeEventData) error {
	log.Printf("Handling Configuration Changed Event: %s", incomingEvent.Context.GetID())

	return nil
}

//
// Handles DeploymentFinishedEventType = "sh.keptn.events.deployment-finished"
// TODO: add in your handler code
//
func HandleDeploymentFinishedEvent(myKeptn *keptn.Keptn, incomingEvent cloudevents.Event, data *keptn.DeploymentFinishedEventData) error {
	log.Printf("Handling Deployment Finished Event: %s", incomingEvent.Context.GetID())

	// capture start time for tests
	// startTime := time.Now()

	// run tests
	// ToDo: Implement your tests here

	// Send Test Finished Event
	// return myKeptn.SendTestsFinishedEvent(&incomingEvent, "", "", startTime, "pass", nil, "ansibletower-service")
	return nil
}

//
// Handles TestsFinishedEventType = "sh.keptn.events.tests-finished"
// TODO: add in your handler code
//
func HandleTestsFinishedEvent(myKeptn *keptn.Keptn, incomingEvent cloudevents.Event, data *keptn.TestsFinishedEventData) error {
	log.Printf("Handling Tests Finished Event: %s", incomingEvent.Context.GetID())

	return nil
}

//
// Handles EvaluationDoneEventType = "sh.keptn.events.evaluation-done"
// TODO: add in your handler code
//
func HandleStartEvaluationEvent(myKeptn *keptn.Keptn, incomingEvent cloudevents.Event, data *keptn.StartEvaluationEventData) error {
	log.Printf("Handling Start Evaluation Event: %s", incomingEvent.Context.GetID())

	return nil
}

//
// Handles DeploymentFinishedEventType = "sh.keptn.events.deployment-finished"
// TODO: add in your handler code
//
func HandleEvaluationDoneEvent(myKeptn *keptn.Keptn, incomingEvent cloudevents.Event, data *keptn.EvaluationDoneEventData) error {
	log.Printf("Handling Evaluation Done Event: %s", incomingEvent.Context.GetID())

	return nil
}

//
// Handles ProblemOpenEventType = "sh.keptn.event.problem.open"
// Handles ProblemEventType = "sh.keptn.events.problem"
// TODO: add in your handler code
//
func HandleProblemEvent(myKeptn *keptn.Keptn, incomingEvent cloudevents.Event, data *keptn.ProblemEventData) error {
	log.Printf("Handling Problem Event: %s", incomingEvent.Context.GetID())

	// Deprecated since Keptn 0.7.0 - use the HandleActionTriggeredEvent instead

	return nil
}

//
// Handles ActionTriggeredEventType = "sh.keptn.event.action.triggered"
// TODO: add in your handler code
//
func HandleActionTriggeredEvent(myKeptn *keptn.Keptn, incomingEvent cloudevents.Event, data *keptn.ActionTriggeredEventData) error {
	log.Printf("Handling Action Triggered Event: %s", incomingEvent.Context.GetID())

	// check if action is supported
	if data.Action.Action == "hello" {
		log.Printf("Supported action: %s", data.Action.Action)

		actionStartedEventData := &keptn.ActionStartedEventData{}
		err := incomingEvent.DataAs(actionStartedEventData)
		if err != nil {
			log.Printf("Got Data Error: %s", err.Error())
			return err
		}

		// myKeptn.SendActionStartedEvent(actionStartedEventData) // TODO: implement the SendActionStartedEvent in keptn/go-utils/pkg/lib/events.go

		// check if everything is ok
		if os.Getenv("ANSIBLETOWER_HOST") == "" {
			return errors.New("Stopping execution of remediation action: ANSIBLETOWER_HOST is empty")
		}
		if os.Getenv("ANSIBLETOWER_TOKEN") == "" {
			return errors.New("Stopping execution of remediation action: ANSIBLETOWER_TOKEN is empty")
		}

		// get environment variables
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
			log.Printf("Expected type map[string]interface{};  got %T", data.Action.Value)
			return errors.New("Could not parse values from actionsOnOpen in the remediation specification")
		}

		var jobTemplateID string

		// check if value is supported
		for k, v := range valuesMap {
			switch k {
			case "JobTemplate":
				// convert numeric type to string
				value, ok := v.(float64)
				if !ok {
					return errors.New("Could not parse values from actionsOnOpen in the remediation specification")
				}
				jobTemplateID = fmt.Sprintf("%.0f", value)
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

		// what about self-signed: skip or import?
		// serviceCertificate := os.Getenv("ANSIBLETOWER_CERTIFICATE")

		log.Printf("requestUrl: %s", requestUrl)
		log.Printf("requestMethod: %s", requestMethod)
		log.Printf("requestAuthorization: %s", requestAuthorization)
		log.Printf("requestContentType: %s", requestContentType)

		// check that everything is ok with the url
		_, err = url.Parse(requestUrl)
		if err != nil {
			log.Printf("URL could not be assembled: %s", err.Error())
			return err
		}

		// create request
		req, err := http.NewRequest(requestMethod, requestUrl, nil)
		req.Header.Add("Authorization", requestAuthorization)
		req.Header.Add("Content-Type", requestContentType)

		if err != nil {
			log.Printf("Error creating the request: %s", err.Error())
			return err
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
		}
		resp, err := client.Do(req)

		if err != nil {
			log.Printf("Response error: %s", err.Error())
			return err
		}

		statusCode := strconv.Itoa(resp.StatusCode)

		// known status codes
		if resp.StatusCode == 201 {
			log.Println("action execute successfully. status code: " + statusCode)

		} else if resp.StatusCode == 400 {
			log.Println("failed to execute service action. status code: " + statusCode)
			log.Println("required passwords were not provided.")
			return errors.New("could not execute service action")

		} else if resp.StatusCode == 401 {
			log.Println("failed to execute service action. status code: " + statusCode)
			log.Println("server is requesting log, check authorization header.")
			return errors.New("could not execute service action")

		} else if resp.StatusCode == 403 {
			log.Println("failed to execute service action. status code: " + statusCode)
			log.Println("provided credential or inventory are not allowed to be used by the user.")
			return errors.New("could not execute service action")

		} else if resp.StatusCode == 405 {
			log.Println("failed to execute service action. status code: " + statusCode)
			log.Println("job cannot be launched.")
			return errors.New("could not execute service action")

		} else {
			// unknown response code
			log.Println("Got unknown response. status code: " + statusCode)
			return errors.New("could not execute service action")
		}

		//myKeptn.SendActionFinishedEvent() TODO: implement the SendActionFinishedEvent in keptn/go-utils/pkg/lib/events.go
	}

	return nil
}

package main

import (
	"errors"
	"log"
	"os"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	tower "github.com/keptn-sandbox/ansibletower-service/ansibletower-provider"
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

	// if the required input is not present, no action should be executed
	if os.Getenv("ANSIBLETOWER_HOST") == "" {
		return errors.New("Stopping execution of remediation action: ANSIBLETOWER_HOST is empty")
	}
	if os.Getenv("ANSIBLETOWER_TOKEN") == "" {
		return errors.New("Stopping execution of remediation action: ANSIBLETOWER_TOKEN is empty")
	}

	// check if action is supported
	if data.Action.Action == "job_template_launch" {
		log.Printf("Supported action: %s", data.Action.Action)

		// populate action started event from incoming event
		actionStartedEventData := &keptn.ActionStartedEventData{}
		err := incomingEvent.DataAs(actionStartedEventData)
		if err != nil {
			log.Printf("Got Data Error: %s", err.Error())
			return err
			// TODO should this send any event?
		}

		err = myKeptn.SendActionStartedEvent(&incomingEvent, data.Labels, "ansibletower-service")
		if err != nil {
			log.Printf("Got Error From SendActionStartedEvent: %s", err.Error())
			return err
		}

		// launch job template
		var jobURL string
		jobURL, err = tower.LaunchJobTemplate(data)

		if err != nil {
			log.Printf("Error launching the template: %s", err.Error())
			return err
		}

		// wait for job to finish
		tower.WaitJobEnd(jobURL)

		log.Println("all done.")

		var actionResult keptn.ActionResult
		actionResult.Result = "pass"
		actionResult.Status = "succeeded"

		myKeptn.SendActionFinishedEvent(&incomingEvent, actionResult, data.Labels, "ansibletower-service")
		if err != nil {
			log.Printf("Got Error From SendActionFinished: %s", err.Error())
			return err
		}
	}

	return nil
}

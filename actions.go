package main

import (
	"fmt"
	"terraform-api/runner"

	"github.com/google/uuid"
)

func plan(path string, r runner.Request) map[string]string {
	job := runner.CreateJob(runner.TF_PLAN, path, &r)

	go runner.ExecuteJob(job)

	return map[string]string{"id": job.ID.String()}
}

func apply(path string, r runner.Request) map[string]string {
	job := runner.CreateJob(runner.TF_APPLY, path, &r)

	go runner.ExecuteJob(job)

	return map[string]string{"id": job.ID.String()}
}

func destroy(path string, r runner.Request) map[string]string {
	job := runner.CreateJob(runner.TF_DESTROY, path, &r)

	go runner.ExecuteJob(job)

	return map[string]string{"id": job.ID.String()}
}

func output(id string) (map[string]interface{}, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	output := runner.GetJobOutput(uuid)

	return output, nil
}

func result(id string) (map[string]string, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	result := runner.GetJobResult(uuid)

	return map[string]string{
		"result": result.Result,
		"errors": fmt.Sprint(result.Errors),
	}, nil
}

func delete(path string) map[string]string {
	r := runner.Request{}
	job := runner.CreateJob(runner.WS_DELETE, path, &r)

	go runner.ExecuteJob(job)

	return map[string]string{"status": "workspace delete job created"}
}

func status(id string) (map[string]string, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	status := runner.GetJobStatus(uuid)

	switch status {
	case runner.JOB_CREATED:
		return map[string]string{"status": "Job Created"}, nil
	case runner.JOB_COMPLETE:
		return map[string]string{"status": "Job Complete"}, nil
	case runner.JOB_RUNNING:
		return map[string]string{"status": "Job Running"}, nil
	case runner.JOB_ERROR:
		return map[string]string{"status": "Job Error"}, nil
	default:
		return map[string]string{"status": "Unknown"}, nil
	}
}

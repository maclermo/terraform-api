package main

import (
	"fmt"
	"terraform-api/runner"
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

func jobs() (map[string]interface{}, error) {
	result := runner.GetJobsList()
	if len(result) == 0 {
		return map[string]interface{}{}, fmt.Errorf("no jobs listed")
	}

	return map[string]interface{}{"jobs": result}, nil
}

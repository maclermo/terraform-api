package runner

import (
	"fmt"
	"log"
	"testing"

	"github.com/google/uuid"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

var Jobs map[uuid.UUID]*Job

func InitJobs() {
	Jobs = make(map[uuid.UUID]*Job)
}

func CreateJob(action int, path string, r *Request) *Job {
	options := terraform.Options{
		Upgrade:       true,
		Lock:          true,
		NoColor:       true,
		Vars:          r.Vars,
		EnvVars:       r.EnvVars,
		BackendConfig: r.BackendConfig,
		TerraformDir:  path,
	}
	job := Job{
		ID:        uuid.New(),
		Workspace: uuid.New(),
		Action:    action,
		Request:   options,
		Status:    JOB_CREATED,
	}

	Jobs[job.ID] = &job

	return &job
}

func GetJobsList() map[uuid.UUID]*Job {
	return Jobs
}

func ExecuteJob(job *Job) {
	t := new(testing.T)

	log.Println("new job dispatch received for", job.ID.String())

	job.Status = JOB_RUNNING

	opts := terraform.WithDefaultRetryableErrors(t, &job.Request)

	switch job.Action {
	case TF_PLAN:
		log.Println("request \"plan\" recieved for", job.ID.String())

		if job.Response.Result, job.Response.Errors = terraform.WorkspaceSelectOrNewE(t, opts, job.Workspace.String()); job.Response.Errors != nil {
			job.Status = JOB_ERROR
			return
		}
		if job.Response.Result, job.Response.Errors = terraform.InitAndPlanE(t, opts); job.Response.Errors != nil {
			job.Status = JOB_ERROR
			return
		}
	case TF_APPLY:
		log.Println("request \"apply\" recieved for", job.ID.String())

		if job.Response.Result, job.Response.Errors = terraform.WorkspaceSelectOrNewE(t, opts, job.Workspace.String()); job.Response.Errors != nil {
			job.Status = JOB_ERROR
			return
		}
		if job.Response.Result, job.Response.Errors = terraform.InitAndApplyE(t, opts); job.Response.Errors != nil {
			job.Status = JOB_ERROR
			return
		}
		if job.Response.Output, job.Response.Errors = terraform.OutputJsonE(t, opts, ""); job.Response.Errors != nil {
			job.Status = JOB_ERROR
			return
		}
	case TF_DESTROY:
		log.Println("request \"destroy\" recieved for", job.ID.String())

		if job.Response.Result, job.Response.Errors = terraform.WorkspaceSelectOrNewE(t, opts, job.Workspace.String()); job.Response.Errors != nil {
			job.Status = JOB_ERROR
			return
		}
		if job.Response.Result, job.Response.Errors = terraform.InitE(t, opts); job.Response.Errors != nil {
			job.Status = JOB_ERROR
			return
		}
		if job.Response.Result, job.Response.Errors = terraform.DestroyE(t, opts); job.Response.Errors != nil {
			job.Status = JOB_ERROR
			return
		}
	default:
		job.Response.Errors = fmt.Errorf("request not recognized for %s", job.ID.String())
		log.Println(job.Response.Errors)

		job.Status = JOB_ERROR
		return
	}

	job.Status = JOB_COMPLETE
}

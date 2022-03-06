package runner

import (
	"github.com/google/uuid"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

type JSONRequest struct {
	Vars          string `form:"vars"`
	EnvVars       string `form:"envvars"`
	BackendConfig string `form:"backendconfig"`
}

type Request struct {
	Vars          map[string]interface{}
	EnvVars       map[string]string
	BackendConfig map[string]interface{}
}

type JobResponse struct {
	Result string
	Errors error
	Output string
}

type Job struct {
	ID        uuid.UUID
	Workspace uuid.UUID
	Action    int
	Request   terraform.Options
	Response  JobResponse
	Status    int
}

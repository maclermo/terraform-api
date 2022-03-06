package runner

import (
	"github.com/google/uuid"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

type JSONRequest struct {
	Workspace     string `form:"workspace" binding:"required"`
	Vars          string `form:"vars"`
	EnvVars       string `form:"envvars"`
	BackendConfig string `form:"backendconfig"`
}

type Request struct {
	Workspace     string
	Vars          map[string]interface{}
	EnvVars       map[string]string
	BackendConfig map[string]interface{}
}

type JobResponse struct {
	Result string
	Errors error
}

type Job struct {
	ID        uuid.UUID
	Workspace string
	Action    int
	Request   terraform.Options
	Response  JobResponse
	Status    int
}

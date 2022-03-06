package runner

const (
	TF_PLAN = iota
	TF_APPLY
	TF_DESTROY
)

const (
	JOB_CREATED = iota
	JOB_RUNNING
	JOB_COMPLETE
	JOB_ERROR
)

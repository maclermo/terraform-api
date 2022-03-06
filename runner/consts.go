package runner

const (
	TF_APPLY = iota
	TF_PLAN
	TF_DESTROY
	WS_DELETE
)

const (
	JOB_CREATED = iota
	JOB_RUNNING
	JOB_COMPLETE
	JOB_ERROR
)
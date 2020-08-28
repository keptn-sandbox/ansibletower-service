package jobs

// Generated by https://quicktype.io

type JobsResponse struct {
	ID                      int64             `json:"id"`
	Type                    string            `json:"type"`
	URL                     string            `json:"url"`
	Related                 map[string]string `json:"related"`
	SummaryFields           SummaryFields     `json:"summary_fields"`
	Created                 string            `json:"created"`
	Modified                string            `json:"modified"`
	Name                    string            `json:"name"`
	Description             string            `json:"description"`
	JobType                 string            `json:"job_type"`
	Inventory               int64             `json:"inventory"`
	Project                 int64             `json:"project"`
	Playbook                string            `json:"playbook"`
	SCMBranch               string            `json:"scm_branch"`
	Forks                   int64             `json:"forks"`
	Limit                   string            `json:"limit"`
	Verbosity               int64             `json:"verbosity"`
	ExtraVars               string            `json:"extra_vars"`
	JobTags                 string            `json:"job_tags"`
	ForceHandlers           bool              `json:"force_handlers"`
	SkipTags                string            `json:"skip_tags"`
	StartAtTask             string            `json:"start_at_task"`
	Timeout                 int64             `json:"timeout"`
	UseFactCache            bool              `json:"use_fact_cache"`
	Organization            int64             `json:"organization"`
	UnifiedJobTemplate      int64             `json:"unified_job_template"`
	LaunchType              string            `json:"launch_type"`
	Status                  string            `json:"status"`
	Failed                  bool              `json:"failed"`
	Started                 string            `json:"started"`
	Finished                string            `json:"finished"`
	CanceledOn              interface{}       `json:"canceled_on"`
	Elapsed                 float64           `json:"elapsed"`
	JobArgs                 string            `json:"job_args"`
	JobCwd                  string            `json:"job_cwd"`
	JobEnv                  map[string]string `json:"job_env"`
	JobExplanation          string            `json:"job_explanation"`
	ExecutionNode           string            `json:"execution_node"`
	ControllerNode          string            `json:"controller_node"`
	ResultTraceback         string            `json:"result_traceback"`
	EventProcessingFinished bool              `json:"event_processing_finished"`
	JobTemplate             int64             `json:"job_template"`
	PasswordsNeededToStart  []interface{}     `json:"passwords_needed_to_start"`
	AllowSimultaneous       bool              `json:"allow_simultaneous"`
	Artifacts               Artifacts         `json:"artifacts"`
	SCMRevision             string            `json:"scm_revision"`
	InstanceGroup           int64             `json:"instance_group"`
	DiffMode                bool              `json:"diff_mode"`
	JobSliceNumber          int64             `json:"job_slice_number"`
	JobSliceCount           int64             `json:"job_slice_count"`
	WebhookService          string            `json:"webhook_service"`
	WebhookCredential       interface{}       `json:"webhook_credential"`
	WebhookGUID             string            `json:"webhook_guid"`
	HostStatusCounts        HostStatusCounts  `json:"host_status_counts"`
	PlaybookCounts          PlaybookCounts    `json:"playbook_counts"`
	CustomVirtualenv        string            `json:"custom_virtualenv"`
}

type Artifacts struct {
}

type HostStatusCounts struct {
	Changed int64 `json:"changed"`
}

type PlaybookCounts struct {
	PlayCount int64 `json:"play_count"`
	TaskCount int64 `json:"task_count"`
}

type SummaryFields struct {
	Organization       JobTemplate      `json:"organization"`
	Inventory          Inventory        `json:"inventory"`
	Project            Project          `json:"project"`
	ProjectUpdate      Project          `json:"project_update"`
	JobTemplate        JobTemplate      `json:"job_template"`
	UnifiedJobTemplate JobTemplate      `json:"unified_job_template"`
	InstanceGroup      InstanceGroup    `json:"instance_group"`
	CreatedBy          CreatedBy        `json:"created_by"`
	UserCapabilities   UserCapabilities `json:"user_capabilities"`
	Labels             Labels           `json:"labels"`
	ExtraCredentials   []interface{}    `json:"extra_credentials"`
	Credentials        []interface{}    `json:"credentials"`
}

type CreatedBy struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type InstanceGroup struct {
	ID              int64  `json:"id"`
	Name            string `json:"name"`
	IsContainerized bool   `json:"is_containerized"`
}

type Inventory struct {
	ID                           int64  `json:"id"`
	Name                         string `json:"name"`
	Description                  string `json:"description"`
	HasActiveFailures            bool   `json:"has_active_failures"`
	TotalHosts                   int64  `json:"total_hosts"`
	HostsWithActiveFailures      int64  `json:"hosts_with_active_failures"`
	TotalGroups                  int64  `json:"total_groups"`
	HasInventorySources          bool   `json:"has_inventory_sources"`
	TotalInventorySources        int64  `json:"total_inventory_sources"`
	InventorySourcesWithFailures int64  `json:"inventory_sources_with_failures"`
	OrganizationID               int64  `json:"organization_id"`
	Kind                         string `json:"kind"`
}

type JobTemplate struct {
	ID             int64   `json:"id"`
	Name           string  `json:"name"`
	Description    string  `json:"description"`
	UnifiedJobType *string `json:"unified_job_type,omitempty"`
}

type Labels struct {
	Count   int64         `json:"count"`
	Results []interface{} `json:"results"`
}

type Project struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Status      string  `json:"status"`
	SCMType     *string `json:"scm_type,omitempty"`
	Failed      *bool   `json:"failed,omitempty"`
}

type UserCapabilities struct {
	Delete bool `json:"delete"`
	Start  bool `json:"start"`
}
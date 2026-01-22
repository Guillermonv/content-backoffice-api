package dto

import "time"

type ExecutionDTO struct {
	ID        uint64      `json:"id"`
	Status    string      `json:"status"`
	CreatedAt time.Time   `json:"created_at"`
	Workflow  WorkflowDTO `json:"workflow"`
}

type WorkflowDTO struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type StepExecutionGroupResponse struct {
	ExecutionID uint64       `json:"execution_id"`
	Execution   ExecutionDTO `json:"execution"`
}

package dto

type StepInputDto struct {
	OrderIndex    int    `json:"orderIndex"`
	Name          string `json:"name"`
	OperationType string `json:"operationType"`
	Prompt        string `json:"prompt"`

	WorkflowID uint64 `json:"workflowId"`
	AgentID    uint64 `json:"agentId"`
}
package model

type Step struct {
	ID            uint64 `gorm:"primaryKey"`
	OrderIndex    int
	Name          string
	OperationType string
	Prompt        string `gorm:"type:mediumtext"`

	WorkflowID uint64 `json:"-" gorm:"column:workflow_id;not null"`
	AgentID    *uint64 `json:"-" gorm:"column:agent_id"`

	Workflow Workflow `gorm:"foreignKey:WorkflowID;references:ID"`
	Agent    Agent    `gorm:"foreignKey:AgentID;references:ID"`
}

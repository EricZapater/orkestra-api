package tasks

type TaskStatus string

const (
	StatusPending    TaskStatus = "Pending"
	StatusTodo       TaskStatus = "ToDo"
	StatusInProgress TaskStatus = "InProgress"
	StatusDone       TaskStatus = "Done"
)

type TaskPriority string

const (
	PriorityHigh   TaskPriority = "A"
	PriorityMedium TaskPriority = "B"
	PriorityLow    TaskPriority = "C"
	NoPriority     TaskPriority = "D"
)

func IsValidStatus(s TaskStatus) bool {
	switch s {
	case StatusPending, StatusTodo, StatusInProgress, StatusDone:
		return true
	default:
		return false
	}
}

func IsValidPriority(p TaskPriority) bool {
	switch p {
	case PriorityHigh, PriorityMedium, PriorityLow, NoPriority:
		return true
	default:
		return false
	}
}
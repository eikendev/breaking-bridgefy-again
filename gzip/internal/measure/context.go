package measure

type RunContext struct {
	SenderName string
	SenderUUID string
}

func DefaultRunContext() *RunContext {
	return &RunContext{
		SenderName: "",
		SenderUUID: "",
	}
}

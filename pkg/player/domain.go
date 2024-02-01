package player

type StateType struct {
	ID      string
	Img     string
	Playing bool
}

type Action struct {
	Type    ActionType
	Payload []string
}

var currentState = StateType{
	ID:      "",
	Img:     "",
	Playing: false,
}

type ActionType string

const (
	PlayType ActionType = "play"
	StopType ActionType = "stop"
)

func Reduce(action Action) {
	switch action.Type {
	case PlayType:
		currentState = StateType{
			Playing: true,
			ID:      action.Payload[0],
			Img:     action.Payload[1],
		}
	case StopType:
		currentState = StateType{
			Playing: false,
			Img:     "",
			ID:      "",
		}
	}
}

func State() StateType {
	return currentState
}

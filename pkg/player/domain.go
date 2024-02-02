package player

type StateType struct {
	ID      string
	Img     string
	Playing bool
	Rooms   []string
}

type Action struct {
	Type    ActionType
	Payload []string
}

var currentState = StateType{
	ID:      "",
	Img:     "",
	Rooms:   []string{},
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
			Rooms:   append(currentState.Rooms, action.Payload[2]),
		}
	case StopType:
		rooms := []string{}
		for _, room := range currentState.Rooms {
			if room != action.Payload[2] {
				rooms = append(rooms, room)
			}
		}

		currentState = StateType{
			Playing: len(rooms) > 0,
			ID:      action.Payload[0],
			Img:     action.Payload[1],
			Rooms:   rooms,
		}
	}
}

func State() StateType {
	return currentState
}

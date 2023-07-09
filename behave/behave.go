package behave

type State uint

const (
	Running State = iota
	Success
	Failure
	Unknown
)

func (s State) String() string {
	switch s {
	case Running:
		return "Running"
	case Success:
		return "Success"
	case Failure:
		return "Failure"
	default:
		return "Unknown"
	}
}

type Behavior interface {
	Reset()
	Execute() State
}

type Action func() State

func (Action) Reset() {}

func (a Action) Execute() State {
	return a()
}

type composite struct {
	nodes []Behavior
	index int
}

func (c *composite) Reset() {
	c.index = 0
	for _, n := range c.nodes {
		n.Reset()
	}
}

type sequence struct {
	composite
}

func Sequence(bs ...Behavior) Behavior {
	return &sequence{composite{nodes: bs}}
}

func (s *sequence) Execute() State {
	for ; s.index < len(s.nodes); s.index++ {
		switch s.nodes[s.index].Execute() {
		case Running:
			return Running
		case Success:
			continue
		case Failure:
			return Failure
		default:
			return Unknown
		}
	}
	return Success
}

type selection struct {
	composite
}

func Selection(bs ...Behavior) Behavior {
	return &selection{composite{nodes: bs}}
}

func (s *selection) Execute() State {
	for ; s.index < len(s.nodes); s.index++ {
		switch s.nodes[s.index].Execute() {
		case Running:
			return Running
		case Success:
			return Success
		case Failure:
			continue
		default:
			return Unknown
		}
	}
	return Failure
}

package behave

import (
	"reflect"
	"testing"
)

func Recorded(states ...State) Behavior {
	i := 0
	return Action(func() State {
		result := states[i%len(states)]
		i++
		return result
	})
}

func TestSequence_Success(t *testing.T) {
	b := Sequence(
		Recorded(Running, Success),
		Recorded(Running, Success),
		Recorded(Success),
	)
	expected := []State{Running, Running, Success}
	actual := make([]State, len(expected))
	for i := range expected {
		actual[i] = b.Execute()
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Error("Sequence produced incorrect states", actual)
	}
}

func TestSequence_Failure(t *testing.T) {
	b := Sequence(
		Recorded(Running, Success),
		Recorded(Running, Running, Failure),
		Recorded(Success),
	)
	expected := []State{Running, Running, Running, Failure}
	actual := make([]State, len(expected))
	for i := range expected {
		actual[i] = b.Execute()
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Error("Sequence produced incorrect states", actual)
	}
}

func TestSelection_Success(t *testing.T) {
	b := Selection(
		Recorded(Running, Failure),
		Recorded(Failure),
		Recorded(Running, Failure),
		Recorded(Success),
		Recorded(Success),
	)
	expected := []State{Running, Running, Success}
	actual := make([]State, len(expected))
	for i := range expected {
		actual[i] = b.Execute()
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Error("Selection produced incorrect states", actual)
	}
}

func TestSelection_Failure(t *testing.T) {
	b := Selection(
		Recorded(Running, Failure),
		Recorded(Failure),
		Recorded(Running, Failure),
		Recorded(Failure),
		Recorded(Failure),
		Recorded(Running, Failure),
	)
	expected := []State{Running, Running, Running, Failure}
	actual := make([]State, len(expected))
	for i := range expected {
		actual[i] = b.Execute()
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Error("Selection produced incorrect states", actual)
	}
}

func TestConditional(t *testing.T) {
	cases := []struct {
		output   bool
		expected State
	}{
		{true, Success},
		{false, Failure},
	}
	for _, c := range cases {
		b := Conditional(func() bool { return c.output })
		if b.Execute() != c.expected {
			t.Errorf("Conditional failed to turn %t into %v", c.output, c.expected)
		}
	}
}

func TestInvert(t *testing.T) {
	b := Invert(Recorded(Running, Failure, Success, Unknown))
	expected := []State{Running, Success, Failure, Unknown}
	actual := make([]State, len(expected))
	for i := range expected {
		actual[i] = b.Execute()
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Error("Selection produced incorrect states", actual)
	}
}

type testBehavior struct {
	base   Behavior
	resets int
}

func (b *testBehavior) Execute() State {
	return b.base.Execute()
}

func (b *testBehavior) Reset() {
	b.base.Reset()
	b.resets++
}

func TestRepeat(t *testing.T) {
	wrapped := &testBehavior{base: Recorded(Running, Failure, Success, Unknown)}
	repeat := Repeat(wrapped)
	expected := []State{Running, Running, Running, Unknown}
	actual := make([]State, len(expected))
	for i := range expected {
		actual[i] = repeat.Execute()
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Error("Repeat produced incorrect states", actual)
	}
	if wrapped.resets != 2 {
		t.Error("Repeat failed to reset wrapped Behavior", wrapped.resets)
	}
}

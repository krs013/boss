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

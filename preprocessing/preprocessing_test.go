package preprocessing

import (
	"slices"
	"testing"
)

const dataSet string = "../dataset"

func TestReadDir(t *testing.T) {
	got, _ := ReadDir(dataSet)
	want := []string{"testing.pdf", "testing_2.pdf"}
	if !slices.Equal(got, want) {
		t.Errorf("got %v wanted %v", got, want)
	}
}

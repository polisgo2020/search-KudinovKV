package index

import (
	"reflect"
	"sync"
	"testing"
)

var (
	i           InvertIndex
	in          InvertIndex
	listOfFiles []int
)

func init() {
	i = NewInvertIndex()
	in = NewInvertIndex()
}

func TestContains(t *testing.T) {
	in := []int{0, -500, 1, 2, -1}
	actual := Contains(in, 1)
	expected := true
	if expected != actual {
		t.Errorf("%v is not equal to expected %v", actual, expected)
	}
	expected = false
	actual = Contains(in, 10)
	if expected != actual {
		t.Errorf("%v is not equal to expected %v", actual, expected)
	}
}

func TestParseIndexFile(t *testing.T) {
	in := string("is:0,1,2\na:2\nbanana:2\nit:0,1,2\nwhat:0,1\n")
	expected := []int{0, 1, 2}
	actual := i.ParseIndexFile(in)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("%v is not equal to expected %v", actual, expected)
	}
	listOfFiles = actual
}
func TestMakeSearch(t *testing.T) {
	in := []string{"banana", "is"}
	expected := []int{2, 1, 1}
	actual := i.MakeSearch(in, listOfFiles)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("%v is not equal to expected %v", actual, expected)
	}
}

func TestPrepareTokens(t *testing.T) {
	in := string("Hello 123 is world !")
	expected := []string{"hello", "world"}
	actual := PrepareTokens(in)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("%v is not equal to expected %v", actual, expected)
	}
}

func TestAddToken(t *testing.T) {
	expected := NewInvertIndex()
	expected["newtoken"] = append(expected["newtoken"], 1)
	in.addToken("newtoken", 1)
	if !reflect.DeepEqual(in["newtoken"], expected["newtoken"]) {
		t.Errorf("%v is not equal to expected %v", in, expected)
	}
}

func TestListener(t *testing.T) {
	dataCh := make(chan []string)
	mutex := &sync.Mutex{}

	expected := NewInvertIndex()
	expected["newtoken"] = append(expected["newtoken"], 1)
	expected["newtoken"] = append(expected["newtoken"], 2)

	go in.Listener(dataCh, mutex)
	dataCh <- []string{"newtoken", "2"}
	close(dataCh)
	mutex.Lock()
	if !reflect.DeepEqual(in["newtoken"], expected["newtoken"]) {
		t.Errorf("%v is not equal to expected %v", in, expected)
	}
	mutex.Unlock()
}

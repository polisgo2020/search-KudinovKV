package index

import (
	"reflect"
	"testing"
)

var (
	i           *InvertIndex
	in          *InvertIndex
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
	expected := []Rate{{2, 2}, {0, 1}, {1, 1}}
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
	(*expected).index["newtoken"] = append((*expected).index["newtoken"], 1)
	in.addToken("newtoken", 1)
	if !reflect.DeepEqual((*in).index["newtoken"], (*expected).index["newtoken"]) {
		t.Errorf("%v is not equal to expected %v", in, expected)
	}
}

func TestListener(t *testing.T) {
	expected := NewInvertIndex()
	(*expected).index["newtoken"] = append((*expected).index["newtoken"], 1)
	(*expected).index["newtoken"] = append((*expected).index["newtoken"], 2)

	in.dataCh <- []string{"newtoken", "2"}
	close(in.dataCh)
	in.mutex.Lock()
	if !reflect.DeepEqual((*in).index["newtoken"], (*expected).index["newtoken"]) {
		t.Errorf("%v is not equal to expected %v", *in, *expected)
	}
	in.mutex.Unlock()

	close(expected.dataCh)
	close(i.dataCh)
}

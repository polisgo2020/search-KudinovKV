package index

import (
	"reflect"
	"testing"
	"time"
)

/*
	go test -coverprofile=cover.out
	go tool cover -html=cover.out -o cover.html
*/

var (
	i           *InvertIndex
	in          *InvertIndex
	listOfFiles []string
)

func init() {
	i = NewInvertIndex()
	in = NewInvertIndex()
}

func TestContains(t *testing.T) {
	in := []string{"abc", "test", "hello", "world"}
	actual := Contains(in, "test")
	expected := true
	if expected != actual {
		t.Errorf("%v is not equal to expected %v", actual, expected)
	}
	expected = false
	actual = Contains(in, "abctest")
	if expected != actual {
		t.Errorf("%v is not equal to expected %v", actual, expected)
	}
}

func TestParseIndexFile(t *testing.T) {
	in := string("is:0.txt,1.txt,2.txt\na:2.txt\nbanana:2.txt\nit:0.txt,1.txt,2.txt\nwhat:0.txt,1.txt\n")
	expected := map[string][]string{
		"is":     []string{"0.txt", "1.txt", "2.txt"},
		"a":      []string{"2.txt"},
		"banana": []string{"2.txt"},
		"it":     []string{"0.txt", "1.txt", "2.txt"},
		"what":   []string{"0.txt", "1.txt"},
	}
	i.ParseIndexFile(in)
	if !reflect.DeepEqual(i.index, expected) {
		t.Errorf("%v is not equal to expected %v", i.index, expected)
	}
}
func TestGet(t *testing.T) {
	in := []string{"banana", "is"}
	expected := []Rate{{"2.txt", 2}, {"0.txt", 1}, {"1.txt", 1}}
	actual, err := i.Get(in)
	if err != nil {
		t.Errorf("%v is not equal to expected %v", actual, expected)
	}
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

func TestAdd(t *testing.T) {
	expected := NewInvertIndex()
	(*expected).index["newtoken"] = append((*expected).index["newtoken"], "1.txt")
	in.Add("newtoken", "1.txt")
	if !reflect.DeepEqual((*in).index["newtoken"], (*expected).index["newtoken"]) {
		t.Errorf("%v is not equal to expected %v", in, expected)
	}
}

func TestListener(t *testing.T) {
	expected := NewInvertIndex()
	(*expected).index["newtoken"] = append((*expected).index["newtoken"], "1.txt")
	(*expected).index["newtoken"] = append((*expected).index["newtoken"], "2.txt")

	in.dataCh <- []string{"newtoken", "2.txt"}
	time.Sleep(1 * time.Second)
	close(in.dataCh)
	if !reflect.DeepEqual((*in).index["newtoken"], (*expected).index["newtoken"]) {
		t.Errorf("%v is not equal to expected %v", *in, *expected)
	}
	close(expected.dataCh)
	close(i.dataCh)
}

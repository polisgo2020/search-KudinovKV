package index

import (
	"reflect"
	"testing"
)

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
	expected := []string{"0.txt", "1.txt", "2.txt"}
	actual := i.ParseIndexFile(in)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("%v is not equal to expected %v", actual, expected)
	}
	listOfFiles = actual
}
func TestMakeSearch(t *testing.T) {
	in := []string{"banana", "is"}
	expected := []Rate{{"2.txt", 2}, {"0.txt", 1}, {"1.txt", 1}}
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
	(*expected).index["newtoken"] = append((*expected).index["newtoken"], "1.txt")
	in.addToken("newtoken", "1.txt")
	if !reflect.DeepEqual((*in).index["newtoken"], (*expected).index["newtoken"]) {
		t.Errorf("%v is not equal to expected %v", in, expected)
	}
}

func TestListener(t *testing.T) {
	expected := NewInvertIndex()
	(*expected).index["newtoken"] = append((*expected).index["newtoken"], "1.txt")
	(*expected).index["newtoken"] = append((*expected).index["newtoken"], "2.txt")

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

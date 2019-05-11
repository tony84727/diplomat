package internal

import (
	"github.com/stretchr/testify/suite"
	"reflect"
	"testing"
)

type FieldSearchTestSuite struct {
	suite.Suite
}

func TestFieldSearcher(t *testing.T) {
	suite.Run(t, &FieldSearchTestSuite{})
}

func (f FieldSearchTestSuite) TestSearch() {
	fake := struct {
		Name string `navigate:"first"`
		Number int `navigate:"second"`
	}{
		"whatever",
		100,
	}
	searcher := FieldSearcher{reflect.ValueOf(fake)}
	i, ok := searcher.Search("first")
	f.Require().True(ok)
	f.Equal([]int{0}, i)
}

func (f FieldSearchTestSuite) TestSearch_Nested() {
	type nested struct {
		Name string `navigate:"first"`
		Number int `navigate:"second"`
	}
	fake := struct {
		embedded nested
	}{
		embedded: nested{Name: "whatever", Number: 100 },
	}
	searcher := FieldSearcher{reflect.ValueOf(fake)}
	i, ok := searcher.Search("second")
	f.Require().True(ok)
	f.Equal([]int{0,1}, i)
}
type DummyInterface interface {
	DoNothing()
}

type DummyStruct struct {
	DummyInterface
}

type DummyImpl struct {
	Name string `navigate:"first"`
	Number int `navigate:"second"`
}

func (d DummyImpl) DoNothing() {}

func (f FieldSearchTestSuite) TestSearch_InterfaceEmbedded() {
	fake := DummyStruct{&DummyImpl{Name:"name",Number: 100}}
	searcher := FieldSearcher{reflect.ValueOf(fake)}
	i, ok := searcher.Search("second")
	f.Require().True(ok)
	f.Equal([]int{0,1}, i)
}

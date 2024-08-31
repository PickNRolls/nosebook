package testlib

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/go-test/deep"
)

type J = map[string]any

type Matcher struct {
	t                 *testing.T
	value             any
	lastExpectedValue any
	diff              []string
	ok                bool
	failOnExpect      bool
	inverted          bool
}

func (this *Matcher) assert(ok bool) {
	this.ok = ok
	if this.inverted {
		this.ok = !this.ok
	}
}

func (this *Matcher) Not() *Matcher {
	this.inverted = !this.inverted
	return this
}

func (this *Matcher) ToBe(value any) *Matcher {
	this.diff = deep.Equal(this.value, value)
	this.assert(len(this.diff) == 0)
	this.lastExpectedValue = value

	if this.failOnExpect {
		this.ElseFail()
	}

	return this
}

func (this *Matcher) ToBeTypeOf(t string) *Matcher {
	kind := reflect.TypeOf(this.value).Kind().String()
	this.assert(kind == t)
	this.diff = append(this.diff, fmt.Sprintf("Value type is %v instead of %v", kind, t))

	if this.failOnExpect {
		this.ElseFail()
	}

	return this
}

func contains(source any, target any) []string {
	output := []string{}

	if reflect.TypeOf(source) != reflect.TypeOf(target) {
		output = append(output, "different types")
		return output
	}

	kind := reflect.TypeOf(source).Kind()
	if kind != reflect.Struct && kind != reflect.Map && kind != reflect.Slice {
		output = deep.Equal(source, target)
		return output
	}

	if kind == reflect.Slice {
		sourceS, _ := source.([]any)
		targetS, _ := target.([]any)

		for i, value := range targetS {
			if len(sourceS) <= i {
				output = append(output, "source has less elements than target\n")
				output = append(output, deep.Equal(sourceS, targetS)...)
				return output
			}
			output = append(output, contains(sourceS[i], value)...)
		}

		return output
	}

	sourceJ, _ := source.(J)
	targetJ, _ := target.(J)

	for key, value := range targetJ {
		if _, has := sourceJ[key]; !has {
			output = append(output, "No key %v\n", key)
			continue
		}

		output = append(output, contains(sourceJ[key], value)...)
	}

	return output
}

func (this *Matcher) ToContain(value J) *Matcher {
	j, ok := this.value.(J)
	if !ok {
		return this
	}

	this.diff = contains(j, value)
	this.assert(len(this.diff) == 0)
	if this.failOnExpect {
		this.ElseFail()
	}
	return this
}

func (this *Matcher) ElseFail() {
	if !this.ok {
		this.t.Error(this.diff)
	}
}

func CreateMatcher(t *testing.T, failOnExpect bool) func(any) *Matcher {
	return func(value any) *Matcher {
		return &Matcher{
			t:            t,
			value:        value,
			ok:           true,
			failOnExpect: failOnExpect,
		}
	}
}

package application_tests

import (
	"reflect"
	"testing"
)

type Matcher struct {
	t                 *testing.T
	value             any
	lastExpectedValue any
	ok                bool
	failOnExpect      bool
}

func (this *Matcher) ToBe(value any) *Matcher {
	this.ok = this.value == value
	this.lastExpectedValue = value
	if this.failOnExpect {
		this.ElseFail()
	}
	return this
}

func (this *Matcher) ToDeepEqual(value any) *Matcher {
	this.ok = reflect.DeepEqual(this.value, value)
	this.lastExpectedValue = value
	if this.failOnExpect {
		this.ElseFail()
	}
	return this
}

func (this *Matcher) ElseFail() {
	if !this.ok {
		this.t.Fatalf("Expected %v, got %v", this.lastExpectedValue, this.value)
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

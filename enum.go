// A struct based golang enum package that handles stringifying, marshalling, and unmarshalling.
// All examples will be built from the following example types
//   type CurrencyCodes struct {
//     enum.Enum
//     USD    enum.Const
//     Custom enum.Const `enum:"CUSTOM"`
//   }
//
//   type Money struct {
//     CurrencyCode CurrencyCode `json:"currency_code"`
//     Amount       int          `json:"amount"`
//   }
package enum

import (
	"fmt"
	"github.com/pkg/errors"
	"reflect"
	"strconv"
)

const invalidEnumErrorMsg = "%s is not a valid enum"
const enumNotConstructedErrorMsg = "cannot set a value on an enum that has not be constructed"
const enumNotNilErrorMsg = "cannot set a value on an enum that has not be constructed"

type Enummer interface {
	// Gets the value stored on the enum
	Get() Const
	// Set the value stored on the enum. Returns an error if value is invalid
	Set(c Const) error
	// Set the value stored on the enum. Panics if the value is invalid
	MustSet(c Const)
	// A list of all possible Consts on the enum
	GetAll() []Const
	// Sets the value without any checks. This is a dangerous method and should pretty
	// much never be used but if you do, USE WITH CAUTION.
	unsafeSet(c Const)
	// Adds an additional valid value. This is a dangerous method and should pretty
	// much never be used but if you do, USE WITH CAUTION.
	unsafeAdd(c Const)
}

// The base type for all Enums. Stores the current value and keeps track of all valid values.
type Enum struct {
	val  Const
	vals []Const
}

// The base value for all Enum fields. The name of the field on the enum struct will be
// the default value for the enum const. For example
// The value of the USD const is "USD". In order to customize this value, add the tag enum:"<name>"
// for example
//   type CurrencyCodes struct {
//     enum.Enum
//     USD enum.Const `enum:"not usd"`
//   }
// The value of the USD const is now "not usd"
type Const string

func (e *Enum) unsafeAdd(c Const) {
	if !contains(e.vals, c) {
		e.vals = append(e.vals, c)
	}
}

func (e *Enum) GetAll() []Const {
	return e.vals
}

func (e *Enum) unsafeSet(c Const) {
	e.val = c
}

func (e Enum) String() string {
	return string(e.Get())
}

// Unmarshalls the string into an Enum. NOTE: you must run enum.Validate
// after unmarshalling a string like so
//   func main() {
//     var money Money
//
//     json.Unmarshal([]byte("{\"currency_code\":\"USD\",\"amount\":5}"), &money)
//     enum.Validate(&money.CurrencyCode) // <-- Must be run after unmarshal
//
//     fmt.Println(money) // Prints "{USD 5}"
//   }
func (e *Enum) UnmarshalJSON(b []byte) error {
	s, _ := strconv.Unquote(string(b))
	c := Const(s)
	e.unsafeSet(c)
	return nil
}

func (e Enum) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(string(e.val))), nil
}

func (e *Enum) Get() Const {
	return e.val
}

func (e *Enum) Set(c Const) error {
	if e.vals != nil {
		if contains(e.GetAll(), c) {
			e.val = c
			return nil
		} else {
			return errors.New(fmt.Sprintf(invalidEnumErrorMsg, c))
		}
	} else {
		return errors.New(enumNotConstructedErrorMsg)
	}
}

func (e *Enum) MustSet(s Const) {
	err := e.Set(s)
	if err != nil {
		panic(err.Error())
	}
}

// Creates a new Enummer with no value set
//   cc := enum.New(new(CurrencyCodes)).(*CurrencyCodes)
func New(e Enummer) Enummer {
	construct(e)
	return e
}

// Instantiates an Enum with the provided value. If the value is invalid, an error is returned
// otherwise, an Enummer is returned with a nil error
//   cc, err := enum.Construct(new(CurrencyCodes), enum.Const("USD"))
//   if err != nil {
//     panic(err)
//   }
//   cc = cc.(*CurrencyCodes)
func Construct(e Enummer, c Const) (Enummer, error) {
	if e == nil {
		return nil, errors.New(enumNotNilErrorMsg)
	}
	construct(e)
	if err := e.Set(c); err != nil {
		return nil, err
	}
	return e, nil
}

// Instantiates an Enum with the provided value. If the value is invalid, a panic occurs
//   cc := enum.MustConstruct(new(CurrencyCodes), enum.Const("USD")).(*CurrencyCodes)
func MustConstruct(e Enummer, c Const) Enummer {
	out, err := Construct(e, c)

	if err != nil {
		panic(err.Error())
	}
	return out
}

// Instantiates the enum if that hasn't been done and validates that its current value is valid.
// Commonly used after unmarshalling an enum like so
//   func main() {
//     var money Money
//
//     json.Unmarshal([]byte("{\"currency_code\":\"USD\",\"amount\":5}"), &money)
//     enum.Validate(&money.CurrencyCode) // <-- Must be run after unmarshal
//
//     fmt.Println(money) // Prints "{USD 5}"
//   }
func Validate(e Enummer) error {
	if e.GetAll() == nil {
		construct(e)
	}

	if !contains(e.GetAll(), e.Get()) {
		return errors.New(fmt.Sprintf(invalidEnumErrorMsg, e.Get()))
	}

	return nil
}

func contains(cs []Const, c Const) bool {
	for _, v := range cs {
		if v == c {
			return true
		}
	}
	return false
}

func construct(e Enummer) {
	v := reflect.ValueOf(e).Elem()
	for i := 0; i < v.NumField(); i++ {
		t := v.Type().Field(i)
		f := v.Field(i)
		if _, ok := f.Interface().(Const); ok {
			s := t.Tag.Get("enum")
			if s == "" {
				s = t.Name
			}
			c := Const(s)
			e.unsafeAdd(c)
			f.Set(reflect.ValueOf(c))
		}
	}
}

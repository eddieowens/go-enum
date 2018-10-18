package tests

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go-enum"
	"testing"
)

type CurrencyCode struct {
	enum.Enum
	USD enum.Const `enum:"ASd"`
	DIA enum.Const
}

func TestMustConstructInvalidValue(t *testing.T) {
	asrt := assert.New(t)
	val := enum.Const("USD")

	asrt.PanicsWithValue("USD is not a valid enum", func() {
		enum.MustConstruct(new(CurrencyCode), val)
	})
}

func TestMustConstruct(t *testing.T) {
	asrt := assert.New(t)

	e := enum.MustConstruct(new(CurrencyCode), enum.Const("ASd")).(*CurrencyCode)

	asrt.Equal(enum.Const("ASd"), e.Get())
	asrt.Equal(enum.Const("ASd"), e.USD)
	asrt.Equal(enum.Const("DIA"), e.DIA)
}

func TestConstructInvalidValue(t *testing.T) {
	asrt := assert.New(t)

	c, err := enum.Construct(new(CurrencyCode), enum.Const("USD"))
	asrt.Nil(c)
	asrt.Equal("USD is not a valid enum", err.Error())
}

func TestConstruct(t *testing.T) {
	asrt := assert.New(t)

	c, err := enum.Construct(new(CurrencyCode), enum.Const("ASd"))
	e := c.(*CurrencyCode)

	asrt.Nil(err)
	asrt.Equal(enum.Const("ASd"), e.Get())
	asrt.Equal(enum.Const("ASd"), e.USD)
	asrt.Equal(enum.Const("DIA"), e.DIA)
}

func TestNew(t *testing.T) {
	asrt := assert.New(t)

	c, err := enum.Construct(new(CurrencyCode), enum.Const("ASd"))
	e := c.(*CurrencyCode)

	asrt.Nil(err)
	asrt.Equal(enum.Const("ASd"), e.Get())
	asrt.Equal(enum.Const("ASd"), e.USD)
	asrt.Equal(enum.Const("DIA"), e.DIA)
}

func TestMarshal(t *testing.T) {
	expected := "{\"currency_code\":\"ASd\",\"another_field\":1}"
	asrt := assert.New(t)

	a := enum.MustConstruct(new(CurrencyCode), enum.Const("ASd")).(*CurrencyCode)

	type test struct {
		CurrencyCode *CurrencyCode `json:"currency_code"`
		AnotherField int           `json:"another_field"`
	}

	c := test{
		CurrencyCode: a,
		AnotherField: 1,
	}

	out, err := json.Marshal(c)
	if err != nil {
		asrt.Fail(err.Error())
	}
	asrt.Equal(expected, string(out))
}

func TestInvalidUnmarshal(t *testing.T) {
	asrt := assert.New(t)

	type test struct {
		CurrencyCode CurrencyCode `json:"currency_code"`
		AnotherField int          `json:"another_field"`
	}

	var a test

	mErr := json.Unmarshal([]byte("{\"currency_code\": \"USD\"}"), &a)
	err := enum.Validate(&a.CurrencyCode)

	asrt.Nil(mErr)
	asrt.Equal("USD is not a valid enum", err.Error())
}

func TestValidUnmarshal(t *testing.T) {
	asrt := assert.New(t)

	type test struct {
		CurrencyCode CurrencyCode `json:"currency_code"`
		AnotherField int          `json:"another_field"`
	}

	var a test

	mErr := json.Unmarshal([]byte("{\"currency_code\": \"ASd\"}"), &a)
	err := enum.Validate(&a.CurrencyCode)

	asrt.Nil(mErr)
	asrt.Nil(err)
	asrt.Equal(a.CurrencyCode.USD, a.CurrencyCode.Get())
}

func TestString(t *testing.T) {
	asrt := assert.New(t)

	type test struct {
		CurrencyCode *CurrencyCode `json:"currency_code"`
		AnotherField int           `json:"another_field"`
	}

	c := enum.MustConstruct(new(CurrencyCode), enum.Const("ASd")).(*CurrencyCode)

	v := test{
		CurrencyCode: c,
		AnotherField: 1,
	}

	asrt.Equal("{ASd 1}", fmt.Sprintf("%v", v))
}

func TestSetInvalid(t *testing.T) {
	asrt := assert.New(t)

	c := enum.MustConstruct(new(CurrencyCode), enum.Const("ASd")).(*CurrencyCode)
	err := c.Set(enum.Const("garbage"))

	asrt.Equal("garbage is not a valid enum", err.Error())
}

func TestSetBeforeConstruct(t *testing.T) {
	asrt := assert.New(t)

	c := new(CurrencyCode)
	err := c.Set(c.USD)

	asrt.Equal("cannot set a value on an enum that has not be constructed", err.Error())
}

func TestGetAll(t *testing.T) {
	asrt := assert.New(t)

	c := enum.MustConstruct(new(CurrencyCode), enum.Const("ASd")).(*CurrencyCode)

	asrt.Equal([]enum.Const{"ASd", "DIA"}, c.GetAll())
}

func TestNewWithValue(t *testing.T) {
	asrt := assert.New(t)

	e := enum.New(new(CurrencyCode)).(*CurrencyCode)

	asrt.Equal(enum.Const(""), e.Get())
	asrt.Equal(enum.Const("ASd"), e.USD)
	asrt.Equal(enum.Const("DIA"), e.DIA)
}

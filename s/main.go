package main

import (
	"fmt"
	"go-enum"
)

type CurrencyCodes struct {
	enum.Enum
	USD    enum.Const
	EUR    enum.Const
	CAD    enum.Const
	Custom enum.Const `enum:"CUSTOM"`
}

func main() {
	cc := enum.MustConstruct(new(CurrencyCodes), enum.Const("CUSTOM")).(*CurrencyCodes)
	fmt.Println(cc)
}

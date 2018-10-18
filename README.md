# go-enum
A struct based golang enum package that handles stringifying, marshalling, and unmarshalling.

## Install
```bash
go get github.com/eddieowens/go-enum
```

## Usage
```go
type CurrencyCodes struct {
	enum.Enum
	USD    enum.Const
	EUR    enum.Const
	CAD    enum.Const
	Custom enum.Const `enum:"CUSTOM"`
}

func main() {
	// Create
	cc := enum.New(new(CurrencyCodes)).(*CurrencyCodes)
	
	// Instantiate then get
	cc = enum.MustConstruct(new(CurrencyCodes), enum.Const("USD")).(*CurrencyCodes)
	fmt.Println(cc.USD == cc.Get()) // Prints true
	
	// Set Valid
	cc.MustSet(Const("CUSTOM"))
	
	// Set invalid
	err := cc.Set(Const("Random"))
	fmt.Println(err.Error()) // Prints "Random is not a valid enum"
	
	// Marshal
    type Money struct {
        CurrencyCode CurrencyCode `json:"currency_code"`
        Amount       int          `json:"amount"`
    }
	
	m := Money{
		CurrencyCode: *cc,
		Amount: 5,
	}
	out, err := json.Marshal(c)
	fmt.Println(string(out)) // Prints "{"currency_code":"USD","amount":5}"
	
	// Unmarshal valid value
    var money Money

    json.Unmarshal([]byte("{\"currency_code\":\"USD\",\"amount\":5}"), &money)
    enum.Validate(&money.CurrencyCode) // <-- Must be run after unmarshal
    
    fmt.Println(money) // Prints "{USD 5}"
    
	// Unmarshal invalid value
    var money Money

    json.Unmarshal([]byte("{\"currency_code\":\"Random\",\"amount\":5}"), &money)
    err := enum.Validate(&money.CurrencyCode) // <-- Must be run after unmarshal
    
    fmt.Println(err.Error()) // Prints "Random is not a valid enum"
    
	// Stringify
	cc = enum.MustConstruct(new(CurrencyCodes), enum.Const("CUSTOM")).(*CurrencyCodes)
	fmt.Println(cc) // Prints "CUSTOM"
}
```

## Const type
The name of the field on the struct will be the default value for the enum const.
For example with `CurrencyCodes.USD` the `Const` value is "USD". In order to customize
this value, add the tag `enum:"<NAME>"` like in `CurrencyCodes.Custom`

## [Docs](https://godoc.org/github.com/eddieowens/go-enum)


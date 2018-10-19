# go-enum
A struct based golang enum package that handles stringifying, marshalling, and unmarshalling.

## Install
```bash
go get github.com/eddieowens/go-enum
```

## Usage
### Basic
The following creates an enum of currency codes and uses them on the `Money` struct
```go
type CurrencyCodes struct {
    enum.Enum
    USD    enum.Const
    EUR    enum.Const
    CAD    enum.Const
    Custom enum.Const `enum:"CUSTOM"`
}

type Money struct {
    CurrencyCode CurrencyCode `json:"currency_code"`
    Amount       int          `json:"amount"`
}
```

To create a currency code you must use one of the following functions
```go
// Instantiate with a provided value and panic if the value is not valid
cc := enum.MustConstruct(new(CurrencyCodes), enum.Const("USD")).(*CurrencyCodes)

// Instantiate with a provided value and return an error if the value is not valid
cc, err := enum.Construct(new(CurrencyCodes), enum.Const("USD")).(*CurrencyCodes)

// Create a new enum with no value set
cc := enum.New(new(CurrencyCodes)).(*CurrencyCodes)
```

To get the value and compare it with another
```go
cc.Get() == cc.USD
```

To set the value, use either of the following
```go
// If the value is invalid,  an error is returned
err := cc.Set(cc.CAD)

// If the value is invalid, a panic occurs
cc.MustSet(cc.CAD)
```



### Complete example
```go
func main() {
    // Create without a value
    cc := enum.New(new(CurrencyCodes)).(*CurrencyCodes)
    
    // Instantiate with a provided value then get it
    cc = enum.MustConstruct(new(CurrencyCodes), enum.Const("USD")).(*CurrencyCodes)
    fmt.Println(cc.USD == cc.Get()) // Prints true
    
    // Set valid value
    cc.MustSet(cc.Custom)
    
    // Set invalid value
    err := cc.Set(Const("Random"))
    fmt.Println(err.Error()) // Prints "Random is not a valid enum"
    
    // Get all possible enums
    consts := cc.GetAll()
    for _, c := range consts{
    	...
    }
    
    // Marshal
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

## To note
### Unmarshalling
All unmarshalling needs to be followed with a call to `enum.Validate(...)`. If it is not
the value placed on the enum may not be valid and the enum will not function as expected.

### Const type
The name of the field on the struct will be the default value for the enum const.
For example with `CurrencyCodes.USD` the `Const` value is "USD". In order to customize
this value, add the tag `enum:"<NAME>"` like in `CurrencyCodes.Custom`

## [Docs](https://godoc.org/github.com/eddieowens/go-enum)

## License
[MIT](https://github.com/eddieowens/go-enum/blob/master/LICENSE)

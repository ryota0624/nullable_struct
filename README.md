# nullable-struct

Nullableなstructを定義するためのジェネレータ

```
//go:generate go run github.com/ryota0624/nullable_struct -package=$GOPACKAGE -type=Data -dest=nullable_$GOFILE

type (
	Data struct {
		A string `json:"a"`
		B int    `json:"b"`
	}
)
```

のように使えます。
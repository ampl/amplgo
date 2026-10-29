# amplgo

Go bindings for the [AMPL API](https://ampl.com/api/), built on the AMPL C API (`libampl`) via cgo.

## Requirements

- Go 1.20 or newer, with a C compiler available for cgo
- An AMPL installation that includes the C API (`amplapi/include` and `amplapi/lib`) and a valid license

## Installation

```bash
go get github.com/ampl/amplgo
```

The C API is not bundled with this module. Tell cgo where to find it, and make the shared library
and the `ampl` executable available at run time:

```bash
export AMPL_DIR=/path/to/ampl.linux-intel64
export CGO_CFLAGS="-I$AMPL_DIR/amplapi/include"
export CGO_LDFLAGS="-L$AMPL_DIR/amplapi/lib"
export LD_LIBRARY_PATH="$AMPL_DIR/amplapi/lib:$LD_LIBRARY_PATH"   # DYLD_LIBRARY_PATH on macOS
export PATH="$AMPL_DIR:$PATH"
```

## Usage

```go
package main

import (
	"fmt"

	"github.com/ampl/amplgo/ampl"
)

func main() {
	model, err := ampl.New()
	if err != nil {
		panic(err)
	}
	if err := model.Read("models/diet/diet.mod"); err != nil {
		panic(err)
	}
	if err := model.ReadData("models/diet/diet.dat"); err != nil {
		panic(err)
	}
	if err := model.Solve("", "highs"); err != nil {
		panic(err)
	}
	variant, err := model.GetValue("Total_Cost")
	if err != nil {
		panic(err)
	}
	totalCost, _ := variant.GetNumericValue()
	fmt.Println("Total Cost:", totalCost)
}
```

More examples are in [examples/](examples/). Run them from the repository root so the relative
model paths resolve:

```bash
go run ./examples/quickstart
go run ./examples/firstexample highs
```

## Development

`dev/updatelib.py` downloads the C API headers and libraries into `ampl/include` and `ampl/lib`
(both git-ignored). Point the cgo variables there to build against them:

```bash
python dev/updatelib.py
export CGO_CFLAGS="-I$PWD/ampl/include"
export CGO_LDFLAGS="-L$PWD/ampl/lib/amd64"
export LD_LIBRARY_PATH="$PWD/ampl/lib/amd64:$LD_LIBRARY_PATH"
go test ./tests
```

## License

MIT, see [LICENSE](LICENSE). The AMPL C API and the AMPL executable are distributed separately under
their own license terms.

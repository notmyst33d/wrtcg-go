//go:build debug

package debug

import (
	"encoding/json"
	"fmt"
)

func Inspect(v any) {
	data, _ := json.Marshal(v)
	fmt.Printf("debug: %s\n", string(data))
}

package schemax

import (
	"fmt"
	"testing"
)

/*
This example demonstrates use of the [Options.Positive] method to
determine whether a specific bit value registers as positive, or
enabled.
*/
func ExampleOptions_Positive() {
	opts := mySchema.Options()
	opts.Shift(AllowOverride) // FYR: AllowOverride is 8
	fmt.Println(opts.Positive(32))
	// Output: false
}

/*
This example demonstrates use of the [Options.Shift] method to
set a particular bit value to a positive (enabled) state.
*/
func ExampleOptions_Shift() {
	opts := mySchema.Options()
	opts.Shift(AllowOverride)
	fmt.Println(opts.Positive(AllowOverride))
	// Output: true
}

/*
This example demonstrates use of the [Options.Unshift] method to
set a particular bit value to a negative (disabled) state.
*/
func ExampleOptions_Unshift() {
	opts := mySchema.Options()
	opts.Shift(AllowOverride)
	opts.Unshift(AllowOverride)
	fmt.Println(opts.Positive(AllowOverride))
	// Output: false
}

func TestImports_codecov(t *testing.T) {
	_ = uitoa(uint64(64))
}

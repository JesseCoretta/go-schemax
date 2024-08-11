package schemax

import (
	"fmt"
	"testing"
)

func ExampleDefinitionMaps_Len() {
	classes := mySchema.ObjectClasses()
	fmt.Printf("%d definitions", classes.Maps().Len())
	// Output: 80 definitions
}

func ExampleDefinitionMap_Len() {
	classes := mySchema.ObjectClasses()
	fmt.Printf("ObjectClass map fields: %d", classes.Maps().Index(0).Len())
	// Output: ObjectClass map fields: 8
}

func ExampleDefinitionMap_Type() {
	classes := mySchema.ObjectClasses()
	fmt.Printf("Map type: %s", classes.Maps().Index(0).Type())
	// Output: Map type: objectClass
}

func ExampleDefinitionMap_Keys() {
	classes := mySchema.ObjectClasses()
	fmt.Printf("ObjectClass map keys: %d", len(classes.Maps().Index(0).Keys()))
	// Output: ObjectClass map keys: 8
}

func ExampleDefinitionMaps_IsZero() {
	classes := mySchema.ObjectClasses()
	fmt.Printf("ObjectClass maps zero: %t", classes.Maps().IsZero())
	// Output: ObjectClass maps zero: false
}

func ExampleDefinitionMap_IsZero() {
	classes := mySchema.ObjectClasses()
	fmt.Printf("ObjectClass map zero: %t", classes.Maps().Index(0).IsZero())
	// Output: ObjectClass map zero: false
}

func ExampleDefinitionMaps_Index() {
	classes := mySchema.ObjectClasses()
	fmt.Printf("%s", classes.Maps().Index(2)[`NUMERICOID`][0])
	// Output: 2.5.20.1
}

func ExampleDefinitionMap_Get() {
	classes := mySchema.ObjectClasses()
	fmt.Printf("%s", classes.Maps().Index(2).Get(`NAME`)[0])
	// Output: subschema
}

func ExampleDefinitionMap_Contains() {
	classes := mySchema.ObjectClasses()
	fmt.Printf("%t", classes.Maps().Index(2).Contains(`NAME`))
	// Output: true
}

func TestDefinitionMaps_codecov(t *testing.T) {
	classes := mySchema.ObjectClasses()
	_ = classes.Maps().Index(222)
}

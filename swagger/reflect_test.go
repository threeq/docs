package swagger

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"reflect"
	"testing"

	"fmt"

	"github.com/stretchr/testify/assert"
)

type Person struct {
	First string
}

type Pet struct {
	Friend      Person    `json:"friend"`
	Friends     []Person  `json:"friends"`
	Pointer     *Person   `json:"pointer" required:"true"`
	Pointers    []*Person `json:"pointers"`
	Int         int
	IntArray    []int
	String      string
	StringArray []string

	unexported string
}

type Empty struct {
	Nope int `json:"-"`
}

type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Inline struct {
	Inline string `json:"expected_outline"`
}

type Outline struct {
	InlineField Inline `json:",inline"`
}

func TestDefine(t *testing.T) {
	v := define("", Pet{})
	obj, ok := v["swaggerPet"]
	assert.True(t, ok)
	assert.False(t, obj.IsArray)
	assert.Equal(t, 8, len(obj.Properties))

	content := map[string]Object{}
	data, err := ioutil.ReadFile("testdata/pet.json")
	assert.Nil(t, err)
	err = json.NewDecoder(bytes.NewReader(data)).Decode(&content)
	assert.Nil(t, err)
	expected := content["swaggerPet"]

	assert.Equal(t, expected.IsArray, obj.IsArray, "expected IsArray to match")
	assert.Equal(t, expected.Type, obj.Type, "expected Type to match")
	assert.Equal(t, expected.Required, obj.Required, "expected Required to match")
	assert.Equal(t, len(expected.Properties), len(obj.Properties), "expected same number of properties")

	for k, v := range obj.Properties {
		e := expected.Properties[k]
		assert.Equal(t, e.Type, v.Type, "expected %v.Type to match", k)
		assert.Equal(t, e.Description, v.Description, "expected %v.Required to match", k)
		assert.Equal(t, e.Enum, v.Enum, "expected %v.Required to match", k)
		assert.Equal(t, e.Format, v.Format, "expected %v.Required to match", k)
		assert.Equal(t, e.Ref, v.Ref, "expected %v.Required to match", k)
		assert.Equal(t, e.Example, v.Example, "expected %v.Required to match", k)
		assert.Equal(t, e.Items, v.Items, "expected %v.Required to match", k)
	}
}

func TestNotStructDefine(t *testing.T) {
	v := define("", int32(1))
	obj, ok := v["int32"]
	assert.True(t, ok)
	assert.False(t, obj.IsArray)
	assert.Equal(t, "integer", obj.Type)
	assert.Equal(t, "int32", obj.Format)

	v = define("", uint64(1))
	obj, ok = v["uint64"]
	assert.True(t, ok)
	assert.False(t, obj.IsArray)
	assert.Equal(t, "integer", obj.Type)
	assert.Equal(t, "int64", obj.Format)

	v = define("", "")
	obj, ok = v["string"]
	assert.True(t, ok)
	assert.False(t, obj.IsArray)
	assert.Equal(t, "string", obj.Type)
	assert.Equal(t, "", obj.Format)

	v = define("", byte(1))
	obj, ok = v["uint8"]
	if !assert.True(t, ok) {
		fmt.Printf("%v", v)
	}
	assert.False(t, obj.IsArray)
	assert.Equal(t, "integer", obj.Type)
	assert.Equal(t, "int32", obj.Format)

	v = define("", []byte{1, 2})
	obj, ok = v["uint8"]
	if !assert.True(t, ok) {
		fmt.Printf("%v", v)
	}
	assert.True(t, obj.IsArray)
	assert.Equal(t, "integer", obj.Type)
	assert.Equal(t, "int32", obj.Format)
}

func TestHonorJsonIgnore(t *testing.T) {
	v := define("", Empty{})
	obj, ok := v["swaggerEmpty"]
	assert.True(t, ok)
	assert.False(t, obj.IsArray)
	assert.Equal(t, 0, len(obj.Properties), "expected zero exposed properties")
}

func TestJsonInline(t *testing.T) {
	v := define("", Outline{})
	obj, ok := v["swaggerOutline"]
	assert.True(t, ok)
	assert.False(t, obj.IsArray)
	assert.Equal(t, 1, len(obj.Properties), "expected 1 exposed properties")
	assert.Contains(t, obj.Properties, "expected_outline", "expected contains expected_outline field")
}

func TestIgnoreUnexported(t *testing.T) {
	type Test struct {
		Exported   string
		unexported string
	}
	v := define("", Test{})
	obj, ok := v["swaggerTest"]
	assert.True(t, ok)
	assert.Equal(t, 1, len(obj.Properties), "expected one exposed properties")
	assert.Contains(t, obj.Properties, "Exported")
	assert.NotContains(t, obj.Properties, "unexported")
}

func TestInterfaceActualType(t *testing.T) {
	v := define("", APIResponse{Data: Pet{}})
	obj, ok := v["swaggerAPIResponse"]
	assert.True(t, ok)
	assert.False(t, obj.IsArray)
	assert.Equal(t, 3, len(obj.Properties), "expected 3 exposed properties")
	data := obj.Properties["data"]
	assert.Equal(t, reflect.TypeOf(Pet{}), data.GoType)

	v = define("", APIResponse{Data: Person{}})
	obj, ok = v["swaggerAPIResponse"]
	assert.True(t, ok)
	assert.False(t, obj.IsArray)
	assert.Equal(t, 3, len(obj.Properties), "expected 3 exposed properties")
	data = obj.Properties["data"]
	assert.Equal(t, reflect.TypeOf(Person{}), data.GoType)

	v = define("", APIResponse{Data: &Person{}})
	obj, ok = v["swaggerAPIResponse"]
	assert.True(t, ok)
	assert.False(t, obj.IsArray)
	assert.Equal(t, 3, len(obj.Properties), "expected 3 exposed properties")
	data = obj.Properties["data"]
	assert.Equal(t, reflect.TypeOf(Person{}), data.GoType)
}

package msg

import (
	"fmt"
)

type TechxDataTest struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	Job  string `json:"job,omitempty"`
}

func (techxDataTest *TechxDataTest) Deserialize(str string) *TechxDataTest {
	return deserialize(str, techxDataTest).(*TechxDataTest)
}

func ExampleMsg() {
	testData := TechxDataTest{
		Name: "mo",
		Age:  20,
		Job:  "xx",
	}
	msg, _ := NewTechxRequestMsg("test", "/convert", testData)

	deserializedMsg := (&TechxMsg{}).Deserialize(msg)

	data := (&TechxDataTest{}).Deserialize(Serialize(deserializedMsg.Data))
	fmt.Println(msg, data)
	// Output: xx
}

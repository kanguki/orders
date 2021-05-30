package msg

import "os"

var NODE_ID = os.Getenv("NODE_ID")
var APPLICATION_NAME = os.Getenv("APPLICATION_NAME")

type Msg interface {
	Deserialize(str string) interface{}
}

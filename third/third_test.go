package third

import "testing"

func Test_String(t *testing.T) {
	New().stringCmd()
}

func Test_List(t *testing.T) {
	New().listCmd()
}

func Test_Set(t *testing.T) {
	New().setCmd()
}

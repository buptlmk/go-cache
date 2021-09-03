package log

import "testing"

func TestPrint(t *testing.T) {

	Logger.Error("dfzgdf")
	Logger.WithField("a", "sdgf").Info("sdgfds")
}

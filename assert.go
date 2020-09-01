package protocmp

import (
	"fmt"
	"path"
	"runtime"
	"strings"
	"testing"

	"github.com/golang/protobuf/proto"
)

func AssertEqual(t *testing.T, expected, actual proto.Message) {
	err := Equal(expected, actual)
	if err != nil {
		frame := getFrame(1)
		fmt.Printf("    %s: %s:%d\n        %s\n", t.Name(), path.Base(frame.File), frame.Line, strings.ReplaceAll(err.Error(), "\n", "\n            "))
		t.Fail()
	}
}

func getFrame(skipFrames int) runtime.Frame {
	// We need the frame at index skipFrames+2, since we never want runtime.Callers and getFrame
	targetFrameIndex := skipFrames + 2

	// Set size to targetFrameIndex+2 to ensure we have room for one more caller than we need
	programCounters := make([]uintptr, targetFrameIndex+2)
	n := runtime.Callers(0, programCounters)

	frame := runtime.Frame{Function: "unknown"}
	if n > 0 {
		frames := runtime.CallersFrames(programCounters[:n])
		for more, frameIndex := true, 0; more && frameIndex <= targetFrameIndex; frameIndex++ {
			var frameCandidate runtime.Frame
			frameCandidate, more = frames.Next()
			if frameIndex == targetFrameIndex {
				frame = frameCandidate
			}
		}
	}

	return frame
}
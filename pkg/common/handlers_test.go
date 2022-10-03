package common

import (
	"fmt"
	"os"
	"testing"
)

func Test_ErrorHandler_NonFatal(t *testing.T) {
	testErr := fmt.Errorf("non fatal error occured")
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("error, unexpected panic in ErrorHandler")
			}
		}()
	}()

	ErrorHandler(testErr, false)
}

func Test_ErrorHandler_Fatal(t *testing.T) {
	testErr := fmt.Errorf("fatal error occured")
	defer func() {
		if r := recover(); r == nil {
			t.Error("error, expected panic in ErrorHandler")
		}
	}()

	ErrorHandler(testErr, true)
}

func Test_OpenFile(t *testing.T) {
	os.Create("/tmp/tempfile.txt")
	defer os.Remove("/tmp/tempfile.txt")

	defer func() {
		if r := recover(); r != nil {
			t.Error("error, unexpected panic in ErrorHandler")
		}
	}()

	file := OpenFile("/tmp/tempfile.txt")
	defer file.Close()
}

func Test_OpenFile_Err(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("error, expected panic in ErrorHandler")
		}
	}()

	file := OpenFile("/tmp/tempfile.txt")
	if file != nil {
		t.Error("error, unexpected non-nil file pointer")
	}
}

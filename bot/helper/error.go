package helper

import "os"

func OnError(str string, err error) {
	prefix := "Error: " + str

	if err != nil {
		os.Stderr.WriteString(prefix + ": " + err.Error() + "\n")
	} else {
		os.Stderr.WriteString(prefix)
	}
}

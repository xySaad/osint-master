package lib

import (
	"fmt"
	"io"
	"os"
)

// clones stdout to outputFile
func CloneOutput(outputFilePath string) (func(), error) {
	outputFile, err := os.OpenFile(outputFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o622)
	if err != nil {
		return nil, fmt.Errorf("Error opening output file: %v", err)
	}
	reader, writer, err := os.Pipe()
	if err != nil {
		return nil, fmt.Errorf("Error creating pipe: %v", err)
	}
	stdout := os.Stdout
	os.Stdout = writer
	mw := io.MultiWriter(stdout, outputFile)

	done := make(chan struct{})
	go func() {
		defer close(done)
		_, err := io.Copy(mw, reader)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error copying stdout: %v\n", err)
		}
	}()

	close := func() {
		os.Stdout = stdout
		writer.Close()
		<-done
		outputFile.Close()
		reader.Close()
	}
	return close, nil
}

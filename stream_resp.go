package gpt3

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/pkg/errors"
)

func StreamOnData(streamData io.ReadCloser, output CompletionResponseInterface, onData func(CompletionResponseInterface)) error {
	reader := newEventStreamReader(streamData, 1<<16)
	defer streamData.Close()

LOOP:
	for {
		event, err := reader.ReadEvent()
		if err != nil {
			if err == io.EOF {
				break LOOP
			}
			return errors.Wrap(err, "ReadEvent")
		}

		// If we get an error, ignore it.
		var msg *Event
		if msg, err = processEvent(event); err != nil {
			return errors.Wrap(err, "ProcessEvent")
		}
		output.Reset()
		if bytes.Equal(msg.Data, doneSequence) {
			break LOOP
		}
		if err := json.Unmarshal(msg.Data, output); err != nil {
			return errors.Errorf("invalid json stream data: %v", err)
		}

		onData(output)
	}
	return nil
}

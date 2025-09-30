package datatrails

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
)

var (
	ErrNoJsonGiven = errors.New("no json given")
)

// eventListFromJson normalises a json encoded event or *list* of events, by
// always returning a list of json encoded events.
//
// This converts events from the following apps:
// - assetsv2
// - eventsv1
//
// NOTE: there is no json validation done on the event or list of events given
// any valid json will be accepted, use validation logic after this function.
func EventListFromJson(data []byte) ([]byte, error) {
	var err error

	doc := struct {
		Events        []json.RawMessage `json:"events,omitempty"`
		NextPageToken json.RawMessage   `json:"next_page_token,omitempty"`
	}{}

	// check for empty json
	// NOTE: also len(nil) == 0, so does the nil check also
	if len(data) == 0 {
		return nil, ErrNoJsonGiven
	}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()

	for {

		err = decoder.Decode(&doc)

		// if we can decode the events json
		//  we know its in the form of a list events json response from
		//  the list events api, so just return data
		if errors.Is(err, io.EOF) {
			return data, nil
		}

		if err != nil {
			break
		}

	}

	// if we get here we know that the given data doesn't represent
	//  a list events json response
	// so we can assume its a single event response from the events api.

	var event json.RawMessage
	err = json.Unmarshal(data, &event)
	if err != nil {
		return nil, err
	}

	// purposefully omit the next page token for response
	listEvents := struct {
		Events []json.RawMessage `json:"events,omitempty"`
	}{}

	listEvents.Events = []json.RawMessage{event}

	events, err := json.Marshal(&listEvents)
	if err != nil {
		return nil, err
	}

	return events, nil
}

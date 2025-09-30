package appentry

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"strings"

	"github.com/forestrie/go-merklelog-datatrails/datatrails"
)

// AppDomain returns the app domain of the given app data
func AppDomain(appData []byte) byte {

	// first attempt to convert the appdata to a list of events
	eventList, err := datatrails.EventListFromJson(appData)
	if err != nil {
		// if we can't return default of assetsv2
		return AssetsV2AppDomain
	}

	// decode into events
	events := struct {
		Events        []json.RawMessage `json:"events,omitempty"`
		NextPageToken json.RawMessage   `json:"next_page_token,omitempty"`
	}{}

	decoder := json.NewDecoder(bytes.NewReader(eventList))
	decoder.DisallowUnknownFields()
	for {
		err = decoder.Decode(&events)

		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			// return default of assetsv2
			return AssetsV2AppDomain
		}
	}

	// decode the first event and find the identity
	event := struct {
		Identity string `json:"identity,omitempty"`
	}{}

	decoder = json.NewDecoder(bytes.NewReader(events.Events[0]))

	for {
		err = decoder.Decode(&event)

		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			// if we can't return default of assetsv2
			return AssetsV2AppDomain
		}
	}

	// find if the event identity is assetsv2 or eventsv1 identity
	if strings.HasPrefix(event.Identity, "assets/") || strings.HasPrefix(event.Identity, "publicassets/") {
		return AssetsV2AppDomain
	} else {
		return EventsV1AppDomain
	}

}

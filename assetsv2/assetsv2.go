package assetsv2

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"sort"
	"strings"

	"github.com/datatrails/go-datatrails-common-api-gen/assets/v2/assets"
	"github.com/forestrie/go-merklelog-datatrails/appentry"
	"github.com/forestrie/go-merklelog-datatrails/v3hash"
	"github.com/google/uuid"
	"google.golang.org/protobuf/encoding/protojson"
)

var (
	ErrInvalidAssetsV2EventJson = errors.New(`invalid assetsv2 event json`)
	ErrNoEvents                 = errors.New(`no events found in events json`)
	ErrNoJsonGiven              = errors.New("no json given")
)

func VerifiableAssetsV2EventsFromData(data []byte) ([]appentry.AppEntry, error) {

	// Accept either the list events response format or a single event. Peak
	// into the json data to pick which.
	eventsJson, err := eventListFromJson(data)
	if err != nil {
		return nil, err
	}

	verifiableEvents, err := NewAssetsV2AppEntries(eventsJson)
	if err != nil {
		return nil, err
	}

	for _, event := range verifiableEvents {
		validationErr := v3hash.Validate(event)
		if validationErr != nil {
			return nil, validationErr
		}
	}

	return verifiableEvents, nil
}

// NewAssetsV2AppEntries takes a list of events JSON (e.g. from the assetsv2 events list API), converts them
// into AssetsV2AppEntries and then returns them sorted by ascending MMR index.
func NewAssetsV2AppEntries(eventsJson []byte) ([]appentry.AppEntry, error) {
	// get the event list out of events
	eventListJson := struct {
		Events []json.RawMessage `json:"events"`
	}{}

	err := json.Unmarshal(eventsJson, &eventListJson)
	if err != nil {
		return nil, err
	}

	events := []appentry.AppEntry{}
	for _, eventJson := range eventListJson.Events {
		verifiableEvent, err := NewAssetsV2AppEntry(eventJson)
		if err != nil {
			return nil, err
		}

		events = append(events, *verifiableEvent)
	}

	// check if we haven't got any events
	if len(events) == 0 {
		return nil, ErrNoEvents
	}

	// Sorting the events by MMR index guarantees that they're sorted in log append order.
	sort.Slice(events, func(i, j int) bool {
		return events[i].MMRIndex() < events[j].MMRIndex()
	})

	return events, nil
}

// NewAssetsV2AppEntry takes a single assetsv2 event JSON and returns an AssetsV2AppEntry,
// providing just enough information to verify the incluson of and identify the event.
func NewAssetsV2AppEntry(eventJson []byte) (*appentry.AppEntry, error) {

	// special care is needed here to deal with uint64 types. json marshal /
	// un marshal treats them as strings because they don't fit in a
	// javascript Number

	// Unmarshal into a generic type to get just the bits we need. Use
	// defered decoding to get the raw merklelog entry as it must be
	// unmarshaled using protojson and the specific generated target type.
	entry := struct {
		Identity       string `json:"identity,omitempty"`
		TenantIdentity string `json:"tenant_identity,omitempty"`
		// Note: the proof_details top level field can be ignored here because it is a 'oneof'
		MerklelogEntry json.RawMessage `json:"merklelog_entry,omitempty"`
	}{}
	err := json.Unmarshal(eventJson, &entry)
	if err != nil {
		return nil, err
	}

	merkleLog := &assets.MerkleLogEntry{}
	err = protojson.Unmarshal(entry.MerklelogEntry, merkleLog)
	if err != nil {
		return nil, err
	}

	if entry.TenantIdentity == "" {
		return nil, ErrInvalidAssetsV2EventJson
	}

	// get the logID from the event log tenant
	logUuid := strings.TrimPrefix(entry.TenantIdentity, "tenant/")
	logId, err := uuid.Parse(logUuid)
	if err != nil {
		return nil, err
	}

	return appentry.NewAppEntry(
		entry.Identity,
		logId[:],
		appentry.NewMMREntryFields(
			byte(0),
			eventJson, // we cheat a bit here, because the eventJson isn't really serialized, but its a log version 0 log entry
		),
		merkleLog.Commit.Index,
	), nil
}

// eventListFromJson normalises a json encoded event or *list* of events, by
// always returning a list of json encoded events.
//
// This converts events from the following apps:
// - assetsv2
// - eventsv1
//
// NOTE: there is no json validation done on the event or list of events given
// any valid json will be accepted, use validation logic after this function.
func eventListFromJson(data []byte) ([]byte, error) {
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

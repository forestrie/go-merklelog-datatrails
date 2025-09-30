package appdata

import (
	"errors"

	"github.com/forestrie/go-merklelog-datatrails/appentry"
	"github.com/forestrie/go-merklelog-datatrails/assetsv2"
	"github.com/forestrie/go-merklelog-datatrails/datatrails"
	"github.com/forestrie/go-merklelog-datatrails/eventsv1"
)

// AppDataToVerifiableLogEntries converts the app data (one or more app entries) to verifiable log entries
func AppDataToVerifiableLogEntries(appData []byte, logTenant string) ([]appentry.AppEntry, error) {

	// first attempt to convert the appdata to a list of events
	eventList, err := datatrails.EventListFromJson(appData)
	if err != nil {
		return nil, err
	}

	// now we have an event list we can decipher if the app is
	//  assetsv2 or eventsv1
	appDomain := appentry.AppDomain(appData)

	verifiableLogEntries := []appentry.AppEntry{}

	switch appDomain {
	case appentry.AssetsV2AppDomain:
		// assetsv2
		verifiableAssetsV2Events, err := assetsv2.NewAssetsV2AppEntries(eventList)
		if err != nil {
			return nil, err
		}

		verifiableLogEntries = append(verifiableLogEntries, verifiableAssetsV2Events...)

	case appentry.EventsV1AppDomain:
		verifiableEventsV1Events, err := eventsv1.NewEventsV1AppEntries(eventList, logTenant)
		if err != nil {
			return nil, err
		}

		verifiableLogEntries = append(verifiableLogEntries, verifiableEventsV1Events...)

	default:
		return nil, errors.New("unknown app domain for given app data")
	}

	return verifiableLogEntries, nil
}

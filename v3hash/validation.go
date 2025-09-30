package v3hash

import (
	"errors"

	"github.com/forestrie/go-merklelog-datatrails/appentry"
)

// Note: We need this logic to detect incomplete JSON unmarshalled into these types. This should
// eventually be replaced by JSON Schema validation. We believe its a problem to solve for the
// entire go codebase through generation. We already describe the structs with json annotations and
// typing information. We don't want to half-cook that solution, as JSON type bugs have bitten us
// before.

var (
	ErrNonEmptyAppIDRequired    = errors.New("app id field is required and must be non-empty")
	ErrNonEmptyEventIDRequired  = errors.New("event identity field is required and must be non-empty")
	ErrNonEmptyTenantIDRequired = errors.New("tenant identity field is required and must be non-empty")
	ErrCommitEntryRequired      = errors.New("merkle log commit field is required")
	ErrIdTimestampRequired      = errors.New("idtimestamp field is required and must be non-empty")
)

// Validate performs basic validation on the AppEntryGetter, ensuring that critical fields
// are present.
func Validate(appEntry appentry.AppEntry) error {
	if appEntry.AppID() == "" {
		return ErrNonEmptyAppIDRequired
	}

	if len(appEntry.LogID()) == 0 {
		return ErrNonEmptyTenantIDRequired
	}

	return nil
}

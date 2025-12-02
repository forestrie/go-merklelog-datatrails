package datatrails

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/forestrie/go-merklelog/massifs/storage"
)

var ErrMassifPathFmt = errors.New("invalid massif path")

// XXX: NOTE: Just staging these functions here while the open sourcing effort is in flight
// LogID from the storage path according to the datatrails massif storage schema.
// The storage path is expected to be in the format:
// /v1/mmrs/tenant/<tenant_uuid>/<log_instance>/massifs/
// or
// /v1/mmrs/tenant/<tenant_uuid>/<log_instance>/massifseals/
func StorageLogID(storagePath string) (storage.LogID, error) {
	logID := TenantID2LogID(storagePath)
	if logID != nil {
		return logID, nil
	}

	return nil, fmt.Errorf("invalid storage path prefix: %s", storagePath)
}

// IsMassifPathLike performs a shallow sanity check on a path to see if it could be a massif log path
// Supports both v1 (v1/mmrs/tenant/...) and v2 (v2/merklelog/massifs/...) formats
func IsMassifPathLike(path string) bool {
	// Check v2 format: v2/merklelog/massifs/...
	if strings.HasPrefix(path, V2MerklelogMassifsPrefix+"/") {
		return strings.HasSuffix(path, storage.V1MMRMassifExt)
	}
	// Check v1 format: v1/mmrs/tenant/...
	if strings.HasPrefix(path, V1MMRTenantPrefix) {
		return strings.HasSuffix(path, storage.V1MMRMassifExt)
	}
	return false
}

// IsSealPathLike performs a shallow sanity check on a path to see if it could be a massif seal path
// Supports both v1 (v1/mmrs/tenant/.../massifseals/...) and v2 (v2/merklelog/checkpoints/...) formats
func IsSealPathLike(path string) bool {
	// Check v2 format: v2/merklelog/checkpoints/...
	if strings.HasPrefix(path, V2MerklelogCheckpointsPrefix+"/") {
		return strings.HasSuffix(path, storage.V1MMRSealSignedRootExt)
	}
	// Check v1 format: v1/mmrs/tenant/.../massifseals/...
	if strings.HasPrefix(path, V1MMRTenantPrefix) {
		return strings.HasSuffix(path, storage.V1MMRSealSignedRootExt)
	}
	return false
}

// ParseMassifPathTenant parse the tenant uuid from a massif storage path
// Performs basic sanity checks
// Supports both v1 and v2 path formats
func ParseMassifPathTenant(path string) (string, error) {
	parts := strings.Split(path, storage.V1MMRPathSep)
	
	// v2 format: v2/merklelog/massifs/{massifHeight}/{uuid}/{index}.log
	// or v2/merklelog/checkpoints/{massifHeight}/{uuid}/{index}.sth
	if len(parts) >= 5 && parts[0] == "v2" && parts[1] == "merklelog" {
		// uuid is at index 4 (after v2/merklelog/massifs/{massifHeight})
		if len(parts) < 5 {
			return "", fmt.Errorf("invalid massif path: %s", path)
		}
		return parts[4], nil
	}

	// v1 format: v1/mmrs/tenant/{uuid}/0/massifs/{index}.log
	if !strings.HasPrefix(path, V1MMRTenantPrefix) {
		return "", fmt.Errorf("invalid massif path: %s", path)
	}

	// the +1 strips the leading /
	path = path[len(V1MMRTenantPrefix)+1:]

	parts = strings.Split(path, storage.V1MMRPathSep)
	if len(parts) == 0 {
		return "", fmt.Errorf("invalid massif path: %s", path)
	}
	// we could parse the uuid, but that seems like over kill
	return parts[0], nil
}

// ParseMassifPathNumberExt parse the log file number and extension from the storage path
// Performs basic sanity checks
// Supports both v1 and v2 path formats
func ParseMassifPathNumberExt(path string) (uint32, string, error) {
	parts := strings.Split(path, storage.V1MMRPathSep)
	if len(parts) == 0 {
		return 0, "", fmt.Errorf("%w: %s", ErrMassifPathFmt, path)
	}

	// Check if it's v2 format: v2/merklelog/massifs/{massifHeight}/{uuid}/{index}.log
	// or v2/merklelog/checkpoints/{massifHeight}/{uuid}/{index}.sth
	if len(parts) >= 5 && parts[0] == "v2" && parts[1] == "merklelog" {
		base := parts[len(parts)-1]
		baseParts := strings.Split(base, storage.V1MMRExtSep)
		if len(baseParts) != 2 {
			return 0, "", fmt.Errorf("%w: base name invalid %s", ErrMassifPathFmt, path)
		}
		if baseParts[1] != storage.V1MMRMassifExt && baseParts[1] != storage.V1MMRSealSignedRootExt {
			return 0, "", fmt.Errorf("%w: extension invalid %s", ErrMassifPathFmt, path)
		}
		// Parse as hex (16 digits) for v2 format
		number, err := strconv.ParseUint(baseParts[0], 16, 32)
		if err != nil {
			return 0, "", fmt.Errorf("%w: log file number invalid %s (%v)", ErrMassifPathFmt, path, err)
		}
		return uint32(number), baseParts[1], nil
	}

	// v1 format: v1/mmrs/tenant/{uuid}/0/massifs/{index}.log
	if !strings.HasPrefix(path, V1MMRTenantPrefix) {
		return 0, "", fmt.Errorf("%w: %s", ErrMassifPathFmt, path)
	}
	base := parts[len(parts)-1]
	baseParts := strings.Split(base, storage.V1MMRExtSep)
	if len(baseParts) != 2 {
		return 0, "", fmt.Errorf("%w: base name invalid %s", ErrMassifPathFmt, path)
	}
	if baseParts[1] != storage.V1MMRMassifExt && baseParts[1] != storage.V1MMRSealSignedRootExt {
		return 0, "", fmt.Errorf("%w: extension invalid %s", ErrMassifPathFmt, path)
	}
	// Parse as decimal for v1 format
	number, err := strconv.ParseUint(baseParts[0], 10, 32)
	if err != nil {
		return 0, "", fmt.Errorf("%w: log file number invalid %s (%v)", ErrMassifPathFmt, path, err)
	}
	return uint32(number), baseParts[1], nil
}

// ParseMassifPathHeight extracts the massifHeight from a storage path.
// Supports both v1 (returns 0 or reads from path) and v2 (extracts from path) formats.
func ParseMassifPathHeight(path string) (uint8, error) {
	parts := strings.Split(path, storage.V1MMRPathSep)
	
	// v2 format: v2/merklelog/massifs/{massifHeight}/{uuid}/{index}.log
	// or v2/merklelog/checkpoints/{massifHeight}/{uuid}/{index}.sth
	if len(parts) >= 5 && parts[0] == "v2" && parts[1] == "merklelog" {
		// massifHeight is at index 3 (after v2/merklelog/massifs or v2/merklelog/checkpoints)
		height, err := strconv.ParseUint(parts[3], 10, 8)
		if err != nil {
			return 0, fmt.Errorf("%w: invalid massifHeight in path %s: %v", ErrMassifPathFmt, path, err)
		}
		return uint8(height), nil
	}

	// v1 format: v1/mmrs/tenant/{uuid}/0/massifs/{index}.log
	// For v1, we return 0 (LogInstanceN) as the height is not encoded
	// In practice, v1 paths with "0" typically mean height 14, but we return 0 here
	// and let callers handle the mapping if needed
	if strings.HasPrefix(path, V1MMRTenantPrefix) {
		return 0, nil
	}

	return 0, fmt.Errorf("%w: unrecognized path format: %s", ErrMassifPathFmt, path)
}

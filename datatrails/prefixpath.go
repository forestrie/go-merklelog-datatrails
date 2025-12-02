package datatrails

import (
	"fmt"

	"github.com/forestrie/go-merklelog/massifs/storage"
	"github.com/google/uuid"
)

func StoragePrefixPath(logID storage.LogID) string {
	// This is the prefix path for the blobs in the datatrails schema
	// It is used to derive the massif and checkpoint paths
	// Returns base format without service-specific prefix: {massifHeight}/{uuid}/
	// Note: This function is deprecated in favor of StorageObjectPrefixWithHeight
	return fmt.Sprintf("%s/%s/", V1MMRPrefix, Log2TenantID(logID))
}

// StorageObjectPrefix returns the old v1 path format for backward compatibility.
// Deprecated: Use StorageObjectPrefixWithHeight instead.
func StorageObjectPrefix(logID storage.LogID, otype storage.ObjectType) (string, error) {
	switch otype {
	case storage.ObjectMassifStart, storage.ObjectMassifData, storage.ObjectPathMassifs:
		return fmt.Sprintf("%s/%s/%d/massifs/", V1MMRPrefix, Log2TenantID(logID), storage.LogInstanceN), nil
	case storage.ObjectCheckpoint, storage.ObjectPathCheckpoints:
		return fmt.Sprintf("%s/%s/%d/massifseals/", V1MMRPrefix, Log2TenantID(logID), storage.LogInstanceN), nil
	default:
		return "", fmt.Errorf("unknown object type %v", otype)
	}
}

// StorageObjectPrefixWithHeight returns the base path format (without service-specific prefix).
// Returns: {massifHeight}/{uuid}/ for both massifs and checkpoints.
// Service implementations (Arbor/Canopy) should add their own prefixes:
// - v2/merklelog/massifs/ for massifs
// - v2/merklelog/checkpoints/ for checkpoints
func StorageObjectPrefixWithHeight(logID storage.LogID, massifHeight uint8, otype storage.ObjectType) (string, error) {
	// Convert LogID to UUID string (without "tenant/" prefix for base format)
	uuidStr := uuid.UUID(logID).String()

	switch otype {
	case storage.ObjectMassifStart, storage.ObjectMassifData, storage.ObjectPathMassifs:
		// Base format: {massifHeight}/{uuid}/
		return fmt.Sprintf("%d/%s/", massifHeight, uuidStr), nil
	case storage.ObjectCheckpoint, storage.ObjectPathCheckpoints:
		// Base format: {massifHeight}/{uuid}/ (same for checkpoints)
		return fmt.Sprintf("%d/%s/", massifHeight, uuidStr), nil
	default:
		return "", fmt.Errorf("unknown object type %v", otype)
	}
}

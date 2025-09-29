# go-merklelog-datatrails

This module provides minimal datatrails-specific support migrated from `go-datatrails-logverification`.

## Packages

### appentry

Contains application entry utilities for merkle log verification, migrated from `go-datatrails-logverification/logverification/app`.

### v3hash

Contains v3 event hashing and processing utilities, including:
- decodedv3event.go (migrated from decodedv4event.go)
- v3toeventresponse.go  
- eventhasher.go
- jsonprincipals.go
- merklelogentry.go

## Development

### Bootstrap

```bash
task bootstrap
```

### Testing

```bash
task test:unit
task test
```

## Dependencies

This module depends on:
- github.com/datatrails/go-datatrails-common  
- github.com/datatrails/go-datatrails-simplehash
- github.com/datatrails/go-datatrails-serialization/eventsv1
- github.com/datatrails/go-datatrails-merklelog/massifs
- github.com/datatrails/go-datatrails-merklelog/mmr
package eventsv1

var (
	eventsv1EventJson = []byte(`
{
	"identity": "events/01947000-3456-780f-bfa9-29881e3bac88",
	"attributes": {
		"foo": "bar"
	},
	"trails": [],
	"origin_tenant": "tenant/112758ce-a8cb-4924-8df8-fcba1e31f8b0",
	"created_by": "2ef471c2-f997-4503-94c8-60b5c929a3c3",
	"created_at": 1737045849174,
	"confirmation_status": "CONFIRMED",
	"merklelog_commit": {
		"index": "1",
		"idtimestamp": "019470003611017900"
	}
}
`)

	singleEventsv1EventJsonList = []byte(`
{
	"events":[
		{
			"identity": "events/01947000-3456-780f-bfa9-29881e3bac88",
			"attributes": {
				"foo": "bar"
			},
			"trails": [],
			"origin_tenant": "tenant/112758ce-a8cb-4924-8df8-fcba1e31f8b0",
			"created_by": "2ef471c2-f997-4503-94c8-60b5c929a3c3",
			"created_at": 1737045849174,
			"confirmation_status": "CONFIRMED",
			"event_type": "",
			"ledger_entry": {
				"index": "18446744073709551614",
				"idtimestamp": "019470003611017900",
				"log_id": "112758ce-a8cb-4924-8df8-fcba1e31f8b0"
			}
		}
	]
}`)
)

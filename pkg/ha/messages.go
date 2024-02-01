package ha

import (
	"net/url"
)

type msgType map[string]interface{}
type eventType map[string]interface{}
type dataType map[string]interface{}
type stateType map[string]interface{}
type attributesType map[string]interface{}

func entityId(d dataType) (string, bool) {
	s, ok := d["entity_id"].(string)
	return s, ok
}
func newState(d dataType) (stateType, bool) {
	s, ok := d["new_state"].(map[string]interface{})
	return s, ok
}

func state(evt stateType) (string, bool) {
	s, ok := evt["state"].(string)
	return s, ok
}

func event(evt msgType) (eventType, bool) {
	e, ok := evt["event"].(map[string]interface{})
	return e, ok
}

func data(evt eventType) (dataType, bool) {
	d, ok := evt["data"].(map[string]interface{})
	return d, ok
}

func attributes(s stateType) (attributesType, bool) {
	a, ok := s["attributes"].(map[string]interface{})
	return a, ok
}

func mediaContentId(attr attributesType) (string, bool) {
	mci, ok := attr["media_content_id"].(string)
	parse, _ := url.Parse(mci)
	id, _ := url.QueryUnescape(parse.Opaque)
	return id, ok
}

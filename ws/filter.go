package ws

// Filter for different values of the message. The values of the arrays will be OR concat and the different arrays will
// be ANDed. Ex: 2 ids are specified and one emitType then the message will be processed when one of the id matches AND
// the emitType.
type Filter struct {
	eventTypes    []EventType
	ids           []string
	resourceTypes []ResourceType
	messageTypes  []MessageType
	uniqueIds     []string
}

func (f Filter) check(message Message) bool {
	return contains(f.eventTypes, message.EventType) &&
		contains(f.ids, message.Id) &&
		contains(f.resourceTypes, message.ResourceType) &&
		contains(f.messageTypes, message.MessageType) &&
		contains(f.uniqueIds, message.UniqueId)
}

func contains[T comparable](content []T, value T) bool {
	if len(content) == 0 {
		return true
	}
	hit := false
	for _, t := range content {
		hit = hit || t == value
	}
	return hit
}

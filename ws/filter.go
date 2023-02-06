package ws

// Filter for different values of the message. The values of the arrays will be OR concat and the different arrays will
// be ANDed. Ex: 2 ids are specified and one emitType then the message will be processed when one of the id matches AND
// the emitType.
type Filter struct {
	EventTypes    []EventType
	Ids           []string
	ResourceTypes []ResourceType
	MessageTypes  []MessageType
	UniqueIds     []string
}

func (f Filter) check(message Message) bool {
	return contains(f.EventTypes, message.EventType) &&
		contains(f.Ids, message.Id) &&
		contains(f.ResourceTypes, message.ResourceType) &&
		contains(f.MessageTypes, message.MessageType) &&
		contains(f.UniqueIds, message.UniqueId)
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

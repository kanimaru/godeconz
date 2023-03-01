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

	// HasState checks if state attribute is set otherwise the message gets blocked
	HasState bool
	// HasAttr checks if attr attribute is set otherwise the message gets blocked
	HasAttr bool
	// HasConfig checks if config attribute is set otherwise the message gets blocked
	HasConfig bool
}

func (f Filter) check(message Message) bool {
	return contains(f.EventTypes, message.EventType) &&
		contains(f.Ids, message.Id) &&
		contains(f.ResourceTypes, message.ResourceType) &&
		contains(f.MessageTypes, message.MessageType) &&
		contains(f.UniqueIds, message.UniqueId) &&
		((f.HasState && message.State != nil) ||
			(f.HasAttr && message.Attr != nil) ||
			(f.HasConfig && message.Config != nil) ||
			(!f.HasState && !f.HasConfig && !f.HasAttr))
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

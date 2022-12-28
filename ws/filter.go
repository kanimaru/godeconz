package ws

// Filter for different values of the message. The values of the arrays will be OR concat and the different arrays will
// be ANDed. Ex: 2 ids are specified and one emitType then the message will be processed when one of the id matches AND
// the emitType.
type Filter struct {
	emitTypes []string
	ids       []string
	roles     []string
	types     []string
	uniqueIds []string
}

func (f Filter) check(message Message) bool {
	return contains(f.emitTypes, message.Emit) &&
		contains(f.ids, message.Id) &&
		contains(f.roles, message.Role) &&
		contains(f.types, message.Type) &&
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

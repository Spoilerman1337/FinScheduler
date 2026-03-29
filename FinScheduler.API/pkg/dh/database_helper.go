package dh

func Reconcile[T comparable](incoming []T, current []T) (toDelete []T, toInsert []T) {
	incomingMap := make(map[T]struct{}, len(incoming))
	for _, id := range incoming {
		incomingMap[id] = struct{}{}
	}

	currentMap := make(map[T]struct{}, len(current))
	for _, id := range current {
		currentMap[id] = struct{}{}
	}

	for _, id := range current {
		if _, exists := incomingMap[id]; !exists {
			toDelete = append(toDelete, id)
		}
	}

	for _, id := range incoming {
		if _, exists := currentMap[id]; !exists {
			toInsert = append(toInsert, id)
		}
	}

	return
}

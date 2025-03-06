package pipe

type Pipe[ItemType any] struct {
	source <-chan ItemType
}

func From[ItemType any](source <-chan ItemType) Pipe[ItemType] {
	return Pipe[ItemType]{source: source}
}

func (pipe Pipe[ItemType]) To(destinations ...chan<- ItemType) {
	go func() {
		defer func() {
			for _, destination := range destinations {
				close(destination)
			}
		}()

		for item := range pipe.source {
			for _, destination := range destinations {
				destination <- item
			}
		}
	}()
}

func (pipe Pipe[ItemType]) Receive() <-chan ItemType {
	return pipe.source
}

// Filter removes items that do not meet the given condition.
func (pipe Pipe[ItemType]) Filter(predicate func(ItemType) bool) Pipe[ItemType] {
	filtered := make(chan ItemType, 1)

	go func() {
		defer close(filtered)

		for item := range pipe.source {
			if predicate(item) {
				filtered <- item
			}
		}
	}()

	return From(filtered)
}

// Transform applies a function to modify items in the pipe.
func (pipe Pipe[ItemType]) Transform(predicate func(ItemType) ItemType) Pipe[ItemType] {
	transformed := make(chan ItemType, 1)

	go func() {
		defer close(transformed)

		for item := range pipe.source {
			transformed <- predicate(item)
		}
	}()

	return From(transformed)
}

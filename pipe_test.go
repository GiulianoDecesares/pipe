package pipe_test

import (
	"testing"

	"github.com/GiulianoDecesares/repository/repository/pipe"
	"github.com/stretchr/testify/assert"
)

func TestFrom(context *testing.T) {
	expectedItems := []int{1, 2, 3}
	source := make(chan int, 3)

	for _, item := range expectedItems {
		source <- item
	}

	close(source)

	currentPipe := pipe.From(source)
	receivedItems := collectAll(currentPipe)

	assert.Equal(context, expectedItems, receivedItems, "From() should correctly create a pipeline from a source channel.")
}

func TestToSingle(context *testing.T) {
	expectedItems := []int{10, 20, 30}
	source := make(chan int, 3)

	for _, item := range expectedItems {
		source <- item
	}

	close(source)

	currentPipe := pipe.From(source)

	destination := make(chan int, 3)
	currentPipe.To(destination)

	receivedItems := make([]int, 0)

	for item := range destination {
		receivedItems = append(receivedItems, item)
	}

	assert.Equal(context, expectedItems, receivedItems, "To() should forward data to the destination channel.")
}

func TestToBroadcast(context *testing.T) {
	expectedItems := []int{10, 20, 30}
	source := make(chan int, 3)

	for _, item := range expectedItems {
		source <- item
	}

	close(source)

	currentPipe := pipe.From(source)

	firstDestination := make(chan int, 3)
	secondDestination := make(chan int, 3)

	currentPipe.To(firstDestination, secondDestination)

	receivedItemsFirst := make([]int, 0)
	receivedItemsSecond := make([]int, 0)

	for item := range firstDestination {
		receivedItemsFirst = append(receivedItemsFirst, item)
	}

	for item := range secondDestination {
		receivedItemsSecond = append(receivedItemsSecond, item)
	}

	assert.Equal(context, expectedItems, receivedItemsFirst, "To() should forward data to the first destination channel.")
	assert.Equal(context, expectedItems, receivedItemsSecond, "To() should forward data to the second destination channel.")
}

func TestFilter(context *testing.T) {
	itemsSent := []int{15, 20, 40, 123}
	expectedItems := []int{20, 40}

	source := make(chan int, len(itemsSent))

	for _, item := range itemsSent {
		source <- item
	}

	close(source)

	// Filter even numbers
	currentPipe := pipe.From(source).Filter(func(item int) bool {
		return item%2 == 0
	})

	receivedItems := collectAll(currentPipe)

	assert.Equal(context, expectedItems, receivedItems, "Filter() should correctly filter the items.")
}

func TestTransform(context *testing.T) {
	sourceChannel := make(chan int, 3)
	sourceChannel <- 1
	sourceChannel <- 2
	sourceChannel <- 3

	close(sourceChannel)

	pipe := pipe.From(sourceChannel).Transform(func(value int) int {
		return value * 10
	})

	receivedItems := collectAll(pipe)

	expectedItems := []int{10, 20, 30}
	assert.Equal(context, expectedItems, receivedItems, "Transform() should apply the function to all items.")
}

func TestFilterAndTransform(context *testing.T) {
	sourceChannel := make(chan int, 5)
	sourceChannel <- 1
	sourceChannel <- 2
	sourceChannel <- 3
	sourceChannel <- 4
	sourceChannel <- 5
	close(sourceChannel)

	pipe := pipe.From(sourceChannel).
		Filter(func(value int) bool {
			return value%2 == 1 // Keep only odd numbers
		}).
		Transform(func(value int) int {
			return value * 100
		})

	receivedItems := collectAll(pipe)

	expectedItems := []int{100, 300, 500}
	assert.Equal(context, expectedItems, receivedItems, "Filter and Transform should work correctly together.")
}

// collectAll is a helper function to read all items from a pipe and return them as a slice.
func collectAll[ItemType any](pipe pipe.Pipe[ItemType]) []ItemType {
	var collectedItems []ItemType

	for item := range pipe.Receive() {
		collectedItems = append(collectedItems, item)
	}

	return collectedItems
}

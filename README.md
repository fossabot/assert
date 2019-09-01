# Assert library for Go.
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fyougg%2Fassert.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fyougg%2Fassert?ref=badge_shield)


## Installation

- install from github

	```bash
	go get github.com/yougg/assert
	```

## Usage

- example

	```go
	package abc

	import (
		"testing"

		"github.com/yougg/assert"
	)

	func TestSomething(t *testing.T) {
		a := assert.New(t)
		// All method's of *testing.T will be included.
		a.Log("testing started")

		// assert for nil (good for errors)
		a.Nil(object)

		// assert for not nil (good when you expect something)
		a.NotNil(object)

		// assert for bool (good when you expect type assertion or check element in map)
		a.True(ok)
		a.False(exist)
		a.Bool(expect, actual, "they should be same boolean")

		// assert equality
		a.Equal(123, 123, "they should be equal")

		// assert inequality
		a.NotEqual(123, 456, "they should not be equal")

		// assert list contains element
		a.Contains(list, element, "the list should contains element")

		// assert list not contains element
		a.NotContains(list, element, "the list should not contains element")
	}
	```

## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fyougg%2Fassert.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fyougg%2Fassert?ref=badge_large)
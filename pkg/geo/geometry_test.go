package geo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoundTripJSON(t *testing.T) {
	assert := assert.New(t)
	cases := []string{
		`{"type":"Point","coordinates":[1,2]}`,
		`{"type":"LineString","coordinates":[[1,2],[3,4]]}`,
		`{"type":"Polygon","coordinates":[[[1,2],[3,4],[5,6],[1,2]]]}`,
	}

	for i, v := range cases {
		var g Geometry

		decodeErr := g.UnmarshalJSON([]byte(v))
		assert.Nil(decodeErr, "expected no decode error for case %d", i)

		data, encodeErr := g.MarshalJSON()
		assert.Nil(encodeErr, "expected no encode error for case %d", i)

		assert.Equal(v, string(data), "expected round trip for case %d", i)
	}
}

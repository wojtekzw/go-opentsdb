package tsdb

import (
	. "launchpad.net/gocheck"
	"testing"
	"github.com/davecgh/go-spew/spew"
)

// Hook up gocheck into the gotest runner
func Test(t *testing.T) { TestingT(t) }

type tsdbSuite struct {
	conn    connection
	dpts    []dataPoint
	reqs    []request
	resps   []response
	queries []query
}

// Tie our test suite into gocheck
var _ = Suite(&tsdbSuite{})

func (s *tsdbSuite) TestDatapointsJSON(c *C) {
}

func (s *tsdbSuite) SetUpSuite(c *C) {
	s.client = *NewClient("", "4242")
}

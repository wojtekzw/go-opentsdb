package tsdb

import (
	. "launchpad.net/gocheck"
	"testing"
	"io/ioutil"
	"strings"
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
)

// Hook up gocheck into the gotest runner
func Test(t *testing.T) { TestingT(t) }

type tsdbSuite struct {
	conn      *connection
	dpts      []dataPoint
	reqs      []request
	resps     []response
	queries   []query
	reqsJSON  [][]byte    // Correct JSON requests. Populated by json files.
	respsJSON [][]byte    // Correct JSON responses. Populated by json files.
}

// Tie our test suite into gocheck
var _ = Suite(&tsdbSuite{})

func (s *tsdbSuite) TestToFromJson(c *C) {
	for i, v := range s.reqsJSON {
		fromJSON := NewEmptyRequest()
		err := json.Unmarshal(v, &fromJSON)
		c.Assert(err, IsNil)
		spew.Dump(fromJSON)
		s.reqs = append(s.reqs, *fromJSON)
		JSONFromReq, err := json.Marshal(s.reqs[i])
		c.Assert(err, IsNil)
		c.Check(JSONFromReq, DeepEquals, v)
	}
}

func (s *tsdbSuite) TestQuery(c *C) {
	for i, v := range s.reqs {
		qResp, err := s.conn.Query(v)
		c.Assert(err, IsNil)
		var JSONFromResp []byte
		err = json.Unmarshal(JSONFromResp, qResp)
		c.Assert(err, IsNil)
		c.Assert(JSONFromResp, DeepEquals, s.respsJSON[i])
		s.resps = append(s.resps, *qResp)
	}
}

func (s *tsdbSuite) SetUpSuite(c *C) {
	var err error

	// Connect to a TSDB server
	s.conn, err = NewConnection("pkc-inftsdb01.ak-networks.com", "4242")
	if err != nil { panic(err) }

	// Load from JSON files
	testFiles, err := ioutil.ReadDir("test-metrics/json")
	if err != nil { panic(err) }
	for _, v := range testFiles {
		if v.IsDir() { continue }
		if strings.Contains(v.Name(), "-request.json") {
			json, err := ioutil.ReadFile("test-metrics/json/"+v.Name())
			if err != nil { panic(err) }
			s.reqsJSON = append(s.reqsJSON, json)
		}
		if strings.Contains(v.Name(), "-response.json") {
			json, err := ioutil.ReadFile("test-metrics/json/"+v.Name())
			if err != nil { panic(err) }
			s.respsJSON = append(s.respsJSON, json)
		}
	}
}

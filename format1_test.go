package main

import (
	"net/http"
	"testing"

	log "github.com/Sirupsen/logrus"

	"github.com/stretchr/testify/suite"
)

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing context
type Format1TestSuite struct {
	suite.Suite
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (suite *Format1TestSuite) SetupSuite() {
	manager := GetQPManager()
	log.Info(" -- Format1 Test Suite --\n")
	log.Debug("Format1TestSuite.SetupSuite() - setting active parsers to format1")
	manager.SetActiveParsers("format1")
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestFormat1TestSuite(t *testing.T) {
	//ts := new(Format1TestSuite)
	log.Info("TestFormat1TestSuite - Running test suite")
	suite.Run(t, new(Format1TestSuite))
	log.Info("TestFormat1TestSuite - Finished test suite")

}

func (suite *Format1TestSuite) TestActiveParsers() {
	log.Info("Format1TestSuite.TestActiveParsers()")
	manager := GetQPManager()
	keys, active := manager.GetActiveParsers()
	log.Infof("active parsers:%s %s", keys, active)
	suite.Equal(len(active), 1, "Should be a single active parser")
}

func (suite *Format1TestSuite) TestParseQueryAndOrStatementRegExpError() {
	testString := "()"
	f := &Format1{}

	tf := f.CanParseString(testString)
	suite.Equal(false, tf, "Should not be able to parse string")

	_, err := f.ParseString(testString)

	expectedError := queryParseError{
		StatusCode: http.StatusBadRequest,
		Message:    "Invalid query format",
	}
	suite.Equal(expectedError, err, "Should be formating error")
}

func (suite *Format1TestSuite) TestParseQueryAndOrStatementFilterFormatError() {
	testString := "and(foo, bar)"
	f := &Format1{}

	tf := f.CanParseString(testString)
	suite.Equal(false, tf, "Should not be able to parse string")

	_, err := f.ParseString(testString)

	expectedError := queryParseError{
		StatusCode: http.StatusBadRequest,
		Message:    "Invalid query filter format",
	}
	suite.Equal(expectedError, err, "Should be formating error")
}

func (suite *Format1TestSuite) TestParseQueryAndOrStatementJsonError() {
	testString := "and([foo, bar,])"
	f := &Format1{}

	tf := f.CanParseString(testString)
	suite.Equal(true, tf, "Should not be able to parse string")

	_, err := f.ParseString(testString)

	expectedError := queryParseError{
		StatusCode: http.StatusBadRequest,
		Message:    "invalid character 'o' in literal false (expecting 'a')",
	}
	suite.Equal(expectedError, err, "Should be formating error")
}

func (suite *Format1TestSuite) TestParseQueryAndOrStatementNotAQuery() {
	testString := "andwell_this_is_bad"
	f := &Format1{}

	tf := f.CanParseString(testString)
	suite.Equal(false, tf, "Should not be able to parse string")

	_, err := f.ParseString(testString)

	expectedError := queryParseError{
		StatusCode: http.StatusBadRequest,
		Message:    "Invalid query format",
	}
	suite.Equal(expectedError, err, "Should be formating error")
}

func (suite *Format1TestSuite) TestParseQueryAndOrStatementSliceOfFilters() {
	testString := `and([["name", "is", "blorg"]])`
	f := &Format1{}

	tf := f.CanParseString(testString)
	suite.Equal(true, tf, "Should not be able to parse string")

	rf, err := f.ParseString(testString)

	rfExpected := newReadFilters()
	rfExpected.AddCondition(newQueryCondition("name", "is", "blorg"))

	suite.Equal(nil, err, "Error should be nil")
	suite.Equal(rfExpected, rf, "Should be the same")
}

func (suite *Format1TestSuite) TestParseQueryAndOrStatementBasicAnd() {
	testString := `and(["name", "is", "blorg"])`
	f := &Format1{}

	tf := f.CanParseString(testString)
	suite.Equal(true, tf, "Should not be able to parse string")

	rf, err := f.ParseString(testString)

	rfExpected := newReadFilters()
	rfExpected.AddCondition(newQueryCondition("name", "is", "blorg"))

	suite.Equal(nil, err, "Error should be nil")
	suite.Equal(rfExpected, rf, "Should be the same")
}

func (suite *Format1TestSuite) TestParseQueryAndOrStatementBasicAndUpper() {
	testString := `AND(["name", "is", "blorg"])`
	f := &Format1{}

	tf := f.CanParseString(testString)
	suite.Equal(true, tf, "Should be able to parse string")

	rf, err := f.ParseString(testString)

	rfExpected := newReadFilters()
	rfExpected.AddCondition(newQueryCondition("name", "is", "blorg"))

	suite.Equal(nil, err, "Error should be nil")
	suite.Equal(rfExpected, rf, "Should be the same")
}

func (suite *Format1TestSuite) TestParseQueryAndOrStatementBasicOr() {
	testString := `or(["name", "is", "blorg"])`
	f := &Format1{}

	tf := f.CanParseString(testString)
	suite.Equal(true, tf, "Should be able to parse string")

	rf, err := f.ParseString(testString)

	rfExpected := newReadFilters()
	rfExpected.LogicalOperator = "or"
	rfExpected.AddCondition(newQueryCondition("name", "is", "blorg"))

	suite.Equal(nil, err, "Error should be nil")
	suite.Equal(rfExpected, rf, "Should be the same")
}

func (suite *Format1TestSuite) TestParseQueryAndOrStatementBasicOrUpper() {
	testString := `OR(["name", "is", "blorg"])`
	f := &Format1{}

	tf := f.CanParseString(testString)
	suite.Equal(true, tf, "Should be able to parse string")

	rf, err := f.ParseString(testString)

	rfExpected := newReadFilters()
	rfExpected.LogicalOperator = "or"
	rfExpected.AddCondition(newQueryCondition("name", "is", "blorg"))

	suite.Equal(nil, err, "Error should be nil")
	suite.Equal(rfExpected, rf, "Should be the same")
}

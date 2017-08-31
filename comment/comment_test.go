package comment

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGithubPrURL(t *testing.T) {
	url := "https://github.com/getgauge/gauge/pull/2"
	assert.True(t, isGithubPrURL(url))
}

func TestGithubPrURLWithInvalidURL(t *testing.T) {
	url := "github.com/getgauge/gauge/pull/2"
	assert.False(t, isGithubPrURL(url))
}

func TestGithubPrURLWithNoPR(t *testing.T) {
	url := "https://github.com/getgauge/gauge"
	assert.False(t, isGithubPrURL(url))
}

func TestGetPrInfo(t *testing.T) {
	url := "https://github.com/getgauge/gauge/pull/2"
	expected := PrInfo{
		Owner:    "getgauge",
		Repo:     "gauge",
		PrNumber: 2,
	}
	pi, _ := getPrInfo(url)
	assert.EqualValues(t, expected, pi)
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

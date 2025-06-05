package fetcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchBattingOrder_ValidMatch(t *testing.T) {
	bo, err := FetchBattingOrder("BOS", "R")
	if err != nil {
		t.Fatalf("FetchBattingOrder failed: %v", err)
	}

	assert.Equal(t, 101, bo.PlayerID1)
	assert.Equal(t, 109, bo.PlayerID9)
}

func TestFetchBattingOrder_CaseInsensitive(t *testing.T) {
	bo, err := FetchBattingOrder("bos", "r")
	if err != nil {
		t.Fatalf("FetchBattingOrder failed: %v", err)
	}

	assert.Equal(t, 101, bo.PlayerID1)
	assert.Equal(t, 109, bo.PlayerID9)
}

func TestFetchBattingOrder_NoMatch(t *testing.T) {
	bo, err := FetchBattingOrder("SEA", "L")
	if err != nil {
		t.Fatalf("FetchBattingOrder failed: %v", err)
	}

	assert.Equal(t, 0, bo.PlayerID1)
	assert.Equal(t, 0, bo.PlayerID9)
}

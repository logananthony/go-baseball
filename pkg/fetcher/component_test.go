package fetcher

import (
	//"reflect"
	"testing"
	"github.com/davecgh/go-spew/spew"
)

func TestComponentFetcher(t *testing.T) {

  result := FetchBatterSwingPercentageLeague()

  spew.Dump(result)

}


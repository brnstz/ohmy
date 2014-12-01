package ohmy_test

import (
	"testing"

	"github.com/brnstz/ohmy"
)

func testRegion(region ohmy.Region, t *testing.T) {
	shows, err := ohmy.GetShows(region, 101)
	if err != nil {
		t.Fatal(err)
	}

	if len(shows) != 101 {
		t.Fatal("Expected 101 shows")
	}

	for _, show := range shows {
		if len(show.Venue.Name) < 1 {
			t.Fatal("Expected non-nil venue name")
		}
		if len(show.Bands) < 1 {
			t.Fatal("Expected non-nil band list")
		}
	}

}

func TestNYC(t *testing.T) {
	testRegion(ohmy.RegionNYC, t)
}

func TestLA(t *testing.T) {
	testRegion(ohmy.RegionLA, t)
}

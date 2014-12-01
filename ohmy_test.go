package ohmy_test

import (
	"testing"

	"github.com/brnstz/ohmy"
)

func testRegion(region ohmy.Region, n int, t *testing.T) {
	shows, err := ohmy.GetShows(region, n)
	if err != nil {
		t.Fatal(err)
	}

	if len(shows) != n {
		t.Fatalf("Expected %d shows", n)
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
	testRegion(ohmy.RegionNYC, 101, t)
}

func TestLA(t *testing.T) {
	testRegion(ohmy.RegionLA, 101, t)
}

func TestSmall(t *testing.T) {
	testRegion(ohmy.RegionNYC, 5, t)
}

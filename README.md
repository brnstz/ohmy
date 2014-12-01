Look at the GoDoc
=================

[![GoDoc](https://godoc.org/github.com/brnstz/ohmy?status.svg)](https://godoc.org/github.com/brnstz/ohmy)

Or here's an example
====================

```go
package main

import (
    "fmt"
    "github.com/brnstz/ohmy"
)

func main() {
    shows, err := ohmy.GetShows(100)

    if err == nil {
        fmt.Println(shows[0].Venue.Name, shows[0].Bands[0].Name)
    }
}

```


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

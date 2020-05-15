package nomad

import (
    "github.com/GPKyte/nomad/common"
)

/*
 * NewTravelMap builds a graph to use as a metaphorical map that can model travel plans
 */
func NewTravelMap(deals []Listing) *graph {
    const waitFare = 0 // Since price is a factor in trips, but not in waiting, weight is 0

    var nodes = make([]TimeAndPlace)
    var edges = make([]Trip)
    var groupByPlace = make(map[string][]TimeAndPlace) /* Unknown size, # keys is # unique locations */

    /* One deal Connects a TimeAndPlace to another via a Trip */
    for _, L := range deals {
        locD := getCode(L.Depart.Location)
        depart := TimeAndPlace{L.Depart.Time, locD}
        locA := getCode(L.Arrive.Location)
        arrive := TimeAndPlace{L.Arrive.Time, locA}

        // Build up nodes and edges
        nodes = append(nodes, depart)
        nodes = append(nodes, arrive)
        edges = append(edges, Trip{depart, arrive, L.price})

        // Build up struct for Waiting Edges to use next
        groupByPlace[locD] = append(groupByPlace[locD], depart)
        groupByPlace[locA] = append(groupByPlace[locA], arrive)
    }

    /* Each Listing resulted in two TimeAndPlace objects, each of those
     * are now grouped by Location and will be sorted by Time
     * Then edges are formed in a line from start to finish in each group. */
	for place, slice := range groupByPlace {
		/* First sort by time to properly model behavior of waiting in place */
		slice = slice.sort(true)

		for waitFrom, until := 0, 1; until < len(slice); waitFrom++ until++ {
			/* Because TimeAndPlace are connected by the ability to wait in place
			 * we make an edge between consecutive times at same place
			 * Also, note that arrival AND Departure connect just the same */
			e := Trip{slice[waitFrom], slice[until], waitFare}
			edges = append(edges, e)
		}
    }

    return NewGraph(nodes, edges)
}

func MakePlans(from, to Location, budget Money, departWhen []DateTime) {
	checkGoalHueristic := func(atLoc Location) int {
		if hasCoordinates(atLoc) && hasCoordinates(to) {
			/* Estimate distance between both */
			estimateDistance := 1
			return estimateDistance
		}
		return -1
	}

	/* Impl A* below */
	return
}

func getCode(L location) string {
    return L.code
}


func convPath(path []edge, mapping map[int]TimeAndPlace) []Trip {
    if len(path) == 0 {
        /* Replace panic with more robust action */
        panic("No path found/provided")
    }

    for _, trip := range path {
        if trip.price == 0 {
            continue // This is time spent waiting until the next departure
        }
        var start TimeAndPlace = mapping[trip.from]
        var end TimeAndPlace = mapping[trip.to]
        
        // append(Trip{start, end, trip.price})
    }
}

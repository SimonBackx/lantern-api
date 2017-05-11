package main

import (
	"encoding/json"
	"fmt"
	"github.com/SimonBackx/lantern-crawler/queries"
	"io/ioutil"
	"net/http"
)

// Elke 10 minuten
var statisticsTimeline = make([]*queries.Stats, 0)

// Elke minuut
var statisticsQueue = make([]*queries.Stats, 0)

/**
 * POST /stats
 */
func newStatsHandler(w http.ResponseWriter, r *http.Request) {
	str, err := ioutil.ReadAll(r.Body)

	if err != nil {
		internalErrorHandler(w, r, err)
		return
	}

	var result queries.Stats
	err = json.Unmarshal(str, &result)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid result.")
		return
	}

	statisticsQueue = append(statisticsQueue, &result)
	if len(statisticsQueue) >= 10 {
		// Combineren
		avg := queries.AverageStats(statisticsQueue)
		statisticsQueue = statisticsQueue[:0]
		statisticsTimeline = append(statisticsTimeline, avg)

		if len(statisticsTimeline) > 24*6 {
			offset := len(statisticsTimeline) - 24*6
			statisticsTimeline = statisticsTimeline[offset:]
		}

		// todo: wissen na 24 uur
	}
	fmt.Fprintf(w, "Success")
}

/**
 * GET /stats
 */
func statsHandler(w http.ResponseWriter, r *http.Request) {
	merged := statisticsTimeline
	if len(statisticsQueue) > 0 {
		merged = append(merged, statisticsQueue[len(statisticsQueue)-1])
	}

	str, err := json.Marshal(merged)

	if err != nil {
		internalErrorHandler(w, r, err)
		return
	}

	fmt.Fprintf(w, "%s", str)
}

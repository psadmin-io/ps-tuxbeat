package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"

	"github.com/psadmin-io/ps-tuxbeat/beater"
)

func main() {
	err := beat.Run("tuxbeat", "", beater.New)
	if err != nil {
		os.Exit(1)
	}
}

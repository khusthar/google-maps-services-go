// Copyright 2015 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package main contains a simple command line tool for Timezone API
// Directions docs: https://developers.google.com/maps/documentation/timezone/
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/kr/pretty"
	"google.golang.org/maps"
)

var (
	apiKey      = flag.String("key", "", "API Key for using Google Maps API.")
	path        = flag.String("path", "", "The path to be snapped. The path parameter accepts a list of latitude/longitude pairs. Latitude and longitude values should be separated by commas. Coordinates should be separated by the pipe character.")
	interpolate = flag.Bool("interpolate", false, "Whether to interpolate a path to include all points forming the full road-geometry.")
)

func usageAndExit(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	fmt.Println("Flags:")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Parse()
	client := &http.Client{}
	if *apiKey == "" {
		usageAndExit("Please specify an API Key.")
	}
	ctx := maps.NewContext(*apiKey, client)
	r := &maps.SnapToRoadRequest{
		Interpolate: *interpolate,
	}

	parsePath(*path, r)

	pretty.Println(r)

	resp, err := r.Get(ctx)
	if err != nil {
		log.Fatalf("error %v", err)
	}

	pretty.Println(resp)
}

// parsePath takes a location argument string and decodes it.
func parsePath(path string, r *maps.SnapToRoadRequest) {
	if path != "" {
		ls := strings.Split(path, "|")
		for _, l := range ls {
			ll := strings.Split(l, ",")
			lat, err := strconv.ParseFloat(ll[0], 64)
			if err != nil {
				usageAndExit(fmt.Sprintf("Could not parse path: %v", err))
			}
			lng, err := strconv.ParseFloat(ll[1], 64)
			if err != nil {
				usageAndExit(fmt.Sprintf("Could not parse path: %v", err))
			}
			r.Path = append(r.Path, maps.Location{Latitude: lat, Longitude: lng})
		}
	} else {
		usageAndExit("Path required")
	}
}

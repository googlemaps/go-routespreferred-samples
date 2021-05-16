// Copyright 2021 Google LLC
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

// [START maps_routespreferred_samples_custom_routes_default]
package main

import (
	"context"
	"crypto/tls"
	"log"
	"os"
	"time"

	"github.com/golang/protobuf/proto"
	v1 "google.golang.org/genproto/googleapis/maps/routes/v1"
	v1alpha "google.golang.org/genproto/googleapis/maps/routes/v1alpha"
	"google.golang.org/genproto/googleapis/type/latlng"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

const (
	serverAddr = "routespreferred.googleapis.com:443"
	// Note that setting the field mask to * is OK for testing, but discouraged in
	// production.
	// For example, for ComputeRoutes, set the field mask to
	// "routes.distanceMeters,routes.duration,routes.polyline.encodedPolyline"
	// in order to get the route distances, durations, and encoded polylines.
	fieldMask = "*"
)

func createWaypoint(lat float64, lng float64) *v1.Waypoint {
	return &v1.Waypoint{LocationType: &v1.Waypoint_Location{
		Location: &v1.Location{
			LatLng: &latlng.LatLng{Latitude: lat, Longitude: lng},
		},
	}}
}

func callComputeCustomRoutes(client v1alpha.RoutesAlphaClient, ctx *context.Context) {
	request := v1.ComputeCustomRoutesRequest{
		Origin:      createWaypoint(37.420761, -122.081356),
		Destination: createWaypoint(37.420999, -122.086894),
		RouteObjective: &v1.RouteObjective{
			Objective: &v1.RouteObjective_RateCard_{
				RateCard: &v1.RouteObjective_RateCard{
					CostPerMinute: &v1.RouteObjective_RateCard_MonetaryCost{Value: 1.1},
				},
			},
		},
		TravelMode:        v1.RouteTravelMode_DRIVE,
		RoutingPreference: v1.RoutingPreference_TRAFFIC_AWARE,
		Units:             v1.Units_METRIC,
		LanguageCode:      "en-us",
		RouteModifiers: &v1.RouteModifiers{
			AvoidTolls:    false,
			AvoidHighways: true,
			AvoidFerries:  true,
		},
		PolylineQuality: v1.PolylineQuality_OVERVIEW,
	}
	marshaler := proto.TextMarshaler{}
	log.Printf("Sending request: \n%s", marshaler.Text(&request))
	result, err := client.ComputeCustomRoutes(*ctx, &request)

	if err != nil {
		log.Fatalf("Failed to call ComputeCustomRoutes: %v", err)
	}
	log.Printf("Result: %s", marshaler.Text(result))
}

func main() {
	config := tls.Config{}
	conn, err := grpc.Dial(serverAddr,
		grpc.WithTransportCredentials(credentials.NewTLS(&config)))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()
	client := v1alpha.NewRoutesAlphaClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	ctx = metadata.AppendToOutgoingContext(ctx, "X-Goog-Api-Key", os.Getenv("GOOGLE_MAPS_API_KEY"))
	ctx = metadata.AppendToOutgoingContext(ctx, "X-Goog-Fieldmask", fieldMask)
	defer cancel()

	callComputeCustomRoutes(client, &ctx)
}

// [END maps_routespreferred_samples_custom_routes_default]

/*
***

# package for business logic

***
*/

package server

import (
	"TestP/pkg/utilities"
	"github.com/gin-gonic/gin"
	"gopkg.in/square/go-jose.v2/json"
	"log"
	"math"
	"net/http"
)

type Evn struct {
	Port string
}

// to parse lines
var lines []struct {
	Line struct {
		Type        string      `json:"type"`
		Coordinates [][]float64 `json:"coordinates"`
	} `json:"line"`
}

type IntersectionResp struct {
	LineId int
	Points Point
}

// Earth's radius in kilometers
const earthRadius = 6371.0

type Point struct {
	Latitude  float64
	Longitude float64
}

// Haversine formula to calculate the distance between two points on the Earth's surface
func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	lat1Rad := degToRad(lat1)
	lon1Rad := degToRad(lon1)
	lat2Rad := degToRad(lat2)
	lon2Rad := degToRad(lon2)

	dLat := lat2Rad - lat1Rad
	dLon := lon2Rad - lon1Rad

	a := math.Pow(math.Sin(dLat/2), 2) + math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Pow(math.Sin(dLon/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := earthRadius * c
	return distance
}

// conversion from degree to radian
func degToRad(deg float64) float64 {
	return deg * (math.Pi / 180)
}

// Find the intersection point of two geolocations
func findIntersection(point1, point2 Point) (Point, bool) {
	lat1 := degToRad(point1.Latitude)
	lon1 := degToRad(point1.Longitude)
	lat2 := degToRad(point2.Latitude)
	lon2 := degToRad(point2.Longitude)

	x := (lon2 - lon1) * math.Cos((lat1+lat2)/2)
	y := lat2 - lat1
	distance := math.Sqrt(math.Pow(x, 2)+math.Pow(y, 2)) * earthRadius

	// Calculate the intersection point's latitude and longitude
	intersectionLat := point1.Latitude + y*(point2.Latitude-point1.Latitude)/distance
	intersectionLon := point1.Longitude + x*(point2.Longitude-point1.Longitude)/distance

	inte := doGeolocationsIntersect(point1, point2)
	return Point{Latitude: intersectionLat, Longitude: intersectionLon}, inte
}

// Check if two points intersect
func doGeolocationsIntersect(point1, point2 Point) bool {
	// Calculate the distance between the two points
	distance := haversine(point1.Latitude, point1.Longitude, point2.Latitude, point2.Longitude)

	// Define a threshold distance (e.g., 1 kilometer) to consider intersection
	thresholdDistance := 1.0

	return distance <= thresholdDistance
}

func (receiver Evn) FindIntersection(c *gin.Context) {
	// Parse the GeoJSON linestring from the request body
	var linestring struct {
		Type         string      `json:"type"`
		Cooardinates [][]float64 `json:"Coordinates"`
	}
	log.Println("called FindIntersection")
	err := c.BindJSON(&linestring)
	if err != nil {
		log.Println("error bind body ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incomplete Request"})
		return
	}
	log.Println("Input Received")
	//reading lines data from file
	data, err := utilities.ReadDataFromFile("data", "lines.json")
	if err != nil {
		log.Println("Read data error ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Some error occurred"})

		return
	}

	err = json.Unmarshal(data, &lines)
	if err != nil {
		log.Println("unmarshal error ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Some error occurred"})

		return

	}

	var intersections []IntersectionResp
	log.Println("Calculating Intersection")
	for id, line := range lines {

		for _, line2 := range linestring.Cooardinates {
			L1 := line.Line.Coordinates

			var l2, l1 Point

			l1.Latitude = L1[0][0]
			l1.Longitude = L1[0][1]

			l2.Latitude = line2[0]
			l2.Longitude = line2[1]

			// Perform intersection check
			inter, isIntersect := findIntersection(l1, l2)

			if isIntersect && !math.IsNaN(inter.Latitude) {

				intersections = append(intersections, IntersectionResp{
					LineId: id,
					Points: inter,
				})
			}
		}

	}
	//sending final output
	log.Println("Returning Response")
	c.JSON(http.StatusOK, gin.H{"data": intersections})
	return
}

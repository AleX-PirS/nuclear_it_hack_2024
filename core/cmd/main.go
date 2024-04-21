package main

import (
	"log"
	"math"
	"os"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geo"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb/quadtree"

	"golang.org/x/exp/maps"
)

type PointData struct{
	geoIndex int
	atrIndex int
	findFlag bool
}

func betweenPoint(n int, l, r orb.Point) []orb.Point{
	res := make([]orb.Point,0)

	for i:=0;i<n;i++{
		lon := math.Min(l.Lon(), r.Lon()) + float64(i)/float64(n-1)*(math.Abs(l.Lon()-r.Lon()))
		lat := math.Min(l.Lat(), r.Lat()) + float64(i)/float64(n-1)*(math.Abs(l.Lat()-r.Lat()))
		res = append(res, orb.Point{lon, lat})
	}

	return res
}

func processRawData(atrGraph *geojson.FeatureCollection) (*quadtree.Quadtree, map[orb.Point]*PointData){
	qTree := quadtree.New(orb.Bound{Min: orb.Point{-180, -90}, Max: orb.Point{180, 90}})
	hash := make(map[orb.Point]*PointData)

	for idx, ft := range atrGraph.Features{
		for _, lineStr := range ft.Geometry.(orb.MultiLineString){
			for pointIdx, p := range lineStr{
				qTree.Add(p)
				hash[p] = &PointData{atrIndex: idx}
				if pointIdx != len(lineStr)-1{
					bps := betweenPoint(5, lineStr[pointIdx], lineStr[pointIdx+1])
					for i, bp := range bps{
						log.Println("1, 2, bp, i:", lineStr[pointIdx], lineStr[pointIdx+1], bp, i)
						qTree.Add(bp)
						hash[bp] = &PointData{atrIndex: idx}
					}
				}
			}
		}
	}

	return qTree, hash
}

func main() {
	// redGEOJson, err := os.ReadFile("D:/dev/github.com/AleX-PirS/nuclear_it_hack_2024/data/kaliningrad_red_WGS84.geojson")
	redGEOJson, err := os.ReadFile("D:/dev/github.com/AleX-PirS/nuclear_it_hack_2024/data/red_new1.geojson")
	if err != nil{
		log.Fatal(err)
	}

	// blueGEOJson, err := os.ReadFile("D:/dev/github.com/AleX-PirS/nuclear_it_hack_2024/data/kaliningrad_blue_WGS84.geojson")
	blueGEOJson, err := os.ReadFile("D:/dev/github.com/AleX-PirS/nuclear_it_hack_2024/data/blue_new.geojson")
	if err != nil{
		log.Fatal(err)
	}

	atrData, err := geojson.UnmarshalFeatureCollection(redGEOJson)

	if err != nil{
		log.Fatal("Error unmarshall RED geojson:, ", err)
	}

	geoData, err := geojson.UnmarshalFeatureCollection(blueGEOJson)

	if err != nil{
		log.Fatal("Error unmarshall BLUE geojson:, ", err)
	}

	qTree, geoHash := processRawData(atrData)

	resGraph := ResolveGraph(100, geoData, atrData, qTree, geoHash)
	log.Println(resGraph)
	data, err := resGraph.MarshalJSON()
	if err != nil {
		log.Fatal("Error marshall json")
	}

	err = os.WriteFile("green.geojson", data, 0644)
	if err != nil {
		log.Fatal("Error save json")
	}

	maps.Copy(geoData.Features[0].Properties, atrData.Features[0].Properties)
}

func ResolveGraph(accarucy int, geoGraph, atrGraph *geojson.FeatureCollection, qTree *quadtree.Quadtree, geoHash map[orb.Point]*PointData) *geojson.FeatureCollection{
	resGraph := geojson.NewFeatureCollection()

	curr := 0
	for ft_idx, ft := range geoGraph.Features{
		newFtFunc := func (id int) *geojson.Feature{
			geometry := ft.Geometry
			properties := ft.Properties
			ftType := ft.Type

			return &geojson.Feature{ID: id, Type: ftType, Geometry: geometry, Properties: properties}
		}
		

		newFt := newFtFunc(curr)
		multyLine := orb.MultiLineString{}
		atrIndexBefore := -1
		atrIndex := -1

		for _, lineStr := range ft.Geometry.(orb.MultiLineString){

			line := orb.LineString{}
			for pointIdx, geoPoint := range lineStr{
				atrIndex = CalculateGraph(accarucy, geoPoint, ft_idx, geoHash, qTree)
				if atrIndex != -1{
					line = append(line, geoPoint)
				}
				
				if pointIdx == len(lineStr)-1{
					multyLine = append(multyLine, line)
					newFt.Geometry = multyLine
					if atrIndex != -1{
						maps.Copy(newFt.Properties, atrGraph.Features[atrIndex].Properties)
					} else if atrIndexBefore != -1{
						maps.Copy(newFt.Properties, atrGraph.Features[atrIndexBefore].Properties)
					}
					resGraph.Features = append(resGraph.Features, newFt)
					curr += 1
					break
				}

				if atrIndex != atrIndexBefore && atrIndexBefore!=-1{
					multyLine = append(multyLine, line)
					newFt.Geometry = multyLine
					if atrIndex != -1{
						maps.Copy(newFt.Properties, atrGraph.Features[atrIndex].Properties)
					} else if atrIndexBefore != -1{
						maps.Copy(newFt.Properties, atrGraph.Features[atrIndexBefore].Properties)
					}
					resGraph.Features = append(resGraph.Features, newFt)

					multyLine = orb.MultiLineString{}
					line = orb.LineString{}
					if pointIdx != len(lineStr)-1{
						line = append(line, geoPoint)
					}
					curr += 1
					newFt = newFtFunc(curr)
				}
				
				atrIndexBefore = atrIndex
			}
			// if len(line) > 0 && atrIndex != -1{
			// 	multyLine = append(multyLine, line)
			// 	newFt.Geometry = multyLine
			// 	maps.Copy(newFt.Properties, atrGraph.Features[atrIndex].Properties)
			// 	resGraph.Features = append(resGraph.Features, newFt)
			// }
		}
	}

	return resGraph
}

func CalculateGraph(accarucy int, point orb.Point, geoFtIdx int, geoHash map[orb.Point]*PointData, qTree *quadtree.Quadtree) int{
	nearestPoint := qTree.Find(point).Point()

	if mh, ok := geoHash[nearestPoint]; ok{
		r := geo.Distance(point, nearestPoint)
		mh.geoIndex = geoFtIdx
		if r>float64(accarucy){
			return -1
		}
		log.Printf("Log. R=%v, nearest:%v", r, nearestPoint)
		mh.findFlag = true

		return mh.atrIndex
	}

	return -1
}
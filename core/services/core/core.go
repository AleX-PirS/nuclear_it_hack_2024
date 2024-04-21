package core

import (
	"log"
	"math"
	"os"
	"time"

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

func betweenPoint(use bool, accarucy int, l, r orb.Point) []orb.Point{
	if !use{
		return []orb.Point{} 
	}

	res := make([]orb.Point,0)

	n := 0
	switch dist := geo.Distance(l, r);{
		case dist < float64(accarucy)/4:
			return []orb.Point{} 
		case dist < float64(accarucy) * 2:
			n = 20
		case dist < float64(accarucy) * 20:
			n = 100
		case dist < float64(accarucy) * 80:
			n = 400
		default:
			n = 1000
	}

	n = int(geo.Distance(l, r)/float64(accarucy))*10

	for i:=0;i<n+2;i++{
		lon := math.Min(l.Lon(), r.Lon()) + float64(i)/float64(n-1)*(math.Abs(l.Lon()-r.Lon()))
		lat := math.Min(l.Lat(), r.Lat()) + float64(i)/float64(n-1)*(math.Abs(l.Lat()-r.Lat()))
		res = append(res, orb.Point{lon, lat})
	}

	return res[1:len(res)-1]
}

func processRawData(accarucy int, atrGraph *geojson.FeatureCollection) (*quadtree.Quadtree, map[orb.Point]*PointData){
	qTree := quadtree.New(orb.Bound{Min: orb.Point{-180, -90}, Max: orb.Point{180, 90}})
	hash := make(map[orb.Point]*PointData)

	for idx, ft := range atrGraph.Features{
		for _, lineStr := range ft.Geometry.(orb.MultiLineString){
			for pointIdx, p := range lineStr{
				qTree.Add(p)
				hash[p] = &PointData{atrIndex: idx}
				if pointIdx != len(lineStr)-1{
					bps := betweenPoint(true, accarucy, lineStr[pointIdx], lineStr[pointIdx+1])
					for _, bp := range bps{
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
	start := time.Now()
	fileName := "green_v1.geojson"
	redGEOJson, err := os.ReadFile("D:/dev/github.com/AleX-PirS/nuclear_it_hack_2024/data/kaliningrad_red_WGS84.geojson")
	// redGEOJson, err := os.ReadFile("D:/dev/github.com/AleX-PirS/nuclear_it_hack_2024/data/red_new1.geojson")
	// redGEOJson, err := os.ReadFile("D:/dev/github.com/AleX-PirS/nuclear_it_hack_2024/data/red_new_new.geojson")
	// redGEOJson, err := os.ReadFile("D:/dev/github.com/AleX-PirS/nuclear_it_hack_2024/data/red_new_new_new.geojson")
	// redGEOJson, err := os.ReadFile("D:/dev/github.com/AleX-PirS/nuclear_it_hack_2024/data/red_1.geojson")
	if err != nil{
		log.Fatal(err)
	}

	blueGEOJson, err := os.ReadFile("D:/dev/github.com/AleX-PirS/nuclear_it_hack_2024/data/kaliningrad_blue_WGS84.geojson")
	// blueGEOJson, err := os.ReadFile("D:/dev/github.com/AleX-PirS/nuclear_it_hack_2024/data/blue_new.geojson")
	// blueGEOJson, err := os.ReadFile("D:/dev/github.com/AleX-PirS/nuclear_it_hack_2024/data/blue_new_new.geojson")
	// blueGEOJson, err := os.ReadFile("D:/dev/github.com/AleX-PirS/nuclear_it_hack_2024/data/blue_new_new_new.geojson")
	// blueGEOJson, err := os.ReadFile("D:/dev/github.com/AleX-PirS/nuclear_it_hack_2024/data/blue_1.geojson")
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
	
	accarucy := 20

	qTree, atrHash := processRawData(accarucy, atrData)

	resGraph := ResolveGraph(accarucy, geoData, atrData, qTree, atrHash)
	// log.Println(resGraph)
	data, err := resGraph.MarshalJSON()
	if err != nil {
		log.Fatal("Error marshall json")
	}

	err = os.WriteFile(fileName, data, 0644)
	if err != nil {
		log.Fatal("Error save json")
	}
	log.Println("Time:", time.Since(start))
}

func ResolveGraph(accarucy int, geoGraph, atrGraph *geojson.FeatureCollection, qTree *quadtree.Quadtree, geoHash map[orb.Point]*PointData) *geojson.FeatureCollection{
	resGraph := geojson.NewFeatureCollection()

	curr := 0
	for ft_idx, ft := range geoGraph.Features{
		newFtFunc := func (id int) *geojson.Feature{
			geometry := orb.MultiLineString{}
			properties := ft.Properties
			ftType := ft.Type

			return &geojson.Feature{ID: id, Type: ftType, Geometry: geometry, Properties: properties}
		}
		

		newFt := newFtFunc(curr)
		multyLine := orb.MultiLineString{}
		atrIndexBefore := -1
		atrIndex := -1


		line := orb.LineString{}
		for pointIdx, geoPoint := range ft.Geometry.(orb.MultiLineString)[0]{
			atrIndex = CalculateGraph(accarucy, geoPoint, ft_idx, geoHash, qTree)
			// log.Println("LOGGING:", atrIndex, atrIndexBefore, geoPoint, pointIdx, len(ft.Geometry.(orb.MultiLineString)[0]))
			// log.Println("Len line start:", len(line))
			if atrIndex != -1{
				maps.Copy(newFt.Properties, atrGraph.Features[atrIndex].Properties)
				line = append(line, geoPoint)
			}
			if pointIdx == len(ft.Geometry.(orb.MultiLineString)[0])-1{
				if len(line) > 1{
					multyLine = append(multyLine, line)
				}
				if len(multyLine) > 0 {
					newFt.Geometry = multyLine
				}
				// if atrIndex != -1{
				// 	maps.Copy(newFt.Properties, atrGraph.Features[atrIndex].Properties)
				// } else if atrIndexBefore != -1{
				// 	maps.Copy(newFt.Properties, atrGraph.Features[atrIndexBefore].Properties)
				// }
				if len(newFt.Geometry.(orb.MultiLineString)) > 0{
					resGraph.Features = append(resGraph.Features, newFt)
					curr += 1
				}
				break
			}
			// log.Println("Len line after stop_check:", len(line))

			if atrIndex != atrIndexBefore && atrIndexBefore!=-1{
				if len(line) > 1{
					multyLine = append(multyLine, line)
				}
				if len(multyLine) > 0 {
					newFt.Geometry = multyLine
				}
				// if atrIndex != -1{
				// 	maps.Copy(newFt.Properties, atrGraph.Features[atrIndex].Properties)
				// } else if atrIndexBefore != -1{
				// 	maps.Copy(newFt.Properties, atrGraph.Features[atrIndexBefore].Properties)
				// }
				if len(newFt.Geometry.(orb.MultiLineString)) > 0{
					resGraph.Features = append(resGraph.Features, newFt)
				}
				multyLine = orb.MultiLineString{}
				
				if len(line) < 2{
					line = orb.LineString{}
				} else {
					line = orb.LineString{}
					line = append(line, geoPoint)
				}
			
				curr += 1
				newFt = newFtFunc(curr)
			}
			// log.Println("Len line end:", len(line))
			
			atrIndexBefore = atrIndex
		}
	}

	return resGraph
}

func CalculateGraph(accarucy int, point orb.Point, geoFtIdx int, atrHash map[orb.Point]*PointData, qTree *quadtree.Quadtree) int{
	nearestPoint := qTree.Find(point).Point()

	if mh, ok := atrHash[nearestPoint]; ok{
		r := geo.Distance(point, nearestPoint)
		mh.geoIndex = geoFtIdx
		if r>float64(accarucy){
			return -1
		}
		// log.Printf("Log. R=%v, nearest:%v", r, nearestPoint)
		mh.findFlag = true

		return mh.atrIndex
	}

	return -1
}
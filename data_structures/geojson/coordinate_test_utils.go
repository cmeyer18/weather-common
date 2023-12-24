package geojson

func makeCoordinate(lat, lng float64) interface{} {
	return []interface{}{lat, lng}
}

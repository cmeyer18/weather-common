package main

import (
	"github.com/cmeyer18/weather-common/v2/generative/generators/golang"
	"github.com/cmeyer18/weather-common/v2/generative/generators/swift"
	"os"
)

func main() {
	print(os.Getwd())
	golang.GenerateSpcOutlookGo()
	golang.GenerateWeatherAlertsGo()
	swift.GenerateSpcOutlookSwift()
	swift.GenerateWeatherAlertsSwift()
}

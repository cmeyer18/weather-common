package main

import (
	"github.com/cmeyer18/weather-common/v3/generative/generators/golang"
	"github.com/cmeyer18/weather-common/v3/generative/generators/swift"
	"os"
)

func main() {
	print(os.Getwd())
	golang.GenerateConvectiveOutlookGo()
	golang.GenerateWeatherAlertsGo()
	swift.GenerateConvectiveOutlookSwift()
	swift.GenerateWeatherAlertsSwift()
}

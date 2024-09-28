package main

import (
	"os"

	"github.com/cmeyer18/weather-common/v6/generative/generators/golang"
	"github.com/cmeyer18/weather-common/v6/generative/generators/swift"
)

func main() {
	print(os.Getwd())
	golang.GenerateConvectiveOutlookGo()
	golang.GenerateWeatherAlertsGo()
	swift.GenerateConvectiveOutlookSwift()
	swift.GenerateWeatherAlertsSwift()
}

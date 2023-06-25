// Generated on 2023-06-24 20:34:38.374247 -0500 CDT m=+0.274684751
// This file was generated by generative/generators/swift/spc_outlooks_enum.go.
// Please do not hand write.

import Foundation

enum SPCOutlookTypes: String, CaseIterable {
	case Day1Categorical = "Day 1 Categorical"
	case Day1Tornado = "Day 1 Tornado"
	case Day1Wind = "Day 1 Wind"
	case Day1Hail = "Day 1 Hail"
	case Day1SignificantTornado = "Day 1 Significant Tornado"
	case Day1SignificantWind = "Day 1 Significant Wind"
	case Day1SignificantHail = "Day 1 Significant Hail"
	case Day2Categorical = "Day 2 Categorical"
	case Day2Tornado = "Day 2 Tornado"
	case Day2Wind = "Day 2 Wind"
	case Day2Hail = "Day 2 Hail"
	case Day2SignificantTornado = "Day 2 Significant Tornado"
	case Day2SignificantWind = "Day 2 Significant Wind"
	case Day2SignificantHail = "Day 2 Significant Hail"
	case Day3Categorical = "Day 3 Categorical"
	case Day3Probabilistic = "Day 3 Probabilistic"
	case Day3SignificantProbabilistic = "Day 3 Significant Probabilistic"
	case Day4Probabilistic = "Day 4 Probabilistic"
	case Day5Probabilistic = "Day 5 Probabilistic"
	case Day6Probabilistic = "Day 6 Probabilistic"
	case Day7Probabilistic = "Day 7 Probabilistic"
	case Day8Probabilistic = "Day 8 Probabilistic"
}

func SPCOutlookTypesOptions() -> Array<Option> {
    var convertAlerts: Array<Option> = []
    let allAlertsEnums: Array<SPCOutlookTypes> = SPCOutlookTypes.allCases
    
    for alert in allAlertsEnums {
        convertAlerts.append(Option(name: alert.rawValue))
    }
    
    return convertAlerts
}

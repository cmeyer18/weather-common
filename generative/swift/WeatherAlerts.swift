// Generated on 2023-06-19 22:46:56.386964 -0500 CDT m=+0.243043042
// This file was generated by generative/generators/swift/weather_alerts_enum_swift.go.
// Please do not hand write.

import Foundation

enum WeatherAlertTypes: String, CaseIterable {
	case _911TelephoneOutageEmergency = "911 Telephone Outage Emergency"
	case AdministrativeMessage = "Administrative Message"
	case AirQualityAlert = "Air Quality Alert"
	case AirStagnationAdvisory = "Air Stagnation Advisory"
	case ArroyoAndSmallStreamFloodAdvisory = "Arroyo And Small Stream Flood Advisory"
	case AshfallAdvisory = "Ashfall Advisory"
	case AshfallWarning = "Ashfall Warning"
	case AvalancheAdvisory = "Avalanche Advisory"
	case AvalancheWarning = "Avalanche Warning"
	case AvalancheWatch = "Avalanche Watch"
	case BeachHazardsStatement = "Beach Hazards Statement"
	case BlizzardWarning = "Blizzard Warning"
	case BlizzardWatch = "Blizzard Watch"
	case BlowingDustAdvisory = "Blowing Dust Advisory"
	case BlowingDustWarning = "Blowing Dust Warning"
	case BriskWindAdvisory = "Brisk Wind Advisory"
	case ChildAbductionEmergency = "Child Abduction Emergency"
	case CivilDangerWarning = "Civil Danger Warning"
	case CivilEmergencyMessage = "Civil Emergency Message"
	case CoastalFloodAdvisory = "Coastal Flood Advisory"
	case CoastalFloodStatement = "Coastal Flood Statement"
	case CoastalFloodWarning = "Coastal Flood Warning"
	case CoastalFloodWatch = "Coastal Flood Watch"
	case DenseFogAdvisory = "Dense Fog Advisory"
	case DenseSmokeAdvisory = "Dense Smoke Advisory"
	case DustAdvisory = "Dust Advisory"
	case DustStormWarning = "Dust Storm Warning"
	case EarthquakeWarning = "Earthquake Warning"
	case EvacuationImmediate = "Evacuation - Immediate"
	case ExcessiveHeatWarning = "Excessive Heat Warning"
	case ExcessiveHeatWatch = "Excessive Heat Watch"
	case ExtremeColdWarning = "Extreme Cold Warning"
	case ExtremeColdWatch = "Extreme Cold Watch"
	case ExtremeFireDanger = "Extreme Fire Danger"
	case ExtremeWindWarning = "Extreme Wind Warning"
	case FireWarning = "Fire Warning"
	case FireWeatherWatch = "Fire Weather Watch"
	case FlashFloodStatement = "Flash Flood Statement"
	case FlashFloodWarning = "Flash Flood Warning"
	case FlashFloodWatch = "Flash Flood Watch"
	case FloodAdvisory = "Flood Advisory"
	case FloodStatement = "Flood Statement"
	case FloodWarning = "Flood Warning"
	case FloodWatch = "Flood Watch"
	case FreezeWarning = "Freeze Warning"
	case FreezeWatch = "Freeze Watch"
	case FreezingFogAdvisory = "Freezing Fog Advisory"
	case FreezingRainAdvisory = "Freezing Rain Advisory"
	case FreezingSprayAdvisory = "Freezing Spray Advisory"
	case FrostAdvisory = "Frost Advisory"
	case GaleWarning = "Gale Warning"
	case GaleWatch = "Gale Watch"
	case HardFreezeWarning = "Hard Freeze Warning"
	case HardFreezeWatch = "Hard Freeze Watch"
	case HazardousMaterialsWarning = "Hazardous Materials Warning"
	case HazardousSeasWarning = "Hazardous Seas Warning"
	case HazardousSeasWatch = "Hazardous Seas Watch"
	case HazardousWeatherOutlook = "Hazardous Weather Outlook"
	case HeatAdvisory = "Heat Advisory"
	case HeavyFreezingSprayWarning = "Heavy Freezing Spray Warning"
	case HeavyFreezingSprayWatch = "Heavy Freezing Spray Watch"
	case HighSurfAdvisory = "High Surf Advisory"
	case HighSurfWarning = "High Surf Warning"
	case HighWindWarning = "High Wind Warning"
	case HighWindWatch = "High Wind Watch"
	case HurricaneForceWindWarning = "Hurricane Force Wind Warning"
	case HurricaneForceWindWatch = "Hurricane Force Wind Watch"
	case HurricaneLocalStatement = "Hurricane Local Statement"
	case HurricaneWarning = "Hurricane Warning"
	case HurricaneWatch = "Hurricane Watch"
	case HydrologicAdvisory = "Hydrologic Advisory"
	case HydrologicOutlook = "Hydrologic Outlook"
	case IceStormWarning = "Ice Storm Warning"
	case LakeEffectSnowAdvisory = "Lake Effect Snow Advisory"
	case LakeEffectSnowWarning = "Lake Effect Snow Warning"
	case LakeEffectSnowWatch = "Lake Effect Snow Watch"
	case LakeWindAdvisory = "Lake Wind Advisory"
	case LakeshoreFloodAdvisory = "Lakeshore Flood Advisory"
	case LakeshoreFloodStatement = "Lakeshore Flood Statement"
	case LakeshoreFloodWarning = "Lakeshore Flood Warning"
	case LakeshoreFloodWatch = "Lakeshore Flood Watch"
	case LawEnforcementWarning = "Law Enforcement Warning"
	case LocalAreaEmergency = "Local Area Emergency"
	case LowWaterAdvisory = "Low Water Advisory"
	case MarineWeatherStatement = "Marine Weather Statement"
	case NuclearPowerPlantWarning = "Nuclear Power Plant Warning"
	case RadiologicalHazardWarning = "Radiological Hazard Warning"
	case RedFlagWarning = "Red Flag Warning"
	case RipCurrentStatement = "Rip Current Statement"
	case SevereThunderstormWarning = "Severe Thunderstorm Warning"
	case SevereThunderstormWatch = "Severe Thunderstorm Watch"
	case SevereWeatherStatement = "Severe Weather Statement"
	case ShelterInPlaceWarning = "Shelter In Place Warning"
	case ShortTermForecast = "Short Term Forecast"
	case SmallCraftAdvisory = "Small Craft Advisory"
	case SmallCraftAdvisoryForHazardousSeas = "Small Craft Advisory For Hazardous Seas"
	case SmallCraftAdvisoryForRoughBar = "Small Craft Advisory For Rough Bar"
	case SmallCraftAdvisoryForWinds = "Small Craft Advisory For Winds"
	case SmallStreamFloodAdvisory = "Small Stream Flood Advisory"
	case SnowSquallWarning = "Snow Squall Warning"
	case SpecialMarineWarning = "Special Marine Warning"
	case SpecialWeatherStatement = "Special Weather Statement"
	case StormSurgeWarning = "Storm Surge Warning"
	case StormSurgeWatch = "Storm Surge Watch"
	case StormWarning = "Storm Warning"
	case StormWatch = "Storm Watch"
	case Test = "Test"
	case TornadoWarning = "Tornado Warning"
	case TornadoWatch = "Tornado Watch"
	case TropicalDepressionLocalStatement = "Tropical Depression Local Statement"
	case TropicalStormLocalStatement = "Tropical Storm Local Statement"
	case TropicalStormWarning = "Tropical Storm Warning"
	case TropicalStormWatch = "Tropical Storm Watch"
	case TsunamiAdvisory = "Tsunami Advisory"
	case TsunamiWarning = "Tsunami Warning"
	case TsunamiWatch = "Tsunami Watch"
	case TyphoonLocalStatement = "Typhoon Local Statement"
	case TyphoonWarning = "Typhoon Warning"
	case TyphoonWatch = "Typhoon Watch"
	case UrbanAndSmallStreamFloodAdvisory = "Urban And Small Stream Flood Advisory"
	case VolcanoWarning = "Volcano Warning"
	case WindAdvisory = "Wind Advisory"
	case WindChillAdvisory = "Wind Chill Advisory"
	case WindChillWarning = "Wind Chill Warning"
	case WindChillWatch = "Wind Chill Watch"
	case WinterStormWarning = "Winter Storm Warning"
	case WinterStormWatch = "Winter Storm Watch"
	case WinterWeatherAdvisory = "Winter Weather Advisory"
}

func WeatherAlertTypeOptions() -> Array<Option> {
	var convertAlerts: Array<Option> = []
	let allAlertsEnums: Array<WeatherAlertTypes> = WeatherAlertTypes.allCases

	for alert in allAlertsEnums {
		convertAlerts.append(Option(name: alert.rawValue))
	}
	return convertAlerts
}
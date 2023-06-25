package collections

import (
	"github.com/cmeyer18/weather-common/v2/custom_mongo"
	"github.com/cmeyer18/weather-common/v2/data_structures"
)

type SPCOutlookCollection struct {
	*custom_mongo.BaseCollection[data_structures.SPCOutlook]
}

func NewSPCOutlookCollection(connection custom_mongo.BaseConnection) (SPCOutlookCollection, bool, bool, error) {
	collection, dbExist, collectExist, err := custom_mongo.GetCollection[data_structures.SPCOutlook](spcOutlookCollectionArgs(), connection)
	if err != nil {
		return SPCOutlookCollection{}, false, false, err
	}

	collectionConv := SPCOutlookCollection{&collection}

	return collectionConv, dbExist, collectExist, nil
}

func spcOutlookCollectionArgs() custom_mongo.BaseCollectionArgs {
	return custom_mongo.BaseCollectionArgs{DatabaseName: "weather", CollectionName: "spcOutlook"}
}

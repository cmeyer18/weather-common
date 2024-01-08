package processing

func ProcessStringArray(arrayData interface{}) []string {
	if arrayData == nil {
		return nil
	}

	interfaceArrayData := arrayData.([]interface{})

	var returnedStringArray []string
	for _, rawData := range interfaceArrayData {
		returnedStringArray = append(returnedStringArray, rawData.(string))
	}

	return returnedStringArray
}

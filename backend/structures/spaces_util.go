package structures

type Space struct {
	Start int32
	End   int32
}

func ConvertToObjects[T any](list []T) []interface{} {
	objects := make([]interface{}, len(list))
	for i, item := range list {
		objects[i] = item
	}
	return objects
}

func getAvailableSpaces(objects []interface{}, start int32, end int32) []Space {
	var occupiedSpaces []Space
	for _, obj := range objects {
		var objStart, objEnd, objSize int32
		switch v := obj.(type) {
		case Partition:
			objStart = v.PartStart
			objSize = v.PartSize
		case EBR:
			objStart = v.PartStart
			objSize = v.PartSize
		default:
			continue
		}
		objEnd = objStart + objSize - 1
		if objStart <= end && objEnd >= start {
			if objStart < start {
				occupiedSpaces = append(occupiedSpaces, Space{objStart, min(objEnd, end)})
			} else {
				occupiedSpaces = append(occupiedSpaces, Space{objStart, min(objEnd, end)})
			}
		}
	}

	var availableSpaces []Space
	currentStart := start

	for _, occ := range occupiedSpaces {
		if occ.Start > currentStart {
			availableSpaces = append(availableSpaces, Space{currentStart, occ.Start - 1})
		}
		if occ.End >= currentStart {
			currentStart = occ.End + 1
		}
	}

	if end >= currentStart {
		availableSpaces = append(availableSpaces, Space{currentStart, end})
	}

	return availableSpaces
}

func FirstFit(objects []interface{}, blockSize int32, start int32, end int32) int32 {
	spaces := getAvailableSpaces(objects, start, end)
	for _, space := range spaces {
		if space.End-space.Start+1 >= blockSize {
			return space.Start
		}
	}
	return -1
}

func WorstFit(objects []interface{}, blockSize int32, start int32, end int32) int32 {
	spaces := getAvailableSpaces(objects, start, end)
	var largestSpace *Space
	for _, space := range spaces {
		if space.End-space.Start+1 >= blockSize {
			if largestSpace == nil || space.End-space.Start > largestSpace.End-largestSpace.Start {
				largestSpace = &space
			}
		}
	}
	if largestSpace != nil {
		return largestSpace.Start
	}
	return -1
}

func BestFit(objects []interface{}, blockSize int32, start int32, end int32) int32 {
	spaces := getAvailableSpaces(objects, start, end)
	var bestFit *Space
	for _, space := range spaces {
		if space.End-space.Start+1 >= blockSize {
			if bestFit == nil || space.End-space.Start < bestFit.End-bestFit.Start {
				bestFit = &space
			}
		}
	}
	if bestFit != nil {
		return bestFit.Start
	}
	return -1
}

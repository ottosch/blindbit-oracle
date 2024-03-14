package dblevel

import (
	"SilentPaymentAppBackend/src/common"
	"SilentPaymentAppBackend/src/common/types"
)

func InsertFilter(pair types.Filter) error {
	err := insertSimple(FiltersDB, &pair)
	if err != nil {
		common.ErrorLogger.Println(err)
		return err
	}
	common.DebugLogger.Println("Taproot Filter inserted")
	return nil
}

func FetchByBlockHashFilter(blockHash string) (types.Filter, error) {
	var pair types.Filter
	err := retrieveByBlockHash(FiltersDB, blockHash, &pair)
	if err != nil {
		common.ErrorLogger.Println(err)
		return types.Filter{}, err
	}
	return pair, nil
}

// FetchAllFilters returns all types.Filter in the DB
func FetchAllFilters() ([]types.Filter, error) {
	pairs, err := retrieveAll(FiltersDB, types.PairFactoryFilter)
	if err != nil {
		common.ErrorLogger.Println(err)
		return nil, err
	}
	if len(pairs) == 0 {
		common.WarningLogger.Println("Nothing returned")
		return nil, NoEntryErr{}
	}

	result := make([]types.Filter, len(pairs))
	// Convert each Pair to a Filter and assign it to the new slice
	for i, pair := range pairs {
		if pairPtr, ok := pair.(*types.Filter); ok {
			result[i] = *pairPtr
		} else {
			common.ErrorLogger.Printf("%+v\n", pair)
			panic("wrong pair struct returned")
		}
	}
	return result, err
}

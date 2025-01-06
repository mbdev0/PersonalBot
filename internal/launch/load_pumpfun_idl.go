package launch

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"pump_fun/internal/logger"
	"pump_fun/internal/models"
	"sync"
)

type IdlMap struct {
	AccountMap map[string]int
}

var (
	idlMapInstance map[string]IdlMap
	once           sync.Once
)

func GetIdlMap() map[string]IdlMap {
	once.Do(func() {
		loadPumpfunIdl()
	})
	fmt.Println(idlMapInstance)
	return idlMapInstance
}

func loadPumpfunIdl() {
	idl, err := loadIdlIntoModel()
	if err != nil {
		logger.Error(err)
		return
	}
	idlMapInstance = generateIdlMap(idl)
}

func loadIdlIntoModel() (*models.PumpFunIdl, error) {
	jsonFile, err := os.Open("internal/launch/pumpfun_idl/pump_fun_idl.json")
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var idl models.PumpFunIdl
	err = json.Unmarshal(byteValue, &idl)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return &idl, nil
}

/*
Due to the fact that we will need to continusly access the idl to find where the mint, user, and bondingCurve accounts are located
we will generate a map/dict to store the index of the accounts in the instruction. This will allow us to easily access the accounts in O(1) time
instead of having to loop through the accounts each time we need to access them.

The dict/map will look in the format of:
{buy : {mint: 0, user: 1, bondingCurve: 2}, create: {mint: 0, user: 1, bondingCurve: 2}}} etc
Which will correspond to the order of how the accounts are stored in the idl.
*/
func generateIdlMap(idl *models.PumpFunIdl) map[string]IdlMap {
	idlMap := make(map[string]IdlMap)
	for _, instruction := range idl.Instructions {
		for num, account := range instruction.Accounts {
			if _, ok := idlMap[instruction.Name]; !ok {
				idlMap[instruction.Name] = IdlMap{
					AccountMap: make(map[string]int),
				}
			}
			idlMap[instruction.Name].AccountMap[account.Name] = num
		}
	}

	return idlMap
}

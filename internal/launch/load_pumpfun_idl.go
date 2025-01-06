package launch

import (
	"encoding/json"
	"io"
	"os"
	"pump_fun/internal/logger"
	"pump_fun/internal/models"
	"sync"
)

type IdlMap struct {
	Instruction string
	AccountMap  map[string]int
}

var (
	idlMapInstance map[string]IdlMap
	once           sync.Once
)

func GetIdlMap() map[string]IdlMap {
	once.Do(func() {
		loadPumpfunIdl()
	})

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

func generateIdlMap(idl *models.PumpFunIdl) map[string]IdlMap {
	idlMap := make(map[string]IdlMap)
	for _, instruction := range idl.Instructions {
		for num, account := range instruction.Accounts {
			if _, ok := idlMap[instruction.Name]; !ok {
				idlMap[instruction.Name] = IdlMap{
					Instruction: instruction.Name,
					AccountMap:  make(map[string]int),
				}
			}
			idlMap[instruction.Name].AccountMap[account.Name] = num
		}
	}
	return idlMap
}

package core

import (
	"encoding/json"
	"onefinger/asset"
	"onefinger/structs"
	"onefinger/utils"
)

var FofaFingerPrints []structs.FofaFingerPrint

// 解析fofa.json
func init() {

	err := json.Unmarshal(asset.FofaFingerPrintString, &FofaFingerPrints)
	if err != nil {
		utils.OptionsError("Json parse error", 1)
	}
}

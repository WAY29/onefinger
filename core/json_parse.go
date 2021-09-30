package core

import (
	"encoding/json"

	"github.com/WAY29/onefinger/asset"
	"github.com/WAY29/onefinger/structs"
	"github.com/WAY29/onefinger/utils"
)

var FofaFingerPrints []structs.FofaFingerPrint

// 解析fofa.json
func init() {

	err := json.Unmarshal(asset.FofaFingerPrintString, &FofaFingerPrints)
	if err != nil {
		utils.OptionsError("Json parse error", 1)
	}
}

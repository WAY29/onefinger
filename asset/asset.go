package asset

import (
	_ "embed"
)

//go:embed fofa.json
var FofaFingerPrintString []byte

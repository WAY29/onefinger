package asset

import (
	_ "embed"
)

//go:embed fofa.min.json
var FofaFingerPrintString []byte

package local

import (
	_ "embed"
)

//go:embed riot.pem
var riotCertificate []byte

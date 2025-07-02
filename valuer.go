package hexid

import (
	"fmt"
	"sync/atomic"

	"github.com/webmafia/hexid/valuer"
)

var valuerType uint32

func SetValuerType(typ valuer.Type) error {
	switch typ {
	case valuer.Int64Valuer, valuer.Uint64Valuer, valuer.StringValuer, valuer.BinaryValuer:
		atomic.StoreUint32(&valuerType, uint32(typ))
		return nil
	}

	return fmt.Errorf("invalid ValuerType: %d", typ)
}

func getValuerType() valuer.Type {
	return valuer.Type(atomic.LoadUint32(&valuerType))
}

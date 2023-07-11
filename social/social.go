package social

import "github.com/lienkolabs/breeze/protocol"

var AxeCode = protocol.Code(1)

func FilterAction(action []byte) bool {
	return action[10] == 1
}

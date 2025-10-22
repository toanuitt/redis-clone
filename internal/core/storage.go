package core

import "redis-clone/internal/data_structure"

var dictStore *data_structure.Dict

func init() {
	dictStore = data_structure.CreateDict()
}

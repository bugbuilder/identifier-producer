package core

import "bennu.cl/identifier-producer/api/types"

type Service interface {
	Save(identifier types.Identifier) (string, error)
}

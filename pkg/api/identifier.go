package api

type IdentifierService interface {
	Save(identifier Identifier) (string, error)
}

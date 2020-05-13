package identifier

import (
	"bennu.cl/identifier-producer/api/server"
	"bennu.cl/identifier-producer/api/types"
	"net/http"
)

func (idr *identifierRouter) postIdentifier(w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	var id types.Identifier

	err := server.FromJSON(r.Body, &id)
	if err != nil {
		return err
	}

	key, err := idr.api.Save(id)
	if err != nil {
		return err
	}

	return server.ToJSON(w, http.StatusCreated, types.Identifier{TransactionId: key})
}

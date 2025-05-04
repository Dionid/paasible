package sqlify

import (
	"github.com/pocketbase/dbx"
)

func Upsert(query *dbx.ModelQuery) error {
	err := query.Insert()
	if err != nil {
		if !UniqueConstraintFailed(err, "") {
			return err
		}

		err = query.Update()
		if err != nil {
			return err
		}
	}

	return err
}

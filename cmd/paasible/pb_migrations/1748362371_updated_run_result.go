package pb_migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_1795773971")
		if err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(9, []byte(`{
			"cascadeDelete": false,
			"collectionId": "pbc_1903973042",
			"hidden": false,
			"id": "relation361630566",
			"maxSelect": 1,
			"minSelect": 0,
			"name": "target_id",
			"presentable": false,
			"required": false,
			"system": false,
			"type": "relation"
		}`)); err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(10, []byte(`{
			"cascadeDelete": false,
			"collectionId": "pbc_1420492469",
			"hidden": false,
			"id": "relation1040386765",
			"maxSelect": 1,
			"minSelect": 0,
			"name": "application_id",
			"presentable": false,
			"required": false,
			"system": false,
			"type": "relation"
		}`)); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_1795773971")
		if err != nil {
			return err
		}

		// remove field
		collection.Fields.RemoveById("relation361630566")

		// remove field
		collection.Fields.RemoveById("relation1040386765")

		return app.Save(collection)
	})
}

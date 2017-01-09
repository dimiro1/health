package health

/*
{
	"check_groups": [
		{
			"name": "Big Companies",
			"check_type": "http",
			"checks: [
				{
					"name": "Microsoft",
					"target": "https://www.microsoft.com/"
				},
				{
					"name":, "Oracle":,
					"tartget": "https://www.oracle.com/"
				},
				{
					"name":, "Google":,
					"tartget": "https://www.google.com/"
				}

			]
		},
		{
			"name": "MySQL Databases",
			"check_type": "db",
			"driver": "mysql",
			"checks": [
				{
					"name": "Primary",
					"target": "db.host.local"
				},
			]
		},
		{
			"name": "PostgreSQL Databases",
			"check_type": "db",
			"driver": "postgres",
			"checks": [
				{
					"name": "Primary",
					"target": "db.host.local"
				},
			]
		}

	]
}
*/

// Check holds the check itself
type Check struct {
	Name   string `json:"name"`
	Target string `json:"target"`
}

// CheckGroup
type CheckGroup struct {
	Name      string  `json:"name"`
	CheckType string  `json:"check_type"`
	Driver    string  `json:"driver"`
	Checks    []Check `json:"checks"`
}

// BulkChecks holds all checks to be imported and setup
type BulkChecks struct {
	CheckGroups []CheckGroup `json:"check_groups"`
}

// Validate makes sure the marhalled JSON contains valid data
func (b *BulkChecks) Validate() error {
	for _, i := range b.CheckGroups {
		if i.Driver != "" {
			if !ValidDBDriver(i.Driver) {
				return ErrInvalidDBDriver
			}
		}
		if !ValidCheck(i.CheckType) {
			return ErrInvalidCheck
		}
	}
	return nil
}

func (b *BulkChecks) Add() error {
	return nil
}

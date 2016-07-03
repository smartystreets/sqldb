package bindsql

import "github.com/smartystreets/sqldb"

type BindingSelectorAdapter struct {
	selector         sqldb.Selector
	panicOnBindError bool
}

func NewBindingSelectorAdapter(selector sqldb.Selector, panicOnBindError bool) *BindingSelectorAdapter {
	return &BindingSelectorAdapter{selector: selector, panicOnBindError: panicOnBindError}
}

func (this *BindingSelectorAdapter) Select(binder Binder, statement string, parameters ...interface{}) error {
	result, err := this.selector.Select(statement, parameters...)
	if err != nil {
		return err
	}

	for result.Next() {
		if err := result.Err(); err != nil {
			result.Close()
			return err
		}

		if err := binder(result); err != nil {
			result.Close()
			if this.panicOnBindError {
				panic(err)
			} else {
				return err
			}
		}
	}

	return result.Close()
}
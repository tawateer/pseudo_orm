package orm


func execByTx(query string, args ...interface{}) error {
	tx, err := txGetter()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(query, args); err != nil {
		return err
	}
	return tx.Commit()
}

func namedExecByTx(query string, arg interface{}) error {
	tx, err := txGetter()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.NamedExec(query, arg); err != nil {
		return err
	}
	return tx.Commit()
}

func getByTx(resp interface{}, query string, args ...interface{}) error {
	tx, err := txGetter()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := tx.Get(resp, query, args); err != nil {
		return err
	}
	return tx.Commit()
}

func namedGetByTx(resp interface{}, query string, arg interface{}) error {
	tx, err := txGetter()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query, args, err := tx.BindNamed(query, arg)
	if err != nil {
		return err
	}

	if err := tx.Get(resp, query, args); err != nil {
		return err
	}
	return tx.Commit()
}

func selectByTx(resp interface{}, query string, args ...interface{}) error {
	tx, err := txGetter()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := tx.Select(resp, query, args); err != nil {
		return err
	}
	return tx.Commit()
}

func namedSelectByTx(resp interface{}, query string, arg interface{}) error {
	tx, err := txGetter()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query, args, err := tx.BindNamed(query, arg)
	if err != nil {
		return err
	}

	if err := tx.Select(resp, query, args); err != nil {
		return err
	}
	return tx.Commit()
}

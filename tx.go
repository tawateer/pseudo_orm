package orm


func execSqlByTx(query string, args ...interface{}) error {
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

func namedExecSqlByTx(query string, arg interface{}) error {
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

func getSqlByTx(resp interface{}, query string, args ...interface{}) error {
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

func namedGetSqlByTx(resp interface{}, query string, arg interface{}) error {
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

func selectSqlByTx(resp interface{}, query string, args ...interface{}) error {
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

func namedSelectSqlByTx(resp interface{}, query string, arg interface{}) error {
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

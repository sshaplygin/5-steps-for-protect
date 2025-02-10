package internal

func (a *App) InitTables() error {
	_, err := a.db.Exec(`
		CREATE SEQUENCE IF NOT EXISTS users_id_sequence START 1;

		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY DEFAULT nextval('users_id_sequence'), 
			login TEXT UNIQUE, 
			password TEXT, 
			is_admin BOOLEAN,
		);
	`)
	if err != nil {
		return err
	}

	_, err = a.db.Exec(`
		CREATE SEQUENCE IF NOT EXISTS posts_id_sequence START 1;

		CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY DEFAULT nextval('posts_id_sequence'), 
			title TEXT, 
			content TEXT, 
			user_id INTEGER, 
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return err
	}

	return nil
}

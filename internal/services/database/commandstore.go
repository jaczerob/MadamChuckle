package database

import "github.com/zekrotja/ken/store"

var _ store.CommandStore = (*Database)(nil)

func (c *Database) Store(cmds map[string]string) (err error) {
	tx, err := c.db.Begin()
	if err != nil {
		return
	}

	stmt, err := tx.Prepare("INSERT OR REPLACE INTO commands(id, command) VALUES(?, ?)")
	if err != nil {
		return
	}

	defer stmt.Close()

	for id, command := range cmds {
		_, err = stmt.Exec(id, command)
		if err != nil {
			return
		}
	}

	return tx.Commit()
}

func (c *Database) Load() (cmds map[string]string, err error) {
	cmds = make(map[string]string)

	rows, err := c.db.Query("SELECT * FROM commands")
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		var id, command string

		if err = rows.Scan(&id, &command); err != nil {
			return
		}

		cmds[id] = command
	}

	err = rows.Err()
	return
}

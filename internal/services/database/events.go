package database

type Event struct {
	ID        int64
	MessageID string
	ChannelID string
}

func (c *Database) RegisterEvent(eventID int64, messageID string, channelID string) (err error) {
	stmt, err := c.db.Prepare("INSERT INTO events(event_id, message_id, channel_id) VALUES(?, ?, ?)")
	if err != nil {
		return
	}

	defer stmt.Close()

	_, err = stmt.Exec(eventID, messageID, channelID)
	return
}

func (c *Database) UnregisterEvent(messageID string, channelID string) (err error) {
	stmt, err := c.db.Prepare("DELETE FROM events WHERE message_id=?, channel_id=?")
	if err != nil {
		return
	}

	defer stmt.Close()

	_, err = stmt.Exec(messageID, channelID)
	return
}

func (c *Database) GetEvents() (events []*Event, err error) {
	events = make([]*Event, 0)

	rows, err := c.db.Query("SELECT * FROM events")
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		var e int64
		var m, c string

		if err = rows.Scan(&e, &m, &c); err != nil {
			return
		}

		events = append(events, &Event{ID: e, MessageID: m, ChannelID: c})
	}

	return
}

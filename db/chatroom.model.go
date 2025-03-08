package db

type Chatroom struct {
	ID   int
	Name string
}

func (db *DB) initChatroomModel() error {
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS chatrooms (
            id INTEGER PRIMARY KEY,
            name TEXT NOT NULL UNIQUE
        );
        INSERT OR IGNORE INTO chatrooms (name) VALUES ('General Chat');
    `)

	return err
}

func (db *DB) GetChatrooms() ([]Chatroom, error) {
	rows, err := db.Query("SELECT id, name FROM chatrooms")
	if err != nil {
		return nil, err
	}
	return scanRows(rows, func(chatroom *Chatroom) error {
		return rows.Scan(&chatroom.ID, &chatroom.Name)
	})
}

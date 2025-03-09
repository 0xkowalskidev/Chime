package db

type Message struct {
	ID         int
	ChatroomID int
	Username   string
	Content    string
}

func (d *DB) initMessagesModel() error {
	_, err := d.Exec(`
        CREATE TABLE IF NOT EXISTS messages (
            id INTEGER PRIMARY KEY,
            chatroom_id INTEGER NOT NULL,
            username TEXT NOT NULL,
            content TEXT NOT NULL,
            FOREIGN KEY (chatroom_id) REFERENCES chatrooms(id)
        );
    `)
	return err
}

func (db *DB) GetMessages(chatroomID int) ([]Message, error) {
	rows, err := db.Query(`
        SELECT id, chatroom_id, username, content 
        FROM messages 
        WHERE chatroom_id = ?
    `, chatroomID)
	if err != nil {
		return nil, err
	}
	return scanRows(rows, func(msg *Message) error {
		return rows.Scan(&msg.ID, &msg.ChatroomID, &msg.Username, &msg.Content)
	})
}

func (db *DB) CreateMessage(chatroomID int, username, content string) (Message, error) {
	result, err := db.Exec(`
        INSERT INTO messages (chatroom_id, username, content) VALUES (?, ?, ?)
    `, chatroomID, username, content)
	if err != nil {
		return Message{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return Message{}, err
	}

	return Message{
		ID:         int(id),
		ChatroomID: chatroomID,
		Username:   username,
		Content:    content,
	}, nil
}

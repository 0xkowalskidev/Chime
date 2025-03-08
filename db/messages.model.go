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
        INSERT OR IGNORE INTO messages (chatroom_id, username, content)
        VALUES ((SELECT id FROM chatrooms WHERE name = 'General Chat'), 'User1', 'Hello everyone!');
        INSERT OR IGNORE INTO messages (chatroom_id, username, content)
        VALUES ((SELECT id FROM chatrooms WHERE name = 'General Chat'), 'User2', 'Hi there!');
    `)
	return err
}

func (d *DB) GetMessages(chatroomID int) ([]Message, error) {
	rows, err := d.Query(`
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

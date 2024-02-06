-- Table Items is a basic table containing mysterious 'items'.
CREATE TABLE IF NOT EXISTS items (
   ID            INTEGER PRIMARY KEY,
   TITLE         TEXT    NOT NULL DEFAULT "",
   DESC          TEXT    NOT NULL DEFAULT "",
   CREATION_DATE INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
);

-- Populate Items with some sample data.
INSERT INTO items(TITLE, DESC) VALUES("Item 1", "This is item 1");
INSERT INTO items(TITLE, DESC) VALUES("Item 2", "This is item 2");
INSERT INTO items(TITLE, DESC) VALUES("Item 3", "This is item 3");
INSERT INTO items(TITLE, DESC) VALUES("Item 4", "This is item 4");


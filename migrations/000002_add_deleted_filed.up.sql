ALTER TABLE notes ADD COLUMN deleted BOOLEAN NOT NULL DEFAULT false;
CREATE INDEX idx_notes_deleted ON notes(deleted);
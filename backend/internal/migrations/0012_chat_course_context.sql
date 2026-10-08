-- Keep course/lesson conversations separate from general assistance.
ALTER TABLE ai_conversations ADD COLUMN program_id uuid REFERENCES programs(id) ON DELETE CASCADE;
ALTER TABLE ai_conversations ADD COLUMN node_id uuid REFERENCES nodes(id) ON DELETE CASCADE;
ALTER TABLE ai_conversations ADD CONSTRAINT ai_conversations_scope_check CHECK (node_id IS NULL OR program_id IS NOT NULL);

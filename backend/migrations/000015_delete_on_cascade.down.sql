-- Down migration: Remove CASCADE delete and restore original foreign keys

-- Sessions
ALTER TABLE sessions DROP CONSTRAINT sessions_user_id_fkey;

ALTER TABLE sessions ADD CONSTRAINT sessions_user_id_fkey 
    FOREIGN KEY (user_id) REFERENCES users(id);

-- Messages  
ALTER TABLE messages DROP CONSTRAINT messages_user_id_fkey;
ALTER TABLE messages DROP CONSTRAINT messages_thread_id_fkey;

ALTER TABLE messages ADD CONSTRAINT messages_user_id_fkey 
    FOREIGN KEY (user_id) REFERENCES users(id);

ALTER TABLE messages ADD CONSTRAINT messages_thread_id_fkey 
    FOREIGN KEY (thread_id) REFERENCES threads(id);

-- Threads
ALTER TABLE threads DROP CONSTRAINT threads_user_id_fkey;

ALTER TABLE threads ADD CONSTRAINT threads_user_id_fkey 
    FOREIGN KEY (user_id) REFERENCES users(id);
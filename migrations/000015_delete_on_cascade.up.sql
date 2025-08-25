-- 1. Sessions (no dependencies)
ALTER TABLE sessions DROP CONSTRAINT sessions_user_id_fkey;
ALTER TABLE sessions ADD CONSTRAINT sessions_user_id_fkey 
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- 2. Threads (before messages, since messages reference threads)
ALTER TABLE threads DROP CONSTRAINT threads_user_id_fkey;
ALTER TABLE threads ADD CONSTRAINT threads_user_id_fkey 
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- 3. Messages (after threads)
ALTER TABLE messages DROP CONSTRAINT messages_user_id_fkey;
ALTER TABLE messages DROP CONSTRAINT messages_thread_id_fkey;
ALTER TABLE messages ADD CONSTRAINT messages_user_id_fkey 
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE messages ADD CONSTRAINT messages_thread_id_fkey 
    FOREIGN KEY (thread_id) REFERENCES threads(id) ON DELETE CASCADE;

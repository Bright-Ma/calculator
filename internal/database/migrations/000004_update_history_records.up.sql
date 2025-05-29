-- 修改 history_records 表结构
ALTER TABLE history_records
CHANGE COLUMN question_content question VARCHAR(255) NOT NULL; 
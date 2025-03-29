-- 创建备份记录表，用于存储每次备份任务的详细记录
CREATE TABLE IF NOT EXISTS backup_records (
    version_id INTEGER PRIMARY KEY, -- 唯一标识每次备份的版本号
    task_id INTEGER, -- 关联的备份任务 ID
    timestamp TEXT, -- 备份任务的时间戳
    task_name TEXT, -- 备份任务的名称
    backup_status TEXT, -- 备份任务的状态（例如成功、失败等）
    backup_file_name TEXT, -- 生成的备份文件名称
    backup_size TEXT, -- 备份文件的大小
    backup_path TEXT, -- 备份文件的存储路径
    version_hash TEXT, -- 备份版本的哈希值，用于校验
    FOREIGN KEY (task_id) REFERENCES backup_tasks (task_id) ON DELETE CASCADE -- 定义外键关系，删除备份任务时自动删除相关备份记录
);

-- 创建备份任务表，用于定义备份任务的基本信息
CREATE TABLE IF NOT EXISTS backup_tasks (
    task_id INTEGER PRIMARY KEY, -- 唯一标识备份任务的 ID
    task_name TEXT, -- 备份任务的名称
    target_directory TEXT, -- 需要备份的目标目录
    backup_directory TEXT, -- 备份文件存放的目标目录
    retention_count INTEGER -- 保留的备份版本数量
);
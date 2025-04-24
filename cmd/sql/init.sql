-- 创建备份任务表，用于定义备份任务的基本信息
CREATE TABLE IF NOT EXISTS backup_tasks (
    task_id INTEGER PRIMARY KEY AUTOINCREMENT, -- 唯一标识备份任务的 ID （自动递增）
    task_name TEXT, -- 备份任务的名称
    target_directory TEXT, -- 需要备份的目标目录
    backup_directory TEXT, -- 备份文件存放的目标目录
    retention_count INTEGER, -- 保留数量
    retention_days INTEGER, -- 保留天数
    no_compression INTEGER,  -- 是否禁用压缩（默认启用压缩，设置为 0 表示启用压缩, 1 表示禁用压缩）
    exclude_rules TEXT -- 用于存储排除规则 （例如: *.txt, *.jpg）"none" 表示不排除任何文件
);

-- 添加索引，用于提高查询效率
CREATE INDEX IF NOT EXISTS idx_backup_tasks_task_name ON backup_tasks (task_name);

-- 创建备份记录表，用于存储每次备份任务的详细记录
CREATE TABLE IF NOT EXISTS backup_records (
    version_id TEXT PRIMARY KEY, -- 唯一标识每次备份的版本号
    task_id INTEGER, -- 关联的备份任务 ID
    timestamp TEXT, -- 备份任务的时间戳
    task_name TEXT, -- 备份任务的名称
    backup_status TEXT, -- 备份任务的状态（例如: true 表示成功 false 表示失败）
    backup_file_name TEXT, -- 生成的备份文件名称
    backup_size TEXT, -- 备份文件的大小
    backup_path TEXT, -- 备份文件的存储路径
    version_hash TEXT-- 备份版本的哈希值，用于校验
);

-- 给备份记录表添加索引，用于提高查询效率 
CREATE INDEX IF NOT EXISTS idx_backup_records_timestamp ON backup_records (timestamp);

-- 给备份记录表添加索引，用于提高查询效率
CREATE INDEX IF NOT EXISTS idx_backup_records_task_id ON backup_records (task_id);
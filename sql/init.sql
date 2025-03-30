-- 创建备份任务表，用于定义备份任务的基本信息
CREATE TABLE IF NOT EXISTS backup_tasks (
    task_id INTEGER PRIMARY KEY AUTOINCREMENT, -- 唯一标识备份任务的 ID （自动递增）
    task_name TEXT, -- 备份任务的名称
    target_directory TEXT, -- 需要备份的目标目录
    backup_directory TEXT, -- 备份文件存放的目标目录
    retention_count INTEGER -- 保留的备份版本数量
);

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
    version_hash TEXT, -- 备份版本的哈希值，用于校验
    data_status TEXT -- 数据状态（例如: 1 表示存在 0 表示不存在）
);

-- 给备份记录表添加索引，用于提高查询效率 
CREATE INDEX IF NOT EXISTS idx_backup_records_timestamp ON backup_records (timestamp);

-- 给备份记录表添加索引，用于提高查询效率
CREATE INDEX IF NOT EXISTS idx_backup_records_task_id ON backup_records (task_id);

-- 创建压缩配置表，用于存储不同操作系统下的压缩工具配置
CREATE TABLE IF NOT EXISTS compress_config (
    id INTEGER PRIMARY KEY AUTOINCREMENT, -- 主键，自增
    os_type TEXT, -- 操作系统类型（如 linux 或 windows）
    compress_tool TEXT, -- 压缩工具名称（如 tar 或 7z）
    compress_args TEXT, -- 压缩工具的参数（如 -czf 或 a -tzip）
    file_extension TEXT -- 文件扩展名（如 .tgz 或 .zip）
);

-- 插入初始的压缩配置数据
INSERT INTO compress_config (os_type, compress_tool, compress_args, file_extension)
VALUES ('linux', 'tar', '-c|-zf', '.tgz');
INSERT INTO compress_config (os_type, compress_tool, compress_args, file_extension)
VALUES ('windows', '7z', 'a|-tzip', '.zip');

-- 创建解压缩配置表，用于存储不同操作系统下的解压缩工具配置
CREATE TABLE IF NOT EXISTS decompress_config (
    id INTEGER PRIMARY KEY AUTOINCREMENT, -- 主键，自增
    os_type TEXT, -- 操作系统类型（如 linux 或 windows）
    decompress_tool TEXT, -- 解压缩工具名称（如 tar 或 7z）
    decompress_args TEXT, -- 解压缩工具的参数（如 -xzf 或 x -tzip）
    file_extension TEXT -- 文件扩展名（如 .tgz 或 .zip）
);

-- 插入初始的解压缩配置数据
INSERT INTO decompress_config (os_type, decompress_tool, decompress_args, file_extension)
VALUES ('linux', 'tar', '-xzf|-C', '.tgz');

INSERT INTO decompress_config (os_type, decompress_tool, decompress_args, file_extension)
VALUES ('windows', '7z', 'x -tzip|-o', '.zip');
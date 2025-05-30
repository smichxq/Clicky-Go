-- 不往 binlog 里写下面的用户创建操作
SET @@SESSION.sql_log_bin = 0;
-- 创建数据库（逻辑库）
CREATE DATABASE IF NOT EXISTS ks1;

CREATE DATABASE IF NOT EXISTS vt_ks1;

-- 创建 _vt 库（Vitess 元数据库），某些镜像会自动创建，可保守确保创建
CREATE DATABASE IF NOT EXISTS _vt;

-- 1. 应用层用户 vt_app（业务访问）
CREATE USER 'vt_app' @'%' IDENTIFIED BY '123';

GRANT SELECT, INSERT , UPDATE, DELETE ON ks1.* TO 'vt_app' @'%';

GRANT SELECT, INSERT , UPDATE, DELETE ON vt_ks1.* TO 'vt_app' @'%';

GRANT ALL ON _vt.* TO 'vt_app' @'%';

-- 2. VReplication 复制用户 vt_filtered（用于 filtered replication / materialized view）
CREATE USER 'vt_filtered' @'%' IDENTIFIED BY '123';

GRANT SELECT ON _vt.* TO 'vt_filtered' @'%';

GRANT SELECT ON vt_ks1.* TO 'vt_filtered' @'%';

-- 3. 主从复制用户 vt_repl（MySQL 层主从用）
CREATE USER 'vt_repl' @'%' IDENTIFIED BY '123';

GRANT REPLICATION SLAVE ON *.* TO 'vt_repl' @'%';

-- 4. 管理员用户 vt_allprivs（用于初始数据迁移、备份恢复等）
CREATE USER 'vt_allprivs' @'%' IDENTIFIED BY '123';

GRANT SUPER,
PROCESS,
REPLICATION SLAVE,
REPLICATION CLIENT,
RELOAD ON *.* TO 'vt_allprivs' @'%';

GRANT ALL PRIVILEGES ON ks1.* TO 'vt_allprivs' @'%';

GRANT ALL PRIVILEGES ON vt_ks1.* TO 'vt_allprivs' @'%';

GRANT ALL ON _vt.* TO 'vt_allprivs' @'%';

-- 授予对系统表的读取权限
GRANT SELECT ON mysql.user TO 'vt_allprivs' @'%';

GRANT SELECT ON mysql.db TO 'vt_allprivs' @'%';

GRANT SELECT ON performance_schema.* TO 'vt_allprivs' @'%';

-- 若用 XtraBackup、分片恢复等功能，可添加更全权限
GRANT
    SUPER,
    PROCESS,
    REPLICATION SLAVE,
    REPLICATION CLIENT,
    RELOAD,
    LOCK TABLES,
    CREATE,
    ALTER,
    DROP,
    INDEX,
    INSERT,
    SELECT,
    UPDATE,
    DELETE,
    EVENT,
    TRIGGER
ON *.* TO 'vt_allprivs'@'%';

-- 5. Vitess 运维工具 vtorc 使用的账号（可选，但推荐）
CREATE USER 'vtorc_user' @'%' IDENTIFIED BY '123';

GRANT PROCESS,
REPLICATION SLAVE,
REPLICATION CLIENT,
SUPER,
RELOAD,
SELECT ON *.* TO 'vtorc_user' @'%';

GRANT SELECT ON mysql.* TO 'vtorc_user' @'%';

GRANT SELECT ON performance_schema.* TO 'vtorc_user' @'%';

-- 6. 完全管理员用户 vt_dba（用于调试、Vitess admin 操作）
CREATE USER 'vt_dba' @'%' IDENTIFIED BY '123';

GRANT ALL PRIVILEGES ON *.* TO 'vt_dba' @'%';

-- 应用所有权限更改
FLUSH PRIVILEGES;

-- 恢复 binlog
SET @@SESSION.sql_log_bin = 1;
-- JXZY AI应用服务数据库初始化脚本

-- 创建数据库
CREATE DATABASE IF NOT EXISTS jxzy_context CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE DATABASE IF NOT EXISTS jxzy_prompt CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE DATABASE IF NOT EXISTS jxzy_rag CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 使用jxzy_context数据库
USE jxzy_context;

-- 上下文表
CREATE TABLE IF NOT EXISTS contexts (
    id VARCHAR(64) PRIMARY KEY COMMENT '上下文ID',
    name VARCHAR(255) NOT NULL COMMENT '上下文名称',
    description TEXT COMMENT '上下文描述',
    type VARCHAR(32) NOT NULL COMMENT '上下文类型: conversation, document',
    user_id VARCHAR(64) NOT NULL COMMENT '用户ID',
    metadata JSON COMMENT '元数据',
    created_at BIGINT NOT NULL COMMENT '创建时间戳',
    updated_at BIGINT NOT NULL COMMENT '更新时间戳',
    INDEX idx_user_id (user_id),
    INDEX idx_type (type),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='上下文表';

-- 消息表
CREATE TABLE IF NOT EXISTS context_messages (
    id VARCHAR(64) PRIMARY KEY COMMENT '消息ID',
    context_id VARCHAR(64) NOT NULL COMMENT '上下文ID',
    role VARCHAR(32) NOT NULL COMMENT '角色: user, assistant, system',
    content TEXT NOT NULL COMMENT '消息内容',
    metadata JSON COMMENT '元数据',
    created_at BIGINT NOT NULL COMMENT '创建时间戳',
    INDEX idx_context_id (context_id),
    INDEX idx_created_at (created_at),
    FOREIGN KEY (context_id) REFERENCES contexts(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='上下文消息表';

-- 使用jxzy_prompt数据库
USE jxzy_prompt;

-- Prompt表
CREATE TABLE IF NOT EXISTS prompts (
    id VARCHAR(64) PRIMARY KEY COMMENT 'Prompt ID',
    name VARCHAR(255) NOT NULL COMMENT 'Prompt名称',
    content TEXT NOT NULL COMMENT 'Prompt内容',
    variables JSON COMMENT '变量定义',
    category VARCHAR(64) COMMENT '分类',
    description TEXT COMMENT '描述',
    user_id VARCHAR(64) NOT NULL COMMENT '用户ID',
    is_public BOOLEAN DEFAULT FALSE COMMENT '是否公开',
    version VARCHAR(32) DEFAULT 'v1.0' COMMENT '版本号',
    tags JSON COMMENT '标签',
    created_at BIGINT NOT NULL COMMENT '创建时间戳',
    updated_at BIGINT NOT NULL COMMENT '更新时间戳',
    INDEX idx_user_id (user_id),
    INDEX idx_category (category),
    INDEX idx_is_public (is_public),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Prompt表';

-- Prompt版本表
CREATE TABLE IF NOT EXISTS prompt_versions (
    id VARCHAR(64) PRIMARY KEY COMMENT '版本ID',
    prompt_id VARCHAR(64) NOT NULL COMMENT 'Prompt ID',
    version VARCHAR(32) NOT NULL COMMENT '版本号',
    content TEXT NOT NULL COMMENT 'Prompt内容',
    variables JSON COMMENT '变量定义',
    change_log TEXT COMMENT '变更日志',
    created_at BIGINT NOT NULL COMMENT '创建时间戳',
    INDEX idx_prompt_id (prompt_id),
    INDEX idx_created_at (created_at),
    FOREIGN KEY (prompt_id) REFERENCES prompts(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Prompt版本表';

-- 使用jxzy_rag数据库
USE jxzy_rag;

-- 文档集合表
CREATE TABLE IF NOT EXISTS document_collections (
    id VARCHAR(64) PRIMARY KEY COMMENT '集合ID',
    name VARCHAR(255) NOT NULL COMMENT '集合名称',
    description TEXT COMMENT '集合描述',
    user_id VARCHAR(64) NOT NULL COMMENT '用户ID',
    document_count INT DEFAULT 0 COMMENT '文档数量',
    vector_store_config JSON COMMENT '向量存储配置',
    metadata JSON COMMENT '元数据',
    created_at BIGINT NOT NULL COMMENT '创建时间戳',
    updated_at BIGINT NOT NULL COMMENT '更新时间戳',
    INDEX idx_user_id (user_id),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文档集合表';

-- 文档表
CREATE TABLE IF NOT EXISTS documents (
    id VARCHAR(64) PRIMARY KEY COMMENT '文档ID',
    title VARCHAR(255) NOT NULL COMMENT '文档标题',
    content LONGTEXT COMMENT '文档内容',
    type VARCHAR(32) NOT NULL COMMENT '文档类型: pdf, txt, md, html, docx',
    source VARCHAR(500) COMMENT '文档来源',
    collection_id VARCHAR(64) NOT NULL COMMENT '集合ID',
    user_id VARCHAR(64) NOT NULL COMMENT '用户ID',
    status VARCHAR(32) DEFAULT 'processing' COMMENT '状态: processing, ready, failed',
    metadata JSON COMMENT '元数据',
    created_at BIGINT NOT NULL COMMENT '创建时间戳',
    updated_at BIGINT NOT NULL COMMENT '更新时间戳',
    INDEX idx_collection_id (collection_id),
    INDEX idx_user_id (user_id),
    INDEX idx_status (status),
    INDEX idx_type (type),
    INDEX idx_created_at (created_at),
    FOREIGN KEY (collection_id) REFERENCES document_collections(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文档表';

-- 文档块表
CREATE TABLE IF NOT EXISTS document_chunks (
    id VARCHAR(64) PRIMARY KEY COMMENT '块ID',
    document_id VARCHAR(64) NOT NULL COMMENT '文档ID',
    content TEXT NOT NULL COMMENT '块内容',
    chunk_index INT NOT NULL COMMENT '块索引',
    start_offset INT COMMENT '开始偏移量',
    end_offset INT COMMENT '结束偏移量',
    vector_id VARCHAR(128) COMMENT '向量数据库中的ID',
    metadata JSON COMMENT '元数据',
    created_at BIGINT NOT NULL COMMENT '创建时间戳',
    INDEX idx_document_id (document_id),
    INDEX idx_chunk_index (chunk_index),
    INDEX idx_vector_id (vector_id),
    FOREIGN KEY (document_id) REFERENCES documents(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文档块表';

-- 处理任务表
CREATE TABLE IF NOT EXISTS processing_tasks (
    id VARCHAR(64) PRIMARY KEY COMMENT '任务ID',
    type VARCHAR(32) NOT NULL COMMENT '任务类型: document_upload, reindex',
    status VARCHAR(32) DEFAULT 'pending' COMMENT '状态: pending, processing, completed, failed',
    progress INT DEFAULT 0 COMMENT '进度: 0-100',
    message TEXT COMMENT '状态消息',
    user_id VARCHAR(64) NOT NULL COMMENT '用户ID',
    target_id VARCHAR(64) COMMENT '目标ID（文档ID或集合ID）',
    details JSON COMMENT '任务详情',
    started_at BIGINT COMMENT '开始时间戳',
    completed_at BIGINT COMMENT '完成时间戳',
    created_at BIGINT NOT NULL COMMENT '创建时间戳',
    updated_at BIGINT NOT NULL COMMENT '更新时间戳',
    INDEX idx_user_id (user_id),
    INDEX idx_status (status),
    INDEX idx_type (type),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='处理任务表';

-- 创建用户并授权（仅开发环境）
CREATE USER IF NOT EXISTS 'jxzy'@'%' IDENTIFIED BY 'jxzy123456';
GRANT ALL PRIVILEGES ON jxzy_context.* TO 'jxzy'@'%';
GRANT ALL PRIVILEGES ON jxzy_prompt.* TO 'jxzy'@'%';
GRANT ALL PRIVILEGES ON jxzy_rag.* TO 'jxzy'@'%';
FLUSH PRIVILEGES;

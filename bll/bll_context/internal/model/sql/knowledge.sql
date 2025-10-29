-- 知识表：对应现实中的一个文件，对应标准语义协议里的完整知识原文，用于管理原始知识文件，可以采用oss辅助存储
CREATE TABLE knowledge_file (
    id BIGINT AUTO_INCREMENT COMMENT '记录ID',
    oss_path VARCHAR(255) NOT NULL COMMENT 'OSS路径',
    file_name VARCHAR(255) NOT NULL COMMENT '文件名称',
    file_size BIGINT NOT NULL COMMENT '文件大小，单位bytes',
    file_type VARCHAR(50) NOT NULL COMMENT '文件类型',
    file_md5 VARCHAR(32) NOT NULL COMMENT '文件MD5',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='知识文件表';

-- 语义段表：对应标准语义协议里的语义段（原文存储于此）
CREATE TABLE knowledge_segment (
    id BIGINT AUTO_INCREMENT COMMENT '记录ID',
    knowledge_file_id BIGINT NOT NULL COMMENT '关联的知识文件ID，对应knowledge_file表',
    segment_text VARCHAR(4096) NOT NULL COMMENT '语义段文本',
    segment_md5 VARCHAR(32) NOT NULL COMMENT '语义段MD5',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='语义段表';

-- 会话表：记录用户创建的会话信息
CREATE TABLE chat_session (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '会话ID',
    name VARCHAR(100) NOT NULL COMMENT '会话名称（用户可自定义）',
    scene_code VARCHAR(50) NOT NULL COMMENT '关联的场景编码，对应llm_scene表',
    user_id VARCHAR(50) NOT NULL COMMENT '会话所属用户ID',
    is_active TINYINT NOT NULL DEFAULT 1 COMMENT '是否为活跃会话（1-是，0-否）',
    last_interact_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '最后交互时间',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '会话创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '会话更新时间',
    INDEX idx_user_id (user_id),
    INDEX idx_scene_code (scene_code),
    INDEX idx_last_interact_time (last_interact_time)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='LLM对话会话表';

-- 会话问答关联表：记录会话内的所有问答记录关联
CREATE TABLE chat_session_qas (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '记录ID',
    session_id BIGINT NOT NULL COMMENT '关联的会话ID，对应chat_session表',
    llm_completion_id BIGINT NOT NULL COMMENT '关联的问答详情ID，对应llm_completion表',
    sequence_num INT NOT NULL COMMENT '会话内的顺序编号（用于保证对话顺序）',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    -- 唯一索引，确保一个问答记录不会重复关联到同一个会话
    UNIQUE KEY uk_session_completion (session_id, llm_completion_id),
    INDEX idx_session_id (session_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='会话与问答记录关联表';

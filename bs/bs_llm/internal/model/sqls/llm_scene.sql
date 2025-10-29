-- 用来封装llm提供商/llm类型/llm说明等，比如上游传递下来一个场景，这个场景需要调用哪个供应商（比如豆包）的哪个llmmodel，这样对上游调用方就是透明的，只需要关心一个场景就行。
CREATE TABLE llm_scene (
    id INT AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    scene_code VARCHAR(50) NOT NULL DEFAULT '' COMMENT '场景编码（上游调用时使用的标识）',
    scene_name VARCHAR(100) NOT NULL DEFAULT '' COMMENT '场景名称',
    provider_code VARCHAR(50) NOT NULL DEFAULT '' COMMENT 'LLM提供商编码（如doubao、openai等）',
    provider_name VARCHAR(100) NOT NULL DEFAULT '' COMMENT 'LLM提供商名称',
    model_code VARCHAR(50) NOT NULL DEFAULT '' COMMENT '模型编码（如doubao-pro、gpt-4等）',
    model_name VARCHAR(100) NOT NULL DEFAULT '' COMMENT '模型名称',
    model_description TEXT COMMENT '模型说明（如适用场景、能力特点等）',
    scene_description TEXT COMMENT '场景说明（描述该场景的业务含义和使用场景）',
    temperature DECIMAL(3,2) DEFAULT 0.70 COMMENT '温度参数（0.00-1.00），控制生成内容的随机性',
    max_tokens INT DEFAULT 1000 COMMENT '最大token数，限制单次生成的最大长度',
    enable_stream TINYINT(1) DEFAULT 1 COMMENT '是否启用流式输出（1-启用，0-禁用）',
    deleted TINYINT NOT NULL DEFAULT 0 COMMENT '是否删除（1-删除，0-未删除）',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY u_scene_code (scene_code),
    INDEX idx_provider_model (provider_code, model_code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='LLM场景映射表（关联场景与对应的LLM提供商及模型）';

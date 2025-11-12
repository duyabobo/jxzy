-- 用来封装embedding提供商/embedding模型名/embedding向量维度，比如上游传递下来一个场景，这个场景需要调用哪个供应商（比如dashvector）的哪个embeddingmodel，这样对上游调用方就是透明的，只需要关心一个场景就行。
CREATE TABLE embedding_scene (
    id INT AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    scene_code VARCHAR(50) NOT NULL DEFAULT '' COMMENT '场景编码（上游调用时使用的标识）',
    scene_name VARCHAR(100) NOT NULL DEFAULT '' COMMENT '场景名称',
    provider_code VARCHAR(50) NOT NULL DEFAULT '' COMMENT 'embedding提供商编码（如bailian等）',
    provider_name VARCHAR(100) NOT NULL DEFAULT '' COMMENT 'embedding提供商名称',
    model_code VARCHAR(50) NOT NULL DEFAULT '' COMMENT 'embedding模型编码（如text-embedding-v4等）',
    model_name VARCHAR(100) NOT NULL DEFAULT '' COMMENT 'embedding模型名称',
    vector_dimension INT NOT NULL DEFAULT 1024 COMMENT 'embedding向量维度',
    deleted TINYINT NOT NULL DEFAULT 0 COMMENT '是否删除（1-删除，0-未删除）',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    unique key uk_scene_code (scene_code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='embedding场景映射表（关联场景与对应的embedding提供商及模型）';

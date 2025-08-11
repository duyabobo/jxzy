-- LLM服务完整初始化脚本
-- 注意：此文件已迁移到 bs/llm-rpc/internal/model/sql/ 目录下
-- 建议使用模块化的SQL文件进行维护
-- 
-- 快速初始化命令：
-- mysql -u root -p < bs/llm-rpc/internal/model/sql/init_all.sql
--
-- 详细文档请参考：bs/llm-rpc/internal/model/sql/README.md

-- 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS jxzy_llm DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE jxzy_llm;

-- 1. LLM供应商表
CREATE TABLE IF NOT EXISTS `llm_providers` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `provider_code` varchar(50) NOT NULL COMMENT '供应商代码：doubao, bailian, openai等',
  `provider_name` varchar(100) NOT NULL COMMENT '供应商名称',
  `api_key` varchar(500) NOT NULL COMMENT 'API密钥，加密存储',
  `base_url` varchar(200) NOT NULL COMMENT 'API基础URL',
  `headers` json DEFAULT NULL COMMENT '额外请求头，JSON格式',
  `timeout_seconds` int(11) DEFAULT 30 COMMENT '请求超时时间（秒）',
  `max_retries` int(11) DEFAULT 3 COMMENT '最大重试次数',
  `retry_delay_ms` int(11) DEFAULT 1000 COMMENT '重试延迟（毫秒）',
  `status` tinyint(4) DEFAULT 1 COMMENT '状态：0-禁用，1-启用',
  `description` text COMMENT '供应商描述',
  `config_extra` json DEFAULT NULL COMMENT '扩展配置，JSON格式',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_provider_code` (`provider_code`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='LLM供应商配置表';

-- 2. LLM模型表
CREATE TABLE IF NOT EXISTS `llm_models` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `model_code` varchar(100) NOT NULL COMMENT '模型代码',
  `model_name` varchar(100) NOT NULL COMMENT '模型名称',
  `provider_id` bigint(20) unsigned NOT NULL COMMENT '供应商ID',
  `endpoint_id` varchar(100) DEFAULT NULL COMMENT '端点ID（豆包等需要）',
  `max_tokens` int(11) DEFAULT 4096 COMMENT '最大token数',
  `support_stream` tinyint(4) DEFAULT 1 COMMENT '是否支持流式：0-否，1-是',
  `support_function_call` tinyint(4) DEFAULT 0 COMMENT '是否支持函数调用：0-否，1-是',
  `input_price_per_1k` decimal(10,6) DEFAULT 0.000000 COMMENT '输入价格（每千token）',
  `output_price_per_1k` decimal(10,6) DEFAULT 0.000000 COMMENT '输出价格（每千token）',
  `default_temperature` decimal(3,2) DEFAULT 0.70 COMMENT '默认温度值',
  `default_top_p` decimal(3,2) DEFAULT 0.90 COMMENT '默认TopP值',
  `description` text COMMENT '模型描述',
  `advantages` text COMMENT '模型优点',
  `disadvantages` text COMMENT '模型缺点',
  `recommended_scenes` json DEFAULT NULL COMMENT '推荐使用场景，JSON数组',
  `model_params` json DEFAULT NULL COMMENT '模型参数配置，JSON格式',
  `status` tinyint(4) DEFAULT 1 COMMENT '状态：0-禁用，1-启用',
  `sort_order` int(11) DEFAULT 0 COMMENT '排序序号',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_model_code` (`model_code`),
  KEY `idx_provider_id` (`provider_id`),
  KEY `idx_status` (`status`),
  KEY `idx_sort_order` (`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='LLM模型配置表';

-- 3. 场景配置表
CREATE TABLE IF NOT EXISTS `llm_scenes` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `scene_code` varchar(100) NOT NULL COMMENT '场景代码',
  `scene_name` varchar(100) NOT NULL COMMENT '场景名称',
  `description` text COMMENT '场景描述',
  `system_prompt` text COMMENT '系统提示词',
  `primary_model_id` bigint(20) unsigned NOT NULL COMMENT '主要模型ID',
  `backup_model_ids` json DEFAULT NULL COMMENT '备用模型ID列表，JSON数组',
  `default_temperature` decimal(3,2) DEFAULT 0.70 COMMENT '默认温度值',
  `default_top_p` decimal(3,2) DEFAULT 0.90 COMMENT '默认TopP值',
  `default_max_tokens` int(11) DEFAULT 2000 COMMENT '默认最大输出token',
  `default_presence_penalty` decimal(3,2) DEFAULT 0.00 COMMENT '默认存在惩罚',
  `default_frequency_penalty` decimal(3,2) DEFAULT 0.00 COMMENT '默认频率惩罚',
  `max_concurrency` int(11) DEFAULT 10 COMMENT '最大并发数',
  `rate_limit_rps` int(11) DEFAULT 50 COMMENT '每秒请求限制',
  `priority` int(11) DEFAULT 5 COMMENT '优先级（1-10）',
  `status` tinyint(4) DEFAULT 1 COMMENT '状态：0-禁用，1-启用',
  `sort_order` int(11) DEFAULT 0 COMMENT '排序序号',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_scene_code` (`scene_code`),
  KEY `idx_primary_model_id` (`primary_model_id`),
  KEY `idx_status` (`status`),
  KEY `idx_priority` (`priority`),
  KEY `idx_sort_order` (`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='LLM场景配置表';

-- 4. 系统配置表
CREATE TABLE IF NOT EXISTS `llm_system_config` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `config_key` varchar(100) NOT NULL COMMENT '配置键',
  `config_value` text NOT NULL COMMENT '配置值',
  `config_type` varchar(20) DEFAULT 'string' COMMENT '配置类型：string, int, float, json, bool',
  `description` varchar(500) DEFAULT NULL COMMENT '配置描述',
  `is_encrypted` tinyint(4) DEFAULT 0 COMMENT '是否加密：0-否，1-是',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_config_key` (`config_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='LLM系统配置表';

-- 插入初始数据

-- 插入豆包供应商
INSERT INTO `llm_providers` (`provider_code`, `provider_name`, `api_key`, `base_url`, `headers`, `timeout_seconds`, `max_retries`, `retry_delay_ms`, `status`, `description`) VALUES 
('doubao', '豆包大模型', 'your-doubao-api-key', 'https://ark.cn-beijing.volces.com/api/v3', '{"Content-Type": "application/json"}', 30, 3, 1000, 1, '字节跳动豆包大模型，支持中文对话'),
('bailian', '百炼大模型', 'your-bailian-api-key', 'https://dashscope.aliyuncs.com/api/v1', '{"Content-Type": "application/json"}', 30, 3, 1000, 0, '阿里云百炼大模型平台（预留）'),
('openai', 'OpenAI', 'sk-your-openai-key', 'https://api.openai.com/v1', '{"Content-Type": "application/json"}', 30, 3, 1000, 0, 'OpenAI官方API（预留）');

-- 插入豆包模型
INSERT INTO `llm_models` (`model_code`, `model_name`, `provider_id`, `endpoint_id`, `max_tokens`, `support_stream`, `support_function_call`, `input_price_per_1k`, `output_price_per_1k`, `default_temperature`, `default_top_p`, `description`, `advantages`, `disadvantages`, `recommended_scenes`, `status`) VALUES 
('doubao_lite', '豆包-lite', 1, 'ep-20241201-xxxx', 4096, 1, 0, 0.0008, 0.002, 0.70, 0.90, '豆包轻量级模型，适合日常对话和简单任务', '响应速度快，成本低廉，适合高频调用', '复杂推理能力相对较弱', '["general_chat", "customer_service"]', 1),
('doubao_pro', '豆包-pro', 1, 'ep-20241201-yyyy', 4096, 1, 1, 0.005, 0.015, 0.70, 0.90, '豆包高级模型，具备强大的推理和函数调用能力', '推理能力强，支持函数调用，适合复杂任务', '成本相对较高，响应时间较长', '["code_generation", "data_analysis", "technical_qa"]', 0);

-- 插入场景配置
INSERT INTO `llm_scenes` (`scene_code`, `scene_name`, `description`, `system_prompt`, `primary_model_id`, `backup_model_ids`, `default_temperature`, `default_top_p`, `default_max_tokens`, `max_concurrency`, `rate_limit_rps`, `priority`, `status`) VALUES 
('general_chat', '通用对话', '日常对话场景，适合一般性问答和闲聊', '你是一个友善、有帮助的AI助手。请以简洁、准确的方式回答用户问题。', 1, '[2]', 0.70, 0.90, 2000, 20, 50, 5, 1),
('customer_service', '客户服务', '客服场景，要求回答准确、专业、有礼貌', '你是一个专业的客服助手。请以礼貌、专业的态度回答客户问题，提供准确的信息和解决方案。', 1, '[2]', 0.30, 0.80, 1500, 30, 100, 8, 1),
('code_generation', '代码生成', '代码生成和编程辅助场景', '你是一个专业的编程助手。请提供高质量的代码，并详细解释代码逻辑。', 2, '[1]', 0.20, 0.80, 3000, 10, 20, 7, 0),
('data_analysis', '数据分析', '数据分析和报告生成场景', '你是一个数据分析专家。请提供专业的数据分析建议和详细的解释。', 2, '[1]', 0.40, 0.80, 2500, 5, 10, 6, 0),
('technical_qa', '技术问答', '技术问题解答场景，要求答案准确且深入', '你是一个技术专家。请提供准确、深入的技术解答，并给出实用的建议。', 2, '[1]', 0.30, 0.80, 2000, 15, 30, 7, 0);

-- 插入系统配置
INSERT INTO `llm_system_config` (`config_key`, `config_value`, `config_type`, `description`) VALUES 
('default_provider', 'doubao', 'string', '默认LLM供应商'),
('default_model', 'doubao_lite', 'string', '默认LLM模型'),
('default_scene', 'general_chat', 'string', '默认使用场景'),
('global_max_tokens', '4000', 'int', '全局最大token数'),
('global_timeout_seconds', '30', 'int', '全局请求超时时间'),
('enable_rate_limit', 'true', 'bool', '是否启用限流'),
('enable_content_filter', 'true', 'bool', '是否启用内容过滤'),
('cache_ttl_seconds', '300', 'int', '缓存TTL时间（秒）');

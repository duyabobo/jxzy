package logic

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"path"
	"strings"

	knowledgepb "jxzy/bll/bll_knowledge/bll_knowledge"
	"jxzy/bll/bll_knowledge/internal/model"
	"jxzy/bll/bll_knowledge/internal/svc"
	bsllm "jxzy/bs/bs_llm/bs_llm"
	"jxzy/bs/bs_rag/bs_rag"
	consts "jxzy/common/const"
	"jxzy/common/logger"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddVectorKnowledgeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddVectorKnowledgeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddVectorKnowledgeLogic {
	// 使用自定义的 ServiceLogger，在日志中显示服务名
	serviceLogger := logger.NewServiceLogger("bll-knowledge").WithContext(ctx)

	return &AddVectorKnowledgeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: serviceLogger,
	}
}

// AddVectorKnowledge 添加知识库到向量数据库
func (l *AddVectorKnowledgeLogic) AddVectorKnowledge(in *knowledgepb.AddVectorKnowledgeRequest) (*knowledgepb.AddVectorKnowledgeResponse, error) {
	l.Logger.Infof("AddVectorKnowledge called with source_type: %v, user_id: %s",
		in.SourceType, in.UserId)

	// 1. 验证输入参数
	if err := l.validateInput(in); err != nil {
		l.Logger.Errorf("Input validation failed: %v", err)
		return &knowledgepb.AddVectorKnowledgeResponse{
			Success: false,
			Message: fmt.Sprintf("输入参数验证失败: %v", err),
		}, nil
	}

	// 2. 根据source_type处理不同的模式
	var fileMd5 string
	var documents []*bs_rag.VectorDocument
	var err error

	switch in.SourceType {
	case knowledgepb.KnowledgeSourceType_SOURCE_TYPE_FILE_URL:
		_, fileMd5, documents, err = l.processFileUrlMode(in)
	case knowledgepb.KnowledgeSourceType_SOURCE_TYPE_SEGMENTS:
		_, fileMd5, documents, err = l.processSegmentsMode(in)
	case knowledgepb.KnowledgeSourceType_SOURCE_TYPE_SUMMARIES:
		_, fileMd5, documents, err = l.processSummariesMode(in)
	default:
		return &knowledgepb.AddVectorKnowledgeResponse{
			Success: false,
			Message: fmt.Sprintf("不支持的source_type: %v", in.SourceType),
		}, nil
	}

	if err != nil {
		// 处理文件已存在的情况
		if strings.HasPrefix(err.Error(), "FILE_EXISTS:") {
			fileMd5 := strings.TrimPrefix(err.Error(), "FILE_EXISTS:")
			return &knowledgepb.AddVectorKnowledgeResponse{
				VectorId: fileMd5,
				Success:  true,
				Message:  "文件已存在，跳过处理",
			}, nil
		}
		return &knowledgepb.AddVectorKnowledgeResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	// 3. 插入向量到RAG服务
	return l.insertVectorsToRAG(fileMd5, documents, in.UserId)
}

// processFileUrlMode 处理文件URL模式
func (l *AddVectorKnowledgeLogic) processFileUrlMode(in *knowledgepb.AddVectorKnowledgeRequest) (int64, string, []*bs_rag.VectorDocument, error) {
	// 下载文件内容
	fileBytes, fileName, fileType, fileSize, err := l.downloadFile(in.FileUrl)
	if err != nil {
		l.Logger.Errorf("Failed to download file: %v", err)
		return 0, "", nil, fmt.Errorf("下载文件失败: %v", err)
	}

	// 计算文件MD5并进行去重
	fileMd5 := l.md5Bytes(fileBytes)
	l.Logger.Infof("Computed file md5: %s (name=%s, size=%d, type=%s)", fileMd5, fileName, fileSize, fileType)

	// 查询是否已存在
	existing, _ := l.svcCtx.KnowledgeFileModel.FindOneByMd5(l.ctx, fileMd5)
	if existing != nil {
		l.Logger.Infof("File already processed, knowledge_file.id=%d", existing.Id)
		// 返回特殊错误，让调用者知道文件已存在
		return 0, "", nil, fmt.Errorf("FILE_EXISTS:%s", fileMd5)
	}

	// 入库 knowledge_file
	kf := &model.KnowledgeFile{
		OssPath:  in.FileUrl,
		FileName: fileName,
		FileSize: fileSize,
		FileType: fileType,
		FileMd5:  fileMd5,
	}
	res, err := l.svcCtx.KnowledgeFileModel.Insert(l.ctx, kf)
	if err != nil {
		l.Logger.Errorf("Insert knowledge_file failed: %v", err)
		return 0, "", nil, fmt.Errorf("保存文件记录失败: %v", err)
	}
	fileId, _ := res.LastInsertId()

	// 利用LLM进行语义段拆分
	segments, err := l.segmentTextWithLLM(string(fileBytes), in.UserId)
	if err != nil {
		l.Logger.Errorf("LLM segmentation failed: %v", err)
		return 0, "", nil, fmt.Errorf("LLM拆分失败: %v", err)
	}

	// 处理segments并生成documents
	documents := l.processSegmentsToDocuments(segments, fileId, in.UserId)
	return fileId, fileMd5, documents, nil
}

// processSegmentsMode 处理segments模式
func (l *AddVectorKnowledgeLogic) processSegmentsMode(in *knowledgepb.AddVectorKnowledgeRequest) (int64, string, []*bs_rag.VectorDocument, error) {
	// 过滤空字符串
	segments := l.filterEmptyStrings(in.Segments)
	if len(segments) == 0 {
		return 0, "", nil, fmt.Errorf("segments列表不能为空")
	}

	// fileId使用0标识，fileMd5使用segments的MD5组合
	fileId := int64(0)
	fileMd5 := l.md5String(strings.Join(segments, "\n"))
	l.Logger.Infof("Using segments mode, segments count: %d, fileMd5: %s", len(segments), fileMd5)

	// 处理segments并生成documents
	documents := l.processSegmentsToDocuments(segments, fileId, in.UserId)
	return fileId, fileMd5, documents, nil
}

// processSummariesMode 处理summaries模式
func (l *AddVectorKnowledgeLogic) processSummariesMode(in *knowledgepb.AddVectorKnowledgeRequest) (int64, string, []*bs_rag.VectorDocument, error) {
	// 过滤空字符串
	summaries := l.filterEmptyStrings(in.Summaries)
	if len(summaries) == 0 {
		return 0, "", nil, fmt.Errorf("summaries列表不能为空")
	}

	// fileId和segId都使用0标识，fileMd5使用summaries的MD5组合
	fileId := int64(0)
	segId := int64(0)
	fileMd5 := l.md5String(strings.Join(summaries, "\n"))
	l.Logger.Infof("Using summaries mode, summaries count: %d, fileMd5: %s", len(summaries), fileMd5)

	// 直接处理摘要句列表：存储并构建RAG文档
	documents := l.processSummariesToDocuments(summaries, fileId, segId, in.UserId)
	return fileId, fileMd5, documents, nil
}

// processSegmentsToDocuments 处理segments并生成documents
func (l *AddVectorKnowledgeLogic) processSegmentsToDocuments(segments []string, fileId int64, userId string) []*bs_rag.VectorDocument {
	var documents []*bs_rag.VectorDocument

	for _, seg := range segments {
		// 存储语义段到数据库
		segMd5 := l.md5String(seg)
		segRec := &model.KnowledgeSegment{
			KnowledgeFileId: fileId,
			SegmentText:     seg,
			SegmentMd5:      segMd5,
		}
		segRes, err := l.svcCtx.KnowledgeSegmentModel.Insert(l.ctx, segRec)
		if err != nil {
			l.Logger.Errorf("Insert knowledge_segment failed: %v", err)
			continue
		}
		segId, _ := segRes.LastInsertId()

		// 为语义段生成多个维度的摘要句
		summaries, err := l.summarizeSegmentWithLLM(seg, userId)
		if err != nil || len(summaries) == 0 {
			l.Logger.Errorf("LLM summary failed for segment %d: %v", segId, err)
			continue
		}

		// 对每个摘要句进行处理：存储并构建RAG文档
		for _, summary := range summaries {
			doc := l.insertSummaryToDocument(summary, fileId, segId, userId)
			if doc != nil {
				documents = append(documents, doc)
			}
		}
	}

	return documents
}

// processSummariesToDocuments 处理summaries并生成documents
func (l *AddVectorKnowledgeLogic) processSummariesToDocuments(summaries []string, fileId, segId int64, userId string) []*bs_rag.VectorDocument {
	var documents []*bs_rag.VectorDocument

	for _, summary := range summaries {
		doc := l.insertSummaryToDocument(summary, fileId, segId, userId)
		if doc != nil {
			documents = append(documents, doc)
		}
	}

	return documents
}

// insertSummaryToDocument 存储摘要句到数据库并构建RAG文档
func (l *AddVectorKnowledgeLogic) insertSummaryToDocument(summary string, fileId, segId int64, userId string) *bs_rag.VectorDocument {
	summary = strings.TrimSpace(summary)
	if summary == "" {
		return nil
	}

	// 存储摘要句到数据库
	sumMd5 := l.md5String(summary)
	sumRec := &model.KnowledgeSummarySentence{
		KnowledgeFileId:     fileId,
		KnowledgeSegmentId:  segId,
		SummarySentenceText: summary,
		SummarySentenceMd5:  sumMd5,
	}
	sumRes, err := l.svcCtx.KnowledgeSummarySentenceModel.Insert(l.ctx, sumRec)
	if err != nil {
		l.Logger.Errorf("Insert knowledge_summary_sentence failed: %v", err)
		return nil
	}
	summaryId, _ := sumRes.LastInsertId()

	// 构建RAG文档
	docId := fmt.Sprintf("%d", summaryId)
	return &bs_rag.VectorDocument{
		Id:   docId,
		Text: summary,
		Metadata: map[string]string{
			"knowledge_file_id":    fmt.Sprintf("%d", fileId),
			"knowledge_segment_id": fmt.Sprintf("%d", segId),
			"user_id":              userId,
		},
		Content: summary,
	}
}

// insertVectorsToRAG 插入向量到RAG服务
func (l *AddVectorKnowledgeLogic) insertVectorsToRAG(fileMd5 string, documents []*bs_rag.VectorDocument, userId string) (*knowledgepb.AddVectorKnowledgeResponse, error) {
	if len(documents) == 0 {
		return &knowledgepb.AddVectorKnowledgeResponse{
			VectorId: fileMd5,
			Success:  false,
			Message:  "没有可插入的摘要文档",
		}, nil
	}

	if l.svcCtx.RagRpc == nil {
		l.Logger.Error("RAG service is not available")
		return &knowledgepb.AddVectorKnowledgeResponse{
			VectorId: fileMd5,
			Success:  false,
			Message:  "RAG服务不可用",
		}, nil
	}

	ragReq := &bs_rag.VectorInsertRequest{
		CollectionName: consts.DefaultCollectionName,
		Documents:      documents,
		UserId:         userId,
	}
	ragResp, err := l.svcCtx.RagRpc.VectorInsert(l.ctx, ragReq)
	if err != nil {
		l.Logger.Errorf("Failed to insert vectors to RAG service: %v", err)
		return &knowledgepb.AddVectorKnowledgeResponse{
			VectorId: fileMd5,
			Success:  false,
			Message:  fmt.Sprintf("插入RAG失败: %v", err),
		}, nil
	}

	l.Logger.Infof("Inserted %d documents to RAG", ragResp.InsertedCount)
	return &knowledgepb.AddVectorKnowledgeResponse{
		VectorId: fileMd5,
		Success:  true,
		Message:  "知识库添加成功",
	}, nil
}

// filterEmptyStrings 过滤空字符串
func (l *AddVectorKnowledgeLogic) filterEmptyStrings(strs []string) []string {
	filtered := make([]string, 0, len(strs))
	for _, s := range strs {
		s = strings.TrimSpace(s)
		if s != "" {
			filtered = append(filtered, s)
		}
	}
	return filtered
}

// validateInput 验证输入参数
func (l *AddVectorKnowledgeLogic) validateInput(in *knowledgepb.AddVectorKnowledgeRequest) error {
	if strings.TrimSpace(in.UserId) == "" {
		return fmt.Errorf("user_id不能为空")
	}
	if in.SourceType == knowledgepb.KnowledgeSourceType_SOURCE_TYPE_FILE_URL {
		if strings.TrimSpace(in.FileUrl) == "" {
			return fmt.Errorf("file_url不能为空")
		}
	} else if in.SourceType == knowledgepb.KnowledgeSourceType_SOURCE_TYPE_SEGMENTS {
		if len(in.Segments) == 0 {
			return fmt.Errorf("segments列表不能为空")
		}
	} else if in.SourceType == knowledgepb.KnowledgeSourceType_SOURCE_TYPE_SUMMARIES {
		if len(in.Summaries) == 0 {
			return fmt.Errorf("summaries列表不能为空")
		}
	} else {
		return fmt.Errorf("source_type必须指定为FILE_URL、SEGMENTS或SUMMARIES")
	}
	return nil
}

// md5String 计算字符串MD5
func (l *AddVectorKnowledgeLogic) md5String(s string) string {
	sum := md5.Sum([]byte(s))
	return hex.EncodeToString(sum[:])
}

// md5Bytes 计算字节MD5
func (l *AddVectorKnowledgeLogic) md5Bytes(b []byte) string {
	sum := md5.Sum(b)
	return hex.EncodeToString(sum[:])
}

// downloadFile 下载远端文件，返回内容及基础元信息
func (l *AddVectorKnowledgeLogic) downloadFile(rawUrl string) ([]byte, string, string, int64, error) {
	req, err := http.NewRequestWithContext(l.ctx, http.MethodGet, rawUrl, nil)
	if err != nil {
		return nil, "", "", 0, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, "", "", 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, "", "", 0, fmt.Errorf("http status: %d", resp.StatusCode)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", "", 0, err
	}
	// 解析文件名
	parsed, _ := url.Parse(rawUrl)
	base := path.Base(parsed.Path)
	// 猜测类型
	ctype := resp.Header.Get("Content-Type")
	if ctype == "" {
		ctype = mime.TypeByExtension(path.Ext(base))
	}
	size := int64(len(data))
	return data, base, ctype, size, nil
}

// segmentTextWithLLM 使用LLM将文本拆分为语义段
func (l *AddVectorKnowledgeLogic) segmentTextWithLLM(text string, userId string) ([]string, error) {
	if l.svcCtx.LlmRpc == nil {
		return []string{strings.TrimSpace(text)}, nil
	}
	req := &bsllm.LLMRequest{
		SceneCode: "knowledge_segmentation",
		UserId:    userId,
		Messages: []*bsllm.ChatMessage{
			{Role: "system", Content: "请将以下内容按语义进行合理拆分，返回为每段一行。不要编号。"},
			{Role: "user", Content: text},
		},
	}
	resp, err := l.svcCtx.LlmRpc.LLM(l.ctx, req)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(resp.Completion, "\n")
	cleaned := make([]string, 0, len(lines))
	for _, ln := range lines {
		ln = strings.TrimSpace(ln)
		if ln != "" {
			cleaned = append(cleaned, ln)
		}
	}
	if len(cleaned) == 0 {
		cleaned = []string{strings.TrimSpace(text)}
	}
	return cleaned, nil
}

// summarizeSegmentWithLLM 使用LLM为语义段生成多个维度的摘要句
func (l *AddVectorKnowledgeLogic) summarizeSegmentWithLLM(segment string, userId string) ([]string, error) {
	if l.svcCtx.LlmRpc == nil {
		return []string{segment}, nil
	}
	req := &bsllm.LLMRequest{
		SceneCode: "knowledge_segment_summary",
		UserId:    userId,
		Messages: []*bsllm.ChatMessage{
			{Role: "system", Content: "请为以下文本从多个维度生成一句话摘要，每个摘要句简洁、完整且可用于检索。每个摘要句占一行，直接返回摘要句，不要编号。"},
			{Role: "user", Content: segment},
		},
	}
	resp, err := l.svcCtx.LlmRpc.LLM(l.ctx, req)
	if err != nil {
		return nil, err
	}
	// 按行分割，去除空行
	lines := strings.Split(resp.Completion, "\n")
	summaries := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			summaries = append(summaries, line)
		}
	}
	// 如果没有生成任何摘要，返回原始段落作为默认摘要
	if len(summaries) == 0 {
		summaries = []string{strings.TrimSpace(segment)}
	}
	return summaries, nil
}

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"jxzy/bs/bs_rag/bs_rag"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 连接到 gRPC 服务器
	conn, err := grpc.Dial("localhost:8082", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// 创建客户端
	client := bs_rag.NewBsRagServiceClient(conn)

	// 设置超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	fmt.Println("=== BS RAG Service Test Client ===")

	// 测试获取集合信息
	fmt.Println("\n1. Testing GetCollectionInfo...")
	infoReq := &bs_rag.CollectionInfoRequest{
		CollectionName: "test",
		UserId:         "test_user",
	}

	infoResp, err := client.GetCollectionInfo(ctx, infoReq)
	if err != nil {
		log.Printf("GetCollectionInfo failed: %v", err)
	} else {
		fmt.Printf("Collection info: Name=%s, Documents=%d, Dimension=%d, Exists=%v\n",
			infoResp.CollectionName, infoResp.DocumentCount, infoResp.VectorDimension, infoResp.Exists)
	}

	// 测试向量插入
	fmt.Println("\n2. Testing VectorInsert...")
	insertReq := &bs_rag.VectorInsertRequest{
		CollectionName: "test",
		Documents: []*bs_rag.VectorDocument{
			{
				Id:       "doc1",
				Text:     "测试文档内容",
				Metadata: map[string]string{"source": "test", "type": "document"},
				Content:  "This is a test document for vector insertion",
			},
		},
		UserId: "test_user",
	}

	insertResp, err := client.VectorInsert(ctx, insertReq)
	if err != nil {
		log.Printf("VectorInsert failed: %v", err)
	} else {
		fmt.Printf("Inserted %d documents: %v\n", insertResp.InsertedCount, insertResp.InsertedIds)
	}

	// 测试向量搜索
	fmt.Println("\n3. Testing VectorSearch...")
	searchReq := &bs_rag.VectorSearchRequest{
		QueryText:      "测试查询文本",
		TopK:           5,
		MinScore:       0.5,
		CollectionName: "test",
		UserId:         "test_user",
	}

	searchResp, err := client.VectorSearch(ctx, searchReq)
	if err != nil {
		log.Printf("VectorSearch failed: %v", err)
	} else {
		fmt.Printf("Found %d results\n", searchResp.TotalCount)
		for i, result := range searchResp.Results {
			fmt.Printf("  Result %d: ID=%s, Score=%.3f\n", i+1, result.Id, result.Score)
		}
	}

	// 测试向量删除
	fmt.Println("\n4. Testing VectorDelete...")
	deleteReq := &bs_rag.VectorDeleteRequest{
		CollectionName: "test",
		DocumentIds:    []string{"doc1"},
		UserId:         "test_user",
	}

	deleteResp, err := client.VectorDelete(ctx, deleteReq)
	if err != nil {
		log.Printf("VectorDelete failed: %v", err)
	} else {
		fmt.Printf("Deleted %d documents: %v\n", deleteResp.DeletedCount, deleteResp.DeletedIds)
	}

	fmt.Println("\n✅ All tests completed!")
}

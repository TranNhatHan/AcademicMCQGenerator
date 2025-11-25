package main

import (
	"context"
	"fmt"
	"os"

	"github.com/tmc/langchaingo/documentloaders"
	embedding_huggingface "github.com/tmc/langchaingo/embeddings/huggingface"
	llm_huggingface "github.com/tmc/langchaingo/llms/huggingface"
	"github.com/tmc/langchaingo/textsplitter"
	"github.com/tmc/langchaingo/vectorstores/weaviate"
)

func main() {
	ctx := context.Background()
	file, err := os.Open("material.pdf")
	if err != nil {
		panic(err)
	}

	stat, err := file.Stat()
	if err != nil {
		panic(err)
	}
	size := stat.Size()

	defer file.Close()

	pdfLoader := documentloaders.NewPDF(file, size)

	splitter := textsplitter.NewRecursiveCharacter(
		textsplitter.WithChunkSize(256),
		textsplitter.WithChunkOverlap(32),
	)

	splittedContent, err := pdfLoader.LoadAndSplit(ctx, splitter)
	if err != nil {
		panic(err)
	}
	for i := 5000; i < len(splittedContent); i++ {
		fmt.Printf("---- Chunk %d ----\n", i+1)
		fmt.Println(splittedContent[i].Metadata)
		fmt.Println(splittedContent[i].PageContent)
		fmt.Println(splittedContent[i].Score)
	}

	embeddingClient, err := llm_huggingface.New(
		llm_huggingface.WithURL(os.Getenv("EMBEDDING_URL")),
	)
	if err != nil {
		panic(err)
	}

	embedder, err := embedding_huggingface.NewHuggingface(
		embedding_huggingface.WithClient(*embeddingClient),
	)
	if err != nil {
		panic(err)
	}

	weaviate.New(
		weaviate.WithScheme("http"),
		weaviate.WithHost("localhost:8080"),
		weaviate.WithEmbedder(embedder),
	)

}

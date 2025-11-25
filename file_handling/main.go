package main

import (
	"context"
	"fmt"
	"os"

	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/textsplitter"
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
		textsplitter.WithChunkSize(512),
		textsplitter.WithChunkOverlap(64),
	)

	splittedContent, err := pdfLoader.LoadAndSplit(ctx, splitter)
	if err != nil {
		panic(err)
	}
	for i := 5000; i < len(splittedContent); i++ {
		fmt.Printf("---- Chunk %d ----\n", i+1)
		fmt.Println(splittedContent[i].PageContent)
	}
}

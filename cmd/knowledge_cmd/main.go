package main

import (
	"DocReader/internal/agent/knowledge_index_pipeline"
	"DocReader/internal/component/loader"
	"DocReader/pkg/client"
	"DocReader/pkg/common"
	"DocReader/pkg/log_call_back"
	"context"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/cloudwego/eino/components/document"
	"github.com/cloudwego/eino/compose"
)

func main() {
	ctx := context.Background()

	// Build indexing pipeline
	r, err := knowledge_index_pipeline.BuildKnowledgeIndexing(ctx)
	if err != nil {
		panic(err)
	}

	// Walk through docs directory
	docsDir := "./docs"
	fmt.Printf("Starting to index documents from: %s\n", docsDir)

	err = filepath.WalkDir(docsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walk dir failed: %w", err)
		}
		if d.IsDir() {
			return nil
		}

		// Only process markdown files
		if !strings.HasSuffix(path, ".md") {
			fmt.Printf("[skip] not a markdown file: %s\n", path)
			return nil
		}

		fmt.Printf("\n[start] indexing file: %s\n", path)

		// Load document to get metadata
		ldr, err := loader.NewFileLoader(ctx)
		if err != nil {
			return err
		}
		docs, err := ldr.Load(ctx, document.Source{URI: path})
		if err != nil {
			return err
		}

		// Delete existing documents with same source
		cli, err := client.NewMilvusClient(ctx)
		if err != nil {
			return err
		}

		expr := fmt.Sprintf(`metadata["_source"] == "%s"`, docs[0].MetaData["_source"])
		queryResult, err := cli.Query(ctx, common.MilvusCollectionName, []string{}, expr, []string{"id"})
		if err != nil {
			fmt.Printf("[warn] query existing data failed: %v\n", err)
		} else if len(queryResult) > 0 {
			var idsToDelete []string
			for _, column := range queryResult {
				if column.Name() == "id" {
					for i := 0; i < column.Len(); i++ {
						id, err := column.GetAsString(i)
						if err == nil {
							idsToDelete = append(idsToDelete, id)
						}
					}
				}
			}
			if len(idsToDelete) > 0 {
				deleteExpr := fmt.Sprintf(`id in ["%s"]`, strings.Join(idsToDelete, `","`))
				err = cli.Delete(ctx, common.MilvusCollectionName, "", deleteExpr)
				if err != nil {
					fmt.Printf("[warn] delete existing data failed: %v\n", err)
				} else {
					fmt.Printf("[info] deleted %d existing records\n", len(idsToDelete))
				}
			}
		}

		// Index the document
		ids, err := r.Invoke(ctx, document.Source{URI: path}, compose.WithCallbacks(log_call_back.LogCallback(nil)))
		if err != nil {
			return fmt.Errorf("invoke index pipeline failed: %w", err)
		}

		fmt.Printf("[done] indexed file: %s, chunks: %d\n", path, len(ids))
		return nil
	})

	if err != nil {
		fmt.Printf("Error during indexing: %v\n", err)
	} else {
		fmt.Println("\n========== Indexing Completed ==========")
	}
}

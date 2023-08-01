package fs_pkg

import (
	"context"
	"fmt"
	"log"

	"firestore_tool/model"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

func Update(ctx context.Context, client *firestore.Client, docPath string, updateStruct []firestore.Update, whereQuery []model.WhereQuery) {
	q := client.Collection(docPath).Query
	for _, v := range whereQuery {
		q = q.Where(v.Path, v.Op, v.Value)
	}
	iter := q.Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
			fmt.Println("Failed Ref: ", doc.Ref)
		}

		_, err = doc.Ref.Update(ctx, updateStruct)
		if err != nil {
			log.Fatalf("Failed to update: %v", err)
			fmt.Println("Failed Ref: ", doc.Ref)
		}
	}
}

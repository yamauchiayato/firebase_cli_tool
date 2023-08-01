package fs_pkg

import (
	"context"
	"errors"

	"cloud.google.com/go/firestore"
)

func Get(ctx context.Context, client *firestore.Client) ([]interface{}, error) {
	// iter := client.Collection("CUSTOMERS/20220124091346358922498/TOWERMASTER").Where("towerName", "==", "3号機").Limit(1).Documents(ctx)
	iter := client.Collection("CUSTOMERS/20220124091346358922498/TOWERMASTER").Documents(ctx)
	docs, err := iter.GetAll() //イテレータを使う必要がなくなり、コードの簡素化
	if err != nil {
		return nil, err
	}

	if len(docs) == 0 {
		return nil, errors.New("Not found in firestore")
	}

	var m []interface{}
	for _, doc := range docs {
		m = append(m, doc.Data())
	}

	return m, nil
}

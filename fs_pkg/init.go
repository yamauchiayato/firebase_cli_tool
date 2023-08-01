package fs_pkg

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

// firebaseのサービスアカウントキーを取得、クライアントの初期化
func FirebaseInit(ctx context.Context, jsonPath string) (*firestore.Client, error) {
	sa := option.WithCredentialsFile(jsonPath)

	// Firebaseアプリケーションの初期化
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		fmt.Printf("Firebaseアプリケーションの初期化に失敗しました", err)
		return nil, err
	}

	// Firestoreクライアントの初期化
	fsClient, err := app.Firestore(ctx)
	if err != nil {
		fmt.Printf("Firestoreクライアントの初期化に失敗しました", err)
		return nil, err
	}

	return fsClient, nil
}

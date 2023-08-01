package main

import (
	"bufio"
	"context"
	"firestore_tool/fs_pkg"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"firestore_tool/model"

	"cloud.google.com/go/firestore"
)

type Setting struct {
	credentialJsonPath string
	mode               string
	docPath            string
	updateStruct       []firestore.Update
	whereQuery         []model.WhereQuery
}

var setting Setting

func main() {
	flag.StringVar(&setting.credentialJsonPath, "credentialJsonPath", "", "エントリーデフォルトパス")
	flag.Parse()

	ctx := context.Background()
	client, err := fs_pkg.FirebaseInit(ctx, setting.credentialJsonPath)
	if err != nil {
		log.Fatalf("FirebaseInit error: %v", err)
	}
	defer client.Close()

	intro()

	// create a channel to indicat e when the program can quit
	doneChan := make(chan bool)

	// start a gorouutin to read user input and run program
	go readUserInput(os.Stdin, doneChan, ctx, client)

	// block until the done chan gets a value
	<-doneChan

	// close the channel
	close(doneChan)

	printSetting()
	fmt.Println("done!")

}

func intro() {
	fmt.Println("Firebaseと接続しました。")
	fmt.Println("メソッドを入力してください。updateまたはdelete")
	prompt()
}

func prompt() {
	fmt.Print("-> ")
}

func readUserInput(in io.Reader, doneChan chan bool, ctx context.Context, client *firestore.Client) {
	scanner := bufio.NewScanner(in)

	for {
		scanner.Scan()
		mode := scanner.Text()
		// invertedIndexからscanner.Text()に対応するデータ取得
		if mode == "exit" || mode == "q" {
			doneChan <- true
			return
		}
		if mode != "delete" && mode != "update" {
			fmt.Println("Invalid mode entered, exiting.")
			doneChan <- true
			return
		}
		setting.mode = mode

		fmt.Println("対象のドキュメントパスを入力してください。")
		prompt()
		scanner.Scan()
		docPath := scanner.Text()
		setting.docPath = docPath

		fmt.Println("Where句を設定しますか？(y/n)")
		prompt()
		scanner.Scan()
		if scanner.Text() == "y" {
			addWhereQuery(scanner)
		}

		if mode == "update" {
			addUpdateStruct(scanner)
		}

		// ユーザーからの入力を取得した後、データを削除または更新します。
		if setting.mode == "delete" {
			// データを削除する関数を呼び出します。
			fs_pkg.Delete(ctx, client, setting.docPath, setting.whereQuery)
		} else if setting.mode == "update" {
			// データを更新する関数を呼び出します。
			fs_pkg.Update(ctx, client, setting.docPath, setting.updateStruct, setting.whereQuery)
		}

		doneChan <- true
		return
	}
}

func addUpdateStruct(scanner *bufio.Scanner) {
	fmt.Println("保存対象のキーと値を入力してください。扱える型は、string,float64,boolean。形式は以下の通りです。")
	fmt.Println("key:type:value 例) shouhinName:string:商品A")
	fmt.Println("複数の場合は,で区切ります。")
	prompt()
	scanner.Scan()
	str := scanner.Text()

	// strをカンマ区切りで分割し、スライスに格納
	dataSlice := splitStr(str, ",")
	for _, v := range dataSlice {
		// dataSliceの要素をコロン区切りで分割し、スライスに格納
		keyAndValue := splitStr(v, ":")
		if len(keyAndValue) != 3 {
			log.Fatalf("Invalid keyAndValue entered, exiting.")
			return
		}
		value, err := convertToTargetType(keyAndValue[2], keyAndValue[1])
		if err != nil {
			log.Fatal(err)
			return
		}

		setting.updateStruct = append(setting.updateStruct, firestore.Update{Path: keyAndValue[0], Value: value})
	}
}

func addWhereQuery(scanner *bufio.Scanner) {
	fmt.Println("Where句を入力してください。扱える型は、string,float64,boolean。形式は以下の通りです。")
	fmt.Println("Path:Op:type:value 例) shouhinName:==:string:商品A")
	fmt.Println("複数の場合は,で区切ります。")
	prompt()
	scanner.Scan()
	str := scanner.Text()

	// strをカンマ区切りで分割し、スライスに格納
	whereQuerySlice := splitStr(str, ",")
	for _, v := range whereQuerySlice {
		// whereQuerySliceの要素をコロン区切りで分割し、スライスに格納
		whereQueryElementSlice := splitStr(v, ":")
		if len(whereQueryElementSlice) != 4 {
			log.Fatalf("Invalid where query entered, exiting.")
			return
		}
		value, err := convertToTargetType(whereQueryElementSlice[3], whereQueryElementSlice[2])
		if err != nil {
			log.Fatal(err)
			return
		}

		setting.whereQuery = append(setting.whereQuery, model.WhereQuery{Path: whereQueryElementSlice[0], Op: whereQueryElementSlice[1], Value: value})
	}
}

func splitStr(str string, delimiter string) []string {
	return strings.Split(str, delimiter)
}

func convertToTargetType(value interface{}, targetType string) (interface{}, error) {
	var err error
	if targetType == "string" {
		return value, nil
	} else if targetType == "float64" {
		value, err = strconv.ParseFloat(value.(string), 64)
		if err != nil {
			return nil, fmt.Errorf("Invalid float64 value entered: %v", err.Error())
		}
		return value, nil
	} else if targetType == "boolean" {
		return value == "true", nil
	} else {
		return nil, fmt.Errorf("Invalid targetType: %s", targetType)
	}
}

// settingの内容を出力する関数
func printSetting() {
	fmt.Println("メソッド: ", setting.mode)
	fmt.Println("ドキュメントパス: ", setting.docPath)
	fmt.Println("クエリ: ", setting.whereQuery)
	fmt.Println("更新データ: ", setting.updateStruct)
}

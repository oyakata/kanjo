package main

import (
	"fmt"
	"log"
	"net/http"
)

func init() {
	// URLパターンに正規表現は渡せない。
	http.HandleFunc("/count", WordCountHandler)

	// 順番に注意。"/"を先頭に指定すると他のPathがマッチしない。
	http.HandleFunc("/", TopPageHandler)
}

func main() {
	port := 8080
	log.Printf("start server, port %v", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

func TopPageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`
	<html>
		<head>
			<meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
			<title>Welcome to kanjo.</title>
		</head>

		<body>
			<h1>Hello, world.</h1>

			文字を入力してください。
			<form action="/count" method="GET">
			<input type="text" name="text" size="32">
			<input type="submit">
			</form>
		</body>
	</html>
	`))
}

func WordCountHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	// 1. テンプレートエンジンを使って処理するよう変更
	// 2. JSONで返却する機能を追加
	// 3. ファイルを読み取って文字数を数える機能を追加
	// 4. 文字数とバイト数を数えるよう変更
	// 5. ユニットテストを追加
	// 6. 最初のプログラムだと文字数ではなくバイト数を返してしまうので直す
	// 7. logをファイルに出力するよう変更

	text := r.FormValue("text")
	count := len(text)

	css := ""
	if count == 0 {
		css = "err"
	}

	content := fmt.Sprintf(`
	<html>
		<head>
			<meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
			<title>文字数カウント結果</title>
			<style type="text/css">
				.err { border: solid 1px red; }
			</style>
		</head>

		<body>
			<h1>文字数カウント結果</h1>

			入力文字: %v<br>
			文字数は: %v<br>

			でした。<br><br>

			文字を入力してください。
			<form action="/count" method="GET">
			<input type="text" name="text" size="32" class="%v">
			<input type="submit">
			</form>
		</body>
	</html>`, text, count, css)

	w.Write([]byte(content))
}

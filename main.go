package main

import (
	"encoding/json"
	"fmt"
	"github.com/oyakata/kanjo/lib/wc"
	htmlTemplate "html/template"
	"log"
	"net/http"
)

func init() {
	// URLパターンに正規表現は渡せない。
	http.HandleFunc("/count/file", FileWordCountHandler)
	http.HandleFunc("/count", WordCountHandler)

	// 順番に注意。"/"を先頭に指定すると他のPathがマッチしない。
	http.HandleFunc("/", TopPageHandler)
}

func main() {
	port := 8080
	log.Printf("start server, port %v", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

// 1. [済]テンプレートエンジンを使って処理するよう変更
// 2. [済]JSONで返却する機能を追加
// 3. [済]ファイルを読み取って文字数を数える機能を追加
// 4. [済]文字数とバイト数を数えるよう変更
// 5. ユニットテストを追加
// 6. [済]最初のプログラムだと文字数ではなくバイト数を返してしまうので直す
// 7. logをファイルに出力するよう変更
// 8. [済]HTMLのエスケープがないので直す

type Context map[string]interface{}

func FileWordCountHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "POSTでアクセスしてください", http.StatusMethodNotAllowed)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusBadRequest)
	}

	defer func() {
		file.Close()
		r.MultipartForm.RemoveAll()
	}()

	count, bc, invalid := wc.WordCountInFile(file)

	tmpl, _ := htmlTemplate.New("file_wc").Parse(`
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

			文字数は: {{.count}}<br>
			バイト数は: {{.bc}}<br>
			不正なバイト数は: {{.invalid}}<br>

			でした。<br><br>

			文字を入力してください。
			<form action="/count" method="GET">
			<input type="text" name="text" size="32">
			<input type="submit">
			</form>

			ファイルを調べたい場合はこちら。
			<form action="/count/file" method="POST" enctype="multipart/form-data">
			<input type="file" name="file">
			<input type="submit">
			</form>

		</body>
	</html>`)

	data := Context{
		"count":   count,
		"bc":      bc,
		"invalid": invalid,
	}
	if err := tmpl.Execute(w, data); err != nil {
		log.Panic(err)
	}
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

			ファイルを調べたい場合はこちら。
			<form action="/count/file" method="POST" enctype="multipart/form-data">
			<input type="file" name="file">
			<input type="submit">
			</form>

		</body>
	</html>
	`))
}

func WordCountHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("format") == "json" {
		JSONWordCountHandler(w, r)
	} else {
		HTMLWordCountHandler(w, r)
	}
}

func JSONWordCountHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	text := r.FormValue("text")
	count, bc, _ := wc.WordCountInString(text)

	// json.Marshalは構造体の公開フィールドしか出力してくれないので注意。
	// 小文字でJSONのキーを出力したければタグを指定する。
	type WordCount struct {
		Text      string `json:"text"`
		Count     int    `json:"count"`
		ByteCount int    `json:"byte_count"`
	}

	result, err := json.Marshal(WordCount{text, count, bc})
	if err != nil {
		log.Panic(err)
	}
	w.Write(result)
}

func HTMLWordCountHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	text := r.FormValue("text")
	count, _, _ := wc.WordCountInString(text)

	tmpl, _ := htmlTemplate.New("wc").Parse(`
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

			入力文字: {{.text}}<br>
			文字数は: {{.count}}<br>

			でした。<br><br>

			文字を入力してください。
			<form action="/count" method="GET">
			<input type="text" name="text" size="32" class="{{.css}}">
			<input type="submit">
			</form>
		</body>
	</html>`)

	css := ""
	if count == 0 {
		css = "err"
	}

	data := Context{"text": text, "count": count, "css": css}
	if err := tmpl.Execute(w, data); err != nil {
		log.Panic(err)
	}
}

## jisyotool

辞書の原本を JSON 形式に変換する機能、辞書の重複をチェックする機能、辞書の要素数を調べる機能があります。

### インストール

```
go install github.com/okikae/jisyotool
```

### 使い方

コマンドラインから次のように使います。

```
jisyotool オプション 入力ファイル
```

オプションは下記です。

- `-n` カラム1 カラム2 の順で標準出力に表示する
- `-r` カラム2 カラム1 の順で標準出力に表示する
- `-j` JSON 形式でファイルに書き出す
- `-t` JSON 形式で ts 拡張子のファイルに書き出す
- `-c` 辞書をチェックする
- `-l` 要素数を調べる

下記は仮名辞書を JSON 形式のファイルに書き出す例です。

```
jisyotool -j kana-jisyo
```

`-j`, `-t` オプションで書き出されるファイルの場所はカレントディレクトリです。
# go-rake-ja
Go言語で実装した、RAKE(Rapid Automatic Keyword Extraction)による、日本語テキストからのキーフレーズ抽出器です。

## 使い方

テストコード(extractor_test.go)を参照してください。

## インストール

### Mecabのインストール

テキストの形態素解析にMecabを使っているので、Mecabをインストールしてください。

### 本パッケージのインストール
`go get`で本パッケージを取得できます。

```
$ go get github.com/kanedazz/go-rake-ja
```

なお、MecabのGoバインディングに[go-mecab](https://github.com/shogo82148/go-mecab)を使っているので、ビルド時には、go-mecabの[README](https://github.com/shogo82148/go-mecab?tab=readme-ov-file#install)に記載の通り、Mecabとリンクするために以下の環境変数を設定してください。

```
$ export CGO_LDFLAGS="`mecab-config --libs`"
$ export CGO_FLAGS="`mecab-config --inc-dir`"
```

## 参照

### RAKE
- Rose, Stuart, et al. "Automatic keyword extraction from individual documents." Text mining: applications and theory (2010): 1-20.

### MecabのGoバインディングについて
- https://github.com/shogo82148/go-mecab/
- https://shogo82148.github.io/blog/2016/02/11/golang-mecab-binding/
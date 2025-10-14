# bpflog（日本語）

bpflog は、eBPF プログラムから perf リングバッファを使ってユーザ空間へメッセージを出力する、シンプルな Go ライブラリとそのサンプルです。サンプルは XDP プログラムをアタッチし、IPv4 の送信元/宛先ペアを記録して可読なログ行を出力します。Go 側はそれを受信して標準出力にプリントします。

## サンプルのビルドと実行（簡易）

1. Go (1.25+) 、clang/llvm、[bpftool](https://github.com/libbpf/bpftool)、[mise](https://mise.jdx.dev) をインストールします。
2. リポジトリのルートからサンプルのビルドと実行コマンドを実行します。

```bash
# vmlinux ヘッダを生成
mise gen-vmlinux

# サンプルを実行
cd example
mise start
```

## リポジトリの構成

- `handler.go` - `cilium/ebpf` の perf リーダーをラップし、Start/Stop を提供する小さなライブラリ。
- `example/` - XDP プログラム（C）、生成された Go バインディング、XDP をアタッチしてログを消費する Go `main` のサンプル。
- `bpf.c`, `logf.h` - サンプルで使う BPF プログラムとヘルパー。

## 注意とヒント

- サンプルは `example/main.go` 内でインターフェース index を `2` にハードコードしています。自分の環境の index に変更してください（例: `ip link` で確認）。インターフェース名を受け取り index に変換するようコードを修正することもできます。
- サンプルの BPF プログラムは 192.168.0.0/16 をフィルタしています。必要に応じて `example/bpf.c` を調整してください。
- `handler.go` の `NewHandler` は `perf.Record` を読み込み、ハンドラコールバックを呼び出す Start/Stop ループを提供します。デフォルトの Read デッドラインは 100ms で、優雅な停止を可能にしています。

## ライセンス

リポジトリ内のライセンス表記（LICENSE ファイルがあればそれを参照）に従います。サンプルの BPF プログラムは GPL を宣言しています。

## 貢献

バグや改善案があれば Issue または PR を開いてください（例: インターフェース名対応、ビルドスクリプトの改善、カーネルヘッダ用の CI など）。

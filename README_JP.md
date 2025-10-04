# bpflog（日本語）

bpflog は、eBPF プログラムから perf リングバッファを使ってユーザ空間へメッセージを出力する、シンプルな Go ライブラリとそのサンプルです。サンプルは XDP プログラムをアタッチし、IPv4 の送信元/宛先ペアを記録して可読なログ行を出力します。Go 側はそれを受信して標準出力にプリントします。

## リポジトリの構成

- `handler.go` - `cilium/ebpf` の perf リーダーをラップし、Start/Stop を提供する小さなライブラリ。
- `example/` - XDP プログラム（C）、生成された Go バインディング、XDP をアタッチしてログを消費する Go `main` のサンプル。
- `bpf.c`, `vmlinux.h`, `logf.h` - サンプルで使う BPF プログラムとヘルパー。

## 要件

- カーネルヘッダと BPF（XDP）をサポートする Linux。
- Go 1.25 以上（`go.mod` にモジュール設定あり）。
- BPF プログラムを手動でビルド/ロードする場合は clang/llvm と `bpftool`/`iproute2`。
- サンプルを実行するには XDP をアタッチするために CAP_NET_ADMIN または root 権限が必要。

## サンプルのビルドと実行（簡易）

1. Go (1.25+) と clang/llvm をインストールします。
2. リポジトリのルートからサンプルのビルドと実行コマンドを実行します。サンプルは `bpf2go` を使って `bpf.c` を Go バインディングに変換します。

```bash
# サンプルバイナリをビルド（bpf2go でバインディング生成）
go generate ./example
go build -o ./example/bin/example ./example

# root（または CAP_NET_ADMIN）で実行し、example/main.go のネットワークインターフェース index を適切に設定してください
sudo ./example/bin/example
```

## 注意とヒント

- サンプルは `example/main.go` 内でインターフェース index を `2` にハードコードしています。自分の環境の index に変更してください（例: `ip link` で確認）。インターフェース名を受け取り index に変換するようコードを修正することもできます。
- サンプルの BPF プログラムは 192.168.0.0/16 をフィルタしています。必要に応じて `example/bpf.c` を調整してください。
- `go generate` が `bpf2go` が見つからず失敗する場合は、以下でインストールしてください：

```bash
GO111MODULE=on go install github.com/cilium/ebpf/cmd/bpf2go@latest
```

- `handler.go` の `NewHandler` は `perf.Record` を読み込み、ハンドラコールバックを呼び出す Start/Stop ループを提供します。デフォルトの Read デッドラインは 100ms で、優雅な停止を可能にしています。

## ライセンス

リポジトリ内のライセンス表記（LICENSE ファイルがあればそれを参照）に従います。サンプルの BPF プログラムは GPL を宣言しています。

## 貢献

バグや改善案があれば Issue または PR を開いてください（例: インターフェース名対応、ビルドスクリプトの改善、カーネルヘッダ用の CI など）。

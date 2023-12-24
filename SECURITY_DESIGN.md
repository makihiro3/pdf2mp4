# Security Design

ユーザーがアップロードしたPDFをMP4をサーバー上で変換するに当たって、セキュリティを厳格に考える必要があった。

PDFはECMAScriptを代表に様々な形式のデータを埋め込めるため扱いに注意を要する。

## Isolate a backend container

RCE脆弱性が存在していた場合に備えてbackendはコンテナ単体でnetwork隔離され他のコンテナへの通信もできないようにしている。(network_mode=none)
frontendとbackend間はunix domain socketを経由して通信しており、frontend->backendへの一方通行になっている。

## Internal container security

コンテナ単体のセキュリティとして、権限昇格を防ぐために以下の3種の制限を明示的に行っている。
- Enable NoNewPrivileged flag
- 管理プロセスが持つ特権(Capabilities)を最小化している。
- Runtimeのコンテナイメージをread onlyにしている。書き換え可能なのはVolumeのみ
- backendコンテナの管理プロセスとJobプロセスは異なるユーザーで実行され、Jobプロセスからは管理プロセスにアクセスできない。

## Limit resouce

coin minerなどのリソース消費系のコードが仕込まれる場合に備えてrequest毎に利用できるCPU timeやmemoryにulimitを用いたJobあたりのリソース制限を行っている。

## Build a minimal FFmpeg

FFmpegはubuntuなどのLinux distroの物は様々なライブラリに依存するコンパイルオプションが有効化しており、動作しているバイナリが大きい。
今回pdftoppmが生成したppmからOpenH264を用いてmp4形式に変換する為の最小構成のffmpegバイナリを用意した。
他のライブラリの脆弱性をなるべく受けないようにしている。

### Compiler Options Hardening Guide for C and C++

OpenSSFから公開されている[Compiler Options Hardening Guide for C and C++
](https://best.openssf.org/Compiler-Hardening-Guides/Compiler-Options-Hardening-Guide-for-C-and-C++) / [Linux日本語訳](https://www.linuxfoundation.jp/openssf/2023/12/compiler-options-hardening-guide-for-c-and-cpp-jp/) 並びに各種Linux DistribuitionのPackage作成時のコンパイルガイドに基づきコンパイルフラグによる脆弱性緩和策を導入している。

# PDF2Image

<p align="center">
  <img src="../build/appicon.png" alt="PDF2Image" width="128">
</p>

<p align="center">
  <a href="../README.md">English</a> | <a href="README.zh-TW.md">繁體中文</a> | <a href="README.zh-CN.md">简体中文</a> | 日本語
</p>

PDF ページを高品質な画像に変換する Windows デスクトップアプリケーションです。[Wails](https://wails.io/)（Go バックエンド + Vue 3 フロントエンド）で構築され、[MuPDF](https://mupdf.com/) による高速・高精度な PDF レンダリングを実現しています。

<h2 id="目次">目次</h2>

- [機能](#機能)
- [クイックスタート](#クイックスタート)
- [使い方](#使い方)
- [前提条件](#前提条件)
- [ソースからビルド](#ソースからビルド)
- [プロジェクト構成](#プロジェクト構成)
- [ライセンス](#ライセンス)

<h2 id="機能">機能 <a href="#目次">⬆</a></h2>

- **出力形式**：JPG または PNG
- **DPI 調整**：72〜600（デフォルト 300）
- **JPEG 品質調整**：10〜100%（デフォルト 90%）
- **柔軟なページ選択**：全ページ、指定ページ、または範囲指定（例：`1-5, 8, 10-12`）
- **並列変換**：1〜20 の worker プロセスを設定可能；各 worker は独立したサブプロセスとして動作し、それぞれ独自の MuPDF インスタンスを持つため完全に分離
- **リアルタイムプレビュー**：ズーム（スクロール）とパン（ドラッグ）に対応；ダブルクリックで表示リセット
- **リアルタイム進捗表示**：ページごとの変換進捗を表示
- **出力先ディレクトリの選択**：保存先を指定可能、デフォルトは PDF と同じディレクトリ
- **多言語 UI**：English、繁體中文 — 右上のドロップダウンから切替可能、設定は自動保存

<h2 id="クイックスタート">クイックスタート <a href="#目次">⬆</a></h2>

<h3 id="方法-aビルド済みリリースをダウンロード推奨">方法 A：ビルド済みリリースをダウンロード（推奨） <a href="#目次">⬆</a></h3>

1. [Releases](https://github.com/alexwudev/go-pdf2image/releases) ページにアクセス
2. 最新の `go-pdf2image.zip` をダウンロード
3. 任意のフォルダに解凍
4. `pdf2image.exe` を実行

> **注意**：`libmupdf.dll` は `pdf2image.exe` と同じディレクトリに配置する必要があります。リリースパッケージに同梱されています。

<h3 id="方法-bソースからビルド">方法 B：ソースからビルド <a href="#目次">⬆</a></h3>

下記の[ソースからビルド](#ソースからビルド)を参照してください。

<h2 id="使い方">使い方 <a href="#目次">⬆</a></h2>

1. `pdf2image.exe` を起動
2. **ファイルを参照**をクリック（またはドラッグ＆ドロップ）して PDF を開く
3. ナビゲーションボタンでページ移動；スクロールでズーム、ドラッグでパン
4. 左パネルで出力設定を調整：
   - **出力形式**：JPG または PNG
   - **DPI**：スライダーで解像度を設定（72〜600）
   - **JPEG 品質**：スライダーで圧縮率を設定（JPG のみ）
   - **同時処理数**：スライダーで並列 worker プロセス数を設定（1〜20）
   - **ページ範囲**：全ページまたはカスタム範囲を指定
   - **出力ディレクトリ**：保存先フォルダを選択
5. **変換**をクリック
6. 変換された画像が出力ディレクトリに保存される

<h2 id="前提条件">前提条件 <a href="#目次">⬆</a></h2>

- **Windows 10/11**（x64）
- **[Microsoft Edge WebView2 Runtime](https://developer.microsoft.com/en-us/microsoft-edge/webview2/)**（ほとんどの Windows 10/11 にプリインストール済み）
- **`libmupdf.dll`** が実行ファイルと同じディレクトリに必要（リリースに同梱）

<h2 id="ソースからビルド">ソースからビルド <a href="#目次">⬆</a></h2>

<h3 id="必要なもの">必要なもの <a href="#目次">⬆</a></h3>

- [Go](https://go.dev/) 1.24+
- [Node.js](https://nodejs.org/)
- [go-winres](https://github.com/tc-hib/go-winres)（アプリアイコン埋め込み用）：`go install github.com/tc-hib/go-winres@latest`

<h3 id="wslからwindowsへクロスコンパイル">WSL（Windows へクロスコンパイル） <a href="#目次">⬆</a></h3>

```bash
./build.sh
```

<h3 id="windowsネイティブ">Windows（ネイティブ） <a href="#目次">⬆</a></h3>

```batch
build.bat
```

<h3 id="開発モード">開発モード <a href="#目次">⬆</a></h3>

[Wails CLI](https://wails.io/docs/gettingstarted/installation) が必要です。

```bash
wails dev
```

<h3 id="libmupdf-dll">libmupdf.dll <a href="#目次">⬆</a></h3>

実行ファイルには `libmupdf.dll`（MuPDF 1.24.9, x64）が同じディレクトリに必要です。WSL からクロスコンパイルする方法：

```bash
# mingw-w64 が必要：sudo apt install gcc-mingw-w64-x86-64
git clone --recursive --branch 1.24.9 --depth 1 https://github.com/ArtifexSoftware/mupdf.git
cd mupdf
make OS=mingw64-cross shared=yes build=release \
  HAVE_X11=no HAVE_GLUT=no HAVE_CURL=no USE_SYSTEM_LIBS=no \
  -j$(nproc)
# 出力：build/shared-release/libmupdf.dll
```

<h2 id="プロジェクト構成">プロジェクト構成 <a href="#目次">⬆</a></h2>

```
go-pdf2image/
├── main.go              # エントリポイント：GUI モードまたは --worker サブプロセスモード
├── app.go               # Go バックエンド：PDF 情報、プレビュー、マルチプロセス変換
├── worker.go            # ヘッドレス worker サブプロセス：ページのレンダリングとエンコード
├── wails.json           # Wails プロジェクト設定
├── winres.json          # go-winres 設定（アイコン＆マニフェスト埋め込み）
├── go.mod / go.sum      # Go 依存関係
├── libmupdf.dll         # MuPDF 共有ライブラリ（ランタイム依存）
├── build.sh             # WSL クロスコンパイルスクリプト
├── build.bat            # Windows ネイティブビルドスクリプト
├── LICENSE              # MIT ライセンス
├── build/
│   ├── appicon.png      # アプリアイコン
│   └── windows/         # Windows マニフェスト＆アイコン
├── docs/                # 翻訳版 README
└── frontend/
    ├── index.html       # メイン HTML
    ├── package.json     # フロントエンド依存関係
    ├── vite.config.ts   # Vite 設定
    └── src/
        ├── main.ts          # Vue アプリ初期化
        ├── App.vue          # ルートレイアウト + 言語切替
        ├── style.css        # グローバルスタイル
        ├── i18n/            # 国際化（en、zh-TW）
        ├── stores/
        │   └── appStore.ts  # Pinia 状態管理
        └── components/
            ├── PdfImport.vue       # PDF ファイル選択
            ├── SettingsPanel.vue   # 出力設定（形式、DPI、品質、並列数、ページ）
            ├── ActionBar.vue       # 変換ボタン＆ステータスメッセージ
            ├── PreviewPanel.vue    # ページプレビュー（ズーム/パン）
            └── ConvertProgress.vue # 変換プログレスバー
```

<h2 id="ライセンス">ライセンス <a href="#目次">⬆</a></h2>

[MIT](../LICENSE)

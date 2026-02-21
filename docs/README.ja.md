# PDF2Image

<p align="center">
  <img src="../build/appicon.png" alt="PDF2Image" width="128">
</p>

<p align="center">
  <a href="../README.md">English</a> | <a href="README.zh-TW.md">繁體中文</a> | <a href="README.zh-CN.md">简体中文</a> | 日本語
</p>

PDF ページを高品質な画像に変換するデスクトップアプリケーションです。[Wails](https://wails.io/)（Go バックエンド + Vue 3 フロントエンド）で構築され、[MuPDF](https://mupdf.com/) による高速・高精度な PDF レンダリングを実現しています。**Windows** と **Linux** に対応（GUI および CLI モード）。

<h2 id="目次">目次</h2>

- [機能](#機能)
- [クイックスタート](#クイックスタート)
- [使い方](#使い方)
- [CLI モード](#cli-モード)
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
- **リアルタイム進捗表示**：ページごとの変換進捗を表示；カスタムタイトルバーが進捗に応じて塗りつぶし；Windows タスクバーボタンにも進捗を表示
- **変換タイマー**：変換中にリアルタイムで経過時間を表示；完了時に合計所要時間を表示
- **ZIP パッケージ**：変換した画像をすべて1つの `.zip` ファイルにまとめるオプション
- **停止ボタン**：変換中にいつでもキャンセル可能；すべての worker サブプロセスを終了し、未完成のファイルをクリーンアップ
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
   - **ZIP パッケージ**：チェックすると1つの `.zip` ファイルとして出力
5. **変換**をクリック（**停止**をクリックで変換をキャンセル）
6. 変換された画像（または ZIP）が出力ディレクトリに保存される

<h2 id="cli-モード">CLI モード <a href="#目次">⬆</a></h2>

GUI を起動せずにコマンドラインから変換を実行：

```bash
pdf2image.exe --cli --pdf INPUT.pdf [オプション]
```

| オプション | デフォルト | 説明 |
|---|---|---|
| `--pdf PATH` | *（必須）* | 入力 PDF ファイル |
| `--format jpg\|png` | `jpg` | 出力画像形式 |
| `--dpi N` | `300` | 解像度（72–600） |
| `--quality N` | `90` | JPEG 品質（10–100、JPG のみ） |
| `--pages SPEC` | 全ページ | ページ選択（例：`1-5,8,10-12`） |
| `--output DIR` | PDF と同じディレクトリ | 出力ディレクトリ |
| `--workers N` | `4` | 並列 worker プロセス数（1–20） |
| `--zip` | オフ | 出力を1つの `.zip` ファイルにパッケージ |

**例：**

```bash
pdf2image.exe --cli --pdf report.pdf --format png --dpi 150 --pages 1-10 --workers 8 --output ./images
```

進捗は stderr に出力されます：

```
PDF: report.pdf (50 pages)
Converting 10 pages | format=png dpi=150 quality=90 workers=8
[10/10] 100% - Page 10 done
Done! 10 files in 5.2s → ./images
```

<h2 id="前提条件">前提条件 <a href="#目次">⬆</a></h2>

**Windows：**

- **Windows 10/11**（x64）
- **[Microsoft Edge WebView2 Runtime](https://developer.microsoft.com/en-us/microsoft-edge/webview2/)**（ほとんどの Windows 10/11 にプリインストール済み）
- **`libmupdf.dll`** が実行ファイルと同じディレクトリに必要（リリースに同梱）

**Linux**（x64）：

- **GTK 3** と **WebKit2GTK 4.0**（GUI モードに必要）
  ```bash
  # Ubuntu/Debian
  sudo apt install libgtk-3-0 libwebkit2gtk-4.0-37
  ```
- **`libmupdf.so`** がダイナミックリンカーの検索パスに必要（リリースに同梱）
  ```bash
  # 同じディレクトリのライブラリで実行
  LD_LIBRARY_PATH=. ./pdf2image
  ```

<h2 id="ソースからビルド">ソースからビルド <a href="#目次">⬆</a></h2>

<h3 id="必要なもの">必要なもの <a href="#目次">⬆</a></h3>

- [Go](https://go.dev/) 1.24+
- [Node.js](https://nodejs.org/)
- [go-winres](https://github.com/tc-hib/go-winres)（Windows ビルドのみ、アプリアイコン埋め込み用）：`go install github.com/tc-hib/go-winres@latest`
- **Linux ビルド**にはさらに：`sudo apt install libgtk-3-dev libwebkit2gtk-4.0-dev`

<h3 id="wslからwindowsへクロスコンパイル">WSL（Windows へクロスコンパイル） <a href="#目次">⬆</a></h3>

```bash
./build.sh            # または：./build.sh windows
```

<h3 id="linuxネイティブ">Linux（ネイティブ） <a href="#目次">⬆</a></h3>

```bash
./build.sh linux
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

<h3 id="libmupdf">libmupdf.dll / libmupdf.so <a href="#目次">⬆</a></h3>

実行ファイルには MuPDF 共有ライブラリ（1.24.9, x64）が同じディレクトリまたはライブラリパスに必要です。

**Windows**（`libmupdf.dll`）— WSL からクロスコンパイル：

```bash
# mingw-w64 が必要：sudo apt install gcc-mingw-w64-x86-64
git clone --recursive --branch 1.24.9 --depth 1 https://github.com/ArtifexSoftware/mupdf.git
cd mupdf
make OS=mingw64-cross shared=yes build=release \
  HAVE_X11=no HAVE_GLUT=no HAVE_CURL=no USE_SYSTEM_LIBS=no \
  -j$(nproc)
# 出力：build/shared-release/libmupdf.dll
```

**Linux**（`libmupdf.so`）— ネイティブビルド：

```bash
git clone --recursive --branch 1.24.9 --depth 1 https://github.com/ArtifexSoftware/mupdf.git
cd mupdf
make shared=yes build=release \
  HAVE_X11=no HAVE_GLUT=no HAVE_CURL=no USE_SYSTEM_LIBS=no \
  -j$(nproc)
# 出力：build/shared-release/libmupdf.so.24.9 → libmupdf.so にリネーム
```

<h2 id="プロジェクト構成">プロジェクト構成 <a href="#目次">⬆</a></h2>

```
go-pdf2image/
├── main.go              # エントリポイント：フレームレス GUI、--cli、または --worker サブプロセスモード
├── app.go               # Go バックエンド：PDF 情報、プレビュー、マルチプロセス変換
├── cli.go               # CLI モード：GUI なしでコマンドラインから変換
├── worker.go            # ヘッドレス worker サブプロセス：ページのレンダリングとエンコード
├── taskbar_windows.go   # Windows タスクバー進捗（ITaskbarList3）とアイコン
├── taskbar_stub.go      # 非 Windows 向けスタブ
├── wails.json           # Wails プロジェクト設定
├── winres.json          # go-winres 設定（アイコン＆マニフェスト埋め込み）
├── go.mod / go.sum      # Go 依存関係
├── libmupdf.dll         # MuPDF 共有ライブラリ — Windows（ランタイム依存）
├── libmupdf.so          # MuPDF 共有ライブラリ — Linux（ランタイム依存）
├── build.sh             # ビルドスクリプト（Windows・Linux 対応）
├── build.bat            # Windows ネイティブビルドスクリプト
├── CHANGELOG.md         # バージョン履歴
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
        ├── App.vue          # ルートレイアウト、カスタムタイトルバー + 言語切替
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

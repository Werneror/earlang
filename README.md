# EarLang

用耳朵学习英语单词。

## 动机

笔者最近阅读了《找对英语学习方法的第一本书》。
这本书讲到学习外语时应该在大脑中建立外语词汇的**发音**和词汇所表达的**概念**之间的联系。
笔者以前学习英语单词时，实际上在建立英语单词的**拼写**和**汉语翻译**之间的联系。
笔者意识到了自己的错误，想要找一个学习单词的工具，来直接建立**发音**和**概念**之间的联系，但没有找到。
笔者决定自己动手写一个，于是有了此项目。

## 使用

若您是 Windows 用户，可在 [releases](https://github.com/werneror/earlang/releases) 页面下载编译好的 exe
文件，下载完成后不用安装，直接双击运行即可。
若您是其他平台用户（未在其他平台测试过，不保证兼容性），需下载源码后自行构建可执行文件，具体方法见后文。

第一次运行程序，启动可能会比较慢，需要一小会才能出现图形界面。请耐心等待，不要重复运行。 图形界面如下图所示：

![EarLang](https://user-images.githubusercontent.com/16622293/202831416-4eac3cef-c416-4c2c-aec7-2dbd4ff37cae.png)

会展示一些图片，并朗读一个英语单词的读音。
点击`左箭头按钮`，或按下快捷键`左方向键`、`A` 或 `P` 可切换到上一个单词。
点击`扬声器按钮`，或按下快捷键`上方向键`、`下方向键`、`W`、`S` 或 `R` 可手动控制朗读一次单词。
点击`右箭头按钮`，或按下快捷键`右方向键`、`D` 或 `N` 可切换到下一个单词。

点击右上角的设置按钮可打开设置页面。设置比较简单，请自行探索。

点击右上角的文件夹按钮，可打开一个文件夹（仅 Windows 平台支持此功能），其中存放着本软件的全部相关文件。以下称这个文件夹为主文件夹。

若想自定义要学习的英语单词，可修改主文件夹中的 `words.txt` 文件，并重启软件。
注意，本软件只适合学习表示具体概念的名词。

若想要重置学习进度，可删除主文件夹中的的 `learned.txt` 和 `progress.txt` 文件，并重启软件。

## 工作原理

本软件会从给定的单词列表中按顺序或随机选择一个单词（默认设置为按顺序），联网下载这个单词相关的图片（默认使用 bing 图片搜索），
联网下载这个单词的读音（默认从剑桥词典下载单词读音）。联网下载的图片和读音都会缓存在主文件夹中，同一单词在配置不变的情况下不会重复下载。

## 构建可执行文件

1. 安装 golang，推荐版本 1.19
2. 安装 fyne 命令

```
$ go install fyne.io/fyne/v2/cmd/fyne@latest
```

3. 下载本项目源码

```
$ git clone https://github.com/werneror/earlang.git
```

4. 进入项目根目录

```
$ cd earlang
```

5. 打包可执行文件

```
$ fyne package -name earlang -icon resource\ear.ico --appID wiki.werner.earlang --release
```

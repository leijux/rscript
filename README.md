# rscript

rscript 是通过 ssh 远程执行命令的工具。具有三种运行模式，目的是简化边缘场景下的运维过程。

<img alt="gui" height="450" width="550" src=".\example\gui_2.png"/>

## 快速开始

### gui/tui

step 1: 编写 yaml 脚本

step 2: 运行 rscript_gui.exe ，选择 yaml 脚本执行

### rscript package

rscript package 的作用是通过 go embed 将资源文件和脚本文件嵌入到 tui ，制作成单一执行文件的包。

step 1: git clone https://github.com/leijux/rscript.git

step 2: internal/app/package 编写 yaml 脚本

step 3: 编译包 go build -ldflags "-s -w" -o upgrade.exe

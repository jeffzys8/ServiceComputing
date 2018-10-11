# CLI命令行程序开发 — Go实现

## 安装 pflag

[链接](https://github.com/spf13/pflag)

安装：
```shell
$ go get github.com/spf13/pflag
```

测试安装：
```shell
$ go test github.com/spf13/pflag
```

## 代码细节

- 对 struct sp_args 中的 page_type 作了修改，由int改为bool
```Go
page_type		bool 
/* false for lines-delimited, true for form-feed-delimited */
```

- 使用pflag包，程序对参数的处理变得非常简便，检测参数数量仍使用os包，但由于使用了pflag, 参数的形式 由 ```-s1```更换为 ```-s 1```类似的形式，因此检测参数数目的时候也应有所改变

```Go
if(len(os.Args) < 5){
		fmt.Fprintf(os.Stderr, "%s: %s\n", progname, "not enough arguments")
		usage()
		os.Exit(1)
	}
```

- 检查是否有文件输入，比起C语言实现要方便多了；当然这个步骤要放在其他参数读取以后
```Go
if flag.NArg() == 1 {
	sa.inFilename = flag.Arg(0)
} else {
	sa.inFilename = ""
}
```

- 处理从Stdin或从File读取input
```Go
func make_reader(sa *sp_args) *bufio.Reader{
	
	/* Stdin by default */
	in_fd := os.Stdin
	if len(sa.in_filename) > 0 {
		in_fd, err = os.Open(sa.in_filename)
		if err != nil && err != io.EOF {
			panic(err)
		}
	}
	return bufio.NewReader(in_fd)
}
```

- 按行或按页读取input的主要区别即为每次读取的截至标识符不一样，对于按行输入：
```Go
line, err = reader.ReadString('\n')
```
对于按页输入：
```Go
line, err = reader.ReadString('\f')
```

- 翻页
```Go
line_ctr++
if line_ctr > sa.page_len {
	page_ctr++
	line_ctr = 1
}
```

- 输出部分 (Stdout)

- 检测输入的参数和实际页数是否匹配


## 测试

- 测试文档
```
df
sdf
dxcf
asdx
xcvsdf
sdfsdfsdf
xcvwefef
sd
xcvsdl
SDFxcvlsd
sdfsdfsdfsdfsdfsd
xcvsdf
ssd
aas
ccx
```

- 测例1
input 
```Shell
$ go run mine.go -s 2 -e 3 -l 2
1
2
3
4
5
6
7
8
9
10
```
output
```shell
3
4
5
6
mine.exe: done
```

- 测例2
input
```Shell
$ go run .\mine.go -s 2 -e 3 -l 2 test.txt
```

output
```Shell
dxcf
asdx
xcvsdf
sdfsdfsdf
mine.exe: done
```

## 其他

- printf 和 fprintf的区别<br>
[引用](https://blog.csdn.net/ysdaniel/article/details/7052956)<br>
printf 将内容发送到Default的输出设备，通常为本机的显示器，fprintf需要指定输出设备，可以为文件，设备。
 

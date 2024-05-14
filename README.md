## generic

golang generics library include stream, gmap, cache etc.

### 泛型使用最佳实践

在初期设计和实现的泛型库中，采用了面向对象的方式实现。在实际使用中发现，这类泛型会导致编译中间文件和最终文件的急剧膨胀。

结合[goweight](https://github.com/jondot/goweight)项目，以及`go tool nm -size`命令，得出下面最佳实践方式：

1. 针对系统现有slice和map设计泛型函数，参见stream包。

2. 将泛型对象拆分成纯泛型数据对象和处理泛型数据对象的泛型方法，参见gmap包

3. 如果需要设计某个泛型对象，其需要支持数据泛型和方法，则建议减少泛型方法的数量，以及方法body的大小，最好将body中的逻辑实现委托给非泛型函数。

### 具体分析方法：

#### 使用`goweight`

1. 安装[goweight](https://github.com/jondot/goweight)
2. 在项目目录下执行`goweight > dist/goweight.log`，查看编译中间文件大小(
   观察泛型包文件大小，以及使用该泛型对象、方法的包的大小)

#### 使用`go tool nm -size`命令

1. 在编译命令中添加`-work -x`和`2>&1 | tee dist/build_output_linux.txt`，将编译中间文件保存下来，并记录包和编译文件的关系
2. 结合`goweight.log`和`build_output_linux.txt`文件找到包对应的中间文件
3. 使用`go tool nm -size <中间文件>`查看符号表，并做适当代码调整，再反复观察

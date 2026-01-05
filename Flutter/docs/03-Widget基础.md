# Widget 基础

在 Flutter 中，**一切皆为 Widget**。Widget 是 Flutter 应用的基本构建单元，类似于 React 中的组件（Component）或 HTML 中的元素。

## Widget 概念

### 什么是 Widget？

Widget 是描述用户界面一部分的不可变配置。Flutter 使用 Widget 树来构建 UI，类似于 React 的虚拟 DOM 树。

```dart
// Widget 树示例
MaterialApp(
  home: Scaffold(
    body: Center(
      child: Text('Hello Flutter'),
    ),
  ),
)
```

### 与 React/HTML 对比

| Flutter Widget | React Component | HTML Element | 说明 |
|----------------|-----------------|--------------|------|
| `Text('Hello')` | `<Text>Hello</Text>` | `<p>Hello</p>` | 文本显示 |
| `Container(...)` | `<div>...</div>` | `<div>...</div>` | 容器 |
| `Column(...)` | `<View style={{flexDirection: 'column'}}>` | `<div style="display: flex; flex-direction: column">` | 垂直布局 |
| `Row(...)` | `<View style={{flexDirection: 'row'}}>` | `<div style="display: flex; flex-direction: row">` | 水平布局 |
| `Image(...)` | `<Image src="..." />` | `<img src="..." />` | 图片 |
| `Button(...)` | `<Button onClick={...}>` | `<button onclick="...">` | 按钮 |

## Widget 类型

### StatelessWidget（无状态组件）

StatelessWidget 是不可变的，一旦创建就不会改变。类似于 React 的函数组件。

```dart
// 无状态组件（类似于 React 函数组件）
class MyText extends StatelessWidget {
  final String text;
  
  const MyText({super.key, required this.text});
  
  @override
  Widget build(BuildContext context) {
    return Text(
      text,
      style: TextStyle(fontSize: 20),
    );
  }
}

// 使用
MyText(text: 'Hello Flutter')

// 与 React 对比
// React:
//   function MyText({ text }) {
//     return <Text style={{fontSize: 20}}>{text}</Text>;
//   }
```

### StatefulWidget（有状态组件）

StatefulWidget 可以在运行时改变状态，类似于 React 的类组件或使用 hooks 的函数组件。

```dart
// 有状态组件（类似于 React 类组件或使用 useState）
class Counter extends StatefulWidget {
  const Counter({super.key});
  
  @override
  State<Counter> createState() => _CounterState();
}

class _CounterState extends State<Counter> {
  int _count = 0;
  
  void _increment() {
    setState(() {
      _count++;
    });
  }
  
  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        Text('Count: $_count'),
        ElevatedButton(
          onPressed: _increment,
          child: Text('Increment'),
        ),
      ],
    );
  }
}

// 与 React 对比
// React:
//   function Counter() {
//     const [count, setCount] = useState(0);
//     return (
//       <>
//         <Text>Count: {count}</Text>
//         <Button onClick={() => setCount(count + 1)}>Increment</Button>
//       </>
//     );
//   }
```

## 常用 Widget

### Text（文本）

```dart
// 基本文本
Text('Hello Flutter')

// 带样式
Text(
  'Hello Flutter',
  style: TextStyle(
    fontSize: 24,
    fontWeight: FontWeight.bold,
    color: Colors.blue,
    decoration: TextDecoration.underline,
  ),
)

// 文本对齐
Text(
  'Hello Flutter',
  textAlign: TextAlign.center,
)

// 最大行数
Text(
  'Very long text...',
  maxLines: 2,
  overflow: TextOverflow.ellipsis,
)

// 与 HTML/CSS 对比
// HTML: <p style="font-size: 24px; color: blue;">Hello</p>
// Flutter: Text('Hello', style: TextStyle(fontSize: 24, color: Colors.blue))
```

### Container（容器）

Container 是最常用的布局和装饰 Widget，类似于 HTML 的 `<div>`。

```dart
Container(
  // 尺寸
  width: 200,
  height: 100,
  
  // 内边距（类似于 CSS padding）
  padding: EdgeInsets.all(16),
  // padding: EdgeInsets.symmetric(horizontal: 16, vertical: 8),
  // padding: EdgeInsets.only(left: 16, top: 8),
  
  // 外边距（类似于 CSS margin）
  margin: EdgeInsets.all(8),
  
  // 装饰（类似于 CSS background, border）
  decoration: BoxDecoration(
    color: Colors.blue,
    borderRadius: BorderRadius.circular(8),
    border: Border.all(color: Colors.black, width: 2),
    boxShadow: [
      BoxShadow(
        color: Colors.grey,
        blurRadius: 4,
        offset: Offset(2, 2),
      ),
    ],
  ),
  
  // 子 Widget
  child: Text('Hello'),
)

// 与 CSS 对比
// CSS:
//   .container {
//     width: 200px;
//     height: 100px;
//     padding: 16px;
//     margin: 8px;
//     background: blue;
//     border-radius: 8px;
//     border: 2px solid black;
//     box-shadow: 2px 2px 4px grey;
//   }
```

### Column（垂直布局）

Column 垂直排列子 Widget，类似于 CSS 的 `flex-direction: column`。

```dart
Column(
  // 主轴对齐（类似于 CSS justify-content）
  mainAxisAlignment: MainAxisAlignment.start,      // flex-start
  // mainAxisAlignment: MainAxisAlignment.center,  // center
  // mainAxisAlignment: MainAxisAlignment.end,      // flex-end
  // mainAxisAlignment: MainAxisAlignment.spaceBetween,  // space-between
  // mainAxisAlignment: MainAxisAlignment.spaceAround,   // space-around
  // mainAxisAlignment: MainAxisAlignment.spaceEvenly,   // space-evenly
  
  // 交叉轴对齐（类似于 CSS align-items）
  crossAxisAlignment: CrossAxisAlignment.start,     // flex-start
  // crossAxisAlignment: CrossAxisAlignment.center, // center
  // crossAxisAlignment: CrossAxisAlignment.end,     // flex-end
  // crossAxisAlignment: CrossAxisAlignment.stretch, // stretch
  
  // 主轴尺寸（类似于 CSS flex）
  mainAxisSize: MainAxisSize.max,                  // 占用最大空间
  // mainAxisSize: MainAxisSize.min,                // 占用最小空间
  
  children: [
    Text('Item 1'),
    Text('Item 2'),
    Text('Item 3'),
  ],
)

// 与 CSS 对比
// CSS:
//   .column {
//     display: flex;
//     flex-direction: column;
//     justify-content: flex-start;
//     align-items: flex-start;
//   }
```

### Row（水平布局）

Row 水平排列子 Widget，类似于 CSS 的 `flex-direction: row`。

```dart
Row(
  mainAxisAlignment: MainAxisAlignment.start,
  crossAxisAlignment: CrossAxisAlignment.center,
  children: [
    Text('Left'),
    Text('Center'),
    Text('Right'),
  ],
)

// 使用 Expanded 让子 Widget 填充空间（类似于 CSS flex: 1）
Row(
  children: [
    Expanded(
      flex: 2,  // 占用 2 份空间
      child: Container(color: Colors.red),
    ),
    Expanded(
      flex: 1,  // 占用 1 份空间
      child: Container(color: Colors.blue),
    ),
  ],
)

// 与 CSS 对比
// CSS:
//   .row {
//     display: flex;
//     flex-direction: row;
//   }
//   .item { flex: 1; }
```

### Stack（层叠布局）

Stack 允许子 Widget 层叠在一起，类似于 CSS 的 `position: absolute`。

```dart
Stack(
  // 对齐方式
  alignment: Alignment.center,
  
  children: [
    // 底层 Widget
    Container(
      width: 200,
      height: 200,
      color: Colors.blue,
    ),
    // 顶层 Widget（使用 Positioned 定位）
    Positioned(
      top: 10,
      right: 10,
      child: Container(
        width: 50,
        height: 50,
        color: Colors.red,
      ),
    ),
  ],
)

// 与 CSS 对比
// CSS:
//   .stack { position: relative; }
//   .overlay { position: absolute; top: 10px; right: 10px; }
```

### Image（图片）

```dart
// 网络图片
Image.network(
  'https://example.com/image.jpg',
  width: 200,
  height: 200,
  fit: BoxFit.cover,  // 类似于 CSS object-fit
)

// 本地图片（需要在 pubspec.yaml 中声明）
Image.asset(
  'assets/images/logo.png',
  width: 200,
  height: 200,
)

// 占位符和错误处理
Image.network(
  'https://example.com/image.jpg',
  loadingBuilder: (context, child, loadingProgress) {
    if (loadingProgress == null) return child;
    return CircularProgressIndicator();
  },
  errorBuilder: (context, error, stackTrace) {
    return Icon(Icons.error);
  },
)

// 与 HTML 对比
// HTML: <img src="..." style="width: 200px; height: 200px; object-fit: cover;" />
```

### Button（按钮）

```dart
// ElevatedButton（凸起按钮，Material Design）
ElevatedButton(
  onPressed: () {
    print('Button clicked');
  },
  child: Text('Click Me'),
  style: ElevatedButton.styleFrom(
    backgroundColor: Colors.blue,
    foregroundColor: Colors.white,
    padding: EdgeInsets.symmetric(horizontal: 20, vertical: 12),
  ),
)

// TextButton（文本按钮）
TextButton(
  onPressed: () {},
  child: Text('Text Button'),
)

// OutlinedButton（轮廓按钮）
OutlinedButton(
  onPressed: () {},
  child: Text('Outlined Button'),
)

// IconButton（图标按钮）
IconButton(
  onPressed: () {},
  icon: Icon(Icons.favorite),
  color: Colors.red,
)

// 与 HTML 对比
// HTML: <button onclick="...">Click Me</button>
```

### TextField（输入框）

```dart
// 基本输入框
TextField(
  decoration: InputDecoration(
    labelText: 'Username',
    hintText: 'Enter your username',
    border: OutlineInputBorder(),
  ),
  onChanged: (value) {
    print('Input: $value');
  },
)

// 带控制器（用于获取/设置值）
final _controller = TextEditingController();

TextField(
  controller: _controller,
  decoration: InputDecoration(
    labelText: 'Email',
    prefixIcon: Icon(Icons.email),
    suffixIcon: IconButton(
      icon: Icon(Icons.clear),
      onPressed: () => _controller.clear(),
    ),
  ),
  keyboardType: TextInputType.emailAddress,
  obscureText: false,  // 密码输入时设为 true
)

// 与 HTML 对比
// HTML: <input type="text" placeholder="Enter username" />
```

### ListView（列表）

```dart
// 简单列表
ListView(
  children: [
    ListTile(title: Text('Item 1')),
    ListTile(title: Text('Item 2')),
    ListTile(title: Text('Item 3')),
  ],
)

// 动态列表（类似于 React 的 map）
ListView.builder(
  itemCount: items.length,
  itemBuilder: (context, index) {
    return ListTile(
      title: Text(items[index]),
      onTap: () {
        print('Tapped: ${items[index]}');
      },
    );
  },
)

// 分隔符列表
ListView.separated(
  itemCount: items.length,
  separatorBuilder: (context, index) => Divider(),
  itemBuilder: (context, index) {
    return ListTile(title: Text(items[index]));
  },
)

// 与 React 对比
// React:
//   {items.map((item, index) => (
//     <ListItem key={index}>{item}</ListItem>
//   ))}
```

### Card（卡片）

```dart
Card(
  elevation: 4,  // 阴影高度
  margin: EdgeInsets.all(8),
  child: Padding(
    padding: EdgeInsets.all(16),
    child: Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          'Card Title',
          style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
        ),
        SizedBox(height: 8),
        Text('Card content goes here'),
      ],
    ),
  ),
)
```

### Icon（图标）

```dart
// Material Icons
Icon(
  Icons.favorite,
  color: Colors.red,
  size: 24,
)

// 自定义图标
Icon(
  Icons.star,
  color: Colors.amber,
  size: 32,
)
```

## Widget 组合示例

### 登录表单

```dart
class LoginForm extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: Text('Login')),
      body: Padding(
        padding: EdgeInsets.all(16),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            // Logo
            Icon(Icons.lock, size: 64, color: Colors.blue),
            SizedBox(height: 32),
            
            // 用户名输入框
            TextField(
              decoration: InputDecoration(
                labelText: 'Username',
                prefixIcon: Icon(Icons.person),
                border: OutlineInputBorder(),
              ),
            ),
            SizedBox(height: 16),
            
            // 密码输入框
            TextField(
              obscureText: true,
              decoration: InputDecoration(
                labelText: 'Password',
                prefixIcon: Icon(Icons.lock),
                border: OutlineInputBorder(),
              ),
            ),
            SizedBox(height: 24),
            
            // 登录按钮
            ElevatedButton(
              onPressed: () {
                // 处理登录逻辑
              },
              style: ElevatedButton.styleFrom(
                minimumSize: Size(double.infinity, 48),
              ),
              child: Text('Login'),
            ),
          ],
        ),
      ),
    );
  }
}
```

### 商品列表卡片

```dart
class ProductCard extends StatelessWidget {
  final String name;
  final double price;
  final String imageUrl;
  
  const ProductCard({
    super.key,
    required this.name,
    required this.price,
    required this.imageUrl,
  });
  
  @override
  Widget build(BuildContext context) {
    return Card(
      margin: EdgeInsets.all(8),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // 商品图片
          Image.network(
            imageUrl,
            width: double.infinity,
            height: 200,
            fit: BoxFit.cover,
          ),
          Padding(
            padding: EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                // 商品名称
                Text(
                  name,
                  style: TextStyle(
                    fontSize: 18,
                    fontWeight: FontWeight.bold,
                  ),
                ),
                SizedBox(height: 8),
                // 价格
                Text(
                  '\$$price',
                  style: TextStyle(
                    fontSize: 20,
                    color: Colors.green,
                    fontWeight: FontWeight.bold,
                  ),
                ),
                SizedBox(height: 16),
                // 购买按钮
                ElevatedButton(
                  onPressed: () {},
                  child: Text('Add to Cart'),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
```

## Widget 生命周期

### StatelessWidget 生命周期

```dart
class MyWidget extends StatelessWidget {
  const MyWidget({super.key});
  
  @override
  Widget build(BuildContext context) {
    // 每次重建时都会调用
    return Container();
  }
}
```

### StatefulWidget 生命周期

```dart
class MyWidget extends StatefulWidget {
  const MyWidget({super.key});
  
  @override
  State<MyWidget> createState() => _MyWidgetState();
}

class _MyWidgetState extends State<MyWidget> {
  @override
  void initState() {
    super.initState();
    // Widget 创建时调用一次（类似于 React 的 componentDidMount）
    print('initState');
  }
  
  @override
  void didChangeDependencies() {
    super.didChangeDependencies();
    // 依赖变化时调用（如 Theme、Localizations）
    print('didChangeDependencies');
  }
  
  @override
  Widget build(BuildContext context) {
    // 构建 Widget 树（类似于 React 的 render）
    return Container();
  }
  
  @override
  void didUpdateWidget(MyWidget oldWidget) {
    super.didUpdateWidget(oldWidget);
    // Widget 更新时调用（类似于 React 的 componentDidUpdate）
    print('didUpdateWidget');
  }
  
  @override
  void dispose() {
    // Widget 销毁时调用（类似于 React 的 componentWillUnmount）
    print('dispose');
    super.dispose();
  }
}
```

## 总结

- ✅ Widget 是 Flutter UI 的基本构建单元
- ✅ StatelessWidget 用于静态 UI，StatefulWidget 用于动态 UI
- ✅ Widget 树类似于 React 的组件树
- ✅ 常用 Widget：Text、Container、Column、Row、ListView 等
- ✅ Widget 可以组合使用，构建复杂的 UI

下一步学习：[布局和样式](./04-布局和样式.md)


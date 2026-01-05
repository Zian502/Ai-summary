# Dart 语言基础

Dart 是 Flutter 的编程语言，由 Google 开发。Dart 语法与 JavaScript/TypeScript 非常相似，如果你熟悉 JavaScript，学习 Dart 会很容易。

## 变量和常量

### 变量声明

```dart
// var - 类型推断（类似于 JavaScript 的 let/var）
var name = 'Flutter';
var age = 25;
var isActive = true;

// 明确类型声明（类似于 TypeScript）
String name = 'Flutter';
int age = 25;
bool isActive = true;

// dynamic - 动态类型（类似于 JavaScript 的变量）
dynamic value = 'Hello';
value = 42;  // 可以改变类型

// final - 运行时常量（类似于 JavaScript 的 const，但只能赋值一次）
final String apiKey = 'abc123';
final currentTime = DateTime.now();  // 运行时计算

// const - 编译时常量（类似于 JavaScript 的 const，但必须是编译时已知）
const String appName = 'MyApp';
const int maxUsers = 100;
```

### 与 JavaScript 对比

| Dart | JavaScript | 说明 |
|------|------------|------|
| `var name = 'Flutter'` | `let name = 'Flutter'` | 类型推断 |
| `String name = 'Flutter'` | `const name: string = 'Flutter'` | 明确类型 |
| `final` | `const` | 运行时常量 |
| `const` | `const`（编译时） | 编译时常量 |
| `dynamic` | 无类型声明 | 动态类型 |

## 数据类型

### 基本类型

```dart
// 数字类型
int age = 25;                    // 整数（类似于 JavaScript 的 Number）
double price = 99.99;           // 浮点数（类似于 JavaScript 的 Number）
num value = 42;                  // 可以是 int 或 double

// 字符串
String name = 'Flutter';         // 单引号（推荐）
String message = "Hello";        // 双引号也可以
String multiLine = '''
  多行
  字符串
''';                             // 多行字符串（类似于 JavaScript 的模板字符串）

// 字符串插值（类似于 JavaScript 的模板字符串）
String greeting = 'Hello, $name!';           // 简单插值
String info = 'Age: ${age + 1}';            // 表达式插值

// 布尔值
bool isActive = true;            // 只能是 true 或 false（不能像 JS 那样用 0/1）

// 列表（类似于 JavaScript 的 Array）
List<String> fruits = ['apple', 'banana', 'orange'];
List<int> numbers = [1, 2, 3];
var mixed = [1, 'hello', true];  // 动态类型列表

// 添加元素
fruits.add('grape');             // 类似于 fruits.push('grape')
fruits.insert(0, 'mango');        // 在索引 0 插入

// 访问元素
print(fruits[0]);                // 类似于 fruits[0]
print(fruits.length);            // 类似于 fruits.length

// Map（类似于 JavaScript 的 Object）
Map<String, dynamic> person = {
  'name': 'Flutter',
  'age': 25,
  'isActive': true,
};

// 访问
print(person['name']);           // 类似于 person.name 或 person['name']
person['email'] = 'test@example.com';  // 添加/更新

// Set（类似于 JavaScript 的 Set）
Set<String> uniqueFruits = {'apple', 'banana', 'orange'};
uniqueFruits.add('apple');       // 不会添加重复项
```

### 空安全（Null Safety）

Dart 2.12+ 引入了空安全，这是 Dart 与 JavaScript 的重要区别。

```dart
// 非空类型（默认）
String name = 'Flutter';         // 不能为 null
// name = null;                  // 编译错误！

// 可空类型（使用 ?）
String? nullableName;            // 可以为 null
nullableName = null;             // 允许
nullableName = 'Flutter';        // 也允许

// 空检查操作符
String? name;
print(name?.length);             // 如果 name 为 null，返回 null（类似于 JS 的 ?.）
print(name ?? 'Default');        // 如果 name 为 null，使用默认值（类似于 JS 的 ??）

// 空断言操作符（确定不为 null 时使用）
String? name = getName();
print(name!.length);             // 告诉编译器 name 不为 null

// 与 JavaScript 对比
// JavaScript: const name = getName(); console.log(name?.length ?? 0);
// Dart:       String? name = getName(); print(name?.length ?? 0);
```

## 函数

### 函数定义

```dart
// 基本函数（类似于 JavaScript 的函数）
void greet() {
  print('Hello');
}

// 带参数
void greetPerson(String name) {
  print('Hello, $name!');
}

// 带返回值
int add(int a, int b) {
  return a + b;
}

// 箭头函数（单表达式）
int multiply(int a, int b) => a * b;

// 可选参数（位置参数）
void greet(String name, [String? title]) {
  if (title != null) {
    print('Hello, $title $name!');
  } else {
    print('Hello, $name!');
  }
}
greet('Flutter');                // Hello, Flutter!
greet('Flutter', 'Mr.');         // Hello, Mr. Flutter!

// 命名参数（推荐）
void greet({required String name, String? title}) {
  if (title != null) {
    print('Hello, $title $name!');
  } else {
    print('Hello, $name!');
  }
}
greet(name: 'Flutter');          // Hello, Flutter!
greet(name: 'Flutter', title: 'Mr.');  // Hello, Mr. Flutter!

// 默认参数值
void greet({String name = 'Guest'}) {
  print('Hello, $name!');
}
greet();                         // Hello, Guest!

// 与 JavaScript 对比
// JavaScript:
//   function greet(name, title) { ... }
//   function greet({ name, title = 'Mr.' }) { ... }
// Dart:
//   void greet(String name, [String? title]) { ... }
//   void greet({required String name, String title = 'Mr.'}) { ... }
```

### 高阶函数

```dart
// 函数作为参数（类似于 JavaScript 的回调）
void processList(List<int> numbers, Function(int) callback) {
  for (var number in numbers) {
    callback(number);
  }
}

processList([1, 2, 3], (n) => print(n * 2));

// 函数作为返回值
Function makeMultiplier(int multiplier) {
  return (int value) => value * multiplier;
}

var doubleIt = makeMultiplier(2);
print(doubleIt(5));              // 10

// 与 JavaScript 对比
// JavaScript: const doubleIt = (multiplier) => (value) => value * multiplier;
// Dart:       Function makeMultiplier(int multiplier) => (int value) => value * multiplier;
```

### 匿名函数和闭包

```dart
// 匿名函数（类似于 JavaScript 的箭头函数）
var numbers = [1, 2, 3, 4, 5];
numbers.forEach((number) {
  print(number * 2);
});

// 箭头函数（单表达式）
numbers.forEach((number) => print(number * 2));

// map、filter、reduce（类似于 JavaScript 的数组方法）
var doubled = numbers.map((n) => n * 2).toList();
var evens = numbers.where((n) => n % 2 == 0).toList();
var sum = numbers.reduce((a, b) => a + b);

// 与 JavaScript 对比
// JavaScript: const doubled = numbers.map(n => n * 2);
// Dart:       var doubled = numbers.map((n) => n * 2).toList();
```

## 类（Class）

### 基本类定义

```dart
// 类定义（类似于 JavaScript ES6 的 class）
class Person {
  // 属性
  String name;
  int age;
  
  // 构造函数
  Person(this.name, this.age);
  
  // 命名构造函数
  Person.fromJson(Map<String, dynamic> json)
      : name = json['name'],
        age = json['age'];
  
  // 方法
  void greet() {
    print('Hello, I am $name, $age years old.');
  }
  
  // Getter
  String get info => '$name ($age)';
  
  // Setter
  set updateAge(int newAge) {
    if (newAge > 0) {
      age = newAge;
    }
  }
}

// 使用
var person = Person('Flutter', 25);
person.greet();                  // Hello, I am Flutter, 25 years old.
print(person.info);              // Flutter (25)

// 与 JavaScript 对比
// JavaScript:
//   class Person {
//     constructor(name, age) {
//       this.name = name;
//       this.age = age;
//     }
//     greet() { console.log(`Hello, I am ${this.name}`); }
//   }
```

### 私有成员

```dart
// 以下划线开头的成员是私有的
class BankAccount {
  String _accountNumber;         // 私有属性
  double _balance = 0;
  
  BankAccount(this._accountNumber);
  
  // 公共方法访问私有属性
  double get balance => _balance;
  
  void deposit(double amount) {
    _balance += amount;
  }
  
  void withdraw(double amount) {
    if (_balance >= amount) {
      _balance -= amount;
    }
  }
}
```

### 继承

```dart
// 继承（类似于 JavaScript 的 extends）
class Animal {
  String name;
  
  Animal(this.name);
  
  void makeSound() {
    print('Some sound');
  }
}

class Dog extends Animal {
  Dog(String name) : super(name);
  
  @override
  void makeSound() {
    print('Woof!');
  }
  
  void fetch() {
    print('$name is fetching!');
  }
}

var dog = Dog('Buddy');
dog.makeSound();                 // Woof!
dog.fetch();                     // Buddy is fetching!

// 与 JavaScript 对比
// JavaScript:
//   class Animal {
//     constructor(name) { this.name = name; }
//     makeSound() { console.log('Some sound'); }
//   }
//   class Dog extends Animal {
//     makeSound() { console.log('Woof!'); }
//   }
```

### 抽象类和接口

```dart
// 抽象类（类似于 TypeScript 的 abstract class）
abstract class Shape {
  double area();                 // 抽象方法
  void draw();                   // 抽象方法
}

class Circle extends Shape {
  double radius;
  
  Circle(this.radius);
  
  @override
  double area() => 3.14 * radius * radius;
  
  @override
  void draw() {
    print('Drawing a circle');
  }
}

// 接口（使用 implements）
class Flyable {
  void fly() {}
}

class Bird implements Flyable {
  @override
  void fly() {
    print('Bird is flying');
  }
}
```

### Mixin

```dart
// Mixin（类似于 JavaScript 的 mixin 或多重继承）
mixin Swimmable {
  void swim() {
    print('Swimming...');
  }
}

mixin Flyable {
  void fly() {
    print('Flying...');
  }
}

class Duck with Swimmable, Flyable {
  void quack() {
    print('Quack!');
  }
}

var duck = Duck();
duck.swim();                     // Swimming...
duck.fly();                      // Flying...
duck.quack();                    // Quack!
```

## 异步编程

### Future（类似于 Promise）

```dart
// Future（类似于 JavaScript 的 Promise）
Future<String> fetchData() {
  return Future.delayed(Duration(seconds: 2), () {
    return 'Data loaded';
  });
}

// 使用 then（类似于 Promise.then）
fetchData().then((data) {
  print(data);
});

// 使用 async/await（类似于 JavaScript 的 async/await）
Future<void> loadData() async {
  var data = await fetchData();
  print(data);
}

// 错误处理
Future<String> fetchDataWithError() async {
  await Future.delayed(Duration(seconds: 1));
  throw Exception('Error occurred');
}

try {
  var data = await fetchDataWithError();
} catch (e) {
  print('Error: $e');
}

// 与 JavaScript 对比
// JavaScript:
//   async function fetchData() {
//     await new Promise(resolve => setTimeout(resolve, 2000));
//     return 'Data loaded';
//   }
// Dart:
//   Future<String> fetchData() async {
//     await Future.delayed(Duration(seconds: 2));
//     return 'Data loaded';
//   }
```

### Stream（类似于 Observable）

```dart
// Stream（类似于 JavaScript 的 Observable 或 EventEmitter）
Stream<int> countStream() async* {
  for (int i = 1; i <= 5; i++) {
    await Future.delayed(Duration(seconds: 1));
    yield i;                     // 类似于 JavaScript 的 yield
  }
}

// 监听 Stream
countStream().listen((value) {
  print(value);
});

// 使用 await for
Future<void> listenToStream() async {
  await for (var value in countStream()) {
    print(value);
  }
}
```

## 枚举（Enum）

```dart
// 枚举（类似于 TypeScript 的 enum）
enum Status {
  pending,
  approved,
  rejected,
}

var status = Status.pending;

// Switch 中使用
switch (status) {
  case Status.pending:
    print('Pending');
    break;
  case Status.approved:
    print('Approved');
    break;
  case Status.rejected:
    print('Rejected');
    break;
}

// 与 TypeScript 对比
// TypeScript: enum Status { Pending, Approved, Rejected }
// Dart:       enum Status { pending, approved, rejected }
```

## 泛型（Generics）

```dart
// 泛型（类似于 TypeScript 的泛型）
class Box<T> {
  T value;
  
  Box(this.value);
  
  T getValue() => value;
}

var stringBox = Box<String>('Hello');
var intBox = Box<int>(42);

// 泛型函数
T first<T>(List<T> items) {
  return items[0];
}

var firstString = first<String>(['a', 'b', 'c']);
var firstInt = first<int>([1, 2, 3]);

// 与 TypeScript 对比
// TypeScript: class Box<T> { value: T; }
// Dart:       class Box<T> { T value; }
```

## 扩展方法（Extension Methods）

```dart
// 扩展方法（类似于 JavaScript 的原型扩展）
extension StringExtension on String {
  bool get isEmail {
    return contains('@');
  }
  
  String capitalize() {
    return '${this[0].toUpperCase()}${substring(1)}';
  }
}

var email = 'test@example.com';
print(email.isEmail);            // true
print('hello'.capitalize());     // Hello
```

## 常用操作符

```dart
// 级联操作符（类似于 JavaScript 的方法链）
var person = Person('Flutter', 25)
  ..age = 26
  ..greet();

// 等同于
var person = Person('Flutter', 25);
person.age = 26;
person.greet();

// 空安全操作符
String? name;
print(name?.length ?? 0);        // 如果 name 为 null，返回 0

// 类型检查
if (value is String) {
  print(value.toUpperCase());
}

// 类型转换
var value = '123';
var number = int.parse(value);
```

## 总结

Dart 与 JavaScript 的主要相似点：
- ✅ 相似的语法结构
- ✅ 支持函数式编程
- ✅ 支持异步编程（async/await）
- ✅ 支持类和继承

主要区别：
- ⚠️ Dart 是强类型语言（类似 TypeScript）
- ⚠️ Dart 有空安全（Null Safety）
- ⚠️ Dart 有编译时常量（const）
- ⚠️ Dart 有 Mixin 支持

掌握这些基础后，就可以开始学习 Flutter Widget 了！


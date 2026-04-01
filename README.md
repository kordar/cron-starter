# cron-starter 使用文档

一个 Go 语言的 cron 任务启动器，基于 **gocron** 和 **goframework-cron**，支持动态配置与结构化日志输出（slog）。

下面通过 `starter_test.go` 测试示例说明如何使用。

------

## 1️⃣ 初始化 Cron Module

```
import (
    cron_starter "github.com/kordar/cron-starter"
)

var handle = cron_starter.NewCronModule("AA", nil, map[string]interface{}{})
```

- 第一个参数 `"AA"` 是模块名称
- 第二个参数可以传入初始化回调函数，这里为 `nil`
- 第三个参数可以传入额外配置 `map[string]interface{}`

------

## 2️⃣ 定义任务配置

```
var cfg = map[string]interface{}{
    "AAA": map[string]interface{}{
        "id":   "AAA",
        "spec": "@every 5s",
    },
    "BBB": map[string]interface{}{
        "id":   "BBB",
        "spec": "@every 10s",
    },
}
```

- 每个任务 ID 对应一个 map 配置
- 可为不同任务自定义 `spec` 等属性
- 配置会被加载到对应的 cron 任务中，用于初始化任务参数

------

## 3️⃣ 创建自定义任务

```
import (
    "github.com/kordar/gocron"
    "log/slog"
)

type TestNameSchedule struct {
    gocron.BaseSchedule
}

func (s TestNameSchedule) GetId() string {
    return "test-name"
}

func (s TestNameSchedule) GetSpec() string {
    return "@every 5s"
}

func (s TestNameSchedule) Execute() {
    config := s.Config()
    slog.Info("--------------test name--------------", "config", config)
}
```

- 继承 `gocron.BaseSchedule`
- 实现三个方法：
  - `GetId()` → 任务唯一 ID
  - `GetSpec()` → cron 表达式（如 `@every 5s`）
  - `Execute()` → 执行逻辑
- 可通过 `s.Config()` 获取任务对应的配置

------

## 4️⃣ 设置任务配置（可覆盖默认）

```
s := &TestNameSchedule{}
s.SetConfig(map[string]string{
    "spec": "@every 5s",
})
```

- 可以为任务单独设置 cron 表达式或其他自定义参数
- 优先级高于全局配置

------

## 5️⃣ 加载配置到 Cron Module

```
handle.Load(cfg)
```

- 把前面定义的 `cfg` 任务配置加载到模块中
- 模块会根据配置自动注册任务

------

## 6️⃣ 添加任务到调度器

```
import "github.com/kordar/goframework-cron"

goframeworkcron.AddJob("BBB", s)
```

- 第一个参数是任务 ID，对应上面 `cfg` 中的任务
- 第二个参数是实现了 `Schedule` 接口的任务对象
- 添加后，任务会根据 `GetSpec()` 周期自动执行

------

## 7️⃣ 启动任务

```
time.Sleep(100 * time.Second)
```

- 测试中用 `Sleep` 模拟任务运行
- 在实际应用中，可将调度器放到后台服务中运行

------

## 8️⃣ 动态 cron 配置示例

测试中还演示了 **动态初始化函数**：

```
var initializeFn gocron.InitializeFunction = func(job gocron.Schedule) map[string]string {
    cfg := map[string]string{}
    cfg["spec"] = "@every 10s"
    slog.Info("Job initialized", "id", job.GetId())
    if job.GetId() == "AAA" {
        cfg["spec"] = "@every 5s"
    }
    return cfg
}
```

- 可以根据任务 ID 动态生成 cron 表达式
- 可用于在运行前动态计算任务配置

------

## 🔹 总结

### 核心步骤

1. 初始化 CronModule
2. 定义任务配置（cfg）
3. 自定义任务（实现 Schedule 接口）
4. 设置任务参数（可选）
5. 加载任务配置
6. 添加任务到调度器
7. 启动任务 / 运行模块

### 特点

- 支持动态 cron 配置
- 使用标准库 slog 进行结构化日志
- 可快速扩展新任务

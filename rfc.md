# RFC 5545 VEVENT 和 VTODO 规范详解

## 概述

本文档基于 RFC 5545 标准，详细分析 VEVENT（事件）和 VTODO（任务）组件在循环规则（RRULE）相关属性上的规范、实例和区别。

## 1. 基础规范对比

### 1.1 时间属性对比

| 属性 | VEVENT | VTODO | 说明 |
|------|--------|-------|------|
| DTSTART | 必需（除非指定 METHOD） | 可选 | 开始时间 |
| DTEND | 可选 | 不支持 | 事件结束时间（非独占） |
| DUE | 不支持 | 可选 | 任务截止时间（独占） |
| DURATION | 可选 | 可选 | 持续时间 |

### 1.2 时间属性约束

**VEVENT 约束：**
- DTEND 和 DURATION 不能同时存在
- 如果没有 DTEND 和 DURATION，事件在 DTSTART 的同一时间结束
- 全天事件：DTSTART 为 DATE 类型时，DTEND 也必须为 DATE 类型

**VTODO 约束：**
- DUE 和 DURATION 不能同时存在
- 如果存在 DURATION，则必须有 DTSTART
- 可以只有 DUE、只有 DTSTART、DTSTART+DUE、DTSTART+DURATION
- 可以完全没有时间属性（表示持续性任务）

## 2. 全天事件/任务的时间表示差异

### 2.1 全天事件（VEVENT）

对于 2025-09-10 的全天事件：

```ics
BEGIN:VEVENT
DTSTART;VALUE=DATE:20250910
DTEND;VALUE=DATE:20250911
SUMMARY:全天事件示例
END:VEVENT
```

**关键点：**
- DTSTART: 2025-09-10（包含）
- DTEND: 2025-09-11（不包含，表示事件在 2025-09-10 结束）
- DTEND 使用非独占结束时间，需要 +1 天

### 2.2 全天任务（VTODO）

对于 2025-09-10 的全天任务：

```ics
BEGIN:VTODO
DTSTART;VALUE=DATE:20250910
DUE;VALUE=DATE:20250910
SUMMARY:全天任务示例
END:VTODO
```

**关键点：**
- DTSTART: 2025-09-10（包含）
- DUE: 2025-09-10（包含，表示任务在 2025-09-10 截止）
- DUE 使用独占截止时间，不需要 +1 天

### 2.3 时间表示总结

| 类型 | 2025-09-10 全天 | 结束时间表示 | 说明 |
|------|-----------------|--------------|------|
| VEVENT | DTEND=20250911 | 非独占 | 需要 +1 天表示结束 |
| VTODO | DUE=20250910 | 独占 | 直接使用当天表示截止 |

## 3. RRULE 循环规则

### 3.1 基本规则

RRULE 属性可用于 VEVENT、VTODO 和 VJOURNAL 组件：

```
RRULE:FREQ=DAILY;COUNT=10
RRULE:FREQ=WEEKLY;UNTIL=20251224T000000Z
RRULE:FREQ=MONTHLY;INTERVAL=2;BYDAY=1MO
```

### 3.2 DTSTART 的重要性

**核心原则：**
- RRULE 的计算基准是 DTSTART 属性
- DTSTART 定义循环集合中的第一个实例
- 缺失的 RRULE 信息从 DTSTART 中推导

**示例：**
```ics
DTSTART:20250101T090000
RRULE:FREQ=WEEKLY
# 结果：每周一 9:00 AM（从 DTSTART 推导出星期和时间）
```

### 3.3 没有 DTSTART 的 VTODO

**规范说明：**
根据 RFC 5545，VTODO 可以没有 DTSTART 和 DUE 属性。在这种情况下：

1. **不支持 RRULE**：没有 DTSTART 就没有循环计算的基准时间
2. **特殊行为**：VTODO 会与每个连续的日历日期相关联，直到完成
3. **实现差异**：不同的日历应用可能有不同的处理方式

**示例：**
```ics
BEGIN:VTODO
UID:todo-without-time@example.com
SUMMARY:持续性任务
DESCRIPTION:没有时间限制的任务
STATUS:NEEDS-ACTION
END:VTODO
```

## 4. UNTIL 属性行为

### 4.1 UNTIL 的通用规则

UNTIL 属性在 VEVENT 和 VTODO 中的行为基本相同：

```
RRULE:FREQ=DAILY;UNTIL=20251224T000000Z
```

**关键点：**
- UNTIL 指定循环的结束时间（包含）
- 必须与 DTSTART 的值类型匹配（DATE 或 DATE-TIME）
- 如果 DTSTART 是本地时间，UNTIL 应该是 UTC 时间

### 4.2 全天事件/任务的 UNTIL

**全天事件：**
```ics
DTSTART;VALUE=DATE:20250101
RRULE:FREQ=DAILY;UNTIL=20251231
```

**全天任务：**
```ics
DTSTART;VALUE=DATE:20250101
RRULE:FREQ=DAILY;UNTIL=20251231
```

**注意：**
- 全天循环的 UNTIL 使用 DATE 格式
- 不需要考虑 DTEND/DUE 的差异，UNTIL 只影响循环生成

## 5. 完整示例对比

### 5.1 每日会议（VEVENT）

```ics
BEGIN:VEVENT
UID:daily-meeting@example.com
DTSTART;TZID=Asia/Shanghai:20250101T090000
DTEND;TZID=Asia/Shanghai:20250101T100000
RRULE:FREQ=DAILY;COUNT=30
SUMMARY:每日站会
END:VEVENT
```

### 5.2 每日任务（VTODO）

```ics
BEGIN:VTODO
UID:daily-task@example.com
DTSTART;TZID=Asia/Shanghai:20250101T090000
DUE;TZID=Asia/Shanghai:20250101T180000
RRULE:FREQ=DAILY;COUNT=30
SUMMARY:每日工作任务
END:VTODO
```

### 5.3 复杂循环示例

**每月最后一个工作日的任务：**
```ics
BEGIN:VTODO
UID:monthly-report@example.com
DTSTART;VALUE=DATE:20250131
DUE;VALUE=DATE:20250131
RRULE:FREQ=MONTHLY;BYDAY=-1MO,-1TU,-1WE,-1TH,-1FR;BYSETPOS=-1
SUMMARY:月度报告
END:VTODO
```

## 6. 实现注意事项

### 6.1 时区处理

1. **本地时间 + 时区**：推荐用于循环事件/任务
2. **UTC 时间**：适用于跨时区场景
3. **浮动时间**：用于全天事件/任务

### 6.2 全天处理

1. **VEVENT**：DTEND = DTSTART + 1 天
2. **VTODO**：DUE = 截止日期（不加 1 天）
3. **时间格式**：使用 VALUE=DATE 参数

### 6.3 循环计算

1. **基准时间**：始终基于 DTSTART
2. **时间推导**：缺失的时间信息从 DTSTART 推导
3. **异常处理**：使用 EXDATE 排除特定实例

## 7. 主要区别总结

| 方面 | VEVENT | VTODO |
|------|--------|-------|
| 结束时间属性 | DTEND（非独占） | DUE（独占） |
| 全天结束时间 | 需要 +1 天 | 使用当天 |
| 时间属性必需性 | DTSTART 必需 | 所有时间属性可选 |
| 无时间属性支持 | 不支持 | 支持（持续性任务） |
| RRULE 支持 | 需要 DTSTART | 需要 DTSTART |
| 循环计算基准 | DTSTART | DTSTART（如果存在） |

## 8. 最佳实践建议

1. **时间属性选择**：优先使用 DTEND/DUE 而非 DURATION
2. **时区一致性**：确保 DTSTART、DTEND/DUE、UNTIL 使用一致的时区
3. **全天事件处理**：注意 DTEND 和 DUE 的不同语义
4. **循环规则验证**：确保 DTSTART 与 RRULE 同步
5. **异常处理**：合理使用 EXDATE 和 RDATE 处理特殊情况

## 9. Recurrence 存储方式建议

### 9.1 主流厂商实现方式

基于对 Google Calendar、Apple Calendar 和 Microsoft Outlook 的调研，推荐采用以下存储方式：

#### Google Calendar API 方式（推荐）
```go
type Event struct {
    // 时间信息 - 单独字段
    DTStart   time.Time `json:"dtstart"`
    DTEnd     time.Time `json:"dtend,omitempty"`
    DUE       time.Time `json:"due,omitempty"`
    
    // 循环规则 - []string 数组
    Recurrence []string `json:"recurrence,omitempty"`
    // 示例：
    // ["RRULE:FREQ=WEEKLY;BYDAY=MO,WE,FR", "EXDATE:20250115T100000Z", "RDATE:20250120T100000Z"]
}
```

### 9.2 Recurrence 数组组成

**正确的分组方式：**
1. `RRULE:...` - 循环规则（包含 COUNT/UNTIL）
2. `RDATE:...` - 附加日期
3. `EXDATE:...` - 排除日期
4. **时间信息单独存储**（不在 recurrence 数组中）

### 9.3 设计原则

1. **符合 RFC 5545 规范**：时间属性和循环规则在 iCalendar 中本就是分离的
2. **与主流厂商一致**：Google/Apple/Microsoft 都将时间信息与循环规则分开存储
3. **便于处理**：
   - 时间信息用于基础事件
   - 循环规则用于生成重复实例
   - 异常处理更清晰

### 9.4 实现示例

```go
// 循环规则解析
func ParseRecurrence(recurrence []string) (*RecurrenceRule, error) {
    var rrule *RRule
    var rdates []time.Time
    var exdates []time.Time
    
    for _, line := range recurrence {
        switch {
        case strings.HasPrefix(line, "RRULE:"):
            rrule = parseRRule(line[6:])
        case strings.HasPrefix(line, "RDATE:"):
            rdates = append(rdates, parseRDate(line[6:]))
        case strings.HasPrefix(line, "EXDATE:"):
            exdates = append(exdates, parseExDate(line[7:]))
        }
    }
    
    return &RecurrenceRule{
        RRule:   rrule,
        RDates:  rdates,
        ExDates: exdates,
    }, nil
}
```

### 9.5 关键要点

1. **时间信息分离**：`DTSTART`/`DTEND`/`DUE` 不放在 `recurrence` 数组中
2. **遵循 RFC 5545**：每个字符串都是完整的 iCalendar 属性行
3. **便于扩展**：未来可以轻松添加其他循环相关属性
4. **兼容性好**：与 Google Calendar API 格式完全兼容

---

*本文档基于 RFC 5545 标准编写，旨在为日历应用开发提供准确的技术参考。*
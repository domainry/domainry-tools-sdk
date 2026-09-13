# Domainry Tools SDK

中立工具执行契约：Authority、Definition、Request、Result、Host、ResultAuthorizer，以及禁止远程引用的 JSON schema 验证。

`ResultReadAuthorizer` 是独立的已保存结果读取策略，用当前资料权限验证原定义、请求、完整结果与来源范围，不要求原工具执行权限，也不得重新调用、重试或对账写操作。包含结果的交付／其他资源由调用方另行授权。Tools 注册需要明确设置 `Registration.AuthorizeResultRead`；没有设置时返回 `tool.result_read_unsupported`，旧 `ResultAuthorizer` 不能作为独立阅读授权。调用方可以继续要求完整的原执行授权及旧来源检查，但来源读取策略明确拒绝时不能降级。邮件、日历和网页读取适配器已接入当前 Integration 账号读取与范围校验；写入适配器仍需原有策略。

Agent SDK 的旧 `ConversationTool*` 类型以兼容别名 / 等价接口连接这里。ConversationID、RunID、worker guard 是可选宿主上下文，不会进入模型自定义参数。

可选 `ConfirmationVerifier` 由执行 owner 实现，用其持久记录验证当前主体、会话、运行、具体调用、定义、参数和执行租约。接收一个 `Confirmation` 对象不等于获得批准。需要具体确认的外部写入适配器应在可信启动装配中取得此端口；SDK 不保存确认、不打开 owner 数据库，也不新增浏览器授权入口。验证的是执行许可，已完成结果的可见性仍由当前 `ResultAuthorizer` 决定。

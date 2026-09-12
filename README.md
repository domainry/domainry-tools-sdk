# Domainry Tools SDK

中立工具执行契约：Authority、Definition、Request、Result、Host、ResultAuthorizer，以及禁止远程引用的 JSON schema 验证。

Agent SDK 的旧 `ConversationTool*` 类型以兼容别名 / 等价接口连接这里。ConversationID、RunID、worker guard 是可选宿主上下文，不会进入模型自定义参数。

可选 `ConfirmationVerifier` 由执行 owner 实现，用其持久记录验证当前主体、会话、运行、具体调用、定义、参数和执行租约。接收一个 `Confirmation` 对象不等于获得批准。需要具体确认的外部写入适配器应在可信启动装配中取得此端口；SDK 不保存确认、不打开 owner 数据库，也不新增浏览器授权入口。验证的是执行许可，已完成结果的可见性仍由当前 `ResultAuthorizer` 决定。

# 贡献指南

感谢你对 Radiko JP Player 的关注！

## 开发环境

- Go 1.18+
- ffmpeg
- Git

## 开发流程

1. Fork 本项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 代码规范

```bash
# 格式化代码
make fmt

# 检查代码
make vet

# 运行测试
make test
```

## 提交信息规范

- `feat`: 新功能
- `fix`: 修复bug
- `docs`: 文档更新
- `style`: 代码格式调整
- `refactor`: 代码重构
- `test`: 测试相关
- `chore`: 构建/工具相关

示例：
```
feat: 添加音量控制功能
fix: 修复播放卡顿问题
docs: 更新安装文档
```

## 报告问题

使用 [GitHub Issues](https://github.com/your-username/radikojp/issues) 报告问题时，请包含：

- 操作系统和版本
- Go 版本
- ffmpeg 版本
- 详细的错误信息
- 复现步骤

## 许可证

提交代码即表示你同意以 MIT 许可证发布你的贡献。

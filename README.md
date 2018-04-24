# goTools
golang小工具

## 1.mail.go
发送邮件

## 2.system.go
一些os方法

## 3.google_authenticator.go

通过秘钥，获取google两步验证的6位验证码


```
// MakeAuth 获取key&t对应的验证码
// key 秘钥
// t unix时间戳/秒
func MakeAuth(key string, t int64) (string, error);


// MakeAuthNow 获取key对应的验证码
func MakeAuthNow(key string) (string, error);
```

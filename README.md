# goTools
golang小工具

## 1.mail.go
发送邮件

## 2.system.go
一些os方法

## 3.google_authenticator.go

通过秘钥，获取google两步验证的6位验证码

>由于不知为何，golang系统库的base32计算出的秘钥base32 decode位数总是不对，所以自己做了个[base32.Decode](http://github.com/sunreaver/goTools/base32)

```
// MakeAuth 获取key&t对应的验证码
// key 秘钥
// t 1970年的秒
func MakeAuth(key string, t int64) (string, error);


// MakeAuthNow 获取key对应的验证码
func MakeAuthNow(key string) (string, error);
```
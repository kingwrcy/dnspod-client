# dnspod gui client

`dnspod`自从被腾讯收购以后,每次想修改/添加 dns记录 的时候就需要用微信等登录,实在是`烦不胜烦`.
登录进去之后还得找dns解析菜单,又费不少时间.然后再添加解析,实在是麻烦.

正好看到`dnspod`原本就提供了[API接口](https://www.dnspod.cn/docs/index.html),个人可以直接[申请](https://docs.dnspod.cn/account/5f2d466de8320f1a740d9ff3/)一个自己的token.

有了这个token之后就可以增删改查域名解析记录了.

为了练手[golang walk](https://github.com/lxn/walk)库,写了这个小工具.

打开`dnspod.exe`会自动读取用户目录下的`.dnspod`文件,文本格式,内容是前面申请到的`login_token`,没有换行.

```
17xxxx,0532xxxxxxxxxxxxxxxxxxxxx
```

读取当前用户的域名和记录,可以完成简单的增删改查.

目前基本的功能已完成,如有其他需求可以提issue,也可以自己改,提pr.

![demo](http://cdn.kingwrcy.cn/demo.png)

[在线工具](https://base46.com)

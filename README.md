# GeoIP2-CN-misakaio
本项目 Fork 自 [Hackl0us/GeoIP2-CN](https://github.com/Hackl0us/GeoIP2-CN), 
项目中所使用的 IP 地址信息来自于 [misakaio/chnroutes2](https://github.com/misakaio/chnroutes2)

### 📥 下载链接
| 📦 项目 | 📃 文件 | 🐙 GitHub RAW | 🔧 适用范围
|  :--:  |  :--:  |     :--:     | ---- |
| IP-CIDR 列表 | chnroutes.txt | [点我下载](https://github.com/nasaboy/GeoIP2-CN-misakaio/raw/release/chnroutes.txt) | 防火墙、较老的代理工具等 | 
| GeoIP2 数据库 | Country.mmdb | [点我下载](https://github.com/nasaboy/GeoIP2-CN-misakaio/raw/release/Country.mmdb) | Surge, Shadowrocket,<br>QuantumultX, Clash<br>等较新的代理工具|

### 🙋🏻‍♂️ 使用方式
#### Surge 
Surge 用户请确保你的软件版本满足以下要求：

* Surge for macOS: `4.0.2 (1215)` 或更高
* Surge for iOS / iPadOS: `4.10.0 (1851)` 或更高

macOS 💻 配置方式：Setting - General - GeoIP Database 处粘贴上方复制的 `Country.mmdb` 下载链接，点击 Update Now 即可。

iOS / iPadOS 📱 配置方式： Home 页面拉至最下 - More Settings - 
GeoIP Database - CUSTOM GEOIP DATABASE URL 处粘贴上方复制的 `Country.mmdb` 下载链接，点击 Update Now 即可。

#### Clash
Clash 及其衍生工具（如 Clash X, Clash for Windows, Clash for Android, OpenClash 等）的用户直接通过上面链接下载 `Country.mmdb` 并替换掉 Clash 配置文件夹下的同名文件即可。

#### ShadowRocket 和 Quantmult X
直接在 Safari 中打开 `Country.mmdb` 下载链接，Safari 下载完毕后页面下方会提示 “在...中打开”，点击完成导入。


## 🏅 版权说明
GeoIP® 商标版权归 MaxMind 公司所有。

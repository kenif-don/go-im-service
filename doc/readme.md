[toc]
## 全局响应对象
```text
如果接口有返回值,比如是User对象，则会将对象进行JSON序列化后返回,请自行解析
```
```protobuf
syntax = "proto3";
message ResultDTOResp{
  uint32 code = 1;
  string msg = 2;
  string data = 3;
}
enum ResultDTOCode {
  DEFAULT = 0;
  SUCCESS = 200;
  ERROR = 500;
  /** 跳转2级密码输入页 */
  TO_INPUT_PWD2 = 2000;
}
```
## 一、 初始化配置 InitConfig
### 参数
```text
请求参数没有默认值，请务必全部传
```
```protobuf
// 模型
syntax = "proto3";
message ConfigReq {
  string baseDir = 1; //配置根目录
  DeviceType deviceType = 2; //设备类型
  LogSwitch logSwitch = 3; //日志开关
  string ApiHost = 4;
  string WsHost = 5;
  enum LogSwitch{
    LogSwitchUNKNOWN = 0;
    CLOSE = 2;        //关闭日志
    CONSOLE = 3;      //控制台日志
    FILE = 4;         //文件日志
    CONSOLE_FILE = 5; //控制台+文件
  }

  enum DeviceType {
    Unknown = 0;
    PC = 1;     // 电脑端
    Android = 2; // 手机端
    IOS = 3;    // IOS
    H5 = 4;     // H5
  }
}
```
### 返回结果
```text
成功-200 失败-500 失败才会有错误消息 否则都是success
```
## 二、 注册 Register
### 参数
```protobuf
syntax = "proto3";
message UserReq{
  string  username = 1;
  string  password = 2;
}
```
### 返回结果
```text
成功-200 失败-500 失败才会有错误消息 否则都是success
```
## 三、 登录 Login
### 参数
```protobuf
syntax = "proto3";
message UserReq{
  string  username = 1;
  string  password = 2;
}
```
### 返回结果
```text
成功-200 失败-500 失败才会有错误消息 否则都是success
code-2000 跳转2级密码输入页
```
## 四、 获取登录者信息 Info
```text
code-200 成功 会携带data数据 内容是user对象的json字符串，请自行解析
```
### 参数
```text
无
```
### 返回结果
```json
{
  "id": 61,
  "username": "test123",
  "password": "",
  "password2": "",
  "burstPassword": "",
  "nickname": "用户655916",
  "email": "",
  "publicKey  ": "209505bb58f5e71b30d216dbc0889e2cf77822",
  "privateKey": "",
  "qrcode": "/upload/qrcode/user61.png",
  "intro": "",
  "headImg": "",
  "versionCode": "",
  "noticeType": 1
}
```
### 结果模型
```go
type User struct {
	Id            uint64 `json:"id"`                //主键
	Username      string `json:"username"`          //账号
	Password      string `json:"password"`          //密码
	Password2     string `json:"password2"`         //二级密码
	BurstPassword string `json:"burstPassword"`     //自毁密码
	Nickname      string `json:"nickname"`          //昵称
	Email         string `json:"email"`             //通知邮箱
	PublicKey     string `json:"publicKey"` 
	PrivateKey    string `json:"privateKey"`        
	Qrcode        string `json:"qrcode"`            //二维码
	Intro         string `json:"intro"`             //个性签名
	HeadImg       string `json:"headImg"`           //头像
	VersionCode   string `json:"versionCode"`       
	NoticeType    int    `json:"noticeType"`
}
```
## 五、 上传 Upload
```text
code-200 成功 会携带访问链接
```
### 参数
```protobuf
message UploadReq{
  string path = 1; // 文件路径 包含路径和文件名 列如 C:\\Users\\Administrator\\Desktop\\result.png
}
```
### 结果
```text
仅有公共结果集 Data中存有访问链接
```
## 六、 修改昵称、签名、邮箱、头像
```text
参数采用通用模型，掉的方法不一样
修改昵称：UpdateNickname
修改签名：UpdateIntro
修改邮箱：UpdateEmail
修改头像：UpdateHeadImg
```
### 参数
```protobuf
message UpdateUserReq{
  uint64 id = 1;//用户ID
  string data = 2;//可以是昵称、签名、邮箱、头像
}
```
### 结果
```text
结果与获取登录者信息info返回一致
```
## 七、 退出登录 Logout
```text
无参 code==200 代表退出成功
```

## 八、 根据昵称搜索 Search
```text
原需求为： 搜索好友、非好友、所在群、非所在群  目前近保留非好友记录
```
### 参数
```protobuf
/** 搜索用户、群聊得到的请求参数 */
message SearchReq{
  string keyword = 1; // 搜索关键字
}
```
### 结果
```text
结果为0-20个用户模型数组 未做分页
```
## 九、 添加好友 AddFriend
### 参数
```protobuf
/** 好友请求参数 */
message FriendApplyReq{
  uint64 id = 1;//添加好友时是用户ID 同意或拒绝好友时是新朋友记录ID
  string remark = 2;//备注 同意或拒绝时留空
  int32 state = 3;//-1拒绝 2同意
}
```
### 结果
```text
无参 code==200 代表请求发送成功
```
## 十、 获取"新朋友"列表 SelectAllFriendApply
```text
无参,默认获取登录者全部新朋友
```
### 结果
```go
type FriendApply struct {
	Id     uint64 `gorm:"unique;<-:create" json:"id"`
	From   uint64 `json:"from"`
	To     uint64 `json:"to"`
	Remark string `json:"remark"`
	State  int    `json:"state"`

	FromUser *User `gorm:"-"` // 用户模型 发起人
}
```
## 十一、 同意或拒绝好友请求 UpdateFriendApply
```protobuf
message FriendApplyReq{
  uint64 id = 1;//添加好友时是用户ID 同意或拒绝好友时是新朋友记录ID
  string remark = 2;//备注 同意或拒绝时留空
  int32 state = 3;//-1拒绝 2同意
}
```
```text
通用返回结果 code==200 代表操作成功
```

## 十二、 获取通讯录列表 SelectAllFriend
```text
无参
```
### 结果
```go
Id    uint64 `gorm:"unique;<-:create" json:"id"`
	Me    uint64 `json:"me"`
	He    uint64 `json:"he"`
	Name  string `json:"name"`
	State int    `json:"state"`

	HeUser *User `gorm:"-"` //用户模型
```
## 十三、好友详情页查询单个好友 SelectOneFriend
### 参数
```protobuf
message FriendReq{
  uint64 id = 1; //好友ID
  string name = 2; //好友备注 查询/删除时留空，修改好友备注时传入
}
```
### 结果
```go
Id    uint64 `gorm:"unique;<-:create" json:"id"`
	Me    uint64 `json:"me"`
	He    uint64 `json:"he"`
	Name  string `json:"name"`
	State int    `json:"state"`

	HeUser *User `gorm:"-"` //用户模型
```
## 十四、修改好友备注 UpdateFriendName
### 参数
```protobuf
message FriendReq{
  uint64 id = 1; //好友ID
  string name = 2; //好友备注 查询/删除时留空，修改好友备注时传入
}
```
### 结果
```text
通用返回结果
```
## 十五、删除好友 DelFriend
### 参数
```protobuf
message FriendReq{
  uint64 id = 1; //好友ID
  string name = 2; //好友备注 查询/删除时留空，修改好友备注时传入
}
```
### 结果
```text
通用返回结果
```
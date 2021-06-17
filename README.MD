# api代码生成器

通过执行`gali gen api`即可生成api/service/dao/model等单表的所有操作

复制[gf-cli](https://github.com/gogf/gf-cli)项目，扩展增加支持创建api接口，从而快速生成api相关接口与文档，方便后续修改


通过`gali gen api`会创建的文档包括

- dao     基于gf-cli的`gen dao`生成
- model   基于gf-cli的`gen dao`生成
- service 创建单表操作的service，包括CRUD的基本功能
- api     接口文档与注释说明，方便后续维护增加




## build

`
gf pack ./templates commands/gen/gen_api_templates.go

go build -o gali.exe main.go


gf build main.go -a amd64 -s linux,windows,darwin -n gali

`
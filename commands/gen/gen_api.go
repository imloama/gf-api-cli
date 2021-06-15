package gen

import (
	"context"
	"github.com/gogf/gf-cli/library/mlog"
	"github.com/gogf/gf-cli/library/utils"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcmd"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/text/gregex"
	"github.com/gogf/gf/text/gstr"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/lib/pq"
	//_ "github.com/mattn/go-oci8"
	//_ "github.com/mattn/go-sqlite3"
)

// generateApiReq is the input parameter for generating dao.
type generateApiReq struct {
	TableName          string // TableName specifies the table name of the table.
	NewTableName       string // NewTableName specifies the prefix-stripped name of the table.
	PrefixName         string // PrefixName specifies the custom prefix name for generated dao and model struct.
	GroupName          string // GroupName specifies the group name of database configuration node for generated DAO.
	ModName            string // ModName specifies the module name of current golang project, which is used for import purpose.
	JsonCase           string // JsonCase specifies the case of generated 'json' tag for model struct, value from gstr.Case* function names.
	DirPath            string // DirPath specifies the directory path for generated files.
	TplDaoIndexPath    string // TplDaoIndexPath specifies the file path for generating dao index files.
	TplDaoInternalPath string // TplDaoInternalPath specifies the file path for generating dao internal files.
	TplModelIndexPath  string // TplModelIndexPath specifies the file path for generating model index content.
	TplModelStructPath string // TplModelStructPath specifies the file path for generating model struct content.
}


func HelpApi() {
	mlog.Print(gstr.TrimLeft(`
USAGE 
    gali gen api [OPTION]
OPTION
    -/--path             directory path for generated files.
    -l, --link           database configuration, the same as the ORM configuration of GoFrame.
    -t, --tables         generate models only for given tables, multiple table names separated with ',' 
    -e, --tablesEx       generate models excluding given tables, multiple table names separated with ',' 
    -g, --group          specifying the configuration group name of database for generated ORM instance,
                         it's not necessary and the default value is "default"
    -p, --prefix         add prefix for all table of specified link/database tables.
    -r, --removePrefix   remove specified prefix of the table, multiple prefix separated with ',' 
    -m, --mod            module name for generated golang file imports.
    -j, --jsonCase       generated json tag case for model struct, cases are as follows:
                         | Case            | Example            |
                         |---------------- |--------------------|
                         | Camel           | AnyKindOfString    | 
                         | CamelLower      | anyKindOfString    | default
                         | Snake           | any_kind_of_string |
                         | SnakeScreaming  | ANY_KIND_OF_STRING |
                         | SnakeFirstUpper | rgb_code_md5       |
                         | Kebab           | any-kind-of-string |
                         | KebabScreaming  | ANY-KIND-OF-STRING |
    -/--tplDaoIndex      template content for Dao index files generating.
    -/--tplDaoInternal   template content for Dao internal files generating.
    -/--tplModelIndex    template content for Model index files generating.
    -/--tplModelStruct   template content for Model internal files generating.
                  
CONFIGURATION SUPPORT
    Options are also supported by configuration file.
    It's suggested using configuration file instead of command line arguments making producing. 
    The configuration node name is "gf.gen.dao", which also supports multiple databases, for example:
    [gfcli]
        [[gfcli.gen.dao]]
            link     = "mysql:root:12345678@tcp(127.0.0.1:3306)/test"
            tables   = "order,products"
            jsonCase = "CamelLower"
        [[gfcli.gen.dao]]
            link   = "mysql:root:12345678@tcp(127.0.0.1:3306)/primary"
            path   = "./my-app"
            prefix = "primary_"
            tables = "user, userDetail"
                  
CONFIGURATION SUPPORT
    Options are also supported by configuration file.
    It's suggested using configuration file instead of command line arguments making producing. 
    The configuration node name is "gf.gen.api", which also supports multiple databases, for example:
    [gali]
        [[gali.gen.api]]
            link     = "mysql:root:12345678@tcp(127.0.0.1:3306)/test"
            tables   = "order,products"
            jsonCase = "CamelLower"
        [[gali.gen.api]]
            link   = "mysql:root:12345678@tcp(127.0.0.1:3306)/primary"
            path   = "./my-app"
            prefix = "primary_"
            tables = "user, userDetail"

EXAMPLES
    gali gen api
    gali gen api -l "mysql:root:12345678@tcp(127.0.0.1:3306)/test"
    gali gen api -path ./model -c config.yaml -g user-center -t user,user_detail,user_login
    gali gen api -r user_
`))
}

// doGenApi implements the "gen api" command.
func doGenApi() {
	parser, err := gcmd.Parse(g.MapStrBool{
		"path":           true,
		"m,mod":          true,
		"l,link":         true,
		"t,tables":       true,
		"e,tablesEx":     true,
		"g,group":        true,
		"c,config":       true,
		"p,prefix":       true,
		"r,removePrefix": true,
		"j,jsonCase":     true,
		"tplDaoIndex":    true,
		"tplDaoInternal": true,
		"tplModelIndex":  true,
		"tplModelStruct": true,
	})
	if err != nil {
		mlog.Fatal(err)
	}
	config := g.Cfg()
	if config.Available() {
		v := config.GetVar(nodeNameGenDaoInConfigFile)
		if v.IsEmpty() && g.IsEmpty(parser.GetOptAll()) {
			mlog.Fatal(`command arguments and configurations not found for generating dao files`)
		}
		if v.IsSlice() {
			for i := 0; i < len(v.Interfaces()); i++ {
				doGenApiForArray(i, parser)
			}
		} else {
			doGenApiForArray(-1, parser)
		}
	} else {
		doGenApiForArray(-1, parser)
	}
	mlog.Print("done!")
}

// doGenApiForArray implements the "gen dao" command for configuration array.
func doGenApiForArray(index int, parser *gcmd.Parser) {
	doGenDaoForArray(index, parser, generateApiContent)

}

// generateApiAndDaoAndModelContentFile generates the dao and model content of given table.
func generateApiContent(db gdb.DB, req generateDaoReq) {
	// Generating table data preparing.
	fieldMap, err := db.TableFields(context.TODO(), req.TableName)
	if err != nil {
		mlog.Fatalf("fetching tables fields failed for table '%s':\n%v", req.TableName, err)
	}
	// Change the `newTableName` if `prefixName` is given.
	newTableName := req.PrefixName + req.NewTableName
	var (
		tableNameCamelCase      = gstr.CaseCamel(newTableName)
		tableNameCamelLowerCase = gstr.CaseCamelLower(newTableName)
		tableNameSnakeCase      = gstr.CaseSnake(newTableName)
		importPrefix            = ""
		dirRealPath             = gfile.RealPath(req.DirPath)
	)
	if dirRealPath == "" {
		dirRealPath = req.DirPath
		importPrefix = dirRealPath
		importPrefix = gstr.Trim(dirRealPath, "./")
	} else {
		importPrefix = gstr.Replace(dirRealPath, gfile.Pwd(), "")
	}
	importPrefix = gstr.Replace(importPrefix, gfile.Separator, "/")
	importPrefix = gstr.Join(g.SliceStr{req.ModName, importPrefix}, "/")
	importPrefix, _ = gregex.ReplaceString(`\/{2,}`, `/`, gstr.Trim(importPrefix, "/"))

	fileName := gstr.Trim(tableNameSnakeCase, "-_.")
	if len(fileName) > 5 && fileName[len(fileName)-5:] == "_test" {
		// Add suffix to avoid the table name which contains "_test",
		// which would make the go file a testing file.
		fileName += "_table"
	}

// service
	view := g.View()
	serviceTpl := string(g.Res().GetContent("templates/gen_api_service.vm"))
	if serviceTpl == "" {
		mlog.Fatalf("获取service template失败！")
		return
	}
	ctx := context.Background()
	service, err := view.ParseContent(ctx, serviceTpl, g.Map{
		"TplImportPrefix":            importPrefix,
		"TplTableName":               req.TableName,
		"TplGroupName":               req.GroupName,
		"TplTableNameCamelCase":      tableNameCamelCase,
		"TplTableNameCamelLowerCase": tableNameCamelLowerCase,
		"TplColumnDefine":            gstr.Trim(generateColumnDefinitionForDao(fieldMap)),
		"TplColumnNames":             gstr.Trim(generateColumnNamesForDao(fieldMap)),
	})
	path := gfile.Join(req.DirPath, "service", fileName+"_service.go")
	if err != nil {
		mlog.Fatalf("gen service content failed: %v", err)
		return
	}
	if err := gfile.PutContents(path, strings.TrimSpace(service)); err != nil {
		mlog.Fatalf("writing content to '%s' failed: %v", path, err)
	} else {
		utils.GoFmt(path)
		mlog.Print("generated:", path)
	}
	// api
	apiTpl := string(g.Res().GetContent("templates/gen_api_api.vm"))
	if apiTpl == "" {
		mlog.Fatalf("获取api template失败！")
		return
	}
	apitpl, err := view.ParseContent(ctx, apiTpl, g.Map{
		"TplImportPrefix":            importPrefix,
		"TplTableName":               req.TableName,
		"TplGroupName":               req.GroupName,
		"TplTableNameCamelCase":      tableNameCamelCase,
		"TplTableNameCamelLowerCase": tableNameCamelLowerCase,
		"TplColumnDefine":            gstr.Trim(generateColumnDefinitionForDao(fieldMap)),
		"TplColumnNames":             gstr.Trim(generateColumnNamesForDao(fieldMap)),
	})
	path = gfile.Join(req.DirPath, "api", fileName+"_api.go")
	if err != nil {
		mlog.Fatalf("gen api content failed: %v", err)
		return
	}
	if err := gfile.PutContents(path, strings.TrimSpace(apitpl)); err != nil {
		mlog.Fatalf("writing content to '%s' failed: %v", path, err)
	} else {
		utils.GoFmt(path)
		mlog.Print("generated:", path)
	}
}

package gen

const templateDaoModelIndexContent = `
// =================================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================
package model
{TplPackageImports}
{TplModelStructs}
`

const templateDaoModelStructContent = `
// {TplTableNameCamelCase} is the golang structure for table {TplTableName}.
{TplStructDefine}
`
package driver

import (
	. "github.com/rickb777/sqlapi/schema"
	. "github.com/rickb777/sqlapi/types"
)

var i64 = Type{Name: "int64", Base: Int64}
var boo = Type{Name: "bool", Base: Bool}
var cat = Type{Name: "Category", Base: Int32}
var str = Type{Name: "string", Base: String}
var spt = Type{Name: "string", IsPtr: true, Base: String}
var phn = Type{Name: "PhoneNumber", IsPtr: true, Base: String}
var ipt = Type{Name: "int32", IsPtr: true, Base: Int32}
var upt = Type{Name: "uint32", IsPtr: true, Base: Uint32}
var fpt = Type{Name: "float32", IsPtr: true, Base: Float32}
var sli = Type{Name: "[]string", Base: Slice}
var bgi = Type{PkgPath: "math/big", PkgName: "big", Name: "Int", Base: Struct}
var bys = Type{Name: "[]byte", Base: Slice}
var scv1 = Type{Name: "Foo", IsScanner: true, IsValuer: true, Base: String}
var scv2 = Type{Name: "Foo", IsScanner: true, IsValuer: true, Base: String, IsPtr: true}
var bar1 = Type{Name: "Bar", Base: String}
var bar2 = Type{Name: "Bar", Base: String, IsPtr: true}
var tim = Type{PkgPath: "time", PkgName: "time", Name: "Time", Base: Struct}

var id = &Field{Node: Node{Name: "Id", Type: i64}, SqlName: "id", Encode: ENCNONE, Tags: &Tag{Primary: true, Auto: true}}
var category = &Field{Node: Node{Name: "Cat", Type: cat}, SqlName: "cat", Encode: ENCNONE, Tags: &Tag{Index: "catIdx"}}
var name = &Field{Node: Node{Name: "Name", Type: str}, SqlName: "username", Encode: ENCNONE, Tags: &Tag{Size: 2048, Name: "username", Unique: "nameIdx"}}
var active = &Field{Node: Node{Name: "Active", Type: boo}, SqlName: "active", Encode: ENCNONE}
var mobile = &Field{Node: Node{Name: "Mobile", Type: phn}, SqlName: "mobile", Encode: ENCNONE}
var qual = &Field{Node: Node{Name: "Qual", Type: spt}, SqlName: "qual", Encode: ENCNONE}
var diff = &Field{Node: Node{Name: "Diff", Type: ipt}, SqlName: "diff", Encode: ENCNONE}
var age = &Field{Node: Node{Name: "Age", Type: upt}, SqlName: "age", Encode: ENCNONE}
var bmi = &Field{Node: Node{Name: "Bmi", Type: fpt}, SqlName: "bmi", Encode: ENCNONE}
var labels = &Field{Node: Node{Name: "Labels", Type: sli}, SqlName: "labels", Encode: ENCJSON, Tags: &Tag{Encode: "json"}}
var fave = &Field{Node: Node{Name: "Fave", Type: bgi}, SqlName: "fave", Encode: ENCJSON, Tags: &Tag{Encode: "json"}}
var avatar = &Field{Node: Node{Name: "Avatar", Type: bys}, SqlName: "avatar", Encode: ENCNONE}
var fooey1 = &Field{Node: Node{Name: "Foo1", Type: scv1}, SqlName: "foo1", Encode: ENCNONE}
var fooey2 = &Field{Node: Node{Name: "Foo2", Type: scv2}, SqlName: "foo2", Encode: ENCNONE}
var barey1 = &Field{Node: Node{Name: "Bar1", Type: bar1}, SqlName: "bar1", Encode: ENCDRIVER, Tags: &Tag{Encode: "driver"}}
var barey2 = &Field{Node: Node{Name: "Bar2", Type: bar2}, SqlName: "bar2", Encode: ENCDRIVER, Tags: &Tag{Encode: "driver"}}
var updated = &Field{Node: Node{Name: "Updated", Type: tim}, SqlName: "updated", Encode: ENCTEXT, Tags: &Tag{Size: 100, Encode: "text"}}

var icat = &Index{Name: "catIdx", Fields: FieldList{category}}
var iname = &Index{Name: "nameIdx", Unique: true, Fields: FieldList{name}}

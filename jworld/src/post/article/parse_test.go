package article

import (
	"post"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestPostParser(t *testing.T) {
	newPost, err := Parse(strings.NewReader(testdata.original))
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(testdata.post, newPost) {
		t.Fatalf("expected %+v, got \n%+v", testdata.post, newPost)
	}
}

var testdata = struct {
	original string
	post     *post.Post
}{
	`<HTML>
<HEAD>
<meta HTTP-EQUIV=Content-Type content="text/html;charset=UTF-8">
<base href="http://www.javaworld.com.tw/jute/">
<TITLE>JWorld@TW</TITLE>
<STYLE TYPE='text/css'>
  BODY      { font-family: Tahoma,Georgia; color: #000000; font-size: 12px }
  P         { font-family: Tahoma,Georgia; color: #000000; font-size: 12px } 
  TD        { font-family: Tahoma,Georgia; color: #000000; font-size: 12px ;table-layout:fixed; word-break :break-all}  
  TEXTAREA  { font-family: Tahoma,Georgia; font-size: 12px }
  A         { text-decoration: underline }
  A:link    { color: #000000; text-decoration: underline }
  A:visited { color: #000080; text-decoration: underline }
  A:active  { color: #FF0000; text-decoration: underline }
  A:hover   { color: #FF0000; text-decoration: underline }
  .nav         { text-decoration: underline; color:#000000 }
  .nav:link    { text-decoration: underline; color:#000000 }
  .nav:visited { text-decoration: underline; color:#000000 }
  .nav:active  { text-decoration: underline; color:#FF0000 }
  .nav:hover   { text-decoration: none; color:#FF0000 }
  .topic         { text-decoration: none }
  .topic:link    { text-decoration: none; color:#000000 }
  .topic:visited { text-decoration: none; color:#000080 }
  .topic:active  { text-decoration: none; color:#FF0000 }
  .topic:hover   { text-decoration: underline; color:#FF0000 }
  .ilink         { text-decoration: underline; color:#0000FF }
  .ilink:link    { text-decoration: underline; color:#0000FF }
  .ilink:visited { text-decoration: underline; color:#004080 }
  .ilink:active  { text-decoration: underline; color:#FF0000 }
  .ilink:hover   { text-decoration: underline; color:#FF0000 }          
</STYLE>
</HEAD>
<BODY BGCOLOR='#FFFFFF' TEXT='#000000' LINK='#000000' VLINK='#000080' ALINK='#FF0000' >
Hi,parkghost <br><br>
您訂閱的話題:<a href="post/view?bid=35&id=310575">Re:jdbmaplite 一個簡易的 orm utils</a>有新的回覆!
<br>
<b>KeepItSimple</b>在<a href="post/view?bid=35&id=311689&sty=3" target=_blank>JWorld@TW</a>裡回覆了該話題，下面是回覆內容：<br>
<br>
---------- <br>
Subject: <b>Re:jdbmaplite 一個簡易的 orm utils</b><br>
Date:       2013-01-19 14:51<br><br>
您的功力真不錯<BR>如果有property file的設定方式，會更加有選擇性，xml有點傷眼睛<BR><BR>或是能提供Fluent API也不錯<BR>比如:<BR><BR>List&lt;Customer&gt; result =query.from(customer)<BR>    .where(customer.lastName.like("%den"), customer.active.eq(true))<BR>    .orderBy(customer.lastName.asc(), customer.firstName.desc())<BR>    .list(customer); <br>
	<br>
---------- <br>
<br>
這封EMAIL為系統郵件,請勿直接回覆到此EMAIL, 請到論壇回覆。 <br>

----------<br>
JWorld@TW http://www.javaworld.com.tw/jute
</BODY>
</HTML>`,
	&post.Post{
		Id:       311689,
		Bid:      35,
		Receiver: "parkghost",
		Sender:   "KeepItSimple",
		Subject:  "Re:jdbmaplite 一個簡易的 orm utils",
		Date:     time.Date(2013, 01, 19, 14, 51, 0, 0, time.UTC),
		Content:  "\n" + `您的功力真不錯如果有property file的設定方式，會更加有選擇性，xml有點傷眼睛或是能提供Fluent API也不錯比如:List<Customer> result =query.from(customer)    .where(customer.lastName.like("%den"), customer.active.eq(true))    .orderBy(customer.lastName.asc(), customer.firstName.desc())    .list(customer); ` + "\n\t",
		Link:     "http://www.javaworld.com.tw/jute/post/view?bid=35&id=311689&sty=3",
	},
}

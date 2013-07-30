package general

import (
	"post"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestPostParser(t *testing.T) {
	for _, item := range testdata {
		newPost, err := Parse(strings.NewReader(item.original))
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(item.post, newPost) {
			t.Fatalf("expected %+v, got \n%+v", item.post, newPost)
		}
	}

}

var testdata = []struct {
	original string
	post     *post.Post
}{
	{ //JWorld@TW新話題通知
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
<BODY BGCOLOR='#FFFFFF' TEXT='#000000' LINK='#000000' VLINK='#000080' ALINK='#FF0000'>
Hi,parkghost <br><br>
您訂閱的版面:Java 工作機會有新話題了!
<br>
<b>Alicecys</b>在<a href="http://www.javaworld.com.tw/jute" target=_blank>JWorld@TW</a>裡發了新話題：<br>
<br>
---------- <br>
Subject: <b>知名外商內湖徵 .NET Developer1名-待遇優</b><br>
Date:       2013-01-27 17:24<br><br>
需求人數：&nbsp;&nbsp;1人<BR>職務類別：&nbsp;&nbsp;軟體設計工程師<BR>職務說明：&nbsp;&nbsp;Responsibilities: <BR>1. Develop Web applications and Web services using ASP.NET MVC, Web Forms, and WCF.<BR>2. GIS development.<BR>管理責任：&nbsp;&nbsp;非管理職，無需負擔管理責任<BR> <BR>工作條件限制<BR>學歷：&nbsp;&nbsp;大學以上<BR>科系：&nbsp;&nbsp;資訊工程相關、資訊管理、數學及電算機科學學科類全部<BR>工作經驗：&nbsp;&nbsp;2年以上<BR>語文條件：&nbsp;&nbsp;英文(聽/中等、說/普通、讀/中等、寫/中等)<BR>擅長工具：&nbsp;&nbsp;程式設計類：<BR>ASP.NET MVC and Web Forms、ADO.NET、DHTML and HTML5, CSS3, and JavaScript.<BR>Google Map API or Bing Map API is a plus.<BR> <BR> <BR><BR>【職務聯絡人】 Alice Chang<BR>【E-mail】 it1@recruitexpress.com.tw<BR>【公司電話】: 02-87806811 <br>
  <br>
---------- <br>
<br>
請勿直接回覆此EMAIL,此為系統系統信件, 請到論壇回覆。 <br>

點選以下LINK可以直接連到此話題：<a href="http://www.javaworld.com.tw/jute/post/view?bid=15&id=311825" target=_blank>http://www.javaworld.com.tw/jute/post/view?bid=15&id=311825</a> <br>

----------<br>
JWorld@TW http://www.javaworld.com.tw/jute
</BODY>
</HTML>`,
		&post.Post{
			Id:       311825,
			Bid:      15,
			Receiver: "parkghost",
			Sender:   "Alicecys",
			Subject:  "知名外商內湖徵 .NET Developer1名-待遇優",
			Date:     time.Date(2013, 01, 27, 17, 24, 0, 0, time.UTC),
			Content:  "\n需求人數：  1人職務類別：  軟體設計工程師職務說明：  Responsibilities: 1. Develop Web applications and Web services using ASP.NET MVC, Web Forms, and WCF.2. GIS development.管理責任：  非管理職，無需負擔管理責任 工作條件限制學歷：  大學以上科系：  資訊工程相關、資訊管理、數學及電算機科學學科類全部工作經驗：  2年以上語文條件：  英文(聽/中等、說/普通、讀/中等、寫/中等)擅長工具：  程式設計類：ASP.NET MVC and Web Forms、ADO.NET、DHTML and HTML5, CSS3, and JavaScript.Google Map API or Bing Map API is a plus.  【職務聯絡人】 Alice Chang【E-mail】 it1@recruitexpress.com.tw【公司電話】: 02-87806811 \n  ",
			Link:     "http://www.javaworld.com.tw/jute/post/view?bid=15&id=311825",
		},
	},
	{ //JWorld@TW話題更新通知
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
您訂閱的主題:&lt;分享&gt;我自己製作的棒球遊戲.更新了!
<br>
<b>gaugyj</b>在<a href="http://www.javaworld.com.tw/jute" target=_blank>JWorld@TW</a>裡修改了文章<b>&lt;分享&gt;我自己製作的棒球遊戲.</b>以下為修改後的內容：<br>
<br>
---------- <br>
Subject: <b>&lt;分享&gt;我自己製作的棒球遊戲.</b><br>
Date:       Sat Jan 08 19:04:12 CST 2005<br><br>
<A HREF="http://p64.cc/w64/baseballGame/bg.jsp" TARGET=_blank class=ilink>http://p64.cc/w64/baseballGame/bg.jsp</A><BR>歡迎大家玩一玩 <br>
  <br>
---------- <br>
<br>
請勿直接回覆此EMAIL, 此為系統信件, 請到論壇回覆。 <br>

點選以下LINK可以直接進入文章：<a href="http://www.javaworld.com.tw/jute/post/view?bid=35&id=88339&sty=2" target=_blank>http://www.javaworld.com.tw/jute/post/view?bid=35&id=88339&sty=2</a> <br>

----------<br>
JWorld@TW http://www.javaworld.com.tw/jute
</BODY>
</HTML>`,
		&post.Post{
			Id:       88339,
			Bid:      35,
			Receiver: "parkghost",
			Sender:   "gaugyj",
			Subject:  "<分享>我自己製作的棒球遊戲.",
			Date:     getTime(time.UnixDate, "Sat Jan 08 19:04:12 CST 2005"),
			Content:  "\nhttp://p64.cc/w64/baseballGame/bg.jsp歡迎大家玩一玩 \n  ",
			Link:     "http://www.javaworld.com.tw/jute/post/view?bid=35&id=88339&sty=2",
		},
	},
	{ //JWorld@TW, 回覆通知
		`<HTML>
<HEAD>
<meta HTTP-EQUIV=Content-Type content="text/html;charset=UTF-8">
<base href="http://www.javaworld.com.tw/jute/">
<TITLE>JWorld@TW</TITLE>
<STYLE TYPE='text/css'>
  BODY      { font-family: Tahoma,Georgia; color: #000000; font-size: 12px }
  P         { font-family: Tahoma,Georgia; color: #000000; font-size: 12px } 
  TD        { font-family: Tahoma,Georgia; color: #000000; font-size: 12px }
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
  .mod         { text-decoration: none; color:#000000 }
  .mod:link    { text-decoration: none; color:#000000 }
  .mod:visited { text-decoration: none; color:#000080 }
  .mod:active  { text-decoration: none; color:#FF0000 }
  .mod:hover   { text-decoration: underline; color:#FF0000 }  
  .thd         { text-decoration: none; color:#808080 }
  .thd:link    { text-decoration: underline; color:#808080 }
  .thd:visited { text-decoration: underline; color:#808080 }
  .thd:active  { text-decoration: underline; color:#FF0000 }
  .thd:hover   { text-decoration: underline; color:#FF0000 }
  .curpage     { text-decoration: none; color:#FFFFFF; font-family: Tahoma; font-size: 9px }
  .page         { text-decoration: none; color:#003063; font-family: Tahoma; font-size: 9px }
  .page:link    { text-decoration: none; color:#003063; font-family: Tahoma; font-size: 9px }
  .page:visited { text-decoration: none; color:#003063; font-family: Tahoma; font-size: 9px }
  .page:active  { text-decoration: none; color:#FF0000; font-family: Tahoma; font-size: 9px }
  .page:hover   { text-decoration: none; color:#FF0000; font-family: Tahoma; font-size: 9px }
  .subject  { font-family: Tahoma,Georgia; font-size: 12px }
  .text     { font-family: Tahoma,Georgia; color: #000000; font-size: 12px }
  .codeStyle {  padding-right: 0.5em; margin-top: 1em; padding-left: 0.5em;  font-size: 9pt; margin-bottom: 1em; padding-bottom: 0.5em; margin-left: 0pt; padding-top: 0.5em; font-family: Courier New; background-color: #000000; color:#ffffff }
  .smalltext   { font-family: Tahoma,Georgia; color: #000000; font-size:11px }
  .verysmalltext  { font-family: Tahoma,Georgia; color: #000000; font-size:4px }          
</STYLE>
</HEAD>
<BODY BGCOLOR='#FFFFFF' TEXT='#000000' LINK='#000000' VLINK='#000080' ALINK='#FF0000' >
Hi,parkghost <br><br>

<b>ctl1690</b>在<a href="http://www.javaworld.com.tw/jute" target=_blank>JWorld@TW</a>裡回覆了您的文章<b>Re:getLastModified GMT+8問題</b>下面是回覆內容<br>
<br>
---------- <br>
Subject: <b>Re:getLastModified GMT+8問題</b><br>
Date:       Tue Dec 25 12:07:30 CST 2012<br><br>
getLastModified 最終是回傳ms的long<BR>所以我想跟DeteFormat應該沒有關係?! <br>
---------- <br>
<br>
請勿回覆此 Email, 此email由系統自動寄出 <br>

點選此連結,直接觀看回覆的文章:<a href="http://www.javaworld.com.tw/jute/post/view?bid=6&id=311305&sty=2" target=_blank>http://www.javaworld.com.tw/jute/post/view?bid=6&id=311305&sty=2</a> <br>

----------<br>
JWorld@TW http://www.javaworld.com.tw/jute
</BODY>
</HTML>`,
		&post.Post{
			Id:       311305,
			Bid:      6,
			Receiver: "parkghost",
			Sender:   "ctl1690",
			Subject:  "Re:getLastModified GMT+8問題",
			Date:     getTime(time.UnixDate, "Tue Dec 25 12:07:30 CST 2012"),
			Content:  "\ngetLastModified 最終是回傳ms的long所以我想跟DeteFormat應該沒有關係?! ",
			Link:     "http://www.javaworld.com.tw/jute/post/view?bid=6&id=311305&sty=2",
		},
	},
}

func getTime(layout, value string) time.Time {
	dt, _ := time.Parse(layout, value)
	return dt
}

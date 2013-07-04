package html

import (
	"testing"
)

func TestExtractText(t *testing.T) {
	result := ExtractText(testdata.original, BannerRemover("----------", 0, 1))
	if testdata.expected != result {
		t.Fatalf("expected cleaned text is %s, got %s", testdata.expected, result)
	}
}

var testdata = struct {
	original string
	expected string
}{`<HTML>
<HEAD>
<meta HTTP-EQUIV=Content-Type content="text/html;charset=UTF-8">
<base href="http://www.javaworld.com.tw/jute/">
<TITLE>JWorld@TW</TITLE>
<STYLE TYPE='text/css'>
  BODY      { font-family: Tahoma,Georgia; color: #000000; font-size: 12px }
  .ilink         { text-decoration: underline; color:#0000FF }
  .ilink:link    { text-decoration: underline; color:#0000FF }
  .ilink:visited { text-decoration: underline; color:#004080 }
  .ilink:active  { text-decoration: underline; color:#FF0000 }
  .ilink:hover   { text-decoration: underline; color:#FF0000 }          
</STYLE>
</HEAD>
<BODY BGCOLOR='#FFFFFF' TEXT='#000000' LINK='#000000' VLINK='#000080' ALINK='#FF0000'>
Hi,parkghost <br><br>
您訂閱的版面:Object Relational Mapping有新話題了!
<br>
<b>sai</b>在<a href="http://www.javaworld.com.tw/jute" target=_blank>JWorld@TW</a>裡發了新話題：<br>
<br>
---------- <br>
Subject: <b>Hibernate 與 spring framework 整合</b><br>
Date:       2013-06-30 17:34<br><br>
當 mybatis 與 spring framework 整合後...
	<br>
---------- <br>
<br>
請勿直接回覆此EMAIL,此為系統系統信件, 請到論壇回覆。 <br>

點選以下LINK可以直接連到此話題：<a href="http://www.javaworld.com.tw/jute/post/view?bid=41&id=314421" target=_blank>http://www.javaworld.com.tw/jute/post/view?bid=41&id=314421</a> <br>

----------<br>
JWorld@TW http://www.javaworld.com.tw/jute
</BODY>
</HTML>`, `Subject: Hibernate 與 spring framework 整合
Date:       2013-06-30 17:34
當 mybatis 與 spring framework 整合後...`,
}

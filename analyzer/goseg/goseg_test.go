package goseg

import (
	"fmt"
	"os"
	"testing"
)

const ps = string(os.PathSeparator)

func TestSeg(t *testing.T) {
	text := `卡納維洛在報告中指出，他已經成功用PEG接合過老鼠和狗的脊髓，但目前若要應用到人類的頭顱移植手術，他還必須面對2大困難：資金和道德。卡納維洛認為，若要2年內達到人頭移植的目標，他至少需要3千萬美金的資金投入研究，以拯救那些為肌肉萎縮症、癱瘓、器官衰竭，甚至癌症的重症病患`

	tk, err := NewTokenizer(".." + ps + ".." + ps + "data" + ps + "dict.txt")

	if err != nil {
		t.Fatal(err)
	}

	words := tk.Cut([]rune(text))
	for _, word := range words {
		fmt.Println(word)
	}
}

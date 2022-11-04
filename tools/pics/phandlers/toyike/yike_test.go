package toyike

import (
	"encoding/json"
	"os"
	"pc-phone-conn-go/tools/pics/pcomm"
	"testing"
)

func TestPrecreate(t *testing.T) {
	bs, err := os.ReadFile("C:/Users/Do/Downloads/金-01.jpg")
	if err != nil {
		t.Fatal(err)
	}
	yk := New(bs, "/4637251374683426/tttt.jpg", 1621090346)
	resp, err := yk.precreate()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v\n", resp)
}

func TestUploadFile(t *testing.T) {
	bs, err := os.ReadFile("C:/Users/Do/Downloads/33112314.png")
	if err != nil {
		t.Fatal(err)
	}

	yk := New(bs, "/tttt.jpg", 0)
	err = yk.UploadFile()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("已上传图集\n")
}

func TestNew(t *testing.T) {
	bs, err := os.ReadFile("C:/Users/Do/Downloads/33112314.png")
	if err != nil {
		t.Fatal(err)
	}
	yk := New(bs, "/test/tttt.jpg", 0)
	t.Logf("分块的 MD5：%v\n", yk.BlockMD5List)

	err = yk.UploadFile()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("上传完成\n")
}

func TestIDStr(t *testing.T) {
	tx := `{"4639797492847721": { "id": "4639797492847721","caption": "热裤长腿，甚是好看!![吃瓜][吃瓜]\nFrom: KgHjvbWOZ", "created": 1621697388, "urls": [ "https://wx2.sinaimg.cn/large/006AfEgvgy1gqrmf64p7ej31vv2w5b2h.jpg", "https://wx3.sinaimg.cn/large/006AfEgvgy1gqrmfk57bdj31yb31rx6w.jpg", "https://wx3.sinaimg.cn/large/006AfEgvgy1gqrmfv9l7nj31u92u87wo.jpg", "https://wx1.sinaimg.cn/large/006AfEgvgy1gqrfb055q4j31vz2x1u15.jpg", "https://wx2.sinaimg.cn/large/006AfEgvgy1gqrmgo3w6zj31xb2zukjt.jpg", "https://wx1.sinaimg.cn/large/006AfEgvgy1gqrouyzwspj31ty2wdnpk.jpg", "https://wx1.sinaimg.cn/large/006AfEgvgy1gqrmh092vkj31zg31vhe1.jpg", "https://wx2.sinaimg.cn/large/006AfEgvgy1gqrmetd621j31vo2x1qvd.jpg", "https://wx4.sinaimg.cn/large/006AfEgvgy1gqrmhe0pljj31wf2z24qy.jpg" ], "urls_m": [ "https://wx2.sinaimg.cn/orj1080/006AfEgvgy1gqrmf64p7ej31vv2w5b2h.jpg", "https://wx3.sinaimg.cn/orj1080/006AfEgvgy1gqrmfk57bdj31yb31rx6w.jpg", "https://wx3.sinaimg.cn/orj1080/006AfEgvgy1gqrmfv9l7nj31u92u87wo.jpg", "https://wx1.sinaimg.cn/orj1080/006AfEgvgy1gqrfb055q4j31vz2x1u15.jpg", "https://wx2.sinaimg.cn/orj1080/006AfEgvgy1gqrmgo3w6zj31xb2zukjt.jpg", "https://wx1.sinaimg.cn/orj1080/006AfEgvgy1gqrouyzwspj31ty2wdnpk.jpg", "https://wx1.sinaimg.cn/orj1080/006AfEgvgy1gqrmh092vkj31zg31vhe1.jpg", "https://wx2.sinaimg.cn/orj1080/006AfEgvgy1gqrmetd621j31vo2x1qvd.jpg", "https://wx4.sinaimg.cn/orj1080/006AfEgvgy1gqrmhe0pljj31wf2z24qy.jpg" ] }}`
	var tasks map[string]pcomm.Album
	err := json.Unmarshal([]byte(tx), &tasks)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v\n", tasks)

	task := tasks["4639797492847721"]
	err = Send(task)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDelAll(t *testing.T) {
	err := DelAll()
	if err != nil {
		t.Fatal(err)
	}
}

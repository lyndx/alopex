package backend

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"mime/multipart"
	"os"
	"strings"

	"alopex/app"
)

type CommonController struct{}

func init() {
	app.CJoin("common", CommonController{})
}

func (ctrl CommonController) Upload(h *app.Http) {
	h.Rep.WriteHeader(200)
	tmp, err := app.String("app").C("upload_path")
	if err != nil {
		bs, _ := json.Marshal(map[string]string{"result": "failed", "message": err.Error()})
		h.Rep.Write(bs)
		return
	}
	udir := tmp.SwitchValue(strings.HasSuffix(tmp.ToString(), "/"), strings.TrimRight(tmp.ToString(), "/"), tmp.ToString()).(app.T).ToString()
	tmp, err = app.String("app").C("page_url")
	if err != nil {
		bs, _ := json.Marshal(map[string]string{"result": "failed", "message": err.Error()})
		h.Rep.Write(bs)
		return
	}
	purl := tmp.SwitchValue(strings.HasSuffix(tmp.ToString(), "/"), strings.TrimRight(tmp.ToString(), "/"), tmp.ToString()).(app.T).ToString()
	platform, path, file := h.P("platform"), h.P("path"), h.P("file").(*multipart.FileHeader)
	source, err := file.Open()
	if err != nil {
		bs, _ := json.Marshal(map[string]string{"result": "failed", "message": err.Error()})
		h.Rep.Write(bs)
		return
	}
	defer source.Close()
	stream, err := ioutil.ReadAll(source)
	if err != nil {
		bs, _ := json.Marshal(map[string]string{"result": "failed", "message": err.Error()})
		h.Rep.Write(bs)
		return
	}
	// 文件格式
	ctmp := strings.Split(file.Header.Get("Content-Type"), "/")
	if (len(ctmp) < 1) || (ctmp[0] != "image") {
		bs, _ := json.Marshal(map[string]string{"result": "failed", "message": "非图片格式，上传失败"})
		h.Rep.Write(bs)
		return
	}
	ctype := "." + ctmp[1]
	// 文件名称（MD5）
	md5 := md5.New()
	md5.Write(stream)
	filename := hex.EncodeToString(md5.Sum(nil))
	// 文件目录
	fpath := "/" + platform.(string) + "/" + path.(string) + "/"
	_, err = os.Stat(udir + fpath)
	if (err != nil) && (os.IsNotExist(err)) {
		err = os.MkdirAll(udir+fpath, os.ModePerm)
		if err != nil {
			bs, _ := json.Marshal(map[string]string{"result": "failed", "message": "文件目录创建失败，上传失败"})
			h.Rep.Write(bs)
			return
		}
	}
	nfile, err := os.Create(udir + fpath + filename + ctype)
	if err != nil {
		bs, _ := json.Marshal(map[string]string{"result": "failed", "message": err.Error()})
		h.Rep.Write(bs)
		return
	}
	defer nfile.Close()
	// 复制
	if _, err := nfile.Write(stream); err != nil {
		bs, _ := json.Marshal(map[string]string{"result": "failed", "message": err.Error()})
		h.Rep.Write(bs)
		return
	}
	bs, _ := json.Marshal(map[string]string{"result": "ok", "url": purl + "/assets/upload" + fpath + filename + ctype})
	h.Rep.Write(bs)
}

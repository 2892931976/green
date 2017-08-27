package service

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/csv"
	"encoding/hex"
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/extrame/xls"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"github.com/inu1255/green/config"
	"github.com/inu1255/green/model"
	"github.com/tealeg/xlsx"
)

type FileService struct {
	Service
}

// @desc 上传文件
func (this *FileService) Upload(user *model.User, file *multipart.FileHeader) (interface{}, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	return FileUpload(this.Db, file.Filename, src, user)
}

// @desc 上传base64文件
// @param filename 文件名
// @param body base64字符串
func (this *FileService) UploadBase64(user *model.User, body io.ReadCloser, filename string) (interface{}, error) {
	// 上传者
	src, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}
	if filename == "" {
		filename = time.Now().Format("2006-01-02 03:04:05")
	}
	// data:image/jpg;base64,
	index := bytes.IndexByte(src, ';')
	if index >= 4 {
		s := src[:index]
		begin := bytes.IndexByte(s, '/')
		ext := string(s[begin+1:])
		if ext == "jpeg" {
			ext = "jpg"
		}
		filename += "." + ext
		src = src[index+8:]
	}
	dst := make([]byte, base64.StdEncoding.DecodedLen(len(src)))
	_, err = base64.StdEncoding.Decode(dst, src)
	if err != nil {
		return nil, err
	}
	return FileUpload(this.Db, filename, bytes.NewReader(dst), user)
}

// @desc 导出csv表格文件
func (this *FileService) ExportCsv(ctx *gin.Context, tables [][]string, filename string) (interface{}, error) {
	if filename == "" {
		filename = "表格"
	}
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Disposition", "attachment; filename="+filename+".csv")
	err := simpleWriteCsv(ctx.Writer, tables)
	return nil, err
}

// @desc 导出xlsx表格文件
func (this *FileService) ExportXlsx(ctx *gin.Context, tables [][]string, filename string) (interface{}, error) {
	if filename == "" {
		filename = "表格"
	}
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Disposition", "attachment; filename="+filename+".xlsx")
	err := simpleWriteExcel(ctx.Writer, tables)
	return nil, err
}

func NewFileServ() *FileService {
	config.Db.Sync2(new(model.File))
	return new(FileService)
}

/*****************************************************************************
 *                                 api above                                 *
 *****************************************************************************/

func FileUpload(Db *xorm.Session, filename string, src io.Reader, user *model.User) (*model.File, error) {
	var err error
	f := new(model.File)
	// 创建用户文件夹
	uid := "0"
	if user != nil {
		uid = strconv.Itoa(user.Id)
		f.OwnerId = user.Id
	}
	dir := strings.Join([]string{config.Cfg.UploadPath, uid}, "/")
	err = os.MkdirAll(dir, 0755)

	bs, err := ioutil.ReadAll(src)
	if err != nil {
		return nil, err
	}
	h := md5.New()
	h.Write(bs)
	f.MD5 = hex.EncodeToString(h.Sum(nil))
	// 保存文件
	f.Place = strings.Join([]string{dir, "/", f.MD5}, "")
	if _, err = os.Stat(f.Place); err != nil {
		err = ioutil.WriteFile(f.Place, bs, 0644)
		if err != nil {
			return nil, err
		}
	}
	//  保存文件
	f.Ext = strings.ToLower(path.Ext(filename))
	f.Filename = filename
	f.Url = f.Place
	if ok, _ := Db.Where("place=? and owner_id=?", f.Place, f.OwnerId).Get(f); ok {
		return f, nil
	}
	_, err = Db.InsertOne(f)
	return f, err
}

func simpleReadExcel(r io.Reader) ([][]string, error) {
	bs, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	// xls 文件
	xlFile, err := xls.OpenReader(bytes.NewReader(bs), "utf-8")
	if err == nil {
		return xlFile.ReadAllCells(10000), nil
	}
	// xlsx 文件
	file, err := xlsx.OpenBinary(bs)
	if err != nil {
		return nil, err
	}
	table, err := file.ToSlice()
	if err != nil {
		return nil, err
	}
	if len(table) < 1 {
		return nil, errors.New("没有Sheet")
	}
	for _, item := range table[1:] {
		table[0] = append(table[0], item...)
	}
	return table[0], nil
}

func simpleWriteExcel(f io.Writer, table [][]string) error {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		return err
	}
	for _, item := range table {
		row = sheet.AddRow()
		for _, col := range item {
			cell = row.AddCell()
			cell.Value = col
		}
	}
	if err = file.Write(f); err != nil {
		return err
	}
	return nil
}

func simpleWriteCsv(f io.Writer, table [][]string) error {
	_, err := f.Write([]byte("\xEF\xBB\xBF")) // 写入UTF-8 BOM
	if err != nil {
		return err
	}
	w := csv.NewWriter(f)
	for _, item := range table {
		err = w.Write(item)
		if err != nil {
			return err
		}
	}
	w.Flush()
	return nil
}

package service

import (
	"io/ioutil"
	"os"
	"path"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/pyihe/go-example/fsnotify/model"
	"github.com/pyihe/go-example/fsnotify/pkg"
	"github.com/pyihe/plogs"
)

type FileService struct {
	mu       sync.Mutex
	domain   string            // 文件服务器域名
	filePath string            // 配置文件路径
	watcher  *fsnotify.Watcher // watcher
	files    []*model.File     // 所有的配置文件信息
	handlers map[string]func() // 每个文件变更时对应的handler
}

func NewFileService(filePath string) *FileService {
	s := &FileService{}
	s.filePath = filePath
	s.handlers = make(map[string]func())

	s.watch()

	s.loadProvince()
	return s
}

func (f *FileService) watch() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		assert(true, err.Error())
	}
	f.watcher = watcher
	go func() {
		for {
			select {
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				plogs.Errorf("监听配置文件发生错误: %v", err)
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op != fsnotify.Write {
					break
				}
				plogs.Infof("监听到配置文件[%v]变动: %v", event.Name, event.Op)
				f.mu.Lock()
				handler := f.handlers[event.Name]
				f.mu.Unlock()
				if handler != nil {
					handler()
				}
			}
		}
	}()
}

func (f *FileService) loadFile(fileName string) error {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		plogs.Errorf("读取文件[%v]出错: %v", fileName, err)
		return err
	}

	fileInfo, err := os.Stat(fileName)
	if err != nil {
		plogs.Errorf("获取文件[%v]属性失败: %v", fileName, err)
		return err
	}

	mf := &model.File{
		Name: fileInfo.Name(),
		Size: fileInfo.Size(),
		Url:  path.Join(f.domain, fileInfo.Name()),
		MD5:  pkg.MD5(data),
	}

	if err = f.watcher.Add(fileName); err != nil {
		plogs.Errorf("Watcher添加文件失败: %v", err)
		return err
	}
	f.mu.Lock()
	find := false
	for i, info := range f.files {
		if info.Name == mf.Name {
			f.files[i] = mf
			find = true
			break
		}
	}
	if !find {
		f.files = append(f.files, mf)
	}
	f.mu.Unlock()
	return nil
}

func (f *FileService) addUpdater(key string, fn func()) {
	if fn != nil {
		f.mu.Lock()
		_, exist := f.handlers[key]
		if !exist {
			f.handlers[key] = fn
		}
		f.mu.Unlock()
	}
}

func (f *FileService) loadProvince() {
	name := path.Join(f.filePath, "province.json")
	if err := f.loadFile(name); err != nil {
		plogs.Fatalf("加载province配置失败: %v", err)
	}
	f.addUpdater(name, f.loadProvince)
}

func (f *FileService) Close() {
	if f.watcher != nil {
		f.watcher.Close()
	}
}

func (f *FileService) GetConfigFilePath() string {
	return f.filePath
}

func (f *FileService) List() (list []*model.File, err error) {
	list = make([]*model.File, len(f.files))
	copy(list, f.files)
	return
}

func assert(b bool, text string) {
	if b {
		panic(text)
	}
}

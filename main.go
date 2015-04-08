package main

import (
    "fmt"
    "net/http"
    "io/ioutil"
    "os"
    "regexp"
    "strings"
    "path/filepath"
    "./utils"
    "github.com/spf13/viper"
)

type novel struct {
    title string
    url string
}

var wenku8url string = "http://www.wenku8.cn/wap/article/packtxt.php?id=%d"

func main() {
    // http://www.wenku8.cn/wap/article/packtxt.php?id=1159
    var id int

    fmt.Print("請輸入小說編號：")
    fmt.Scan(&id)
    if (id <= 0) {
        fmt.Println("錯誤的小說編號")
        os.Exit(1)
    }
    getList(id)
}

// func getList(id int) (string, string) {
func getList(id int) {

    workingDir := getCurrentDir()
    viper.SetConfigName("config")
    viper.AddConfigPath(workingDir)
    err := viper.ReadInConfig()

    if err != nil {
        fmt.Println("No configuration file loaded - using defaults")
    }

    viper.SetDefault("SaveRootPath", "/Users/chrisliu/Downloads")
    saveRootPath := viper.GetString("SaveRootPath")

    url := fmt.Sprintf(wenku8url, id)
    content := getContent(url)
    r, _ := regexp.Compile(`《<a href="articleinfo\.php\?id=[0-9]+">(.*)?</a>`)
    title := utils.S2T(strings.TrimSpace(r.FindStringSubmatch(content)[1]))

    if (strings.Contains(title, "(")) {
        r, _ := regexp.Compile(`[^\(]+\(([^\)]+)\)`)
        title = r.FindStringSubmatch(content)[1]
    }

    r, _ = regexp.Compile(`(.*)<br/>(\r\n|\r|\n)\(`)
    categories := r.FindAllStringSubmatch(content, -1)

    r, _ = regexp.Compile(`<a href="([^"]+)">繁体`)
    urls := r.FindAllStringSubmatch(content, -1)

    savePath := fmt.Sprintf("%s/%s", saveRootPath, title)
    fmt.Printf("標題：%s\n儲存路徑：%s\n", title, savePath)

    createFolder(savePath)

    novelSavePath := ""
    stringBytes := []byte("")

    for k, v := range categories {
        fmt.Printf("%s：%s\n", v[1], strings.Replace(urls[k][1], "&amp;", "&", -1))
        novelSavePath = fmt.Sprintf("%s/%02d - %s.txt", savePath, k+1, utils.S2T(strings.TrimSpace(v[1])))
        if (fileExists(novelSavePath) == true) {
            continue
        }
        stringBytes = []byte(getContent(strings.Replace(urls[k][1], "&amp;", "&", -1)))
        ioutil.WriteFile(novelSavePath, stringBytes, 0644)
    }
}

func getContent(url string) string {

    msg(fmt.Sprintf("準備抓取 %s 的網頁內容", url))

    response, err := http.Get(url)
    if err != nil {
        fmt.Printf("%s", err)
        os.Exit(1)
    } else {
        defer response.Body.Close()
        contents, err := ioutil.ReadAll(response.Body)
        if err != nil {
            fmt.Printf("%s", err)
            os.Exit(1)
        }
        return string(contents)
    }

    return ""
}

func fileExists(path string) bool {
    if _, err := os.Stat(path); os.IsNotExist(err) {
        return false
    }
    return true
}

func createFolder(path string) {
    if _, err := os.Stat(path); os.IsNotExist(err) {
        msg(fmt.Sprintf("建立資料夾：%s", path))
        os.Mkdir(path, 0755)
    }
}

func getCurrentDir() string {
    dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
    if err != nil {
        fmt.Printf("%s", err)
    }
    return dir
}

func msg(msg string) {
    fmt.Printf("%s\n", msg)
}
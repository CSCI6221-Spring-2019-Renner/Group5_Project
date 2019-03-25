//package main
//
//import (
//	"fmt"
//	"github.com/gocolly/colly"
//	"strings"
//)
//
//func main() {
//	// 初始化 colly
//	//var result []string
//	c := colly.NewCollector(
//		// 只采集规定的域名下的内容
//		colly.AllowedDomains("stackoverflow.com"),
//	)
//	c.Limit(&colly.LimitRule{DomainGlob: "*", Delay: 3})
//
//	// 任何具有 href 属性的标签都会触发回调函数
//	// 第一个参数其实就是 jquery 风格的选择器
//	c.OnHTML("html", func(e *colly.HTMLElement) {
//		goquerySelection := e.DOM
//		link, _ := goquerySelection.Find("div[id=questions] > div > div > h3 > a[href]").Attr("href")
//		//link := e.Attr("href")
//		//fmt.Printf("Link found: %q -> \n", e.Text)
//		fmt.Printf("Link found: -> %s\n", link)
//		// 访问该网站
//		//goquerySelection := e.DOM
//		title := goquerySelection.Find("title").Text()
//		fmt.Println("title ", title)
//		if strings.HasPrefix(title, "Newest") {
//			//result = append(result, link)
//			c.Visit(e.Request.AbsoluteURL(link))
//			fmt.Println("Question link: ", link)
//		} else {
//			//fmt.Println(e.Text)
//			//result := goquerySelection.Find("body").Text();
//			//fmt.Println("result: ", result)
//		}
//	})
//	// 在请求发起之前输出 url
//	c.OnRequest(func(r *colly.Request) {
//		//fmt.Println("Visiting", r.URL.String())
//	})
//
//	//c.OnResponse(func (resp *colly.Response) {
//	//
//	//	//result = appe nd(result, string(resp.Body))
//	//	fmt.Println("body", string(resp.Body))
//	//
//	//})
//
//
//	//从以下地址开始抓起
//	//for i := 1; i <= 1000; i++ {
//	//	c.Visit("https://stackoverflow.com/questions?sort=newest&page=" + strconv.Itoa(i))
//	//}
//	c.Visit("https://stackoverflow.com/questions?sort=newest&page=1")
//
//}





package main


import (
	"fmt"
	"github.com/gocolly/colly"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"unicode/utf8"
)


type AtomicMap struct { // the map struct for multithread
	data   map[string]int
	rwLock sync.RWMutex
}

func (self *AtomicMap) Init(keywords []string) { // init the KEYWORDS LIST before using it!!
	self.data = make(map[string]int)
	for _, keyword := range keywords {
		self.data[keyword] = 0
	}
}

func (self *AtomicMap) Get(key string) (int) { // read a value
	self.rwLock.RLock()
	defer self.rwLock.RUnlock()
	val, found := self.data[key]
	if found {
		return val
	} else {
		return 0
	}
}

func (self *AtomicMap) PlusOne(key string) { // add 1 for the key
	self.rwLock.Lock()
	defer self.rwLock.Unlock()
	val, found := self.data[key]
	if found {
		self.data[key] = val + 1
	}
}

var frequencyForWord = new(AtomicMap) // the global map of record word frequencies on websites
var appearanceForWord = new(AtomicMap) // the global map of record word appeances on websites

func updateFrequencies(input string) { // analyze website with THIS!!!
	re, _ := regexp.Compile("[^0-9A-Za-z_\\-\\+#]")
	input = re.ReplaceAllString(input, " ")
	appear := make(map[string]bool)
	//fmt.Println(input)
	words := strings.Fields(input)
	for _, word := range words {
		if len(word) > utf8.UTFMax || utf8.RuneCountInString(word) > 1 {
			low_word := strings.ToLower(word)
			frequencyForWord.PlusOne(low_word)
			appear[low_word] = true
		}
	}
	for word, _ := range appear {
		appearanceForWord.PlusOne(word)
	}
}

func reportByWords(frequencyForWord map[string]int) { // show result, not necessary
	words := make([]string, 0, len(frequencyForWord))
	wordWidth, frequencyWidth := 0, 0
	for word, _ := range frequencyForWord {
		words = append(words, word)
		if width := utf8.RuneCountInString(word); width > wordWidth {
			wordWidth = width
		}
	}
	sort.Strings(words)
	gap := wordWidth - len("Word")
	fmt.Printf("Word %*s%s\n", gap, " ", "Frequency")
	for _, word := range words {
		fmt.Printf("%-*s %*d\n", wordWidth, word, frequencyWidth,
			frequencyForWord[word])
	}
}



func main() {
	keywords := []string{"javascript", "swift", "java", "c", "python", "c#", "php", "android", "c++", "sql", "objective-c", "matlab", "perl", "r", "ruby", "groovy", "go", "delphi", "visual", "assembly"}
	frequencyForWord.Init(keywords)
	appearanceForWord.Init(keywords)
	//str := "[c++] (c#) .. ,, --  python"
	//updateFrequencies(str)
	reportByWords(frequencyForWord.data)
	reportByWords(appearanceForWord.data)

	c := colly.NewCollector(
		colly.AllowedDomains("stackoverflow.com"),
	)
	c.Limit(&colly.LimitRule{DomainGlob: "*", Delay: 3})

	c.OnHTML("div[id=questions] > div > div > h3 > a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		c.Visit(e.Request.AbsoluteURL(link))
		fmt.Println("Question link: ", link)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnResponse(func (resp *colly.Response) {
		//fmt.Println("body", string(resp.Body))
		updateFrequencies(string(resp.Body))

	})

	//to find 1000 pages
	for i := 25; i <= 36; i++ {
		c.Visit("https://stackoverflow.com/questions?sort=newest&page=" + strconv.Itoa(i))
	}
	reportByWords(frequencyForWord.data)
	reportByWords(appearanceForWord.data)
}

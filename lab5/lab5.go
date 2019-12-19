package main

import (
	"fmt"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

type newsObject struct {
	URL     string
	IMGhref string
	Tag     string
	Header  string
	//Text    string
}

func (nobj newsObject) printIt() {
	println(nobj.URL)
	println(nobj.IMGhref)
	println(nobj.Tag)
	println(nobj.Header)
	//println(nobj.Text)
}

func newsObjectFromMap(m map[string]interface{}) newsObject {
	return newsObject{
		URL:     m["URL"].(string),
		IMGhref: m["IMGhref"].(string),
		Tag:     m["Tag"].(string),
		Header:  m["Header"].(string),
		//Text:    m["Text"].(string),
	}
}

func main() {
	// Start a Selenium WebDriver server instance (if one is not already
	// running).
	const (
		// These paths will be different on your system.
		seleniumPath     = "vendor/selenium-server-standalone-3.141.59.jar"
		chromeDriverPath = "vendor/chromedriver.exe"
		port             = 8080
	)
	opts := []selenium.ServiceOption{
		//selenium.StartFrameBuffer(),             // Start an X frame buffer for the browser to run in.
		selenium.ChromeDriver(chromeDriverPath), // Specify the path to GeckoDriver in order to use Firefox.
		//selenium.Output(os.Stderr),              // Output debug information to STDERR.
	}
	selenium.SetDebug(false)
	service, err := selenium.NewSeleniumService(seleniumPath, port, opts...)
	if err != nil {
		panic(err) // panic is used only as an example and is not otherwise recommended.
	}
	defer service.Stop()

	// Connect to the WebDriver instance running locally.
	caps := selenium.Capabilities{"browserName": "chrome"}
	chromeCaps := chrome.Capabilities{
		Path: "",
		Args: []string{
			"--headless", // <<<
			"--no-sandbox",
			"--log-level=3",
		},
	}
	caps.AddChrome(chromeCaps)
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		panic(err)
	}
	defer wd.Quit()
	// Navigate to the simple playground interface.
	if err := wd.Get("https://meduza.io/razbor"); err != nil {
		panic(err)
	}

	res, err := wd.ExecuteScript("let objs = []; let sections = document.querySelectorAll('article.RichBlock-root'); sections.forEach(function (s) { let IMGobj = s.firstChild.firstChild.lastChild; let content = s.lastChild.firstChild.children; objs.push({ URL: content[1].firstChild.firstChild.getAttribute('href'), IMGhref: IMGobj.getAttribute('src'), Tag: content[0].firstChild.innerText, Header:  content[1].firstChild.firstChild.firstChild.innerText }) }); return objs;", nil)
	if err != nil {
		panic(err)
	}
	vals := res.([]interface{})
	println(len(vals))
	for i := 0; i < len(vals); i++ {
		newsObjectFromMap(vals[i].(map[string]interface{})).printIt()
		println()
	}
}

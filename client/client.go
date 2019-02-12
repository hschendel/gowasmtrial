package main

import (
	"fmt"
	"syscall/js"
	"net/http"
	"html/template"
	"io/ioutil"
	"strings"
	"github.com/hschendel/wasmtrial/shared"
	"encoding/json"
)

func main() {
	fmt.Println("Hello Gophers at grandcentrix!")

	// fetch template tmpl/main.html

	mainTmpl, err := fetchTemplate("main")
	if err != nil {
		fmt.Println(err)
		return
	}

	// execute template and fill body with output

	page := struct { Heading string } {
		Heading: "Hello Gophers at grandcentrix!",
	}
	body, err := executeTemplate(mainTmpl, page)
	if err != nil {
		fmt.Println(err)
		return
	}
	setInnerHTML("body", body)

	// set up callbacks

	doc := js.Global().Get("document")
	inputName := doc.Call("getElementById", "name")
	buttonUpper := doc.Call("getElementById", "buttonUpper")
	// somehow this only works with js.PreventDefault, otherwise the page reloads
	onButtonUpperClick := js.NewEventCallback(js.PreventDefault, func(ev js.Value) {
		inputName.Set("value", strings.ToUpper(inputName.Get("value").String()))
	})
	buttonUpper.Call("addEventListener", "click", onButtonUpperClick)

	buttonRpc := doc.Call("getElementById", "buttonRpc")
	onButtonRpcClick := js.NewEventCallback(js.PreventDefault, func(ev js.Value) {
		fmt.Println("rpcButton click")
		xhrGet("/entity", func(status int, body string) {
			fmt.Println("rpcButton click callback")
			fmt.Printf("status %d and body %s\n", status, body)
			if status != 200 {
				fmt.Printf("GET /entity failed with status %d\n", status)
				return
			}
			var entity shared.SomeEntity
			if decErr := json.Unmarshal([]byte(body), &entity); decErr != nil {
				fmt.Printf("cannot decode JSON from %q: %s", body, decErr)
				return
			}
			inputName.Set("value", entity.A)
		})
	})
	buttonRpc.Call("addEventListener", "click", onButtonRpcClick)

	// now block main goroutine so Go does not exit, and the handlers are still callable

	fmt.Println("waiting forever")
	select {}
	fmt.Println("exiting")
}

func setInnerHTML(elemId, innerHTML string) {
	js.Global().Get("document").Call("getElementById", elemId).Set("innerHTML", innerHTML)
}

func xhrGet(u string, callback func(status int, body string)) {
	xhr := js.Global().Get("XMLHttpRequest").New()
	var onLoad js.Callback
	done := xhr.Get("DONE").Int()
	onLoad = js.NewEventCallback(js.PreventDefault, func(ev js.Value) {
		fmt.Println("onload called")
		target := ev.Get("target")

		if target.Get("readyState").Int() != done {
			return
		}
		fmt.Println("xhr done")
		status := target.Get("status").Int()
		responseText := target.Get("responseText").String()
		callback(status, responseText)
		onLoad.Release()
	})
	xhr.Call("open", "GET", u, true)
	xhr.Call("addEventListener", "load", onLoad)
	xhr.Call("send", nil)
	fmt.Println("xhr sent")

}

func fetchTemplate(name string) (tmpl *template.Template, err error) {
	// http.Get only seems to work during this "setup" phase, at a later stage it leads to deadlock,
	//  so we have to use XHR
	tmplUrl := fmt.Sprintf("tmpl/%s.html", name)
	resp, err := http.Get(tmplUrl)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	tmplBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	tmpl, err = template.New(name).Parse(string(tmplBytes))
	return
}

func executeTemplate(tmpl *template.Template, pageData interface{}) (output string, err error) {
	var sb strings.Builder
	if err = tmpl.Execute(&sb, pageData); err != nil {
		return
	}
	output = sb.String()
	return
}


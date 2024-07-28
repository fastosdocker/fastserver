package mw

import (
	"bytes"
	"fast/web/base"
	"io/ioutil"
	"log"
)

func Log(c *base.Ctx) {
	data, _ := c.GetRawData()
	if len(data) > 0 {
		log.Println(string(data))
		c.Request.Body = ioutil.NopCloser(bytes.NewReader(data))
	}
}

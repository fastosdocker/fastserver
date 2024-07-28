package repo

import "fast/web/base"

func GetDockerRepo(c *base.Ctx) {
	cid, _ := c.GetPostForm("cid")
	buf, ok := GetBuffer(cid)
	if !ok {
		return
	}
	c.Success(buf.data)
}

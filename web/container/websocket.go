package container

import (
	"context"
	"fast/web/base"
	"io"
	"log"
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gorilla/websocket"
)

// Ws websocket服务
func Ws(c *base.Ctx) {
	// websocket握手
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.Error(err.Error())
		return
	}
	defer conn.Close()

	// 获取容器ID或name
	container := c.Query("container")
	shell := c.Query("shell")
	user := c.Query("user")

	if container == "" {
		c.Error("参数错误,容器不能为空")
		return
	}

	if shell == "" {
		shell = "/bin/sh"
	}

	if user == "" {
		user = "root"
	}

	// 执行exec，获取到容器终端的连接
	hr, err := containerExec(container, shell, user)
	if err != nil {
		log.Println("containerExec" + err.Error())
		return
	}
	// 关闭I/O流
	defer hr.Close()
	// 退出进程
	defer func() {
		hr.Conn.Write([]byte("exit\r"))
	}()

	// 转发输入/输出至websocket
	go func() {
		wsWriterCopy(hr.Conn, conn)
	}()
	wsReaderCopy(conn, hr.Conn)
}

// containerExec 执行命令
func containerExec(container, shell, user string) (hr types.HijackedResponse, err error) {
	// 执行/bin/bash命令
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	cli.NegotiateAPIVersion(ctx)
	if err != nil {
		return types.HijackedResponse{}, err
	}
	ir, err := cli.ContainerExecCreate(ctx, container, types.ExecConfig{
		User:         user,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          []string{shell},
		Tty:          true,
	})
	if err != nil {
		return
	}

	// 附加到上面创建的/bin/bash进程中
	hr, err = cli.ContainerExecAttach(ctx, ir.ID, types.ExecStartCheck{Detach: false, Tty: true})
	if err != nil {
		return
	}
	return
}

// 将终端的输出转发到前端
func wsWriterCopy(reader io.Reader, writer *websocket.Conn) {
	buf := make([]byte, 8192)
	for {
		nr, err := reader.Read(buf)
		if err != nil {
			log.Println("reader.Read" + err.Error())
			return
		}
		if nr > 0 {
			err := writer.WriteMessage(websocket.BinaryMessage, buf[0:nr])
			if err != nil {
				log.Println(err)
				return
			}
		}
		if err != nil {
			log.Println(err)
			return
		}
	}
}

// 将前端的输入转发到终端
func wsReaderCopy(reader *websocket.Conn, writer io.Writer) {
	for {
		messageType, p, err := reader.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		if messageType == websocket.TextMessage {
			writer.Write(p)
		}
	}
}

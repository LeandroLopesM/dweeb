package main

import (
	"fmt"
	"net"
	"time"
)

func (ctx *AppCtx) submitMessage(message string) {
	ctx.mu.Lock()
    defer ctx.mu.Unlock()

    ctx.addMsgSync("SYSTEM", "Pushing " + message)

    for _,peer := range ctx.peers {
		peerDialer := net.Dialer{
			Timeout: 1 * time.Second,
		}

		conn, err := peerDialer.Dial("tcp", peer)

        if err != nil {
            ctx.addErrSync(err.Error())
			continue
		}

		defer conn.Close()

        _,err = conn.Write([]byte(message))
        if err != nil {
            ctx.addErrSync(err.Error())
			continue
        }

        // continue

        // errBlk:
    }
}

func (ctx *AppCtx) pushMsgBox(msg Message) {
    ctx.msgBox.SetText(ctx.msgBox.GetText(false) + fmt.Sprintf("\n(%s) -> %s", msg.author, msg.content))
}

func (ctx *AppCtx) addMsgSync(author, msg string) {
	ctx.messages = append(ctx.messages, Message{
		author: author,
		content: msg,
	})

	ctx.pushMsgBox(ctx.messages[len(ctx.messages)-1])
}

func (ctx *AppCtx) addMsg(author, msg string) {
    ctx.mu.Lock()
    defer ctx.mu.Unlock()

	ctx.addMsgSync(author, msg)
}

func (ctx *AppCtx) addErrSync(msg string) {
    ctx.addMsgSync("SYSTEM", fmt.Sprintf("[red]ERROR: %s[white]", msg))
}

func (ctx *AppCtx) addErr(msg string) {
    ctx.addMsg("SYSTEM", fmt.Sprintf("[red]ERROR: %s[white]", msg))
}

func receiverThread(ctx *AppCtx) {
	srv,err := net.Listen("tcp", "localhost:-1")

    if err != nil {
        ctx.addErr(err.Error())
        return
    }
    defer srv.Close()

	for {
		conn, err := srv.Accept()

		if err != nil {
			ctx.addErr(err.Error())
			continue
		}
		defer conn.Close()

		// TODO:
		// Read bytes from remote conn and parse into Message{}
		// Then display

		// var bytes []byte
		// for {
		// 	var data []byte = make([]byte, 128)
		// 	n, err := conn.Read(&data)
		// 	if err != nil {
		// 		ctx.addErr(err)
		// 		break
		// 	}

		// 	if n < 128 {
		// 		break
		// 	}

		// 	bytes = slices.Concat(bytes, data)
		// }

		// if len(data) != 0 {
		// 	var msg Message

		// 	if err := json.Unmarshal(data); err != nil {
		// 		ctx.addErr(fmt.Sprintf("Failed to receive from %s: %s", conn.RemoteAddr().String(), err))
		// 		continue
		// 	} else {

		// 	}
		// }
	}
}
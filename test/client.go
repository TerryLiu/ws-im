package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"runtime"
	"strings"
	"time"

	"nhooyr.io/websocket"
	"ws-im/test/xrand"
)

func cloneMessages(msgs map[string]struct{}) map[string]struct{} {
	msgs2 := make(map[string]struct{}, len(msgs))
	for m := range msgs {
		msgs2[m] = struct{}{}
	}
	return msgs2
}

func randMessages(n, maxMessageLength int) map[string]struct{} {
	msgs := make(map[string]struct{})
	for i := 0; i < n; i++ {
		m := randString(randInt(maxMessageLength))
		if _, ok := msgs[m]; ok {
			i--
			continue
		}
		msgs[m] = struct{}{}
	}
	return msgs
}

func assertSuccess( err error) {
	if err != nil {
		panic(err)
	}
}

type client struct {
	url string
	path string
	c   *websocket.Conn
}

func newClient(ctx context.Context, url,path string) (*client, error) {
	u := url + path
	c, _, err := websocket.Dial(ctx, u, nil)
	if err != nil {
		sprintf := fmt.Sprintf("ws拨号失败,url:%v,err:%v", u, err.Error())
		panic(sprintf)
	}

	cl := &client{
		url: url,
		path: path,
		c:   c,
	}

	return cl, nil
}

func Do(url,path string, s time.Duration ) *client  {
	ctx, _ := context.WithTimeout(context.Background(), s)

	cl, err := newClient(ctx, url,path)
	assertSuccess( err)
	go func() {
		for {
			select {
			case <- ctx.Done():
				cl.Close()
			default:
				// expMsg := randString(52)
				msg := xrand.Bytes(xrand.Int(255))
				err2 := cl.c.Write(ctx, websocket.MessageText, msg)
				if err2 == nil {
					// fmt.Printf("\r%v发送消息:%v字节\n",cl.path,len(msg))
				}


				time.Sleep(time.Second*2)
				runtime.Gosched()
			}
		}
	}()
	go func() {
		for {
			select {
			case <- ctx.Done():
				cl.Close()
			default:

				// actType, act, err2 := cl.c.Read(ctx)
				actType, _, err2 := cl.c.Read(ctx)
				if err2 == nil && actType == websocket.MessageText {
					// fmt.Printf("\r%v收到消息:%v字节\n",cl.path,len(act))
				}

				time.Sleep(time.Second*1)
				runtime.Gosched()
			}
		}
	}()
	return cl
}



func (cl *client) nextMessage() (string, error) {
	typ, b, err := cl.c.Read(context.Background())
	if err != nil {
		return "", err
	}

	if typ != websocket.MessageText {
		cl.c.Close(websocket.StatusUnsupportedData, "expected text message")
		return "", fmt.Errorf("expected text message but got %v", typ)
	}
	return string(b), nil
}

func (cl *client) Close() error {
	fmt.Printf("<<<<=========== closed of %v",cl.path)
	return cl.c.Close(websocket.StatusNormalClosure, "")
}

// randString generates a random string with length n.
func randString(n int) string {
	b := make([]byte, n)
	_, err := rand.Reader.Read(b)
	if err != nil {
		panic(fmt.Sprintf("failed to generate rand bytes: %v", err))
	}

	s := strings.ToValidUTF8(string(b), "_")
	s = strings.ReplaceAll(s, "\x00", "_")
	if len(s) > n {
		return s[:n]
	}
	if len(s) < n {
		// Pad with =
		extra := n - len(s)
		return s + strings.Repeat("=", extra)
	}
	return s
}

// randInt returns a randomly generated integer between [0, max).
func randInt(max int) int {
	x, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		panic(fmt.Sprintf("failed to get random int: %v", err))
	}
	return int(x.Int64())
}

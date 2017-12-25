package main

import (
	"fmt"
	"github.com/tsaridas/salt-event-listener-golang/zmqapi"
	"github.com/vmihailenco/msgpack"
	"io/ioutil"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	dat, err := ioutil.ReadFile("/var/cache/salt/master/.root_key")
	check(err)

	var tgt [1]string
	tgt[0] = "salt-minion-01"
	var arg [0]string
	delimiter := map[string]interface{}{"delimiter": ":", "show_timeout": true, "show_jid": false}
	load := map[string]interface{}{"tgt_type": "list", "jid": "", "cmd": "publish", "tgt": tgt, "key": dat, "arg": arg, "fun": "test.ping", "kwargs": delimiter, "ret": "", "user": "root"}
	msg := map[string]interface{}{"load": load, "enc": "clear"}

	b, err := msgpack.Marshal(msg)
	check(err)

	var verbose bool
	if len(os.Args) > 1 && os.Args[1] == "-v" {
		verbose = true
	}
	session, _ := mdapi.NewMdcli("tcp://127.0.0.1:4506", verbose)
	defer session.Close()
	s := string(b)
	ret, err := session.Send(s)
	check(err)

	if len(ret) == 0 {
		fmt.Println("Did not get a return.")
		return
	}
	byte_result := []byte(ret[0])
	var item map[string]interface{}
	err = msgpack.Unmarshal(byte_result, &item)
	check(err)
	fmt.Println(item)

}
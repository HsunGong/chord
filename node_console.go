package chord

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

const (
	helpInfo  string = "There is no help info, you are on your own"
	startInfo string = "Node start on ip: "
	pingHelp  string = "ping <IP Address>:<Port>"
	joinHelp  string = "join <IP Address>:<Port>"
	putHelp   string = "put <key> <value>"
	getHelp   string = "get <key>"
)

type NodeConsole struct {
	ipt        []string
	node       *RingNode
	currentMsg CtrlMessage
	stopSigPNI chan uint8
}

func NewNodeConsole(port int32) *NodeConsole {
	cmd := new(NodeConsole)
	cmd.node = NewNode(port)
	return cmd
}

func (c *NodeConsole) PrintHelp() {
	fmt.Println(helpInfo)
}

func (c *NodeConsole) processNodeInfo(wg *sync.WaitGroup) {
	var nodeMsg CtrlMessage
	var ok bool
	for {
		nodeMsg, ok = <-c.node.NodeMessageQueueOut
		if ok && len(nodeMsg.name) > 0 {
			PrintLog("[NODE INFO]" + nodeMsg.name[0])
		}
		if len(c.node.IfStop) > 0 {
			break
		}
	}
	wg.Done()
}

func (c *NodeConsole) processInput(ipt []string) int {
	mmsg := *NewCtrlMsg(ipt, 1)
	if len(c.node.IfStop) > 0 {
		return 2
	}
	if len(mmsg.name) != 0 && mmsg.name[0] != "" {
		switch mmsg.name[0] {
		case "create":
			c.node.UserMessageQueueIn <- mmsg
			return 1
			break
		case "ping":
			if len(mmsg.name) == 2 {
				c.node.UserMessageQueueIn <- mmsg
			} else {
				fmt.Println(pingHelp + "\n")
			}
			return 1
		case "help":
			c.PrintHelp()
			return 1
		case "join":
			if len(mmsg.name) == 2 {
				c.node.UserMessageQueueIn <- mmsg
			} else {
				fmt.Println(joinHelp)
			}
			return 1
		case "quit":
			c.node.UserMessageQueueIn <- mmsg
			return 2
		case "dump":
			c.node.UserMessageQueueIn <- mmsg
			return 1
		case "dumpsucc":
			c.node.UserMessageQueueIn <- mmsg
			return 1
		case "nf":
			c.node.UserMessageQueueIn <- mmsg
			return 1
		case "put":
			if len(mmsg.name) != 3 {
				fmt.Println(putHelp)
				return 1
			}
			c.node.UserMessageQueueIn <- mmsg
			return 1
		case "get":
			if len(mmsg.name) != 2 {
				fmt.Println(getHelp)
				return 1
			}
			c.node.UserMessageQueueIn <- mmsg
			return 1
		case "dumpdata":
			c.node.UserMessageQueueIn <- mmsg
			return 1
		default:
			return 0
		}
	}
	return 0
}

func (c *NodeConsole) Run() int {
	var wg sync.WaitGroup
	wg.Add(1)
	go c.node.Run()
	wg.Add(1)
	go c.processNodeInfo(&wg)
	reader := bufio.NewReader(os.Stdin)
	var ipt string
	for len(c.node.IfStop) == 0 {
		ipt, _ = reader.ReadString('\n')
		ipt = strings.TrimSpace(ipt)
		ipt = strings.Replace(ipt, "\n", "", -1)
		s := strings.Fields(ipt)
		switch c.processInput(s) {
		case 0:
			PrintLog("Wrong Command")
			break
		case 2:
			wg.Wait()
			return 0
		default:
			break
		}
	}
	return 0
}

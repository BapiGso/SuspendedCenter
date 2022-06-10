package blogfunc

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"log"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func GetCpuPercent() uint8 {
	percent, _ := cpu.Percent(time.Second, false) //false是所有cpu
	return uint8(percent[0])
}

func GetMemPercent() float64 {
	memInfo, _ := mem.VirtualMemory()
	return memInfo.UsedPercent
}

func GetDiskPercent() uint8 {
	parts, _ := disk.Partitions(false) //false是物理分区
	diskInfo, _ := disk.Usage(parts[0].Mountpoint)
	return uint8(diskInfo.UsedPercent)
}

func Ws(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Error during connection upgradation:", err)
		return
	}
	defer conn.Close()

	for {
		message3 := "1"
		message2 := string(time.Now().Format("2006-01-02 15:04:05"))
		conn.WriteMessage(websocket.TextMessage, []byte(message3+message2))
		time.Sleep(time.Duration(10) * time.Second)
	}
	// The event loop
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error during message reading:", err)
			break
		}
		log.Printf("Received: %s", message)
		message1 := "123"
		conn.WriteMessage(websocket.TextMessage, []byte(message1))
		err = conn.WriteMessage(messageType, message)
		if err != nil {
			log.Println("Error during message writing:", err)
			break
		}
	}

}

func Tz(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	pla, fam, ver, _ := host.PlatformInformation()
	_, _ = cpu.Info()
	physicalCnt, _ := cpu.Counts(false)
	logicalCnt, _ := cpu.Counts(true)
	fmt.Fprint(w,
		"<!DOCTYPE html>\n"+
			"<meta charset=\"utf-8\">\n"+
			"<title>Go探针</title>\n"+
			"<body>\n")
	fmt.Fprintln(w, "主机名"+pla+"-"+fam+"-"+ver+"\n")
	fmt.Fprintf(w, "物理核:%d 逻辑核:%d\n", physicalCnt, logicalCnt)
	fmt.Println(GetCpuPercent())
	fmt.Println(GetMemPercent())
	fmt.Println(GetDiskPercent())
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
	fmt.Fprint(w, "</body>")

}

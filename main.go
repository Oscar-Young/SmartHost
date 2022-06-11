package main

import (
	"fmt"
	"github.com/be-ys/pzem-004t-v3/pzem"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"os"
	"time"
)

type PowerMeterData struct {
	Voltage     float32
	Current     float32
	Power       float32
	Frequency   float32
	Energy      float32
	PowerFactor float32
}

type HostInfo struct {
	TotalCpu     []float64
	PercentCpu   []float64
	MemAvailable uint64
	MemUsed      uint64
}

func main() {

	brokerHost := ""
	brokerPassword := ""
	brokerUsername := ""

	tty := os.Args[1]
	hostname := os.Args[2]
	// Coonnect Serial
	p, err := pzem.Setup(
		pzem.Config{
			Port:  tty,
			Speed: 9600,
		})
	if err != nil {
		panic(err)
	}

	err = p.ResetEnergy()
	if err != nil {
		panic(err)
	}

	//Connect MQTT
	var broker = brokerHost
	var port = 1883
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID(hostname)
	opts.SetUsername(brokerUsername)
	opts.SetPassword(brokerPassword)
	//opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	opts.AutoReconnect = true
	opts.ConnectRetryInterval = 60 * time.Second
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	t := time.NewTicker(10 * time.Second)

	for {
		<-t.C

		// Get PowerNeterData
		powerMeterData := getPzem004Info(p)
		//powerMeterDataJson, _ := json.Marshal(&powerMeterData)

		voltage := fmt.Sprintf("%f", powerMeterData.Voltage)
		current := fmt.Sprintf("%f", powerMeterData.Current)
		power := fmt.Sprintf("%f", powerMeterData.Power)
		powerFactor := fmt.Sprintf("%f", powerMeterData.PowerFactor)
		energy := fmt.Sprintf("%f", powerMeterData.Energy)
		frequency := fmt.Sprintf("%f", powerMeterData.Frequency)
		ms := time.Now().UTC().Format(time.RFC3339)

		powerMeterDataCsv := ms + "," + voltage + "," + current + "," + power + "," + powerFactor + "," + energy + "," + frequency
		//meterDataJson, _ := json.Marshal(&powerMeterData)

		fmt.Println(time.Now().UTC().Format(time.RFC3339) + " " + string(powerMeterDataCsv))

		//client.Publish("v1/devices/me/telemetry", 0, false, meterDataJson).Wait()
		client.Publish("iotdata/powermeter/"+hostname, 0, false, powerMeterDataCsv).Wait()

		// Get CPU Info & Mem Info
		virtualMemoryStat := getMemInfo()

		hostInfo := HostInfo{
			TotalCpu:     getTotalCpuInfo(),
			PercentCpu:   getPercentCpuInfo(),
			MemUsed:      virtualMemoryStat.Used,
			MemAvailable: virtualMemoryStat.Available,
		}
		//hostInfoJson, _ := json.Marshal(&hostInfo)
		totalCpu := fmt.Sprintf("%f", hostInfo.TotalCpu[0])
		memUsed := fmt.Sprintf("%d", hostInfo.MemUsed)
		memAvailable := fmt.Sprintf("%d", hostInfo.MemAvailable)

		ms = time.Now().UTC().Format(time.RFC3339)
		hostInfoCsv := ms + "," + totalCpu + "," + memUsed + "," + memAvailable
		fmt.Println(time.Now().UTC().Format(time.RFC3339) + " " + hostInfoCsv)

		client.Publish("iotdata/host/"+hostname, 0, false, hostInfoCsv).Wait()
	}

}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	t := time.Now().UTC()
	fmt.Println(t.Format(time.RFC3339) + " Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	t := time.Now().UTC()
	fmt.Println(t.Format(time.RFC3339) + " Connect lost")
}

func getPzem004Info(p pzem.Probe) PowerMeterData {

	voltage, err := p.Voltage()
	if err != nil {
		panic(err)
	}
	intensity, err := p.Intensity()
	if err != nil {
		panic(err)
	}
	power, err := p.Power()
	if err != nil {
		panic(err)
	}
	frequency, err := p.Frequency()
	if err != nil {
		panic(err)
	}
	energy, err := p.Energy()
	if err != nil {
		panic(err)
	}
	powerFactor, err := p.PowerFactor()
	if err != nil {
		panic(err)
	}

	powerMeterData := PowerMeterData{
		Voltage:     voltage,
		Current:     intensity,
		Power:       power,
		Frequency:   frequency,
		Energy:      energy,
		PowerFactor: powerFactor,
	}

	return powerMeterData
}

func getMemInfo() *mem.VirtualMemoryStat {
	v, _ := mem.VirtualMemory()
	return v
}

func getTotalCpuInfo() []float64 {
	totalPercent, _ := cpu.Percent(3*time.Second, false)
	return totalPercent
}

func getPercentCpuInfo() []float64 {
	perPercents, _ := cpu.Percent(3*time.Second, true)
	return perPercents
}
func getNetInfo() []net.IOCountersStat {
	info, _ := net.IOCounters(true)
	return info
}

func getDiskInfo() *disk.UsageStat {
	info, _ := disk.Usage("C:/")
	return info
}

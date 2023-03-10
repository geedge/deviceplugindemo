package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"

	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

const (
	// SysfsDevices   = "/sys/bus/pci/devices"
	SysfsDevices   = "/root/demo"
	MgmtPrefix     = "/dev/xclmgmt"
	UserPrefix     = "/dev/dri"
	QdmaPrefix     = "/dev/xfpga"
	QDMASTR        = "dma.qdma.u"
	UserPFKeyword  = "drm"
	DRMSTR         = "renderD"
	ROMSTR         = "rom"
	SNSTR          = "xmc.u."
	DSAverFile     = "VBNV"
	DSAtsFile      = "timestamp"
	InstanceFile   = "instance"
	MgmtFile       = "mgmt_pf"
	UserFile       = "user_pf"
	VendorFile     = "vendor"
	DeviceFile     = "device"
	SNFile         = "serial_num"
	VtShell        = "xilinx_u30"
	U30CommonShell = "ama_u30"
	XilinxVendorID = "0x10ee"
	ADVANTECH_ID   = "0x13fe"
	AWS_ID         = "0x1d0f"
	AristaVendorID = "0x3475"
)

type Pairs struct {
	Mgmt string
	User string
	Qdma string
}

type Device struct {
	index    string
	deviceID string //devid of the user pf
	Healthy  string
	SN       string
}

func GetInstance(DBDF string) (string, error) {
	strArray := strings.Split(DBDF, ":")
	domain, err := strconv.ParseUint(strArray[0], 16, 16)
	if err != nil {
		return "", fmt.Errorf("strconv failed: %s\n", strArray[0])
	}
	bus, err := strconv.ParseUint(strArray[1], 16, 8)
	if err != nil {
		return "", fmt.Errorf("strconv failed: %s\n", strArray[1])
	}
	strArray = strings.Split(strArray[2], ".")
	dev, err := strconv.ParseUint(strArray[0], 16, 8)
	if err != nil {
		return "", fmt.Errorf("strconv failed: %s\n", strArray[0])
	}
	fc, err := strconv.ParseUint(strArray[1], 16, 8)
	if err != nil {
		return "", fmt.Errorf("strconv failed: %s\n", strArray[1])
	}
	ret := domain*65536 + bus*256 + dev*8 + fc
	return strconv.FormatUint(ret, 10), nil
}

func GetFileNameFromPrefix(dir string, prefix string) (string, error) {
	userFiles, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf("Can't read folder %s", dir)
	}
	for _, userFile := range userFiles {
		fname := userFile.Name()

		if !strings.HasPrefix(fname, prefix) {
			continue
		}
		return fname, nil
	}
	return "", nil
}

func GetFileContent(file string) (string, error) {
	if buf, err := ioutil.ReadFile(file); err != nil {
		return "", fmt.Errorf("Can't read file %s", file)
	} else {
		ret := strings.Trim(string(buf), "\n")
		return ret, nil
	}
}

func FileExist(fname string) bool {
	if _, err := os.Stat(fname); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func IsMgmtPf(pciID string) bool {
	fname := path.Join(SysfsDevices, pciID, MgmtFile)
	return FileExist(fname)
}

func IsUserPf(pciID string) bool {
	fname := path.Join(SysfsDevices, pciID, UserFile)
	return FileExist(fname)
}

func GetDevices() ([]Device, error) {
	var devices []Device
	pciFiles, err := ioutil.ReadDir(SysfsDevices)
	if err != nil {
		return nil, fmt.Errorf("Can't read folder %s", SysfsDevices)
	}

	for _, pciFile := range pciFiles {
		pciID := pciFile.Name()
		// get device id
		fname := path.Join(SysfsDevices, pciID, DeviceFile)
		content, err := GetFileContent(fname)
		if err != nil {
			return nil, err
		}
		devid := content
		fname = path.Join(SysfsDevices, pciID, "SN")
		content, err = GetFileContent(fname)
		if err != nil {
			return nil, err
		}
		sn := content
		healthy := pluginapi.Healthy
		devices = append(devices, Device{
			index:    strconv.Itoa(len(devices) + 1),
			deviceID: devid,
			Healthy:  healthy,
			SN:       sn,
		})
	}
	return devices, nil
}

// func main() {
// 	devices, err := GetDevices()
// 	if err != nil {
// 		fmt.Printf("%s !!!\n", err)
// 		return
// 	}
// 	//SNFolder, err := GetFileNameFromPrefix(path.Join(SysfsDevices, "0000:e3:00.1"), SNSTR)
// 	//fname := path.Join(SysfsDevices, "0000:e3:00.1", SNFolder, SNFile)
// 	//content, err := GetFileContent(fname)
// 	//SN := content
// 	//fmt.Printf("SN: %v \n", SN)
// 	for _, device := range devices {
// 		fmt.Printf("Device: %v \n", device)
// 		fmt.Printf("ID: %s  \n\n", device.deviceID)
// 		fmt.Printf("SN: %s  \n\n", device.SN)
// 		fmt.Printf("Heathy: %s  \n\n", device.Healthy)
// 	}
// }

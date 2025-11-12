package render

import (
	"fmt"
	"io"
	"os"
	"strings"

	c "github.com/chigaji/diag-cli/internal/collectors"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/viper"
)

func Print(w io.Writer, data any) error {
	out := viper.GetString("output")
	switch strings.ToLower(out) {
	case "json":
		return writeJSON(w, data)
	case "yaml":
		return writeYaml(w, data)
	default:
		return writeTable(w, data)
	}
}

func writeTable(w io.Writer, data any) error {
	t := table.NewWriter()
	t.SetOutputMirror(w)
	noColor := viper.GetBool("ui.no_color")
	if noColor {
		t.Style().Color = table.ColorOptions{}
	}

	switch v := data.(type) {
	case c.SystemData:
		t.AppendHeader(table.Row{"HOST", "OS", "PLATFORM", "KERNEL", "UPTIME(s)"})
		t.AppendRow(table.Row{v.Host.Hostname, v.Host.OS, v.Host.Platform, v.Host.Kernel, v.Host.Uptime})
		t.Render()
		fmt.Fprintln(w)
		t = table.NewWriter()
		t.SetOutputMirror(w)
		if noColor {
			t.Style().Color = table.ColorOptions{}
		}
		t.AppendHeader(table.Row{"CPU MODEL", "CORES", "MHZ", "USER", "SYSTEM", "IDLE"})
		for _, cinfo := range v.CPU {
			t.AppendRow(table.Row{cinfo.Model, cinfo.Cores, int(cinfo.Mhz), int(cinfo.User), int(cinfo.System), int(cinfo.Idle)})
		}
		t.Render()
		fmt.Fprintln(w)
		t = table.NewWriter()
		t.SetOutputMirror(w)
		if noColor {
			t.Style().Color = table.ColorOptions{}
		}
		t.AppendHeader(table.Row{"DEVICE", "MOUNT", "FSTYPE", "USED%", "USED", "FREE"})
		for _, d := range v.Disks {
			t.AppendRow(table.Row{d.Device, d.Mountpoint, d.Fstype, int(d.UsedPercent), d.Used, d.Free})
		}
		t.Render()
	case c.NetworkData:
		t.AppendHeader(table.Row{"IFACE", "HWADDR", "MTU", "ADDRS"})
		for _, i := range v.Interfaces {
			t.AppendRow(table.Row{i.Name, i.HardwareAddr, i.MTU, strings.Join(i.Addrs, ", ")})
		}
		t.Render()
		fmt.Fprintln(w)
		t = table.NewWriter()
		t.SetOutputMirror(w)
		if noColor {
			t.Style().Color = table.ColorOptions{}
		}
		t.AppendHeader(table.Row{"IFACE", "BYTES SENT", "BYTES RECV", "PKTS SENT", "PKTS RECV"})
		for _, io := range v.IO {
			t.AppendRow(table.Row{io.Name, io.BytesSent, io.BytesRecv, io.PacketsSent, io.PacketsRecv})
		}
		t.Render()
	case c.ProcData:
		t.AppendHeader(table.Row{"PID", "NAME", "CPU%", "MEM%"})
		for _, r := range v.Rows {
			t.AppendRow(table.Row{r.PID, r.Name, fmt.Sprintf("%.1f", r.CPU), fmt.Sprintf("%.1f", r.MEM)})
		}
		t.Render()
	default:
		fmt.Fprintln(os.Stderr, "unsupported type for table rendering")
	}
	return nil
}

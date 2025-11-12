package gopsutil

import (
	"context"
	"sort"

	c "github.com/chigaji/diag-cli/internal/collectors"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"
	"github.com/shirou/gopsutil/v4/process"
	"github.com/shirou/gopsutil/v4/sensors"
)

type psCollector struct{}

func New() *psCollector {
	return &psCollector{}
}

func (p *psCollector) System(ctx context.Context, all bool) (c.SystemData, error) {
	var out c.SystemData
	if hi, err := host.InfoWithContext(ctx); err == nil {
		out.Host = c.HostInfo{
			Hostname: hi.Hostname,
			OS:       hi.OS,
			Platform: hi.Platform,
			Kernel:   hi.KernelVersion,
			Uptime:   hi.Uptime,
			BootTime: hi.BootTime,
		}
	}
	if times, err := cpu.TimesWithContext(ctx, false); err == nil && len(times) > 0 {
		if infos, err2 := cpu.InfoWithContext(ctx); err2 == nil && len(infos) > 0 {
			row := c.CPUInfo{
				Model:  infos[0].ModelName,
				Cores:  int(infos[0].Cores),
				Mhz:    infos[0].Mhz,
				User:   times[0].User,
				System: times[0].System,
				Idle:   times[0].Idle,
			}
			out.CPU = []c.CPUInfo{row}
		}
	}
	if vm, err := mem.VirtualMemoryWithContext(ctx); err == nil {
		out.Memory = c.MemoryInfo{
			Total:       vm.Total,
			Used:        vm.Used,
			Free:        vm.Free,
			UsedPercent: uint64(vm.UsedPercent),
		}
	}
	if parts, err := disk.PartitionsWithContext(ctx, true); err == nil {
		for _, p := range parts {
			if st, err2 := disk.UsageWithContext(ctx, p.Mountpoint); err2 == nil {
				out.Disks = append(out.Disks, c.DiskInfo{
					Device:      p.Device,
					Mountpoint:  p.Mountpoint,
					Fstype:      p.Fstype,
					Total:       st.Total,
					Used:        st.Used,
					Free:        st.Free,
					UsedPercent: st.UsedPercent,
				})
			}
		}
	}
	if all {
		if temps, err := sensors.SensorsTemperatures(); err == nil {
			for _, t := range temps {
				out.Temperatures = append(out.Temperatures, c.Temperature{
					Sensor:  t.SensorKey,
					Celsius: t.Temperature,
				})
			}
		}
		// if temps, err := host.SensorsTemperaturesWithContext(ctx); err == nil {
		// 	for _, t := range temps {
		// 		out.Temperatures = append(out.Temperatures, c.Temperature{
		// 			Sensor:  t.SensorKey,
		// 			Celsius: t.Temperature,
		// 		})
		// 	}
		// }
	}
	return out, nil
}

func (p *psCollector) Network(ctx context.Context, iface string) (c.NetworkData, error) {
	var out c.NetworkData
	if list, err := net.InterfacesWithContext(ctx); err == nil {
		for _, it := range list {
			if iface != "" && it.Name != iface {
				continue
			}
			addrs := make([]string, 0, len(it.Addrs))
			for _, a := range it.Addrs {
				addrs = append(addrs, a.Addr)
			}
			out.Interfaces = append(out.Interfaces, c.NetIF{
				Name:         it.Name,
				HardwareAddr: it.HardwareAddr,
				MTU:          it.MTU,
				Addrs:        addrs,
			})
		}
	}
	if io, err := net.IOCountersWithContext(ctx, true); err == nil {
		for _, i := range io {
			if iface != "" && i.Name != iface {
				continue
			}
			out.IO = append(out.IO, c.NetIO{
				Name:        i.Name,
				BytesSent:   i.BytesSent,
				BytesRecv:   i.BytesRecv,
				PacketsSent: i.PacketsSent,
				PacketsRecv: i.PacketsRecv,
			})
		}
	}
	return out, nil
}

func (p *psCollector) Processes(ctx context.Context, top int, sortBy string) (c.ProcData, error) {
	procs, _ := process.ProcessesWithContext(ctx)
	rows := make([]c.ProcessRow, 0, len(procs))
	for _, pr := range procs {
		name, _ := pr.NameWithContext(ctx)
		cpuPct, _ := pr.CPUPercentWithContext(ctx)
		memPct, _ := pr.MemoryPercentWithContext(ctx)
		rows = append(rows, c.ProcessRow{
			PID:  pr.Pid,
			Name: name,
			CPU:  cpuPct,
			MEM:  memPct,
		})
	}
	sort.Slice(rows, func(i, j int) bool {
		if sortBy == "mem" {
			return rows[i].MEM > rows[j].MEM
		}
		return rows[i].CPU > rows[j].CPU
	})
	if top > 0 && top < len(rows) {
		rows = rows[:top]
	}
	return c.ProcData{Rows: rows}, nil
}

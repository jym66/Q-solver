package audio

import (
	"Q-Solver/pkg/logger"
	"sync"

	"github.com/gen2brain/malgo"
)

var (
	targetTime = 40 // 攒多ms发送一次
	buffsize   = 16000 * 2 * targetTime / 1000
)

// LoopbackCapture 扬声器音频采集 (实时)
type LoopbackCapture struct {
	ctx         *malgo.AllocatedContext
	device      *malgo.Device
	onData      func([]byte)
	mu          sync.Mutex
	running     bool
	audioBuffer []byte
}

// NewLoopbackCapture 创建 Loopback 采集器
func NewLoopbackCapture(onData func([]byte)) (*LoopbackCapture, error) {
	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, nil)
	if err != nil {
		return nil, err
	}
	return &LoopbackCapture{
		ctx:         ctx,
		onData:      onData,
		audioBuffer: make([]byte, 0,buffsize*2), //多申请一些防止频繁扩容
	}, nil
}

// Start 开始采集扬声器输出
func (c *LoopbackCapture) Start() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.running {
		return nil
	}

	deviceConfig := malgo.DefaultDeviceConfig(malgo.Loopback)
	deviceConfig.Capture.Format = malgo.FormatS16
	deviceConfig.Capture.Channels = 1
	deviceConfig.SampleRate = 16000

	// 数据回调 - 实时发送
	onRecv := func(_, pInput []byte, frameCount uint32) {
		if len(pInput) > 0 && c.onData != nil {
			c.audioBuffer = append(c.audioBuffer, pInput...)
			for len(c.audioBuffer) >= buffsize {
				dataCopy := make([]byte, buffsize)
				copy(dataCopy, c.audioBuffer[:buffsize])
				c.onData(dataCopy)
				c.audioBuffer = c.audioBuffer[buffsize:]
			}
		}
	}

	callbacks := malgo.DeviceCallbacks{Data: onRecv}

	device, err := malgo.InitDevice(c.ctx.Context, deviceConfig, callbacks)
	if err != nil {
		return err
	}

	if err := device.Start(); err != nil {
		device.Uninit()
		return err
	}

	c.device = device
	c.running = true

	logger.Println("Loopback 采集已启动")
	return nil
}

// Stop 停止采集
func (c *LoopbackCapture) Stop() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.device != nil {
		c.device.Stop()
		c.device.Uninit()
		c.device = nil
	}
	c.running = false
	logger.Println("Loopback 采集已停止")
}

// Close 释放资源
func (c *LoopbackCapture) Close() {
	c.Stop()
	if c.ctx != nil {
		_ = c.ctx.Uninit()
		c.ctx.Free()
		c.ctx = nil
	}
}

// IsRunning 是否正在采集
func (c *LoopbackCapture) IsRunning() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.running
}

package volume

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/k8shuginn/event-collector/exporter"
	"github.com/k8shuginn/event-collector/pkg/logger"
	"go.uber.org/zap"
)

var _ exporter.Exporter = (*VolumeExporter)(nil)

type VolumeExporter struct {
	currentFile *os.File
	fileName    string
	filePath    string

	dataChan chan []byte

	currentCount int
	maxFileSize  int
	maxFileCount int
}

// NewVolumeExporter volume exporter 생성
// fileName: 파일명
// filePath: 파일 경로
// opts: exporter option
func NewVolumeExporter(fileName, filePath string, opts ...Option) (*VolumeExporter, error) {
	c := fromOptions(opts...)
	e := &VolumeExporter{
		fileName:     fileName,
		filePath:     filePath,
		dataChan:     make(chan []byte, 200),
		maxFileSize:  c.maxFileSize,
		maxFileCount: c.maxFileCount,
	}

	// directory 생성
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			return nil, fmt.Errorf("failed to create directory: %w", err)
		}
	}

	files, err := e.getSortFileList()
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		e.currentCount = 0
	} else {
		lastCount := extractNumber(files[len(files)-1])
		e.currentCount = lastCount
	}

	return e, nil
}

// Start exporter 시작
func (e *VolumeExporter) Start(ctx context.Context, wg *sync.WaitGroup) error {
	logger.Info("[volume exporter] started")
	defer func() {
		e.shutdown()
		logger.Info("[volume exporter] stopped")
		wg.Done()
	}()

	wg.Add(1)
	for {
		select {
		case <-ctx.Done():
			return nil
		case data := <-e.dataChan:
			if err := e.writeData(data); err != nil {
				logger.Error("[volume exporter] failed to write data", zap.Error(err))
			}
		}
	}
}

// Write 데이터 기록
func (e *VolumeExporter) Write(data []byte) {
	e.dataChan <- data
}

// shutdown exporter 종료
func (e *VolumeExporter) shutdown() {
	if e.currentFile != nil {
		e.currentFile.Close()
	}

	close(e.dataChan)
}

// getFileList 파일 목록을 반환
func (e *VolumeExporter) getFileList() ([]string, error) {
	entries, err := os.ReadDir(e.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var files []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if strings.Contains(entry.Name(), e.fileName) {
			files = append(files, entry.Name())
		}
	}

	return files, nil
}

// getSortFileList 파일 목록을 정렬해서 반환
func (e *VolumeExporter) getSortFileList() ([]string, error) {
	files, err := e.getFileList()
	if err != nil {
		return nil, err
	}

	SortByNumericSuffix(files)

	return files, nil
}

// removeFile 파일 삭제
func (e *VolumeExporter) removeFile(file string) error {
	if err := os.Remove(filepath.Join(e.filePath, file)); err != nil {
		return err
	}

	return nil
}

// checkAndRemove 파일 개수가 maxFileCount를 넘으면 파일을 삭제
func (e *VolumeExporter) checkAndRemove() error {
	files, err := e.getSortFileList()
	if err != nil {
		return err
	}

	if len(files) >= e.maxFileCount {
		count := len(files) - e.maxFileCount + 1
		for i := 0; i < count; i++ {
			if err := e.removeFile(files[i]); err != nil {
				logger.Error("[volume exporter] failed to remove file", zap.Error(err), zap.String("file", files[i]))
			}
		}
	}

	return nil
}

// writeData 파일에 데이터를 기록
func (e *VolumeExporter) writeData(data []byte) error {
	// 파일이 없으면 새로 생성
	if e.currentFile == nil {
		file, err := os.Create(filepath.Join(e.filePath, fmt.Sprintf("%s_%d", e.fileName, e.currentCount)))
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		e.currentFile = file
	}

	// 파일 사이즈 체크
	// 파일의 최대 개수가 넘으면 파일을 닫고 새로 생성
	fInfo, _ := e.currentFile.Stat()
	if fInfo.Size() > int64(e.maxFileSize) {
		if err := e.checkAndRemove(); err != nil {
			return fmt.Errorf("failed to check and remove: %w", err)
		}

		if err := e.currentFile.Close(); err != nil {
			return fmt.Errorf("failed to close file: %w", err)
		}

		e.currentCount++
		file, err := os.Create(filepath.Join(e.filePath, fmt.Sprintf("%s_%d", e.fileName, e.currentCount)))
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		e.currentFile = file
	}

	// 데이터 기록
	if _, err := e.currentFile.Write(data); err != nil {
		return fmt.Errorf("failed to write data: %w", err)
	}

	logger.Debug("[volume exporter] data written", zap.String("file", e.currentFile.Name()), zap.Int("size", len(data)))

	return nil
}

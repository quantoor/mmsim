package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

//
type LimitOrder struct {
	Id              string
	Symbol          string
	Timestamp       int64
	Price           float64
	Size            float64
	IsBuy           bool
	FilledTimestamp int64
	RealizedPnl     float64
}

func NewLimitOrder(symbol string, price float64, amount float64, isBuy bool) *LimitOrder {
	return &LimitOrder{
		Symbol:          symbol,
		Price:           price,
		Size:            amount,
		IsBuy:           isBuy,
		FilledTimestamp: 0,
		RealizedPnl:     0,
	}
}

func (l *LimitOrder) ComputeRealizedPnl(position *Position) {
	realizedPnl := 0.0

	if position.Size < 0 && l.IsBuy {
		realizedPnl = (position.AveragePrice - l.Price) * l.Size

	} else if position.Size > 0 && !l.IsBuy {
		realizedPnl = (l.Price - position.AveragePrice) * l.Size
	}

	l.RealizedPnl = realizedPnl
}

func (l *LimitOrder) String() string {
	side := "sell"
	if l.IsBuy {
		side = "buy"
	}
	return fmt.Sprintf("symbol %s, price %f, size %f, side %s, realized pnl %f", l.Symbol, l.Price, l.Size, side, l.RealizedPnl)
}

//
type Position struct {
	Symbol       string
	Size         float64
	AveragePrice float64
}

//
type Tick struct {
	Timestamp int64
	Price     float64
}

type TickData struct {
	Symbol string
	Data   []Tick
}

func (d *TickData) readFromFile(path string) {
	file, err := os.Open(path)
	failOnError(err, fmt.Sprintf("Could not open file %s", path))
	defer file.Close()

	d.Data = make([]Tick, 0)
	scanner := bufio.NewScanner(file)
	scanner.Scan() // skip header
	for scanner.Scan() {
		line := scanner.Text()
		values := strings.Split(line, ",")

		timestamp, err := strconv.ParseFloat(values[0], 64)
		failOnError(err, "Error parsing timestamp")
		price, err := strconv.ParseFloat(values[1], 64)
		failOnError(err, "Error parsing price")

		d.Data = append(d.Data, Tick{Timestamp: int64(timestamp), Price: price})
	}
	err = scanner.Err()
	failOnError(err, "Scanner error")
}

func (d *TickData) append(symbolData *TickData) {
	d.Data = append(d.Data, symbolData.Data...)
}

func (d *TickData) writeToFile(filePath string) error {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	datawriter := bufio.NewWriter(file)
	_, err = datawriter.WriteString("Timestamp,Price\n")
	if err != nil {
		log.Error("Error writing header of file")
	}

	for _, di := range d.Data {
		_, err = datawriter.WriteString(fmt.Sprintf("%d,%f\n", di.Timestamp, di.Price))
		if err != nil {
			log.Error("Error writing status to result file")
		}
	}
	datawriter.Flush()
	return nil
}

func NewSymbolDataFromProcessedFile(filePath string) *TickData {
	tickData := &TickData{}
	tickData.readFromFile(filePath)
	return tickData
}

func NewSymbolDataFromTickDataFolder(folderPath string) *TickData {
	var files []string
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".csv" {
			files = append(files, path)
		}
		return nil
	})
	failOnError(err, fmt.Sprintf("Could not get files in folder %s", folderPath))

	tickData := &TickData{}
	for _, file := range files {
		tickData.append(processTickData(file))
	}

	tickData.writeToFile("../datasets/DOGE_1s.csv")
	return tickData
}

func processTickData(filePath string) *TickData {
	log.Infof("Processing tick data file: %s", filePath)

	file, err := os.Open(filePath)
	failOnError(err, fmt.Sprintf("Could not open file %s", filePath))
	defer file.Close()

	symbolData := &TickData{}
	symbolData.Data = make([]Tick, 0)
	scanner := bufio.NewScanner(file)
	scanner.Scan() // skip header

	lastTimestamp := 0
	for scanner.Scan() {
		line := scanner.Text()
		values := strings.Split(line, ",")

		timestamp, err := strconv.Atoi(values[5])
		timestamp /= 1000 // convert to seconds
		failOnError(err, "Error parsing timestamp")
		price, err := strconv.ParseFloat(values[2], 64)
		failOnError(err, "Error parsing price")

		if timestamp >= (lastTimestamp + 1) { // append new value only if 1 second has passed
			lastTimestamp = timestamp
			symbolData.Data = append(symbolData.Data, Tick{Timestamp: int64(lastTimestamp), Price: price})
		}
	}
	err = scanner.Err()
	failOnError(err, "Scanner error")
	return symbolData
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panic(msg)
	}
}

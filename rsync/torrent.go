package main
/*
import (
	"crypto/md5"
	"fmt"
	"hash/adler32"
	"math"
	"os"
)

type chunklist struct {
	adler uint32
	md5   [16]byte
}
func main() {
	var finalresult []chunklist
	fileToBeChunked := "./mandala/mandala.war"
	file, err := os.Open(fileToBeChunked)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()
	fileInfo, _ := file.Stat()
	var fileSize int64 = fileInfo.Size()
	const fileChunk = 1 * (1 << 10) // 1 kb, change this to your requirement
	// calculate total number of parts the file will be chunked into
	totalPartsNum := uint64(math.Ceil(float64(fileSize) / float64(fileChunk)))
	fmt.Printf("Splitting to %d pieces.\n", totalPartsNum)
	for i := uint64(0); i < totalPartsNum; i++ {
		partSize := int(math.Min(fileChunk, float64(fileSize-int64(i*fileChunk))))
		partBuffer := make([]byte, partSize)
		file.Read(partBuffer)
		//fmt.Println("Split to : ", i)
		result := md5.Sum(partBuffer)
		//fmt.Printf("%x\n", result)
		resultt := adler32.Checksum(partBuffer)
		finalresult = append(finalresult, chunklist{resultt, result})
		//fmt.Printf("%x\n", resultt)
	}
	fmt.Println(finalresult[324])
}
*/

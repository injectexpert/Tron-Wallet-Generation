package main

import (
	"bufio"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/antlabs/strsim"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mr-tron/base58"
)

var runCount uint64 = 0
var mutex sync.Mutex
var wg sync.WaitGroup

func s256(s []byte) []byte {
	h := sha256.New()
	h.Write(s)
	bs := h.Sum(nil)
	return bs
}

func IsValuableAddress(str string, length int) bool {
	sub := str[len(str)-length : len(str)] //后缀
	if strings.Count(sub, string(sub[0])) != length {
		return false
	} else {
		return true
	}

}

func PrintStatus() {
	var currentCount uint64 = 0
	start := time.Now()
	log.Println("程序开始运行")

	if logInterval == 0 {
		return
	}

	for {
		elapsed := time.Now().Sub(start).Seconds()

		if currentCount != 0 {
			log.Printf("计算地址次数:%d 算力:%d/s 已运行时间:%vs", runCount, (runCount-currentCount)/uint64(logInterval), elapsed)
		}
		currentCount = runCount
		time.Sleep(time.Duration(logInterval) * time.Second)
	}
}

// 计算能力测试
func ComputeAbilityTest() {
	for i := 0; i < 10000000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			privateKey, err := crypto.GenerateKey()
			if err != nil {
				log.Fatal(err)
			}
			publicKey := privateKey.Public()
			publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
			if !ok {
				log.Fatal("error casting public key to ECDSA")
			}

			address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
			address = "41" + address[2:]
			//fmt.Println("address hex: ", address)
			addb, _ := hex.DecodeString(address)
			hash1 := s256(s256(addb))
			secret := hash1[:4]
			for _, v := range secret {
				addb = append(addb, v)
			}
			mutex.Lock()
			runCount++
			mutex.Unlock()

		}(i)

	}
	wg.Wait()
	time.Sleep(3 * time.Second)
}

// 暴力生成符合条件的前缀或后缀
func BruteAddress(sPrefix string, sSuffix string, generateNum int) {
	var i int = 0
	var isEnd bool = false
	var totalGenerateCount int = 0
	var mutexPrint sync.Mutex
	for {
		if isEnd {
			break
		}

		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			privateKey, err := crypto.GenerateKey()
			if err != nil {
				log.Fatal(err)
			}
			publicKey := privateKey.Public()
			publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
			if !ok {
				log.Fatal("error casting public key to ECDSA")
			}

			address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
			address = "41" + address[2:]

			addb, _ := hex.DecodeString(address)
			hash1 := s256(s256(addb))
			secret := hash1[:4]
			for _, v := range secret {
				addb = append(addb, v)
			}

			base58Address := base58.Encode(addb)

			if (sPrefix != "" && strings.HasPrefix(base58Address, sPrefix) || sPrefix == "") && (sSuffix != "" && strings.HasSuffix(base58Address, sSuffix) || sSuffix == "") {
				privateKeyBytes := crypto.FromECDSA(privateKey)

				mutexPrint.Lock()
				fmt.Println()
				fmt.Println("私钥:", hexutil.Encode(privateKeyBytes)[2:])
				fmt.Println("地址:", base58Address)
				fmt.Println()
				totalGenerateCount++
				if totalGenerateCount >= generateNum {
					isEnd = true
				}
				mutexPrint.Unlock()
			}

			mutex.Lock()
			runCount++
			mutex.Unlock()

		}(i)

	}
	wg.Wait()
	time.Sleep(3 * time.Second)

}

func SimilarAddressGenerate(sourceAddress string) {
	var maxSimilarity float64 = 0
	var i int = 0
	var isEnd bool = false
	var mutexPrint sync.Mutex
	for {
		if isEnd {
			break
		}

		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			privateKey, err := crypto.GenerateKey()
			if err != nil {
				log.Fatal(err)
			}

			publicKey := privateKey.Public()
			publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
			if !ok {
				log.Fatal("生成私钥出错!")
			}

			address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
			address = "41" + address[2:]

			addb, _ := hex.DecodeString(address)
			hash1 := s256(s256(addb))
			secret := hash1[:4]
			for _, v := range secret {
				addb = append(addb, v)
			}

			base58Address := base58.Encode(addb)
			currentSimilarity := strsim.Compare(sourceAddress, base58Address, strsim.Cosine())
			if maxSimilarity < currentSimilarity {
				maxSimilarity = currentSimilarity
				privateKeyBytes := crypto.FromECDSA(privateKey)

				mutexPrint.Lock()
				fmt.Println("相似度:", currentSimilarity*100)
				fmt.Println("私钥:", hexutil.Encode(privateKeyBytes)[2:])
				fmt.Println("地址:", base58Address)
				fmt.Println()
				mutexPrint.Unlock()
			}

			mutex.Lock()
			runCount++
			mutex.Unlock()

		}(i)

	}
	wg.Wait()
	time.Sleep(3 * time.Second)
}

func GenerateAddress(num int, saveFilename string) {
	file, err := os.OpenFile(saveFilename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal("文件打开失败", err)
	}
	//及时关闭file句柄
	defer file.Close()

	write := bufio.NewWriter(file)

	for i := 0; i < num; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			privateKey, err := crypto.GenerateKey()
			if err != nil {
				log.Fatal(err)
			}
			//crypto.FromECDSA(privateKey)
			privateKeyBytes := crypto.FromECDSA(privateKey)
			//fmt.Println("privateKey:", hexutil.Encode(privateKeyBytes)[2:])
			publicKey := privateKey.Public()
			publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
			if !ok {
				log.Fatal("error casting public key to ECDSA")
			}

			address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
			address = "41" + address[2:]
			//fmt.Println("address hex: ", address)
			addb, _ := hex.DecodeString(address)
			hash1 := s256(s256(addb))
			secret := hash1[:4]
			for _, v := range secret {
				addb = append(addb, v)
			}
			mutex.Lock()
			tmp := ""
			tmp += "地址:"
			base58Address := base58.Encode(addb)
			tmp += base58Address
			tmp += "\n私钥:"
			tmp += hexutil.Encode(privateKeyBytes)[2:]
			tmp += "\n\n"
			write.WriteString(tmp)
			write.Flush() //Flush将缓存的文件真正写入到文件中
			runCount++
			mutex.Unlock()

		}(i)

	}
	wg.Wait()
	time.Sleep(3 * time.Second)
}

func ValuableAddressGenerate(length int, generateNum int) { //靓号生成
	var i int = 0
	var isEnd bool = false
	var totalGenerateCount int = 0
	var mutexPrint sync.Mutex
	for {
		if isEnd {
			break
		}

		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			privateKey, err := crypto.GenerateKey()
			if err != nil {
				log.Fatal(err)
			}

			publicKey := privateKey.Public()
			publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
			if !ok {
				log.Fatal("error casting public key to ECDSA")
			}

			address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
			address = "41" + address[2:]

			addb, _ := hex.DecodeString(address)
			hash1 := s256(s256(addb))
			secret := hash1[:4]
			for _, v := range secret {
				addb = append(addb, v)
			}

			base58Address := base58.Encode(addb)

			if IsValuableAddress(base58Address, length) {
				privateKeyBytes := crypto.FromECDSA(privateKey)

				mutexPrint.Lock()
				fmt.Println()
				fmt.Println("私钥:", hexutil.Encode(privateKeyBytes)[2:])
				fmt.Println("地址:", base58Address)
				fmt.Println()
				totalGenerateCount++
				if totalGenerateCount >= generateNum {
					isEnd = true
				}
				mutexPrint.Unlock()
			}

			mutex.Lock()
			runCount++
			mutex.Unlock()

		}(i)

	}
	wg.Wait()
	time.Sleep(3 * time.Second)

}

var (
	runMode                         int
	logInterval                     int
	help                            bool
	argsPrefix                      string //前缀
	argsSuffix                      string //后缀
	argsSourceAddress               string
	argsSaveFilename                string //保存文件名
	argsAddressGenerateNum          int
	argsValuableAddressSuffixLength int //靓号后缀长度
	argsCpuCoreUse                  int //用多少分子一的cpu
)

func init() {
	flag.BoolVar(&help, "h", false, "显示帮助")
	flag.IntVar(&runMode, "m", 0, "模式 1算力测试 2前缀后缀爆破 3相似度爆破 4批量地址生成 5靓号生成")
	flag.StringVar(&argsPrefix, "pf", "", "生成地址前缀，T开头")
	flag.StringVar(&argsSuffix, "sf", "", "生成地址后缀")
	flag.IntVar(&logInterval, "t", 1, "日志输出间隔 为0时不输出")
	flag.StringVar(&argsSourceAddress, "a", "", "相似度爆破模式中的源地址")
	flag.StringVar(&argsSaveFilename, "o", "地址.txt", "批量地址生成中保存生成地址的文件名")
	flag.IntVar(&argsAddressGenerateNum, "n", 1, "要生成的个数")
	flag.IntVar(&argsValuableAddressSuffixLength, "l", 5, "靓号后缀长度")
	flag.IntVar(&argsCpuCoreUse, "c", 0, "用多少个cpu.为0则不限制")
}

func main() {
	fmt.Println("Telegram:@inject_exp")
	//fmt.Println("测试版本")
	fmt.Println("本程序全程不联网，但为了您的安全性考虑，建议在断网的虚拟机中运行!")
	fmt.Println()
	flag.Parse()
	if help {
		flag.Usage()
	} else {
		go PrintStatus()
		cpuTotal := runtime.NumCPU()
		fmt.Println("当前总共cpu核心个数", cpuTotal)
		if argsCpuCoreUse != 0 {
			fmt.Println("使用cpu核心个数为", argsCpuCoreUse)
			runtime.GOMAXPROCS(argsCpuCoreUse)
		}

		switch runMode {
		case 0:
			flag.Usage()
			break
		case 1:
			fmt.Println("当前模式:算力测试")
			ComputeAbilityTest()
			break
		case 2:
			if argsPrefix == "" && argsSuffix == "" {
				flag.Usage()
				os.Exit(0)
			}

			if argsPrefix != "" {
				fmt.Println("前缀为:", argsPrefix)
			}

			if argsSuffix != "" {
				fmt.Println("后缀为:", argsSuffix)
			}

			BruteAddress(argsPrefix, argsSuffix, argsAddressGenerateNum)
			break
		case 3:
			if argsSourceAddress == "" {
				flag.Usage()
				os.Exit(0)
			}
			fmt.Println("当前模式:相似度爆破")
			SimilarAddressGenerate(argsSourceAddress)
			break
		case 4:
			GenerateAddress(argsAddressGenerateNum, argsSaveFilename)
			break
		case 5:
			fmt.Println("当前模式:靓号生成")
			fmt.Println("靓号后缀长度:", argsValuableAddressSuffixLength)
			ValuableAddressGenerate(argsValuableAddressSuffixLength, argsAddressGenerateNum)
			break
		default:
			flag.Usage()
			break
		}

	}
}

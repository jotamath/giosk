package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"giosk/internal/engine"
	"giosk/internal/models"
	"giosk/internal/utils"
)

const (
	Version = "1.0.0"
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Cyan    = "\033[36m"
)

func main() {
	// Configuração do Helper
	flag.Usage = func() {
		printBanner()
		fmt.Fprintf(os.Stderr, "Giosk %s - High-performance Port Scanner for SOC/Pentest\n\n", Version)
		fmt.Fprintf(os.Stderr, "USAGE:\n  giosk -t <target> [options]\n\n")
		fmt.Fprintf(os.Stderr, "OPTIONS:\n")
		fmt.Fprintf(os.Stderr, "  -t string      Target IP or CIDR range\n")
		fmt.Fprintf(os.Stderr, "  -p string      Port range (e.g. '80,443' or '1-1024') (default \"1-1024\")\n")
		fmt.Fprintf(os.Stderr, "  -c int         Number of concurrent workers (default 100)\n")
		fmt.Fprintf(os.Stderr, "  -to duration   Timeout per connection (default 500ms)\n")
		fmt.Fprintf(os.Stderr, "  -v             Verbose mode: shows ALL connection attempts (Open/Closed)\n")
		fmt.Fprintf(os.Stderr, "  -o string      Output file to save results (e.g. 'report.txt')\n")
		fmt.Fprintf(os.Stderr, "  -version       Show version information\n")
		fmt.Fprintf(os.Stderr, "\nEXAMPLES:\n")
		fmt.Fprintf(os.Stderr, "  giosk -t 192.168.1.1 -p 80,443 -o results.txt\n")
		fmt.Fprintf(os.Stderr, "  giosk -t 10.0.0.0/24 -v -o verbose_scan.txt\n\n")
	}

	target := flag.String("t", "", "")
	concurrency := flag.Int("c", 100, "")
	timeout := flag.Duration("to", 500*time.Millisecond, "")
	portsRange := flag.String("p", "1-1024", "")
	outputFile := flag.String("o", "", "")
	showVersion := flag.Bool("version", false, "")
	verbose := flag.Bool("v", false, "")

	flag.Parse()

	if *showVersion {
		fmt.Printf("Giosk version %s\n", Version)
		os.Exit(0)
	}

	if *target == "" {
		flag.Usage()
		os.Exit(1)
	}

	// 1. Parsing de IPs e Portas
	ips, err := utils.ParseCIDR(*target)
	if err != nil {
		fmt.Printf("%s[!] Error parsing target: %v%s\n", Red, err, Reset)
		os.Exit(1)
	}

	ports, err := parsePorts(*portsRange)
	if err != nil {
		fmt.Printf("%s[!] Error parsing ports: %v%s\n", Red, err, Reset)
		os.Exit(1)
	}

	// 2. Preparação do Arquivo de Relatório (se solicitado)
	var file *os.File
	if *outputFile != "" {
		file, err = os.Create(*outputFile)
		if err != nil {
			fmt.Printf("%s[!] Could not create file: %v%s\n", Red, err, Reset)
			os.Exit(1)
		}
		defer file.Close()
		fmt.Fprintf(file, "GIOSK SCAN REPORT - %s\n", time.Now().Format(time.RFC1123))
		fmt.Fprintf(file, "Target: %s | Ports: %s\n", *target, *portsRange)
		fmt.Fprintf(file, strings.Repeat("-", 50)+"\n")
	}

	// 3. Canais e Workers
	jobs := make(chan int, 100)
	results := make(chan models.ScanResults)
	var wg sync.WaitGroup

	printBanner()
	fmt.Printf("%s[*] Target:%s %s | %s[*] Hosts:%s %d | %s[*] Concurrency:%s %d\n",
		Cyan, Reset, *target, Cyan, Reset, len(ips), Cyan, Reset, *concurrency)
	if *outputFile != "" {
		fmt.Printf("%s[*] Saving results to:%s %s\n", Cyan, Reset, *outputFile)
	}
	fmt.Println(strings.Repeat("-", 65))

	start := time.Now()

	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			engine.Worker(ips, jobs, results, *timeout)
		}()
	}

	go func() {
		for _, p := range ports {
			jobs <- p
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	// 4. Consumer TUI e Persistência
	openPorts := 0
	for res := range results {
		statusLine := ""

		if res.Status == "open" {
			openPorts++
			statusLine = fmt.Sprintf("[+] %-15s | Port: %-5d | Banner: %q", res.IP, res.Port, res.Banner)
			// Imprime na tela limpando a linha do scanner
			fmt.Printf("\r\033[K%s%s%s\n", Green, statusLine, Reset)
			// Salva no arquivo sempre que for aberta
			if file != nil {
				fmt.Fprintln(file, statusLine)
			}
		} else {
			// Lógica Verbose: Mostrar fechadas apenas se -v estiver ativo
			if *verbose {
				statusLine = fmt.Sprintf("[-] %-15s | Port: %-5d | Status: %s", res.IP, res.Port, res.Status)
				fmt.Printf("\r\033[K%s%s%s\n", Red, statusLine, Reset)
				if file != nil {
					fmt.Fprintln(file, statusLine)
				}
			} else {
				// Feedback dinâmico para o usuário não achar que travou
				fmt.Printf("\r%s[*] Scanning... IP: %s | Port: %d%s", Yellow, res.IP, res.Port, Reset)
			}
		}
	}

	duration := time.Since(start).Round(time.Millisecond)
	finalMsg := fmt.Sprintf("\n[✓] Scan complete in %v. %d open ports found.", duration, openPorts)
	fmt.Printf("%s%s%s\n", Green, finalMsg, Reset)

	if file != nil {
		fmt.Fprintln(file, strings.Repeat("-", 50))
		fmt.Fprintln(file, finalMsg)
	}
}

func parsePorts(pRange string) ([]int, error) {
	var ports []int
	pRange = strings.TrimSpace(pRange)
	if strings.Contains(pRange, "-") {
		parts := strings.Split(pRange, "-")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid range format")
		}
		start, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
		end, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
		for i := start; i <= end; i++ {
			ports = append(ports, i)
		}
	} else {
		for _, p := range strings.Split(pRange, ",") {
			p = strings.TrimSpace(p)
			if p == "" {
				continue
			}
			val, err := strconv.Atoi(p)
			if err == nil {
				ports = append(ports, val)
			}
		}
	}
	return ports, nil
}

func printBanner() {
	banner := `
   /$$$$$$  /$$$$$$  /$$$$$$   /$$$$$$  /$$   /$$
  /$$__  $$|_  $$_/ /$$__  $$ /$$__  $$| $$  /$$/
 | $$  \__/  | $$  | $$  \ $$| $$  \__/| $$ /$$/ 
 | $$ /$$$$  | $$  | $$  | $$|  $$$$$$ | $$$$$/  
 | $$|_  $$  | $$  | $$  | $$ \____  $$| $$  $$  
 | $$  \ $$  | $$  | $$  | $$ /$$  \ $$| $$\  $$ 
 |  $$$$$$/ /$$$$$$|  $$$$$$/|  $$$$$$/| $$ \  $$
  \______/ |______/ \______/  \______/ |__/  \__/`
	fmt.Printf("%s%s%s\n", Cyan, banner, Reset)
}

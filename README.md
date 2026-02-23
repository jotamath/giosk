# GIOSK: Fast Service Discovery & Recon Engine

```bash
      /$$$$$$  /$$$$$$  /$$$$$$   /$$$$$$  /$$   /$$
     /$$__  $$|_  $$_/ /$$__  $$ /$$__  $$| $$  /$$/
    | $$  \__/  | $$  | $$  \ $$| $$  \__/| $$ /$$/ 
    | $$ /$$$$  | $$  | $$  | $$|  $$$$$$ | $$$$$/  
    | $$|_  $$  | $$  | $$  | $$ \____  $$| $$  $$  
    | $$  \ $$  | $$  | $$  | $$ /$$  \ $$| $$\  $$ 
    |  $$$$$$/ /$$$$$$|  $$$$$$/|  $$$$$$/| $$ \  $$
     \______/ |______/ \______/  \______/ |__/  \__/                                                    
```

<div align="center">
	<img width="50" src="https://raw.githubusercontent.com/marwin1991/profile-technology-icons/refs/heads/main/icons/go.png" alt="Go" title="Go"/>
</div>

Giosk é um scanner de portas de alta performance desenvolvido em Go, focado em reconhecimento ativo e identificação de banners de serviço. Diferente de scanners genéricos, o Giosk prioriza a velocidade de execução via worker pools e a extração imediata de metadados de serviços (Banner Grabbing).


## Por que o Giosk?
No cenário de Segurança da Informação, o tempo entre o scan e a exploração é crítico. O Giosk foi construído para resolver três problemas principais:
1. **Velocidade de Recon**: Minimiza o overhead de rede usando concorrência nativa (Goroutines).
2. **Identificação de Shadow IT**: Localiza serviços rodando em portas não convencionais através da captura de banners.
3. **Redução de Ruído**: Interface limpa que separa o que é "lixo" (portas fechadas) do que é superfície de ataque real.

## Implicações de Portas Abertas
Uma porta aberta é um ponto de entrada. No mundo do Pentesting e SOC, elas significam:
- **Exposição de Superfície**: Cada porta é um socket ouvindo conexões, potencialmente vulnerável a exploits de buffer overflow ou falhas de lógica.
- **Divulgação de Versão**: Banners expõem a versão exata do software (ex: OpenSSH 9.6p1), permitindo a busca por CVEs (Common Vulnerabilities and Exposures) específicas.
- **Vetor de Autenticação**: Portas como 22 (SSH) ou 3389 (RDP) são alvos constantes de ataques de força bruta.

## Instalação (FUTURO)
```bash
sudo snap install giosk
```
### Configuração de Permissões

Por ser um scanner de rede distribuído sob isolamento de segurança (*strict confinement*), você precisa conceder permissões para que o Giosk acesse a rede:

```bash
sudo snap connect giosk:network
sudo snap connect giosk:network-bind

```

---

## Exemplo de uso
```bash
giosk -t 192.168.1.1 -p 1-1024 -v -o scan_report.txt
```
<img width="804" height="363" alt="image" src="https://github.com/user-attachments/assets/239ddea3-2514-47ef-bb1f-be166cda928e" />

Seu README está com uma base técnica excelente! Como o Giosk agora é um **Snap**, você pode substituir aquela seção de "Instalação (FUTURO)" por algo real e profissional.

Aqui estão as adições sugeridas para o seu `README.md`, incluindo a seção de instalação e os comandos necessários para que o scanner funcione com permissões totais no Linux:

---



## 🛠️ Opções e Comandos

| Flag | Descrição | Exemplo |
| --- | --- | --- |
| `-t` | Alvo (IP ou Hostname) | `-t 192.168.1.1` |
| `-p` | Range de portas | `-p 1-1024` ou `-p 80,443` |
| `-c` | Concorrência (Workers) | `-c 500` |
| `-to` | Timeout de conexão | `-to 300ms` |
| `-o` | Arquivo de saída | `-o result.txt` |
| `-v` | Modo Verboso | `-v` |

---

## 📖 Exemplo de Uso Rápido

**Scan padrão em portas comuns:**

```bash
giosk -t scanme.nmap.org -p 1-1000 -v

```

**Scan agressivo em rede local:**

```bash
giosk -t 192.168.15.1 -p 1-65535 -c 1000 -to 200ms
```

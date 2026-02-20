# GIOSK: Fast Service Discovery & Recon Engine

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
apt install -y giosk
``` 

## Exemplo de uso
```bash
giosk -t 192.168.1.1 -p 1-1024 -v -o scan_report.txt
```
<img width="804" height="363" alt="image" src="https://github.com/user-attachments/assets/239ddea3-2514-47ef-bb1f-be166cda928e" />


## Workflow de Análise
1. **Fase de Descoberta**: Use o Giosk com um range amplo (`-p 1-65535`) em um único alvo para mapear serviços ocultos.
2. **Fase de Identificação**: Analise os Banners retornados. Se o banner estiver vazio, tente aumentar o `-to` (timeout), pois serviços lentos podem demorar a responder.
3. **Exportação**: Utilize a flag `-o` para manter logs de auditoria.

## Dicas de Performance
- Para redes locais: `-c 1000 -to 200ms`
- Para redes externas (Internet): `-c 100 -to 500ms` (evita descarte de pacotes por latência).

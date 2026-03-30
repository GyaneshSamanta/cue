package builtins

import "github.com/GyaneshSamanta/cue/internal/macro"

func registerNetworkMacros() {
	reg(&macro.Macro{
		Name: "port-kill", Category: "network", Dangerous: true,
		Description: "Kill the process using a specific port",
		Commands: []macro.Step{
			{OS: "linux", Command: `fuser -k $1/tcp 2>/dev/null || lsof -ti:$1 | xargs kill -9`},
			{OS: "darwin", Command: `lsof -ti:$1 | xargs kill -9`},
			{OS: "windows", Command: `FOR /F "tokens=5" %a IN ('netstat -aon ^| findstr :$1') DO taskkill /F /PID %a`},
		},
		Explanation: `
✔ Done. Process on the specified port was killed.
─────────────────────────────────────────────────────
The port is now free. If a service was running on it,
you may need to restart it.
─────────────────────────────────────────────────────`,
		BuiltIn: true,
	})

	reg(&macro.Macro{
		Name: "port-check", Category: "network",
		Description: "Check what's running on a specific port",
		Commands: []macro.Step{
			{OS: "linux", Command: "ss -tlnp | grep :$1 || echo 'Port $1 is free'"},
			{OS: "darwin", Command: "lsof -i :$1 || echo 'Port $1 is free'"},
			{OS: "windows", Command: `netstat -aon | findstr :$1 || echo Port $1 is free`},
		},
		Explanation: `
✔ Port check complete.
─────────────────────────────────────────────────────
If a process is shown, it's currently using that port.
Use 'cue port-kill <port>' to free it.
─────────────────────────────────────────────────────`,
		BuiltIn: true,
	})

	reg(&macro.Macro{
		Name: "kill-port", Category: "network", Dangerous: true,
		Description: "Interactive port-to-process killer",
		Commands: []macro.Step{
			{OS: "linux", Command: `ss -tlnp | head -20 && echo "Use: cue port-kill <port>"`},
			{OS: "darwin", Command: `lsof -iTCP -sTCP:LISTEN | head -20 && echo "Use: cue port-kill <port>"`},
			{OS: "windows", Command: `netstat -aon | findstr LISTENING | head -20 & echo Use: cue port-kill <port>`},
		},
		Explanation: `
✔ Listed all listening ports and their processes.
─────────────────────────────────────────────────────
Use 'cue port-kill <port>' to kill a specific process.
─────────────────────────────────────────────────────`,
		BuiltIn: true,
	})

	reg(&macro.Macro{
		Name: "ip-info", Category: "network",
		Description: "Show local and public IP addresses",
		Commands: []macro.Step{
			{OS: "linux", Command: `echo "Local:" && hostname -I && echo "Public:" && curl -s ifconfig.me 2>/dev/null || echo "No internet"`},
			{OS: "darwin", Command: `echo "Local:" && ipconfig getifaddr en0 && echo "Public:" && curl -s ifconfig.me 2>/dev/null || echo "No internet"`},
			{OS: "windows", Command: `echo Local: & ipconfig | findstr /i "IPv4" & echo Public: & curl -s ifconfig.me 2>nul || echo No internet`},
		},
		Explanation: `
✔ IP information displayed.
─────────────────────────────────────────────────────
Local IP: Your address on this network.
Public IP: How the internet sees you (via ifconfig.me).
─────────────────────────────────────────────────────`,
		BuiltIn: true,
	})

	reg(&macro.Macro{
		Name: "cert-check", Category: "security",
		Description: "Check SSL certificate expiry for a domain",
		Commands: []macro.Step{
			{OS: "linux", Command: `echo | openssl s_client -servername $1 -connect $1:443 2>/dev/null | openssl x509 -noout -dates`},
			{OS: "darwin", Command: `echo | openssl s_client -servername $1 -connect $1:443 2>/dev/null | openssl x509 -noout -dates`},
			{OS: "windows", Command: `echo | openssl s_client -servername $1 -connect $1:443 2>nul | openssl x509 -noout -dates`},
		},
		Explanation: `
✔ Certificate dates displayed.
─────────────────────────────────────────────────────
notBefore = when the cert became valid
notAfter  = when the cert expires
Renew before notAfter to avoid downtime.
─────────────────────────────────────────────────────`,
		BuiltIn: true,
	})
}

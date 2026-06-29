package shellio

import "io"

const Prompt = "$ "

func redrawLine(w io.Writer, line string) {
	io.WriteString(w, "\r\033[K")
	io.WriteString(w, Prompt)
	io.WriteString(w, line)
}

func writePrompt(w io.Writer, rawMode bool) {
	if rawMode {
		io.WriteString(w, "\r"+Prompt)
		return
	}
	io.WriteString(w, Prompt)
}

func ringBell(w io.Writer) {
	io.WriteString(w, "\a")
}

func writeCRLF(w io.Writer) {
	io.WriteString(w, "\r\n")
}

func writeBackspace(w io.Writer) {
	io.WriteString(w, "\b \b")
}

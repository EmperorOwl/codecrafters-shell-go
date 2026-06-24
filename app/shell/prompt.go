package shell

import "io"

const Prompt = "$ "

func WritePrompt(w io.Writer) {
	io.WriteString(w, Prompt)
}

import { Terminal as XTerminal } from "@xterm/xterm"
import { useEffect, useRef, useState } from "react"
import "@xterm/xterm/css/xterm.css"
import { sendCommand } from "./socket/socketHandler";

export const term = new XTerminal({
  cursorBlink: true,
  allowProposedApi: true,
  convertEol: true, 
});

export default function Terminal() {
  const terminalRef = useRef()
  const isRendered = useRef(false);
  let lastChar = ''
  let command = ""

  useEffect(() => {
  });

  useEffect(() => {
    if (isRendered.current) return;
    isRendered.current = true;

    term.open(terminalRef.current);
    term.write("$ ")

    term.onData((e) => {
      if (lastChar === '\r' || lastChar === '\n') {
        term.write('\n$ ')
      }
      switch (e) {
        case '\u0003': // Ctrl+C
          term.write('^C');
          prompt(term);
          break;
        case '\r': // Enter
          term.write('\n')
          sendCommand(command)
          command = '';
          break;
        case '\u007F': // Backspace (DEL)
          // Do not delete the prompt
          if (term._core.buffer.x > 2) {
            term.write('\b \b');
            if (command.length > 0) {
              command = command.substr(0, command.length - 1);
            }
          }
          break;
        default: // Print all other characters for demo
          if (e >= String.fromCharCode(0x20) && e <= String.fromCharCode(0x7E) || e >= '\u00a0') {
            command += e;
            term.write(e);
          }
      }
      lastChar = e
    })

  }, [terminalRef])

  return (
    <div ref={terminalRef} id="terminal"></div>
  )
}

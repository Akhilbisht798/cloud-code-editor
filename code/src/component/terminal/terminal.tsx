import { Terminal as Xterminal } from "@xterm/xterm";
import "@xterm/xterm/css/xterm.css";
import { FC, useEffect, useRef, useContext } from "react";
import { sendCommand } from "../socket/socketHandler";
import SocketProvider from "@/context/socketContextProvider";

export const term = new Xterminal({
  cursorBlink: true,
  allowProposedApi: true,
  convertEol: true,
  theme: {
    background: "#1e1e1e",
  },
  fontSize: 16
});

const Terminal: FC = () => {
  const terminalRef = useRef<HTMLDivElement | null>(null);
  const isRendered = useRef(false);
  const { ws } = useContext(SocketProvider);
  let lastChar = "";
  let command = "";

  useEffect(() => {
    console.log("rendered");
    if (isRendered.current || !terminalRef.current || !ws) return;
    isRendered.current = true;

    term.open(terminalRef.current);
    term.writeln("Connected to cloud machine.");
    term.writeln("Happy Coding.");
    term.write('\x1b[1;32m' + "$ ");

    term.onData((e) => {
      if (lastChar === "\r" || lastChar === "\n") {
        term.write('\n$ ');
      }
      switch (e) {
        case "\u0003": // Ctrl+C 
          term.write("^C");
          prompt(term);
          break;
        case "\r": // Enter
          term.write("\n");
          sendCommand(ws, command);
          command = "";
          break;
        case "\u007F": // Backspace (DEL)
          // Do not delete the prompt
          if (term._core.buffer.x > 2) {
            term.write("\b \b");
            if (command.length > 0) {
              command = command.substr(0, command.length - 1);
            }
          }
          break;
        default: // Print all other characters for demo
          if (
            (e >= String.fromCharCode(0x20) &&
              e <= String.fromCharCode(0x7e)) ||
            e >= "\u00a0"
          ) {
            command += e;
            term.write(e);
          }
      }
      lastChar = e;
    });
  }, [terminalRef, ws]);

  return <div ref={terminalRef} id="terminal"></div>;
};

export default Terminal;

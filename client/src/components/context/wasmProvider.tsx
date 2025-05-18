import React, { useEffect, useState } from "react";
import { IWasmContext } from "../../interfaces/IWasmContext";
import { WasmContext } from "./wasmContext";
import { loadWasm } from "../../utils/wasm";
export const WasmProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [wasmAPI, setWasmAPI] = useState<IWasmContext>({});

  useEffect(() => {

    const wasmLoader = async () => {

      await loadWasm();

      setWasmAPI({
        joinSession: () => (window as any).joinSession(),
        leaveSession: () => (window as any).leaveSession(),
        getAllPeers: () => (window as any).getAllPeers(),
        getPeerMode: () => (window as any).getPeerMode(),
      });
    };

    wasmLoader();
  }, []);

  return (
    <WasmContext.Provider value={wasmAPI}>{children}</WasmContext.Provider>
  );
};

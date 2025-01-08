export const loadWasm = async () => {
  const go = new window.Go();
  const response = await fetch("../../public/wasm/crdt_conn.wasm");
  const buffer = await response.arrayBuffer();
  const wasm = await WebAssembly.instantiate(buffer, go.importObject);
  go.run(wasm.instance);
};

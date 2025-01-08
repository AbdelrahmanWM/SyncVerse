declare module "*.wasm" {
    const value: any;
    export default value;
  }
  
  interface Go {
    new(): {
    importObject: any;
    run(instance: WebAssembly.Instance): void;
    }
  }
  
  interface Window {
    // crdt_conn.wasm
    connectToSignalingServer:()=>void
    disconnectFromSignalingServer:()=>void
    getAllPeerIDs:()=>void
    newPeerConnection:()=>void
    sendToAll:()=>void
    clearLog:()=>void
    sindIdentifySelfMessage:()=>void
    Go: Go;
  }
  
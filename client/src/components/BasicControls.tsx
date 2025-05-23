import { useWasm } from "./context/wasmContext";

export function BasicControls(){
    const wasmAPI = useWasm();
    return (
      <div>
        <button onClick={wasmAPI.joinSession}>Join Session</button>
        <button onClick={wasmAPI.leaveSession}>Leave Session</button>
        <button onClick={wasmAPI.getAllPeers}>Get Connected Peers</button>
        <button onClick={wasmAPI.getPeerMode}>Get Peer Mode</button>
      </div>
    );
}
export interface IWasmContext {
  joinSession?: () => void;
  leaveSession?: () => void;
  getAllPeers?: () => void;
  getPeerMode?: () => void;
}

// import "./App.css";
import { loadWasm } from "./utils/wasm";
import { useEffect } from "react";

function App() {
  useEffect(() => {
    loadWasm();
  }, []);
  return (
    <div>
      <h1>WebRTC Mesh Network Client</h1>
      <div style={{ display: "flex", justifyContent: "space-between" }}>
        <div style={{ flex: "1" }}>
          <h2>Signaling Server Client</h2>
          <div style={{ marginBottom: "10px" }}>
            <button onClick={() => window.connectToSignalingServer()}>
              Connect to Signaling Server
            </button>
            <button onClick={()=>window.disconnectFromSignalingServer()}>
              Disconnect
            </button>
            <button onClick={()=>window.getAllPeerIDs()}>Get All Peer IDs</button>
            <button onClick={()=>window.sindIdentifySelfMessage()}>Identify self</button>

          </div>
        </div>

        <div
          style={{
            flex: "1",
            backgroundColor: "#f9f9f9",
            border: "1px solid #ddd",
            padding: "10px",
          }}
        >
          <h3>Notes</h3>
          <p>
            <strong>Step 1:</strong> Connect to the signaling server to exchange
            WebRTC signaling messages.
          </p>
          <p>
            <strong>Step 2:</strong> Get the list of available peer IDs to
            establish peer connections.
          </p>
          <p>
            <strong>Step 3:</strong> Once a WebRTC peer connection is
            established, you can disconnect from the signaling server.
          </p>
        </div>
      </div>

      <hr />

      <div style={{ display: "flex", justifyContent: "space-between" }}>
        <div style={{ flex: "1", paddingRight: "20px" }}>
          <h2>WebRTC Peer Connections</h2>
          <div>
            <input
              type="text"
              id="peerIDInput"
              placeholder="Enter Peer ID"
              style={{ width: "80%" }}
            />
            <button onClick={()=>window.newPeerConnection()}>
              New Peer Connection
            </button>
          </div>
          <div style={{ paddingTop: "30px" }}>
            <button onClick={()=>window.sendToAll()}>
              Send to all peer connections
            </button>
          </div>
          <h3>Peer Connections:</h3>
          <div id="peerConnections"></div>
        </div>

        <div
          style={{
            flex: "1",
            paddingLeft: "20px",
            backgroundColor: "#f9f9f9",
            border: "1px solid #ddd",
            padding: "10px",
          }}
        >
          <h3>Notes</h3>
          <p>
            <strong>Note 1:</strong> This client is designed for WebRTC
            peer-to-peer communication.
          </p>
          <p>
            <strong>Note 2:</strong> Peer connections are established via a
            signaling server. Make sure the signaling server is running and
            accessible.
          </p>
          <p>
            <strong>Note 3:</strong> Use the message input field to send data
            messages between peers once the connection is established.
          </p>
          <p>
            <strong>Note 4:</strong> If you face any issues with the connection,
            check the browser's console for error logs.
          </p>
        </div>
      </div>

      <div
        style={{ backgroundColor: "#eee", padding: "10px", marginTop: "20px" }}
      >
        <h2>Send Message via Peer Connection</h2>
        <textarea
          id="message"
          rows={5}
          cols={50}
          placeholder="Message to be sent via one of the WebRTC peer connections"
          style={{ width: "100%", resize: "vertical" }}
        ></textarea>
      </div>

      <div>
        <h3>Logs</h3>
        <button onClick={()=>window.clearLog()}>Clear Logs</button>
        <div id="logArea" style={{ marginTop: "10px" }}></div>
      </div>
    </div>
  );
}
export default App;

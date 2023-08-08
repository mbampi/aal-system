import { useEffect, useState } from 'react';
import Events from './Findings';
import Charts from './Observations';
import './App.css';
import { TFinding, TObservation, TPackage } from './Types';

function App() {
  const [findings, setFindings] = useState<TFinding[]>([]);
  const [observations, setObservations] = useState<TObservation[]>([]);
  let running = false;

  useEffect(() => {
    if (running) {
      const ws = new WebSocket("ws://localhost:8080/ws");

      ws.onopen = (event) => {
        console.log("WebSocket is open now.", event);
      }

      ws.onmessage = (event) => {
        try {
          console.log("Received package: ", event.data);
          const pkg = JSON.parse(event.data) as TPackage;

          switch (pkg.type) {
            case "finding":
              let finding = pkg.data as TFinding;
              setFindings((prevEvents) => [...prevEvents, finding]);
              break;
            case "observation":
              let observation = pkg.data as TObservation;
              setObservations((prevEvents) => [...prevEvents, observation]);
              break;
            default:
              console.error("Unknown package type: ", pkg.type);
          }

        } catch (error) {
          console.error("Error parsing event data:", error);
        }
      };

      ws.onclose = (event) => {
        console.log("WebSocket closed.", event);
      };

      ws.onerror = (error) => {
        console.error("WebSocket error: ", error);
      };

      // Close the WebSocket connection when the component unmounts
      return () => {
        console.log("unmounting...");
        ws.close()
      };
    } else {
      running = true;
    }
  }, []);

  return (
    <div className="App">
      <div className="events-container">
        <Events findings={findings} />
      </div>
      <div className="charts-container">
        <Charts observations={observations} />
      </div>
    </div>
  );
}

export default App;

import React, { useEffect, useState } from 'react';
import Events from './Events';
import Charts from './Charts';
import './App.css';
import EventData from './EventData';



function App() {
  const [events, setEvents] = useState<EventData[]>([]);
  let running = false;

  useEffect(() => {
    if (running) {
      const ws = new WebSocket("ws://localhost:8080/ws");

      ws.onopen = (event) => {
        console.log("WebSocket is open now.", event);
      }

      ws.onmessage = (event) => {
        try {
          console.log("Received event: ", event.data);
          const eventData = JSON.parse(event.data) as EventData;
          setEvents((prevEvents) => [...prevEvents, eventData]);
        } catch (error) {
          console.error("Error parsing event data:", error);
        }
      };

      ws.onclose = (event) => {
        console.log("WebSocket is closed now.", event);
      };

      ws.onerror = (error) => {
        console.error("WebSocket encountered error: ", error);
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
        <Events events={events} />
      </div>
      <div className="charts-container">
        <Charts data={events} />
      </div>
    </div>
  );
}

export default App;

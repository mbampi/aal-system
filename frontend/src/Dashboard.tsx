import React, { useEffect, useState } from "react";
import "./Dashboard.css";

type EventData = {
    name: string;
    patient: string;
    sensor: string;
    value: string;
    timestamp: string;
};

const Dashboard: React.FC = () => {
    const [events, setEvents] = useState<EventData[]>([]);

    useEffect(() => {
        const ws = new WebSocket("ws://localhost:8080/ws");

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
        return () => ws.close();
    }, []);

    return (
        <div>
            <h1>Dashboard</h1>
            {events.length === 0 ? (
                <div>Waiting for events...</div>
            ) : (
                events.map((event, index) => (
                    <div key={index} className="event-card">
                        <p className="name">{event.name}</p>
                        <p className="value">{event.value}</p>
                        <p className="label">Patient:</p> <p className="patient">{event.patient}</p>
                        <p className="label">Sensor:</p> <p className="sensor">{event.sensor}</p>
                        <p className="timestamp">{event.timestamp}</p>
                    </div>

                ))
            )}
        </div>
    );
};

export default Dashboard;

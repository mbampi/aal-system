import React, { useEffect, useState } from "react";
import "./Events.css";
import EventData from "./EventData";


function Events({ events }: { events: EventData[] }) {
    return (
        <div>
            <h2>Events</h2>
            {events.length === 0 ? (
                <div>Waiting for events...</div>
            ) : (
                events.map((event, index) => (
                    <div key={index} className="event-card">
                        <p className="name">{event.name}</p>
                        <p className="value">{event.value}</p>
                        <div className="card-field">
                            <p className="label">Patient:</p>
                            <p className="patient">{event.patient}</p>
                        </div>
                        <p></p>
                        <div className="card-field">
                            <p className="label">Sensor:</p>
                            <p className="sensor">{event.sensor}</p>
                        </div>
                        <p></p>
                        <p className="timestamp">{event.timestamp}</p>
                    </div>

                ))
            )}
        </div>
    );
}

export default Events;

import React from "react";
import "./Findings.css";
import { TFinding } from "../Types";

interface FindingsProps {
    findings: TFinding[];
}

const Findings: React.FC<FindingsProps> = ({ findings }) => {
    return (
        <div>
            <h2>Findings</h2>
            {findings.length === 0 ? (
                <div>Waiting for events...</div>
            ) : (
                findings.map((finding: TFinding, index: number) => (
                    <div key={index} className="event-card">
                        <p className="name">{finding.name}</p>
                        <p className="value">{finding.value}</p>
                        <div className="card-field">
                            <p className="label">Patient:</p>
                            <p className="patient">{finding.patient}</p>
                        </div>
                        <p></p>
                        <div className="card-field">
                            <p className="label">Sensor:</p>
                            <p className="sensor">{finding.sensor}</p>
                        </div>
                        <p></p>
                        <p className="timestamp">{finding.timestamp}</p>
                    </div>

                ))
            )}
        </div>
    );
}

export default Findings;

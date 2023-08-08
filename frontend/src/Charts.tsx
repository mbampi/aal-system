import React, { useState } from 'react';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip } from 'recharts';
import { TObservation } from './Types';
import './Charts.css';

interface ChartsProps {
    observations: TObservation[];
}

const Charts: React.FC<ChartsProps> = ({ observations }) => {
    const [expandedProperty, setExpandedProperty] = useState<null | string>(null);

    const observationTypes = Array.from(new Set(observations.map((obs: TObservation) => obs.name)));

    return (
        <div>
            <h2>Observations</h2>

            {observationTypes.map((type) => {
                const filteredObservations = observations.filter((obs: TObservation) => obs.name === type);
                return <PropertyObservations
                    property={type}
                    observations={filteredObservations}
                    expanded={expandedProperty === type}
                    onToggleExpand={() => setExpandedProperty(type === expandedProperty ? null : type)}
                />;
            })}
        </div>
    );
}

interface PropertyObservationsProps {
    property: string;
    observations: TObservation[];
    expanded: boolean;
    onToggleExpand: () => void;
}

const PropertyObservations: React.FC<PropertyObservationsProps> = ({ property, observations, expanded, onToggleExpand }) => {
    const chartData = observations.map((obs: TObservation) => ({
        timestamp: new Date(obs.timestamp).toLocaleTimeString(),
        value: parseFloat(obs.value)
    }));

    return <div key={property} className='property-obs' onClick={onToggleExpand}>
        <h2 >{property}</h2>
        {expanded && <>

            <div className="info-container">
                <p><strong>Sensor:</strong> {observations[0].sensor} </p>
                <p><strong>Number of observations:</strong> {observations.length}</p>
                <p><strong>Min value:</strong> {Math.min(...observations.map((obs: TObservation) => parseFloat(obs.value)))}</p>
                <p><strong>Max value:</strong> {Math.max(...observations.map((obs: TObservation) => parseFloat(obs.value)))}</p>
            </div>

            <LineChart
                width={700}
                height={300}
                data={chartData}
                margin={{ top: 5, right: 30, left: 20, bottom: 5 }}
            >
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="timestamp" />
                <YAxis dataKey="value" />
                <Tooltip />
                <Line type="monotone" dataKey="value" stroke="#82ca9d" />
            </LineChart>

            {/* List with all observations of current type */}
            <ObservationsTable observations={observations} />
        </>}
    </div>
}

const ObservationsTable: React.FC<ChartsProps> = ({ observations }) => {
    return (
        <div>
            <table className='observations-table'>
                <thead>
                    <tr>
                        <th>Date</th>
                        <th>Time</th>
                        <th>Value</th>
                    </tr>
                </thead>
                <tbody>
                    {observations.map((obs: TObservation) => (
                        <tr key={obs.timestamp}>
                            <td>{new Date(obs.timestamp).toLocaleDateString()}</td>
                            <td>{new Date(obs.timestamp).toLocaleTimeString()}</td>
                            <td>{obs.value}</td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    );
}

export default Charts;

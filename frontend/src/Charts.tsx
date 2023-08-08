import React from 'react';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip } from 'recharts';
import { TObservation } from './Types';
import './Charts.css';

interface ChartsProps {
    observations: TObservation[];
}

const Charts: React.FC<ChartsProps> = ({ observations }) => {
    const observationTypes = Array.from(new Set(observations.map((obs: TObservation) => obs.name)));

    return (
        <div>
            <h2>Observations</h2>

            {observationTypes.map((type) => {
                const filteredObservations = observations.filter((obs: TObservation) => obs.name === type);
                return <PropertyObservations
                    property={type}
                    observations={filteredObservations}
                />;
            })}
        </div>
    );
}

interface PropertyObservationsProps {
    property: string;
    observations: TObservation[];
}

const PropertyObservations: React.FC<PropertyObservationsProps> = ({ property, observations }) => {
    const chartData = observations.map((obs: TObservation) => ({
        timestamp: new Date(obs.timestamp).toLocaleTimeString(),
        value: parseFloat(obs.value)
    }));

    return <div key={property} className='property-header'>

        <h2>{property}</h2>
        <p>Sensor: {observations[0].sensor} </p>
        <p>Number of observations: {observations.length}</p>
        <p>Min value: {Math.min(...observations.map((obs: TObservation) => parseFloat(obs.value)))}</p>
        <p>Max value: {Math.max(...observations.map((obs: TObservation) => parseFloat(obs.value)))}</p>

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

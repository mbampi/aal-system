import React from 'react';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip } from 'recharts';
import { TObservation } from './Types';
import './Charts.css';

interface ChartsProps {
    observations: TObservation[];
}

const Charts: React.FC<ChartsProps> = ({ observations }) => {
    // Convert the data to the format needed by recharts
    const chartData = observations.map((obs: TObservation) => ({
        timestamp: new Date(obs.timestamp).toLocaleTimeString(),
        value: parseFloat(obs.value)
    }));

    return (
        <div>
            <h2>Observations</h2>

            <h4>Heart Rate</h4>

            <LineChart
                width={500}
                height={300}
                data={chartData}
                margin={{ top: 5, right: 30, left: 20, bottom: 5 }}
            >
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="timestamp" />
                <YAxis />
                <Tooltip />
                <Line type="monotone" dataKey="value" stroke="#82ca9d" />
            </LineChart>

            {/* List with all obvervations */}
            <table className="observations-table">
                <thead>
                    <tr>
                        <th>Date</th>
                        <th>Time</th>
                        <th>Value</th>
                    </tr>
                </thead>
                <tbody>
                    {observations.map((obs: TObservation, index: number) => (
                        <tr key={index}>
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

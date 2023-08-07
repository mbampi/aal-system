import React from 'react';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip } from 'recharts';
import EventData from './EventData';

interface ChartsProps {
    data: EventData[];
}

const Charts: React.FC<ChartsProps> = ({ data }) => {
    // Convert the data to the format needed by recharts
    const chartData = data.map(event => ({
        timestamp: new Date(event.timestamp).toLocaleTimeString(),
        value: parseFloat(event.value)
    }));

    return (
        <div>
            <h2>Charts</h2>

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
        </div>
    );
}

export default Charts;

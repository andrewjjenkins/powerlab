import React from 'react';
import {
  useQuery
} from '@tanstack/react-query'
import { ServerSensorReadings, ServersResponse } from './generated/api';

function Servers() {
  const { isPending, error, data } = useQuery({
    queryKey: ['repoData'],
    queryFn: () =>
      fetch('/api/servers').then((res) =>
        res.json() as Promise<ServersResponse>
      ).then((servers) => {
        return Promise.all(servers.map((server, i) => {
          const sensorUrl = `/api/servers/${server.name}/sensors`;
          return fetch(sensorUrl).then((res) =>
            res.json() as Promise<ServerSensorReadings>).then((readings) => {
              return {
                name: server.name,
                ...readings,
              };
            })
        }))
      }),
  })

  if (isPending) return (
    <div>Loading...</div>
  );

  if (error) return (
    <div>An error has occurred: {error.message}</div>
  );

    return (
        <div className="overflow-x-auto">
            <table className="table">
                <thead>
                <tr>
                    <th></th>
                    <th>Name</th>
                    <th>CPU Temp</th>
                    <th>Chassis Temp</th>
                    <th>Fan Speed</th>
                    <th>Power (W)</th>
                </tr>
                </thead>
                <tbody>
                {data.map((server, index) => {
                    return (
                        <tr key={index}>
                            <th>{index}</th>
                            <td>{server.name}</td>
                            <td>{server.CpuTemp}</td>
                            <td>{server.ChassisTemp}</td>
                            <td>{server.FanSpeed}</td>
                            <td>{server.PowerWatts}</td>
                        </tr>
                    );
                })
                }
                </tbody>
            </table>
            </div>
    );
}

export default Servers;

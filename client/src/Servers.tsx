import React from 'react';
import {
  useQuery
} from '@tanstack/react-query'
import { ServersResponse } from './generated/api';
import { Server } from 'http';

function Servers() {
  const { isPending, error, data } = useQuery({
    queryKey: ['repoData'],
    queryFn: () =>
      fetch('http://shaftoe.lan:8080/api/servers').then((res) =>
        res.json() as Promise<ServersResponse>
      ),
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
                    <th>Power Status</th>
                    <th>Power (W)</th>
                </tr>
                </thead>
                <tbody>
                {data.map((server, index) => {
                    return (
                        <tr key={index}>
                            <th>{index}</th>
                            <td>{server.name}</td>
                            <td>{server.power_status}</td>
                            <td>{server.power_watts}</td>
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

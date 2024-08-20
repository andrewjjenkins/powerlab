import React from 'react';
import Servers from './Servers';
import './App.css';

function App() {
  return (
    <div className="min-h-full">
      <header className="bg-white shadow">
        <div className="mx-auto max-w-7xl px-4 py-6 sm:px-6 lg:px-8">
          <h1 className="text-3xl font-bold tracking-tight text-gray-900">Powerlab Servers</h1>
        </div>
      </header>
      <main>
        <div className="mx-auto max-w-7xl px-4 py-6 sm:px-6 lg:px-8">
          {<Servers></Servers>}
        </div>
      </main>
    </div>
  );
}


export default App;

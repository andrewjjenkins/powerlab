import React from 'react';
import {
    StopIcon,
    PlayIcon,
} from '@heroicons/react/24/solid'

interface PowerControllerProps {
    powerState: string;
}

function PowerController(props: PowerControllerProps) {
    const onStartClick = function() {
        console.log("Start clicked");
    }
    const onStopClick = function() {
        console.log("Stop clicked");
    }
    return (
        <div className="flex justify-center overflow-x-hidden join">
            <button className="btn btn-sm join-item btn-outline btn-success max-h-[2em]" onClick={onStartClick}>
                <div className="flex flex-row gap-x-1">
                    <PlayIcon className="min-h-[1em]"/>
                    Start
                </div>
            </button>
            <button className="btn btn-sm join-item btn-error" onClick={onStopClick}>
                <div className="flex flex-row gap-x-1">
                    <StopIcon className="min-h-[1em]"/>
                    Stop
                </div>
            </button>
        </div>
    );
}

export default PowerController;

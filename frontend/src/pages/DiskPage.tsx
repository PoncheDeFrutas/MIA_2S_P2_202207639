import Disk from "../components/Disk";
import React, { useState, useEffect } from "react";
import {GET} from "../services/API.ts";

type DiskData = { id: string; name: string; };

const FileExplorer: React.FC = () => {
    const [disksData, setDisksData] = useState<DiskData[]>([]);

    useEffect(() => {
        (async () => {
            try {
                const response = await GET<{ result: DiskData[] }>("disks");
                if (!response.result) {
                    setDisksData([]);
                } else {
                    setDisksData(response.result);
                }
            } catch (error) {
                alert(`Error: ${error}`);
            }
        })();
    }, []);


    return (
        <div className="flex-grow flex items-center justify-center p-44">
            <div className="w-full max-w-3xl p-8 bg-white rounded-lg shadow-md">
                <h2 className="text-2xl font-bold mb-4 text-gray-800">Disks</h2>
                <div className="flex flex-wrap gap-4">
                    {
                        disksData.length > 0 ? (
                        disksData.map((disk) => (
                            <Disk key={disk.id} id={disk.id} name={disk.name} />
                        ))
                        ) : (
                            <p className="text-gray-500"> No disks found</p>
                        )
                    }
                </div>
            </div>
        </div>
    );
}

export default FileExplorer;
